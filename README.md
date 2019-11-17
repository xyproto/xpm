# xpm [![Build Status](https://travis-ci.com/xyproto/xpm.svg?branch=master)](https://travis-ci.com/xyproto/xpm) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/xpm)](https://goreportcard.com/report/github.com/xyproto/xpm) [![GoDoc](https://godoc.org/github.com/xyproto/xpm?status.svg)](https://godoc.org/github.com/xyproto/xpm)

Encode images to the X PixMap (XPM3) image format.

The resulting images are smaller than the one from GIMP, since the question mark character is also used, while at the same time avoiding double question marks, which could result in a trigraph (like `??=`, which has special meaning in C).


Includes a `png2xpm` utility.

## Example use

Converting from a PNG to an XPM file:

```go
// Create a new XPM encoder
enc := xpm.NewEncoder(imageName)

// Open the PNG file
f, err := os.Open(inputFilename)
if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s\n", err)
    os.Exit(1)
}
m, err := png.Decode(f)
if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s\n", err)
    os.Exit(1)
}
f.Close()

// Prepare to output the XPM data to either stdout or to file
if outputFilename == "-" {
    f = os.Stdout
} else {
    f, err = os.Create(outputFilename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s\n", err)
        os.Exit(1)
    }
    defer f.Close()
}

// Generate and output the XPM data
err = enc.Encode(f, m)
if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s\n", err)
    os.Exit(1)
}
```

## General info

* Version: 2.1.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
