package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

func main() {
	//rand.Seed(time.Now().UnixNano())
	text := "hi"
	binaryText := TextToBinary(text)
	fmt.Println(binaryText)

	file, err := os.Open("input3.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Декодируем изображение
	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	// Создаем новое изображение с теми же размерами
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}
	//fmt.Println("min: ", bounds.Min.Y, bounds.Min.X)

	//fmt.Println(len(binaryText), bounds.Dx()*bounds.Dy())
	//fmt.Println(getRandomPixels(len(binaryText), bounds.Dx()*bounds.Dy()))

	randomPixels := getRandomPixels(len(binaryText), bounds.Dx()*bounds.Dy())

	sort.Ints(randomPixels)

	//bounds image.Rectangle, newImg *image.RGBA, img image.Image, randomPixels []int, binaryText []string
	newImg = encrypt(bounds, newImg, img, randomPixels, binaryText)

	// Сохраняем измененное изображение
	//bounds image.Rectangle, newImg *image.RGBA, randomPixels []int, binaryText []string
	decrypt(bounds, newImg, randomPixels)
	outFile, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	jpeg.Encode(outFile, newImg, nil)
}

func TextToBinary(s string) []string {
	array := make([]string, len(s))
	for i, c := range s {
		array[i] = fmt.Sprintf("%s%.8b", "", c)
	}
	return array
}

func getRandomPixels(keyLen int, sizeOfPicture int) []int {
	var array []int
	for i := 0; i < keyLen*8; i++ {
		tmp := rand.Intn(sizeOfPicture)
		if isNotInSlice(tmp, array) {
			array = append(array, tmp)
		} else {
			for {
				tmp = rand.Intn(sizeOfPicture)
				if isNotInSlice(tmp, array) {
					array = append(array, tmp)
					break
				}
			}
		}
	}
	return array
}

func isNotInSlice(x int, slice []int) bool {
	for _, v := range slice {
		if v == x {
			return false // x is found in the slice
		}
	}
	return true // x is not found in the slice
}

func encrypt(bounds image.Rectangle, newImg *image.RGBA, img image.Image, randomPixels []int, binaryText []string) *image.RGBA {

	var newBlueBrightnes uint8
	for i := 0; i < len(randomPixels); i++ {
		r, g, b, a := img.At(randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy()).RGBA()

		rBr := uint8(r >> 8)
		gBr := uint8(g >> 8)
		bBr := uint8(b >> 8)
		aBr := uint8(a >> 8)
		//fmt.Println(rBr, gBr, bBr)
		Y := 0.3*float64(rBr) + 0.59*float64(gBr) + 0.11*float64(bBr)
		//Y := 0.2989*float64(r) + 0.58662*float64(g) + 0.11448*float64(b)
		str := binaryText[i/8]
		lambda := 0.25
		fmt.Println(string(str[i%8]))

		if string(str[i%8]) == "1" {
			newBlueBrightnes = uint8(float64(bBr) + lambda*Y)
			fmt.Println("+", newBlueBrightnes, float64(bBr)+lambda*Y)
		} else if string(str[i%8]) == "0" {
			newBlueBrightnes = uint8(float64(bBr) - lambda*Y)
			fmt.Println("-", newBlueBrightnes, float64(bBr)-lambda*Y)
		}

		//newBlueBr := int(newBlueBrightnes) << 8
		//fmt.Println("1", newBlueBrightnes, bBr)
		newColor := color.RGBA{rBr, gBr, newBlueBrightnes, aBr}
		fmt.Println(rBr, gBr, newBlueBrightnes, aBr)
		//newColor := color.RGBA{uint8(r), uint8(g), uint8(newBlueBr), uint8(a)}
		newImg.Set(randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy(), newColor)
	}
	//}
	return newImg
}

func decrypt(bounds image.Rectangle, newImg *image.RGBA, randomPixels []int) string {
	array := make([]string, len(randomPixels), len(randomPixels))
	var bOriginalBrightnes uint32
	beta := 6

	for i := 0; i < len(randomPixels); i++ {
		bOriginalBrightnes = 0
		var bBr255 uint8
		//fmt.Println(randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy())
		for j := -beta; j <= beta; j++ {
			_, _, bBrYplusI, _ := newImg.At(randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy()+j).RGBA()
			//_, _, bBrYminusI, _ := newImg.At(randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy()-j).RGBA()
			_, _, bBrXplusI, _ := newImg.At(randomPixels[i]%bounds.Dx()+j, randomPixels[i]/bounds.Dy()).RGBA()

			bBrXPlysI := uint8(bBrXplusI >> 8)
			bBrYPlusI := uint8(bBrYplusI >> 8)

			//_, _, bBrXminusI, _ := newImg.At(randomPixels[i]%bounds.Dx()-j, randomPixels[i]/bounds.Dy()).RGBA()
			//bOriginalBrightnes += uint32(bBrXminusI) + uint32(bBrXplusI) + uint32(bBrYminusI) + uint32(bBrYplusI)

			bOriginalBrightnes += uint32(bBrYPlusI + bBrXPlysI)

			//bOriginalBrightnes += bBrYplusI + bBrXplusI
			//fmt.Println("Синий: ", uint8(bBrYplusI>>8), uint8(bBrXplusI>>8), bOriginalBrightnes, uint8(bOriginalBrightnes))
			// fmt.Println("1", randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy()+j)
			// fmt.Println("2", randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy()-j)
			// fmt.Println("3", randomPixels[i]%bounds.Dx()+j, randomPixels[i]/bounds.Dy())
			// fmt.Println("4", randomPixels[i]%bounds.Dx()-j, randomPixels[i]/bounds.Dy())
		}
		_, _, bBr, _ := newImg.At(randomPixels[i]%bounds.Dx(), randomPixels[i]/bounds.Dy()).RGBA()
		bBr255 = uint8(bBr >> 8)
		bOriginalBrightnes -= uint32(2 * bBr255)
		//fmt.Println("bBr255: ", uint32(2*bBr255), uint8(2*bBr255), uint8(bOriginalBrightnes))

		//bOriginalBrightnesInt8 := uint8(bOriginalBrightnes)
		//bOriginalBrightnesInt8 -= 2 * bBr255
		bOriginalBrightnes = bOriginalBrightnes / uint32((4 * beta))
		//fmt.Println("1", bOriginalBrightnes)
		bBr255 = uint8(bBr >> 8)
		//fmt.Println("2", bBr255, bOriginalBrightnes)

		//bOriginalBrightnes255 := uint8(bOriginalBrightnes >> 8)
		//fmt.Println("2", bBr255, bOriginalBrightnes255)
		if uint8(bOriginalBrightnes) < uint8(bBr255) {
			array[i/8] += "1"
		} else if uint8(bOriginalBrightnes) >= uint8(bBr255) {
			array[i/8] += "0"
		}
	}
	fmt.Println(array)
	var result string
	for i := range array {
		num, _ := strconv.ParseInt(array[i], 2, 64)
		char := string(rune(num))
		result += char
	}
	//result := strings.Join(arrayB, " ")
	fmt.Println("result: ", result)
	return result
}
