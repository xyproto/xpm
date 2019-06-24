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
	// These are used when encoding the color ID as ASCII
	FirstLetter uint8 // default: 'a'
	LastLetter  uint8 // default: 'z'

	// The internal image name
	ImageName string

	// The alpha threshold
	AlphaThreshold float64
}

// hexify converts a slice of bytes to a slice of hex strings on the form 0x00
func hexify(data []byte) (r []string) {
	for _, b := range data {
		r = append(r, fmt.Sprintf("0x%02x", b))
	}
	return r
}

// c2hex converts from a color.Color to a XPM friendly color
// (either on the form #000000 or as a color, like "black")
// XPM only supports 100% or 0% alpha, represented as the None color
func c2hex(c color.Color, threshold float64) string {
	byteColor := color.NRGBAModel.Convert(c).(color.NRGBA)
	r, g, b, a := byteColor.R, byteColor.G, byteColor.B, byteColor.A
	if a < uint8(256.0*threshold) {
		return "None"
	}
	// "black" and "red" are shorter than the hex codes
	if r == 0 {
		if g == 0 {
			if b == 0 {
				return "black"
			} else if b == 0xff {
				return "blue"
			}
		} else if g == 0xff && b == 0 {
			return "green"
		}
	} else if r == 0xff {
		if g == 0xff && b == 0xff {
			return "white"
		} else if g == 0 && b == 0 {
			return "red"
		}
	}

	// return hex color code on the form #000000
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// inc will advance to the next string. Uses a-z. From "a" to "b", "zz" to "aaa" etc
func inc(s string, firstLetter, lastLetter uint8) string {
	if s == "" {
		return string(firstLetter)
	}
	if len(s) == 1 {
		if s[len(s)-1] != lastLetter { // one digit, not the last letter
			return string(s[0] + 1) // return the next letter
		}
		// one digit, and it is the last letter
		return string(firstLetter) + string(firstLetter)
	}
	if s[len(s)-1] == lastLetter { // two or more digits, the last digit is z
		return inc(s[:len(s)-1], firstLetter, lastLetter) + string(firstLetter) // increase next to last digit with one + "a"
	}
	// two or more digits, the last digit is not z
	return s[:len(s)-1] + inc(string(s[len(s)-1]), firstLetter, lastLetter) // first digit + last digit increases with one
}

// num2charcode converts a number to ascii letters, like 0 to "a", and 1 to "b".
// Can output multiple letters for higher numbers.
func num2charcode(num int, firstLetter, lastLetter uint8) string {
	// This is not the efficient way, but it's only called once per image conversion
	d := string(firstLetter)
	for i := 0; i < num; i++ {
		d = inc(d, firstLetter, lastLetter)
	}
	return d
}

// Encode will encode the given image as XBM, using a custom image name from
// the Encoder struct. The colors are first converted to grayscale, and then
// with a 50% cutoff they are converted to 1-bit colors.
func (enc *Encoder) Encode(w io.Writer, m image.Image) error {
	width := m.Bounds().Dx()
	height := m.Bounds().Dy()

	paletteMap := make(map[string]color.Color) // hexstring -> color, unordered
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			c := m.At(x, y)
			paletteMap[c2hex(c, enc.AlphaThreshold)] = c
		}
	}

	var paletteSlice []string // hexstrings, ordered
	for hexColor := range paletteMap {
		paletteSlice = append(paletteSlice, hexColor)
	}

	// Find the character code of the highest index
	highestCharCode := num2charcode(len(paletteSlice)-1, enc.FirstLetter, enc.LastLetter)
	charsPerPixel := len(highestCharCode)
	colors := len(paletteSlice)

	// Write the header, now that we know the right values
	fmt.Fprint(w, "/* XPM */\n")
	fmt.Fprintf(w, "static char * %s[] = {\n", enc.ImageName)
	fmt.Fprint(w, "/* <Values> */\n")
	fmt.Fprintf(w, "\"%d %d %d %d\",\n", width, height, colors, charsPerPixel)

	// Write the colors of paletteSlice, and generate a lookup table from hexColor to charcode
	fmt.Fprint(w, "/* <Colors> */\n")
	lookup := make(map[string]string) // hexcolor -> paletteindexchars, unordered
	charcode := strings.Repeat("a", charsPerPixel)
	for index, hexColor := range paletteSlice {
		fmt.Fprintf(w, "\"%s c %s\",\n", charcode, hexColor)
		lookup[hexColor] = charcode
		index++
		charcode = inc(charcode, enc.FirstLetter, enc.LastLetter)
	}

	// Now write the pixels, as character codes
	fmt.Fprint(w, "/* <Pixels> */\n")
	lastY := m.Bounds().Max.Y - 1
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		fmt.Fprintf(w, "\"")
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			c := m.At(x, y)
			charcode := lookup[c2hex(c, enc.AlphaThreshold)]
			// Now write the id code for the hex color to the file
			fmt.Fprint(w, charcode)
		}
		if y < lastY {
			fmt.Fprintf(w, "\",\n")
		} else {
			// Don't output a final comma
			fmt.Fprintf(w, "\"\n")
		}
	}

	fmt.Fprintf(w, "};\n")

	return nil
}

// Encode will encode the image as XBM, using "img" as the image name
func Encode(w io.Writer, m image.Image) error {
	e := &Encoder{'a', 'z', "img", 0.5}
	return e.Encode(w, m)
}
