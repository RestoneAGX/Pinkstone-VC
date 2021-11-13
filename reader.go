package main

import (
	"fmt"
	"github.com/makiuchi-d/gozxing/oned"
	"image"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
)

var path = "barcodes [test]/"

func main() {
	// open and decode image file
	//path := "barcode.png"
	file, _ := os.Open(path + "test 1.png")
	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	barReader := oned.NewCodaBarReader()
	result, _ := barReader.Decode(bmp, nil)
	fmt.Println(result)
}