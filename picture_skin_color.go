package main

import (
    "image"
    "image/color"
    "image/jpeg"
    "math"
    "os"

    "golang.org/x/image/draw"
)

// RGBColor RBG Color Type
type RGBColor struct {
    Red   uint8
    Green uint8
    Blue  uint8
}

// HSVColor HSV Color Type
type HSVColor struct {
    Hue        float64
    Saturation float64
    Value      float64
}

type YCbCrColor struct {
    Y  uint8
    Cb uint8
    Cr uint8
}

func main() {
    i, err := getImage("11.jpeg")
    if err != nil {
        return
    }

    width := i.Bounds().Dx()
    height := i.Bounds().Dy()

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            tmp := i.RGBAAt(x, y)
            c := RGBColor{
                Red:   tmp.R,
                Green: tmp.G,
                Blue:  tmp.B,
            }
            // Д╫©Г■╗hsvЕ┬╓Ф√╜Г ╝Х┌╓
            hc := rgbToHSV(c)
            if hc.Hue > 0 && hc.Hue < 35 && hc.Saturation > 0.23 && hc.Saturation < 0.68 {
                nc := color.RGBA{0, 0, 255, 255}
                i.SetRGBA(x, y, nc)
            }

            // Д╫©Г■╗YCbCrЕ┬╓Ф√╜Г ╝Х┌╓
            yc := rgbToYCbCr(c)
            // if (97.5 <= float64(yc.Cb) && float64(yc.Cb) <= 142.5) &&
            //  (134 <= float64(yc.Cr) && float64(yc.Cr) <= 176) {
            //  nc := color.RGBA{0, 0, 255, 255}
            //  i.SetRGBA(x, y, nc)
            // }
            if yc.Cb > 77 && yc.Cb < 127 && yc.Cr > 133 && yc.Cr < 173 {
                nc := color.RGBA{0, 0, 255, 255}
                i.SetRGBA(x, y, nc)
            }

            // Д╫©Г■╗rgbЕ┬╓Ф√╜Г ╝Х┌╓
            if c.Red > 95 && c.Green > 40 && c.Green < 100 &&
                c.Blue > 20 && max(c)-min(c) > 15 &&
                math.Abs(float64(c.Red-c.Green)) > 15 &&
                c.Red > c.Green && c.Red > c.Blue {
                nc := color.RGBA{0, 0, 255, 255}
                i.SetRGBA(x, y, nc)
            }

        }
    }
    writeImageToJpeg(i, "12.jpg")

}

func writeImageToJpeg(img image.Image, name string) error {
    fso, err := os.Create(name)
    if err != nil {
        return err
    }
    defer fso.Close()

    return jpeg.Encode(fso, img, &jpeg.Options{Quality: 100})
}

func getImage(input string) (*image.RGBA, error) {
    file, err := os.Open(input)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    src, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }
    img := toRGBA(src)

    return img, nil
}

func max(c RGBColor) float64 {
    max := math.Max(float64(c.Red), float64(c.Green))
    max = math.Max(max, float64(c.Blue))
    return max
}

func min(c RGBColor) float64 {
    min := math.Min(float64(c.Red), float64(c.Green))
    min = math.Min(min, float64(c.Blue))
    return min
}

// toRGBA converts an image.Image to an image.RGBA
func toRGBA(img image.Image) *image.RGBA {
    switch img.(type) {
    case *image.RGBA:
        return img.(*image.RGBA)
    }
    out := image.NewRGBA(img.Bounds())
    draw.Copy(out, image.Pt(0, 0), img, img.Bounds(), draw.Src, nil)
    return out
}

func rgbToHSV(c RGBColor) HSVColor {
    r := float64(c.Red) / 255
    g := float64(c.Green) / 255
    b := float64(c.Blue) / 255
    max := math.Max(r, g)
    max = math.Max(max, b)
    min := math.Min(r, g)
    min = math.Min(min, b)
    var h float64
    if r == max {
        h = (g - b) / (max - min)
    }

    if g == max {
        h = 2 + (b-r)/(max-min)
    }

    if b == max {
        h = 4 + (r-g)/(max-min)
    }

    h *= 60
    if h < 0 {
        h += 360
    }
    return HSVColor{h, (max - min) / max, max}

}

func rgbToYCbCr(c RGBColor) YCbCrColor {
    y, cb, cr := color.RGBToYCbCr(c.Red, c.Green, c.Blue)
    return YCbCrColor{y, cb, cr}
}
