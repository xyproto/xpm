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
	AllowedLetters []rune

	// The internal image name
	ImageName string

	// The alpha threshold
	AlphaThreshold float64

	// With comments?
	Comments bool
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
func inc(s string, allowedLetters []rune) string {
	firstLetter := allowedLetters[0]
	lastLetter := allowedLetters[len(allowedLetters)-1]
	if s == "" {
		return string(firstLetter)
	}
	var lastRuneOfString rune
	for i, r := range s {
		if i == len(s)-1 {
			lastRuneOfString = r
		}
	}
	if len(s) == 1 {
		if lastRuneOfString != lastLetter { // one digit, not the last letter
			// lastRuneOfString is the only rune in this string
			pos := strings.IndexRune(string(allowedLetters), lastRuneOfString)
			pos++
			if pos == len(allowedLetters) {
				pos = 0
			}
			return string(allowedLetters[pos]) // return the next letter
		}
		// one digit, and it is the last letter
		return string(firstLetter) + string(firstLetter)
	}
	if lastRuneOfString == lastLetter { // two or more digits, the last digit is z
		return inc(s[:len(s)-1], allowedLetters) + string(firstLetter) // increase next to last digit with one + "a"
	}
	// two or more digits, the last digit is not z
	return s[:len(s)-1] + inc(string(lastRuneOfString), allowedLetters) // first digit + last digit increases with one
}

// num2charcode converts a number to ascii letters, like 0 to "a", and 1 to "b".
// Can output multiple letters for higher numbers.
func num2charcode(num int, allowedLetters []rune) string {
	// This is not the efficient way, but it's only called once per image conversion
	d := string(allowedLetters[0])
	for i := 0; i < num; i++ {
		d = inc(d, allowedLetters)
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
	highestCharCode := num2charcode(len(paletteSlice)-1, enc.AllowedLetters)
	charsPerPixel := len(highestCharCode)
	colors := len(paletteSlice)

	// Write the header, now that we know the right values
	fmt.Fprint(w, "/* XPM */\n")
	fmt.Fprintf(w, "static char * %s[] = {\n", enc.ImageName)
	if enc.Comments {
		fmt.Fprint(w, "/* <Values> */\n")
	}
	fmt.Fprintf(w, "\"%d %d %d %d\",\n", width, height, colors, charsPerPixel)

	// Write the colors of paletteSlice, and generate a lookup table from hexColor to charcode
	if enc.Comments {
		fmt.Fprint(w, "/* <Colors> */\n")
	}
	lookup := make(map[string]string) // hexcolor -> paletteindexchars, unordered
	charcode := strings.Repeat(string(enc.AllowedLetters[0]), charsPerPixel)
	for index, hexColor := range paletteSlice {
		fmt.Fprintf(w, "\"%s c %s\",\n", charcode, hexColor)
		lookup[hexColor] = charcode
		index++
		charcode = inc(charcode, enc.AllowedLetters)
	}

	// Now write the pixels, as character codes
	if enc.Comments {
		fmt.Fprint(w, "/* <Pixels> */\n")
	}
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
	// 93 allowed ASCII letters is correct, ref. https://en.wikipedia.org/wiki/X_PixMap
	// these letters have been moved to the front, for prettier representation
	// of images with few colors, in text editors: *o+
	allowedLetters := []rune(" +*o!#$%&'(),-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnpqrstuvwxyz{|}~")
	e := &Encoder{allowedLetters, "img", 0.5, false}
	return e.Encode(w, m)
}
