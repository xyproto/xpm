package xpm

import (
	"fmt"
	"image"
	"image/color"
	"io"
	//	"strings"
)

// Encoder contains encoding configuration that is used by the Encode method
type Encoder struct {
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

func c2hex(c color.Color) string {
	return "#2080ff"
}

// lastIs checks if the last n letters of the given string is the given letter
func lastIs(s string, n int, letter byte) bool {
	if n >= len(s) {
		panic("n is too large")
	}
	for i := n; i >= len(s)-n; i-- {
		if s[i] != letter {
			return false
		}
	}
	return true
}

// inc will advance to the next string. Uses a-z. From "a" to "b", "zz" to "aaa" etc
func inc(s string) string {
	if s == "" {
		return "a"
	}
	if len(s) == 1 {
		if s[len(s)-1] != 'z' { // one digit, not z
			return string(s[0] + 1) // return the next letter
		}
		// one digit, and it is 'z'
		return "aa"
	}
	if s[len(s)-1] == 'z' { // two or more digits, the last digit is z
		return inc(s[:len(s)-1]) + "a" // increase next to last digit with one + "a"
	}
	// two or more digits, the last digit is not z
	return s[:len(s)-1] + inc(string(s[len(s)-1])) // first digit + last digit increases with one
}

// num2charcode converts a number to ascii letters, like 0 to "a", and 1 to "b".
// Can output multiple letters for higher numbers.
func num2charcode(num int) string {
	// This is not the efficient way, but it's only called once per image conversion
	d := "a"
	for i := 0; i < num; i++ {
		d = inc(d)
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
			paletteMap[c2hex(c)] = c
		}
	}

	var paletteSlice []string // hexstrings, ordered
	for hexColor := range paletteMap {
		paletteSlice = append(paletteSlice, hexColor)
	}

	// Find the character code of the highest index
	highestCharCode := num2charcode(len(paletteSlice) - 1)
	charsPerPixel := len(highestCharCode)
	colors := len(paletteSlice)

	// Write the header, now that we know the right values
	fmt.Fprint(w, "/* XPM */\n")
	fmt.Fprint(w, "static char * XFACE[] = {\n")
	fmt.Fprint(w, "/* <Values> */\n")
	fmt.Fprintf(w, "\"%d %d %d %d\",\n", width, height, colors, charsPerPixel)

	// Write the colors of paletteSlice, and generate a lookup table from hexColor to charcode
	fmt.Fprint(w, "/* <Colors> */\n")
	lookup := make(map[string]string) // hexcolor -> paletteindexchars, unordered
	charcode := ""
	for index, hexColor := range paletteSlice {
		charcode := inc(charcode)
		fmt.Fprintf(w, "\"%s c %s\",\n", charcode, hexColor)
		lookup[hexColor] = charcode
		index++
	}

	// Now write the pixels, as character codes
	fmt.Fprint(w, "/* <Pixels> */\n")
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		fmt.Fprintf(w, "\"")
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			c := m.At(x, y)
			charcode := lookup[c2hex(c)]
			fmt.Fprint(w, charcode)
			// Now write c2hex(c) to the file
		}
		fmt.Fprintf(w, "\",\n")
	}

	fmt.Fprintf(w, "};\n")

	return nil
}

// Encode will encode the image as XBM, using "img" as the image name
func Encode(w io.Writer, m image.Image) error {
	e := &Encoder{}
	return e.Encode(w, m)
}
