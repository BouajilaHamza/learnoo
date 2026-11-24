// Package media validates uploaded lesson images and stores them as
// immutable, content-addressed files.
//
// Process accepts PNG, JPEG and WebP only. The type is detected from magic
// bytes, never from a filename or client Content-Type. Size and dimension
// limits are checked before the full decode, and the decoded pixels are
// re-encoded so that stored bytes carry no trailing data, metadata (EXIF, GPS,
// text chunks, ICC profiles) or polyglot payloads. WebP is decoded with
// golang.org/x/image/webp and stored as lossless PNG because the standard
// library has no WebP encoder. See docs/DECISIONS.md ("Media upload limits").
package media

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"strings"

	"golang.org/x/image/webp"
)

// Limits. Changing them is a recorded decision, not a tuning knob.
const (
	// MaxUploadBytes is the largest accepted upload (5 MiB). Larger input is
	// rejected, never truncated.
	MaxUploadBytes = 5 << 20
	// MaxSide is the largest accepted width or height in pixels.
	MaxSide = 4096
	// MaxPixels is the largest accepted width*height (16 megapixels). Together
	// with MaxSide it bounds decode memory to about 64 MB of RGBA.
	MaxPixels = 16_000_000
	// JPEGQuality is the fixed quality used when re-encoding JPEG uploads.
	JPEGQuality = 90
)

// Content types. Uploads may be any of the three; stored assets are only
// ContentTypePNG or ContentTypeJPEG (WebP is stored as PNG).
const (
	ContentTypePNG  = "image/png"
	ContentTypeJPEG = "image/jpeg"
	ContentTypeWebP = "image/webp"
)

// Rejection reasons. Process wraps one of these, so callers map them with
// errors.Is (for example to HTTP 413 or 415 in the later API slice).
var (
	ErrEmpty           = errors.New("media: empty upload")
	ErrTooLarge        = errors.New("media: upload exceeds size limit")
	ErrUnsupportedType = errors.New("media: unsupported image type")
	ErrTypeMismatch    = errors.New("media: declared content type does not match image content")
	ErrDimensions      = errors.New("media: image dimensions outside limits")
	ErrCorrupt         = errors.New("media: image is corrupt or truncated")
)

// Asset is a validated, canonical image ready to be stored.
type Asset struct {
	// Data is the canonical re-encoded image. SHA256, ByteSize and the storage
	// key all describe these bytes, not the original upload.
	Data []byte
	// SHA256 is the lowercase hex SHA-256 of Data.
	SHA256 string
	// ContentType is the stored type: ContentTypePNG or ContentTypeJPEG.
	ContentType string
	// SourceType is the detected upload type (PNG, JPEG or WebP).
	SourceType string
	ByteSize   int64
	Width      int
	Height     int
}

// Key returns the storage key for the asset: "<sha256>.<ext>".
func (a *Asset) Key() string {
	return keyFor(a.SHA256, a.ContentType)
}

type format struct {
	contentType  string
	decodeConfig func(io.Reader) (image.Config, error)
	decode       func(io.Reader) (image.Image, error)
}

var (
	formatPNG  = format{ContentTypePNG, png.DecodeConfig, png.Decode}
	formatJPEG = format{ContentTypeJPEG, jpeg.DecodeConfig, jpeg.Decode}
	formatWebP = format{ContentTypeWebP, webp.DecodeConfig, webp.Decode}
)

// sniff identifies an accepted format from magic bytes only.
func sniff(b []byte) (format, bool) {
	switch {
	case bytes.HasPrefix(b, []byte("\x89PNG\r\n\x1a\n")):
		return formatPNG, true
	case bytes.HasPrefix(b, []byte{0xFF, 0xD8, 0xFF}):
		return formatJPEG, true
	case len(b) >= 12 && string(b[0:4]) == "RIFF" && string(b[8:12]) == "WEBP":
		return formatWebP, true
	}
	return format{}, false
}

// Process reads at most MaxUploadBytes from r, validates the image and returns
// its canonical form. declaredType is the client's Content-Type; when it is not
// empty it must name the detected type, otherwise it is ignored.
func Process(r io.Reader, declaredType string) (*Asset, error) {
	raw, err := io.ReadAll(io.LimitReader(r, MaxUploadBytes+1))
	if err != nil {
		return nil, fmt.Errorf("media: read upload: %w", err)
	}
	if len(raw) > MaxUploadBytes {
		return nil, fmt.Errorf("%w: more than %d bytes", ErrTooLarge, MaxUploadBytes)
	}
	if len(raw) == 0 {
		return nil, ErrEmpty
	}

	f, ok := sniff(raw)
	if !ok {
		return nil, fmt.Errorf("%w: only PNG, JPEG and WebP are accepted", ErrUnsupportedType)
	}
	if declaredType != "" {
		mt, _, err := mime.ParseMediaType(declaredType)
		if err != nil || !strings.EqualFold(mt, f.contentType) {
			return nil, fmt.Errorf("%w: declared %q, detected %s", ErrTypeMismatch, declaredType, f.contentType)
		}
	}

	// Header-only check before allocating pixels (decompression-bomb defence).
	cfg, err := f.decodeConfig(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorrupt, err)
	}
	if err := checkDimensions(cfg.Width, cfg.Height); err != nil {
		return nil, err
	}

	img, err := f.decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorrupt, err)
	}
	b := img.Bounds()
	if b.Dx() != cfg.Width || b.Dy() != cfg.Height {
		return nil, fmt.Errorf("%w: header and decoded dimensions differ", ErrCorrupt)
	}

	var out bytes.Buffer
	stored := ContentTypePNG
	switch f.contentType {
	case ContentTypeJPEG:
		stored = ContentTypeJPEG
		err = jpeg.Encode(&out, img, &jpeg.Options{Quality: JPEGQuality})
	default: // PNG, and WebP stored as lossless PNG.
		err = (&png.Encoder{CompressionLevel: png.DefaultCompression}).Encode(&out, img)
	}
	if err != nil {
		return nil, fmt.Errorf("media: re-encode: %w", err)
	}

	sum := sha256.Sum256(out.Bytes())
	return &Asset{
		Data:        out.Bytes(),
		SHA256:      hex.EncodeToString(sum[:]),
		ContentType: stored,
		SourceType:  f.contentType,
		ByteSize:    int64(out.Len()),
		Width:       b.Dx(),
		Height:      b.Dy(),
	}, nil
}

func checkDimensions(w, h int) error {
	if w <= 0 || h <= 0 {
		return fmt.Errorf("%w: %dx%d has a zero dimension", ErrDimensions, w, h)
	}
	if w > MaxSide || h > MaxSide {
		return fmt.Errorf("%w: %dx%d exceeds %d px per side", ErrDimensions, w, h, MaxSide)
	}
	if int64(w)*int64(h) > MaxPixels {
		return fmt.Errorf("%w: %dx%d exceeds %d pixels", ErrDimensions, w, h, MaxPixels)
	}
	return nil
}
