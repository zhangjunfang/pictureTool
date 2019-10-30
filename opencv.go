package main

import (
    "fmt"
    "image"
    "image/color"
    "os"

    "gocv.io/x/gocv"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("How to run:\n\t%s [imageFile]\n", os.Args[0])
        return
    }

    filename := os.Args[1]

    img := gocv.IMRead(filename, gocv.IMReadColor)
    gray := gocv.NewMat()
    defer gray.Close()

    // 1Ц─│Х╫╛Е▄√Ф┬░Г│╟Е╨╕Е⌡╬
    gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

    // 2Ц─│Е╫╒Ф─│Е╜╕Е▐≤Ф█╒Г └И╒└Е╓└Г░├О╪▄Е╬≈Е┬╟Е▐╞Д╩╔Ф÷╔Ф┴╬Г÷╘Е╫╒Г └Е⌡╬Г┴┤
    dilation := preprocess(gray)
    defer dilation.Close()

    // 3Ц─│Ф÷╔Ф┴╬Е▓▄Г╜⌡И─┴Ф√┤Е╜≈Е▄╨Е÷÷
    rects := findTextRegion(dilation)

    // 4Ц─│Г■╗Г╩©Г╨©Г■╩Е┤╨Х©≥Д╨⌡Ф┴╬Е┬╟Г └Х╫╝Е╩⌠
    for _, rect := range rects {
        tmpRects := make([][]image.Point, 0)
        tmpRects = append(tmpRects, rect)
        gocv.DrawContours(&img, tmpRects, 0, color.RGBA{0, 255, 0, 255}, 2)
    }

    // 5.Ф≤╬Г╓╨Е╦╕Х╫╝Е╩⌠Г └Е⌡╬Е┐▐
    gocv.IMWrite("imgDrawRect.jpg", img)
    // window := gocv.NewWindow("img")
    // for {
    //  window.IMShow(img)
    //  if window.WaitKey(1) >= 0 {
    //      break
    //  }
    // }
}

func preprocess(gray gocv.Mat) gocv.Mat {
    // 1. SobelГ╝≈Е╜░О╪▄xФ√╧Е░▒Ф╠┌Ф╒╞Е╨╕
    sobel := gocv.NewMat()
    defer sobel.Close()
    gocv.Sobel(gray, &sobel, int(gocv.MatTypeCV8U), 1, 0, 3, 1, 0, gocv.BorderDefault)

    // Д╨▄Е─╪Е▄√
    binary := gocv.NewMat()
    defer binary.Close()
    gocv.Threshold(sobel, &binary, 0, 255, gocv.ThresholdOtsu+gocv.ThresholdBinary)

    // 3. Х├╗Х┐─Е▓▄Х┘░Х ─Ф⌠█Д╫°Г └Ф═╦Е┤╫Ф∙╟
    // Е▐╞Д╩╔Х╟┐Х┼┌
    element1 := gocv.GetStructuringElement(gocv.MorphRect, image.Point{30, 9})
    // Е▐╞Д╩╔Х╟┐Х┼┌
    element2 := gocv.GetStructuringElement(gocv.MorphRect, image.Point{24, 4})

    // 4. Х├╗Х┐─Д╦─Ф╛║О╪▄Х╝╘Х╫╝Е╩⌠Г╙│Е┤╨
    dilation := gocv.NewMat()
    defer dilation.Close()
    gocv.Dilate(binary, &dilation, element2)

    // 5. Х┘░Х ─Д╦─Ф╛║О╪▄Е▌╩Ф▌┴Г╩├Х┼┌О╪▄Е╕┌Х║╗Ф═╪Г╨©Г╜┴Ц─┌ФЁ╗Ф└▐Х©≥И┤▄Е▌╩Ф▌┴Г └Ф≤╞Г╚√Г⌡╢Г └Г╨©
    erosion := gocv.NewMat()
    defer erosion.Close()
    gocv.Erode(dilation, &erosion, element1)

    // 6. Е├█Ф╛║Х├╗Х┐─О╪▄Х╝╘Х╫╝Е╩⌠Ф≤▌Ф≤╬Д╦─Д╨⌡
    dilation2 := gocv.NewMat()
    // defer dilation2.Close()
    gocv.Dilate(erosion, &dilation2, element2)

    // Е╜≤Е┌╗Д╦╜И≈╢Е⌡╬Г┴┤
    gocv.IMWrite("binary.png", binary)
    gocv.IMWrite("dilation.png", dilation)
    gocv.IMWrite("erosion.png", erosion)
    gocv.IMWrite("dilation2.png", dilation2)

    return dilation2
}

func findTextRegion(img gocv.Mat) [][]image.Point {
    // 1. Ф÷╔Ф┴╬Х╫╝Е╩⌠
    rects := make([][]image.Point, 0)
    contours := gocv.FindContours(img, gocv.RetrievalTree, gocv.ChainApproxSimple)

    for _, cnt := range contours {
        // Х╝║Г╝≈Х╞╔Х╫╝Е╩⌠Г └И²╒Г╖╞
        area := gocv.ContourArea(cnt)
        // И²╒Г╖╞Е╟▐Г └И┐╫Г╜⌡И─┴Ф▌┴
        // Е▐╞Д╩╔Х╟┐Х┼┌ 1000
        if area < 700 {
            continue
        }

        // Х╫╝Е╩⌠Х©▒Д╪╪О╪▄Д╫°Г■╗Е╬┬Е╟▐
        epsilon := 0.001 * gocv.ArcLength(cnt, true)
        // approx := gocv.ApproxPolyDP(cnt, epsilon, true)
        _ = gocv.ApproxPolyDP(cnt, epsilon, true)

        // Ф┴╬Е┬╟Ф°─Е╟▐Г÷╘Е╫╒О╪▄Х╞╔Г÷╘Е╫╒Е▐╞Х┐╫Ф°┴Ф√╧Е░▒
        rect := gocv.MinAreaRect(cnt)

        // Х╝║Г╝≈И╚≤Е▓▄Е╝╫
        mWidth := float64(rect.BoundingRect.Max.X - rect.BoundingRect.Min.X)
        mHeight := float64(rect.BoundingRect.Max.Y - rect.BoundingRect.Min.Y)

        // Г╜⌡И─┴И┌ёД╨⌡Е╓╙Г╩├Г └Г÷╘Е╫╒О╪▄Г∙≥Д╦▀Ф┴│Г └
        // Е▐╞Д╩╔Х╟┐Х┼┌ mHeight > (mWidth * 1.2)
        if mHeight > (mWidth * 0.9) {
            continue
        }

        // Г╛╕Е░┬Ф²║Д╩╤Г └rectФ╥╩Е┼═Е┬╟rectsИ⌡├Е░┬Д╦╜
        rects = append(rects, rect.Contour)
    }

    return rects
}
