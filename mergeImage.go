package main

import (
    "flag"
    "image"
    "image/draw"
    "image/png"
    "os"
    "path"
)

var (
    bg      = flag.String("bg", "bg.png", "背景图片")
    pt      = flag.String("pt", "pt.png", "前景图片")
    offsetX = flag.Int("offsetX", 0, "x轴偏移值")
    offsetY = flag.Int("offsetY", 0, "y轴偏移值")
    prefix  = flag.String("prefix", "test_", "文件名前缀")
)

func main() {
    flag.Parse()
    mergeImage(*pt)
}

func mergeImage(file string) {
    imgb, _ := os.Open(*bg)
    img, _ := png.Decode(imgb)
    defer imgb.Close()
    b1 := img.Bounds()

    wmb, _ := os.Open(file)
    watermark, _ := png.Decode(wmb)
    b2 := watermark.Bounds()
    defer wmb.Close()

    offset := image.Pt((b1.Max.X-b2.Max.X)/2+*offsetX, (b1.Max.Y-b2.Max.Y)/2+*offsetY)
    b := img.Bounds()
    m := image.NewRGBA(b)
    draw.Draw(m, b, img, image.ZP, draw.Src)
    draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

    imgw, _ := os.Create(*prefix + path.Base(file))
    png.Encode(imgw, m)
    defer imgw.Close()
}
