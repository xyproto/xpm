package xpm

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

// Encoder contains encoding configuration that is used by the Encode method
type Encoder struct {
	// The internal image name, when encoding the data. The default is "img".
	Name string

	// The threshold, from 0 to 1, for when grayscale colors should appear as
	// white or black. The default is 0.5.
	Threshold float64
}

// hexify converts a slice of bytes to a slice of hex strings on the form 0x00
func hexify(data []byte) (r []string) {
	for _, b := range data {
		hexdigits := fmt.Sprintf("%x", b)
		if len(hexdigits) == 1 {
			r = append(r, "0x0"+hexdigits)
		} else {
			r = append(r, "0x"+hexdigits)
		}
	}
	return r
}

// Encode will encode the given image as XBM, using a custom image name from
// the Encoder struct. The colors are first converted to grayscale, and then
// with a 50% cutoff they are converted to 1-bit colors.
func (enc *Encoder) Encode(w io.Writer, m image.Image) error {
    panic("not implemented yet")

	imageName := enc.Name

	width := m.Bounds().Dx()
	height := m.Bounds().Dy()

	fmt.Fprintf(w, "/* XBM X11 format */\n")

	maskIndex := 0
	masks := []uint8{
		0x1,
		0x2,
		0x4,
		0x8,
		0x10,
		0x20,
		0x40,
		0x80,
	}

	var pixels []byte
	var pixel uint8
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			c := m.At(x, y)
			grayColor := color.GrayModel.Convert(c).(color.Gray)
			value := grayColor.Y
			if value <= byte(float64(256)*enc.Threshold) {
				// white
				pixel |= masks[maskIndex]
			} else {
				// black, skip
			}
			// Prepare to write the next bit
			maskIndex++
			if maskIndex == len(masks) {
				// Filled up an entire byte with pixel bits, flush and reset
				maskIndex = 0
				pixels = append(pixels, pixel)
				pixel = 0
			}
		}
	}

	fmt.Fprintf(w, "#define %s_width %d\n", imageName, width)
	fmt.Fprintf(w, "#define %s_height %d\n", imageName, height)
	fmt.Fprintf(w, "static unsigned char %s_bits[] = {\n", imageName)
	fmt.Fprintf(w, "  %s\n", strings.Join(hexify(pixels), ", "))
	fmt.Fprintf(w, "};\n")

	return nil
}

// Encode will encode the image as XBM, using "img" as the image name
func Encode(w io.Writer, m image.Image) error {
	e := &Encoder{"img", 0.5}
	return e.Encode(w, m)
}
