package graphics

import (
	"errors"
	"image"
	"os"

	"github.com/nfnt/resize"
)

func ImageCopy(src image.Image, x, y, w, h int) (image.Image, error) {
	var tmpImg image.Image
	if rgbImg, ok := src.(*image.RGBA); ok {
		tmpImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		tmpImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.YCbCr); ok {
		tmpImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //裁剪x0 y0 x1 y1
	} else {

		return tmpImg, errors.New("解码失败")
	}

	return tmpImg, nil
}

func ImageCopyFromFile(p string, x, y, w, h int) (src image.Image, err error) {
	file, err := os.Open(p)
	defer file.Close()
	if err != nil {
		return src, err
	}
	src, _, err = image.Decode(file)

	return ImageCopy(src, x, y, w, h)
}

func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}
func ImageResizeSaveFile(src image.Image, width, height int, p string) error {
	dst := resize.Resize(uint(width), uint(height), src, resize.Lanczos3)
	return SaveImage(p, dst)
}
