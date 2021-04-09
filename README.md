# xpm

![Build Status](https://github.com/xyproto/xpm/workflows/Build/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/xpm)](https://goreportcard.com/report/github.com/xyproto/xpm) [![GoDoc](https://godoc.org/github.com/xyproto/xpm?status.svg)](https://godoc.org/github.com/xyproto/xpm) [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/xyproto/xpm/main/LICENSE)

Encode images to the X PixMap (XPM3) image format.

The resulting images are smaller than the one from GIMP, since the question mark character is also used, while at the same time avoiding double question marks, which could result in a trigraph (like `??=`, which has special meaning in C).


The `png2xpm` utility is included.

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

## Reference documentation

* The [XPM pdf](https://www.xfree86.org/current/xpm.pdf)

## General info

* Version: 2.2.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
