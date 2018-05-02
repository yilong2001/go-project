package imgutils

import (
	"errors"
	//"gopkg.in/gographics/imagick.v2/imagick"
	//"github.com/Terry-Mao/paint"
	"github.com/Terry-Mao/paint/wand" //"math"
	"log"
)

// func ResizeImg(path, newpath string, destWid, destHeight uint, keepRatio bool) error {
// 	imagick.Initialize()
// 	defer imagick.Terminate()

// 	mw := imagick.NewMagickWand()

// 	err := mw.ReadImage(path)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	// Get original logo size
// 	width := mw.GetImageWidth()
// 	height := mw.GetImageHeight()

// 	// Calculate half the size
// 	wRatio := float64(width * 1.0 / destWid)
// 	hRatio := float64(height * 1.0 / destHeight)

// 	destWRatio := wRatio
// 	destHRatio := hRatio

// 	if keepRatio {
// 		if wRatio > hRatio {
// 			destWRatio = hRatio
// 			destHRatio = hRatio
// 		} else {
// 			destWRatio = wRatio
// 			destHRatio = wRatio
// 		}
// 	}

// 	if destWRatio < 1 || destHRatio < 1 {
// 		log.Println("current width or height is more little : ", destWRatio, destHRatio)
// 		return errors.New("current width or height is more little ")
// 	}

// 	hWidth := uint(float64(width) * destWRatio)
// 	hHeight := uint(float64(height) * destHRatio)

// 	// Resize the image using the Lanczos filter
// 	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
// 	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
// 	if err != nil {
// 		log.Println("ResizeImage failed : ", err)
// 		return err
// 	}

// 	// Set the compression quality to 95 (high quality = low compression)
// 	err = mw.SetImageCompressionQuality(95)
// 	if err != nil {
// 		log.Println("SetImageCompressionQuality", err)
// 		return err
// 	}

// 	err = mw.WriteImage(newpath)
// 	if err != nil {
// 		log.Println("SaveImage", err)
// 		return err
// 	}

// 	return nil
// }

func ResizeImg(path, newpath string, destWid, destHeight uint, keepRatio bool) error {
	wand.Genesis()
	defer wand.Terminus()
	w := wand.NewMagickWand()
	defer w.Destroy()

	if err := w.ReadImage(path); err != nil {
		log.Println("ReadImage", err)
		return err
	}

	// Get original logo size
	width := w.ImageWidth()
	height := w.ImageHeight()

	// Calculate half the size
	wRatio := float64(width * 1.0 / destWid)
	hRatio := float64(height * 1.0 / destHeight)

	destWRatio := wRatio
	destHRatio := hRatio

	if keepRatio {
		if wRatio > hRatio {
			destWRatio = hRatio
			destHRatio = hRatio
		} else {
			destWRatio = wRatio
			destHRatio = wRatio
		}
	}

	if destWRatio < 1 || destHRatio < 1 {
		log.Println("current width or height is more little : ", destWRatio, destHRatio)
		return errors.New("current width or height is more little ")
	}

	hWidth := uint(float64(width) * destWRatio)
	hHeight := uint(float64(height) * destHRatio)

	if err := w.ResizeImage(hWidth, hHeight, wand.GaussianFilter, 1.0); err != nil {
		log.Println("ResizeImage", err)
		return err
	}

	if err := w.WriteImage(newpath); err != nil {
		log.Println("WriteImage", err)
		return err
	}

	return nil
}
