package media

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func orientationEXIF(order binary.ByteOrder, orientation uint16) []byte {
	b := make([]byte, 32)
	copy(b, "Exif\x00\x00")
	tiff := b[6:]
	copy(tiff, "MM")
	if order == binary.LittleEndian {
		copy(tiff, "II")
	}
	order.PutUint16(tiff[2:], 42)
	order.PutUint32(tiff[4:], 8)
	order.PutUint16(tiff[8:], 1)
	order.PutUint16(tiff[10:], 0x112)
	order.PutUint16(tiff[12:], 3)
	order.PutUint32(tiff[14:], 1)
	order.PutUint16(tiff[18:], orientation)
	return b
}

func TestOrientationPixels(t *testing.T) {
	// Non-zero source origin and six distinct pixels expose flips as well as rotations.
	src := image.NewRGBA(image.Rect(3, 4, 6, 6))
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			src.SetRGBA(x+3, y+4, color.RGBA{R: uint8(y*3 + x + 1), A: 255})
		}
	}
	cases := []struct {
		w, h   int
		pixels []uint8
	}{
		{3, 2, []uint8{1, 2, 3, 4, 5, 6}}, {3, 2, []uint8{3, 2, 1, 6, 5, 4}},
		{3, 2, []uint8{6, 5, 4, 3, 2, 1}}, {3, 2, []uint8{4, 5, 6, 1, 2, 3}},
		{2, 3, []uint8{1, 4, 2, 5, 3, 6}}, {2, 3, []uint8{4, 1, 5, 2, 6, 3}},
		{2, 3, []uint8{6, 3, 5, 2, 4, 1}}, {2, 3, []uint8{3, 6, 2, 5, 1, 4}},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprint(i+1), func(t *testing.T) {
			dst := applyOrientation(src, i+1)
			b := dst.Bounds()
			if b.Dx() != tc.w || b.Dy() != tc.h {
				t.Fatalf("bounds %v", b)
			}
			for y := 0; y < tc.h; y++ {
				for x := 0; x < tc.w; x++ {
					r, _, _, _ := dst.At(b.Min.X+x, b.Min.Y+y).RGBA()
					if uint8(r>>8) != tc.pixels[y*tc.w+x] {
						t.Fatalf("wrong pixel at %d,%d", x, y)
					}
				}
			}
		})
	}
}

func TestProcessExifOrientation(t *testing.T) {
	jpg := jpegBytes(t, 16, 32)
	source, err := jpeg.Decode(bytes.NewReader(jpg))
	if err != nil {
		t.Fatal(err)
	}
	for _, order := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		for orientation := 1; orientation <= 8; orientation++ {
			t.Run(fmt.Sprintf("%s/%d", order, orientation), func(t *testing.T) {
				input := insertJPEGSegment(jpg, 0xe1, orientationEXIF(order, uint16(orientation)))
				a, err := process(input, "image/jpeg")
				if err != nil {
					t.Fatal(err)
				}
				expected := applyOrientation(source, orientation)
				var encoded bytes.Buffer
				if err := jpeg.Encode(&encoded, expected, &jpeg.Options{Quality: JPEGQuality}); err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(a.Data(), encoded.Bytes()) {
					t.Fatal("stored pixels do not match expected orientation")
				}
				if a.Width() != expected.Bounds().Dx() || a.Height() != expected.Bounds().Dy() {
					t.Fatal("incorrect stored dimensions")
				}
				if bytes.Contains(a.Data(), []byte("Exif")) {
					t.Fatal("EXIF retained")
				}
				if err := checkJPEGEnd(a.Data()); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestMalformedOrientationIgnored(t *testing.T) {
	for _, order := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		valid := orientationEXIF(order, 6)[6:]
		for _, change := range []func([]byte){
			func(b []byte) { order.PutUint16(b[12:], 4) }, // LONG, not SHORT
			func(b []byte) { order.PutUint32(b[14:], 2) }, // multiple values
			func(b []byte) { order.PutUint32(b[4:], 0xffffffff) },
			func(b []byte) { order.PutUint16(b[18:], 9) },
		} {
			b := bytes.Clone(valid)
			change(b)
			if got := parseTIFFOrientation(b); got != 1 {
				t.Fatalf("invalid field yielded %d", got)
			}
		}
		for n := 0; n < 22; n++ {
			if got := parseTIFFOrientation(valid[:n]); got != 1 {
				t.Fatalf("truncated TIFF %d yielded %d", n, got)
			}
		}
	}
}

func TestProcessRejectsDisguisedPhoneTrailers(t *testing.T) {
	jpg := jpegBytes(t, 16, 16)
	for name, trailer := range map[string][]byte{
		"fake mp4":    []byte("\x00\x00\x00\x08ftyp<script>alert(1)</script>"),
		"fake jpeg":   []byte("\xff\xd8<script>alert(1)</script>\xff\xd9"),
		"second jpeg": jpegBytes(t, 4, 4),
		"padding":     make([]byte, 128),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := process(append(bytes.Clone(jpg), trailer...), "image/jpeg")
			if !errors.Is(err, ErrTrailingData) {
				t.Fatalf("got %v, want ErrTrailingData", err)
			}
		})
	}
}
