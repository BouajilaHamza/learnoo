package media

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
	"testing"
)

// WebP fixtures: a 4x3 gradient PNG generated in Go (see testImage) converted
// with libwebp's cwebp 1.6.0:
//
//	cwebp -quiet -lossless -metadata none in.png -o ll.webp
//	cwebp -quiet -q 80 -metadata none in.png -o lossy.webp
const (
	webpLosslessHex = "5249464638000000574542505650384c2b0000002f038000005f40906d33c422dcdfe30ec333c8f3df830cb05c00820012a51452186e8515468de87f50970300"
	webpLossyHex    = "524946465a00000057454250565038204e0000003002009d012a0400030001402625b002744c807e8001d9a493c400fec2ad3c6da06240fff26b1fa7d2200fd774f0d6ebf36dd31e1ef53cdfb2ffbfff2bb53eb7b3fb6d51fd3ba0c7dfff92ee0000"
)

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func testImage(w, h int) *image.NRGBA {
	m := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			m.Set(x, y, color.NRGBA{uint8(60 * x), uint8(80 * y), 200, 255})
		}
	}
	return m
}

func pngBytes(t *testing.T, w, h int) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := png.Encode(&b, testImage(w, h)); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func jpegBytes(t *testing.T, w, h int) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := jpeg.Encode(&b, testImage(w, h), nil); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// pngChunk builds a PNG chunk with a valid CRC.
func pngChunk(typ string, data []byte) []byte {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.BigEndian, uint32(len(data)))
	b.WriteString(typ)
	b.Write(data)
	_ = binary.Write(&b, binary.BigEndian, crc32.ChecksumIEEE(append([]byte(typ), data...)))
	return b.Bytes()
}

// pngHeaderOnly is a PNG signature and IHDR claiming w x h with no pixel data:
// DecodeConfig succeeds, a full decode would fail.
func pngHeaderOnly(w, h uint32) []byte {
	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:], w)
	binary.BigEndian.PutUint32(ihdr[4:], h)
	ihdr[8] = 8 // bit depth
	ihdr[9] = 6 // RGBA
	return append([]byte("\x89PNG\r\n\x1a\n"), pngChunk("IHDR", ihdr)...)
}

// insertBeforeIEND inserts extra chunks before the final IEND chunk.
func insertBeforeIEND(p []byte, chunks ...[]byte) []byte {
	cut := len(p) - 12
	out := append([]byte{}, p[:cut]...)
	for _, c := range chunks {
		out = append(out, c...)
	}
	return append(out, p[cut:]...)
}

// jpegWithDims patches the SOF0 frame header of a baseline JPEG.
func jpegWithDims(t *testing.T, w, h uint16) []byte {
	t.Helper()
	b := jpegBytes(t, 8, 8)
	i := bytes.Index(b, []byte{0xFF, 0xC0})
	if i < 0 {
		t.Fatal("no SOF0 marker")
	}
	binary.BigEndian.PutUint16(b[i+5:], h)
	binary.BigEndian.PutUint16(b[i+7:], w)
	return b
}

func TestProcessAccepts(t *testing.T) {
	tests := []struct {
		name, declared, source, stored string
		data                           []byte
		w, h                           int
	}{
		{"png", "image/png", ContentTypePNG, ContentTypePNG, pngBytes(t, 4, 3), 4, 3},
		{"png no declared type", "", ContentTypePNG, ContentTypePNG, pngBytes(t, 4, 3), 4, 3},
		{"jpeg", "image/jpeg", ContentTypeJPEG, ContentTypeJPEG, jpegBytes(t, 40, 30), 40, 30},
		{"jpeg declared with params", "Image/JPEG; charset=binary", ContentTypeJPEG, ContentTypeJPEG, jpegBytes(t, 2, 2), 2, 2},
		{"webp lossless", "image/webp", ContentTypeWebP, ContentTypePNG, mustHex(t, webpLosslessHex), 4, 3},
		{"webp lossy", "image/webp", ContentTypeWebP, ContentTypePNG, mustHex(t, webpLossyHex), 4, 3},
		{"max side", "", ContentTypePNG, ContentTypePNG, pngBytes(t, MaxSide, 1), MaxSide, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, err := Process(bytes.NewReader(tc.data), tc.declared)
			if err != nil {
				t.Fatalf("Process: %v", err)
			}
			if a.ContentType != tc.stored || a.SourceType != tc.source || a.Width != tc.w || a.Height != tc.h {
				t.Fatalf("got %s from %s %dx%d", a.ContentType, a.SourceType, a.Width, a.Height)
			}
			if a.ByteSize != int64(len(a.Data)) || len(a.SHA256) != 64 || !ValidKey(a.Key()) {
				t.Fatalf("inconsistent asset: size %d len %d key %q", a.ByteSize, len(a.Data), a.Key())
			}
			// Stored bytes are a valid image of the stored type.
			f, _ := sniff(a.Data)
			if f.contentType != tc.stored {
				t.Fatalf("stored bytes sniff as %q", f.contentType)
			}
			if _, err := f.decode(bytes.NewReader(a.Data)); err != nil {
				t.Fatalf("stored bytes do not decode: %v", err)
			}
		})
	}
}

func TestProcessDeterministic(t *testing.T) {
	in := pngBytes(t, 5, 5)
	a, err := Process(bytes.NewReader(in), "")
	if err != nil {
		t.Fatal(err)
	}
	b, err := Process(bytes.NewReader(append(append([]byte{}, in...), "trailing"...)), "")
	if err != nil {
		t.Fatal(err)
	}
	if a.SHA256 != b.SHA256 || !bytes.Equal(a.Data, b.Data) {
		t.Fatal("same pixels produced different canonical bytes")
	}
}

func TestProcessRejects(t *testing.T) {
	var gifBuf bytes.Buffer
	if err := gif.Encode(&gifBuf, testImage(2, 2), nil); err != nil {
		t.Fatal(err)
	}
	truncatedPNG := pngBytes(t, 64, 64)
	truncatedPNG = truncatedPNG[:len(truncatedPNG)/2]
	corruptJPEG := jpegBytes(t, 32, 32)
	for i := len(corruptJPEG) / 2; i < len(corruptJPEG)-2; i++ {
		corruptJPEG[i] = 0xFF
	}
	webpLossless := mustHex(t, webpLosslessHex)
	tooBig := append(pngBytes(t, 1, 1), make([]byte, MaxUploadBytes)...)

	tests := []struct {
		name, declared string
		data           []byte
		want           error
	}{
		{"empty", "", nil, ErrEmpty},
		{"svg", "", []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`), ErrUnsupportedType},
		{"svg with xml prolog", "image/svg+xml", []byte(`<?xml version="1.0"?><svg/>`), ErrUnsupportedType},
		{"gif", "image/gif", gifBuf.Bytes(), ErrUnsupportedType},
		{"bmp", "", append([]byte("BM"), make([]byte, 64)...), ErrUnsupportedType},
		{"html", "text/html", []byte("<!doctype html><html><script>alert(1)</script></html>"), ErrUnsupportedType},
		{"text", "", []byte("just some text"), ErrUnsupportedType},
		{"png declared as jpeg", "image/jpeg", pngBytes(t, 2, 2), ErrTypeMismatch},
		{"png declared as svg", "image/svg+xml", pngBytes(t, 2, 2), ErrTypeMismatch},
		{"webp declared as png", "image/png", webpLossless, ErrTypeMismatch},
		{"jpeg declared garbage", "not a type", jpegBytes(t, 2, 2), ErrTypeMismatch},
		{"over 5 MiB", "image/png", tooBig, ErrTooLarge},
		{"png too wide (header only)", "", pngHeaderOnly(MaxSide+1, 1), ErrDimensions},
		{"png too tall (header only)", "", pngHeaderOnly(1, MaxSide+1), ErrDimensions},
		{"png too many pixels (header only)", "", pngHeaderOnly(MaxSide, MaxSide), ErrDimensions},
		{"png huge (header only)", "", pngHeaderOnly(100_000, 100_000), ErrDimensions},
		{"jpeg too wide (header only)", "", jpegWithDims(t, MaxSide+1, 8), ErrDimensions},
		{"jpeg too many pixels (header only)", "", jpegWithDims(t, 4000, 4001), ErrDimensions},
		{"jpeg zero height", "", jpegWithDims(t, 8, 0), ErrDimensions},
		{"png zero width", "", pngHeaderOnly(0, 8), ErrCorrupt},
		{"truncated png", "", truncatedPNG, ErrCorrupt},
		{"corrupt jpeg", "", corruptJPEG, ErrCorrupt},
		{"truncated webp", "", webpLossless[:30], ErrCorrupt},
		{"signature only png", "", []byte("\x89PNG\r\n\x1a\n"), ErrCorrupt},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, err := Process(bytes.NewReader(tc.data), tc.declared)
			if !errors.Is(err, tc.want) {
				t.Fatalf("got asset=%v err=%v, want %v", a != nil, err, tc.want)
			}
		})
	}
}

// The PNG decoder rejects a zero dimension itself (ErrCorrupt above); the JPEG
// header path reaches checkDimensions (ErrDimensions above).
func TestCheckDimensions(t *testing.T) {
	for _, d := range [][2]int{{0, 1}, {1, 0}, {-1, 5}, {MaxSide + 1, 1}, {MaxSide, MaxSide}} {
		if err := checkDimensions(d[0], d[1]); !errors.Is(err, ErrDimensions) {
			t.Errorf("%v: got %v", d, err)
		}
	}
	for _, d := range [][2]int{{1, 1}, {MaxSide, 1}, {4000, 4000}} {
		if err := checkDimensions(d[0], d[1]); err != nil {
			t.Errorf("%v: got %v", d, err)
		}
	}
}

// countingReader yields an endless stream of zero bytes and counts reads.
type countingReader struct{ n int64 }

func (r *countingReader) Read(p []byte) (int, error) {
	clear(p)
	r.n += int64(len(p))
	return len(p), nil
}

func TestProcessReadsAtMostLimit(t *testing.T) {
	r := &countingReader{}
	if _, err := Process(io.MultiReader(bytes.NewReader(pngBytes(t, 1, 1)), r), ""); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("got %v", err)
	}
	if r.n > MaxUploadBytes+1 {
		t.Fatalf("read %d bytes from an unbounded stream", r.n)
	}
}

func TestProcessExactlyAtLimit(t *testing.T) {
	in := pngBytes(t, 3, 3)
	in = append(in, make([]byte, MaxUploadBytes-len(in))...)
	a, err := Process(bytes.NewReader(in), "image/png")
	if err != nil {
		t.Fatalf("upload of exactly %d bytes rejected: %v", MaxUploadBytes, err)
	}
	if a.ByteSize > 1024 {
		t.Fatalf("padding kept: %d bytes stored", a.ByteSize)
	}
}

func TestProcessStripsPolyglotAndMetadata(t *testing.T) {
	const script = "<script>alert(document.cookie)</script>"
	zipData := append([]byte("PK\x03\x04"), []byte("payload.html<html>")...)
	zipData = append(zipData, []byte("PK\x05\x06")...)

	pngWithTrailer := append(pngBytes(t, 4, 4), []byte(script)...)
	pngWithTrailer = append(pngWithTrailer, zipData...)

	pngWithChunks := insertBeforeIEND(pngBytes(t, 4, 4),
		pngChunk("tEXt", []byte("Comment\x00"+script)),
		pngChunk("eXIf", []byte("MM\x00*GPSSECRET")),
		pngChunk("iCCP", []byte("ICCSECRET\x00\x00garbage")),
	)

	exif := []byte("Exif\x00\x00MM\x00*GPSSECRET" + script)
	app1 := append([]byte{0xFF, 0xE1, 0, 0}, exif...)
	binary.BigEndian.PutUint16(app1[2:], uint16(len(exif)+2))
	jpg := jpegBytes(t, 16, 16)
	jpegWithExif := append(append(append([]byte{}, jpg[:2]...), app1...), jpg[2:]...)
	jpegWithExif = append(jpegWithExif, []byte(script)...)

	webpWithTrailer := append(mustHex(t, webpLosslessHex), []byte(script)...)

	tests := []struct {
		name   string
		data   []byte
		banned []string
	}{
		{"png trailing script and zip", pngWithTrailer, []string{script, "PK\x03\x04", "PK\x05\x06", "payload.html"}},
		{"png text exif iccp chunks", pngWithChunks, []string{script, "tEXt", "eXIf", "iCCP", "GPSSECRET", "ICCSECRET"}},
		{"jpeg exif app1 and trailer", jpegWithExif, []string{"Exif", "GPSSECRET", script, "\xFF\xE1"}},
		{"webp trailer", webpWithTrailer, []string{script, "RIFF", "WEBP"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, s := range tc.banned[:1] {
				if !bytes.Contains(tc.data, []byte(s)) {
					t.Fatalf("fixture lacks %q", s)
				}
			}
			a, err := Process(bytes.NewReader(tc.data), "")
			if err != nil {
				t.Fatalf("Process: %v", err)
			}
			for _, s := range tc.banned {
				if bytes.Contains(a.Data, []byte(s)) {
					t.Errorf("stored bytes contain %q", s)
				}
			}
			if strings.Contains(strings.ToLower(string(a.Data)), "<script") {
				t.Error("stored bytes contain a script tag")
			}
		})
	}
}
