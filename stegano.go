package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

//use this instead of standard utf-8 encodings to keep deltas small
func initCharacterMap() (map[string]int, map[int]string) {
	charMap := map[string]int{
		"a":  1,
		"b":  2,
		"c":  3,
		"d":  4,
		"e":  5,
		"f":  6,
		"g":  7,
		"h":  8,
		"i":  9,
		"j":  10,
		"k":  11,
		"l":  12,
		"m":  13,
		"n":  14,
		"o":  15,
		"p":  16,
		"q":  17,
		"r":  18,
		"s":  19,
		"t":  20,
		"u":  21,
		"v":  22,
		"w":  23,
		"x":  24,
		"y":  25,
		"z":  26,
		" ":  27,
		".":  28,
		",":  29,
		"!":  30,
		"\"": 31,
		"'":  32,
		"?":  33,
		"\n": 34,
		"\t": 35,
	}
	intMap := map[int]string{
		1:  "a",
		2:  "b",
		3:  "c",
		4:  "d",
		5:  "e",
		6:  "f",
		7:  "g",
		8:  "h",
		9:  "i",
		10: "j",
		11: "k",
		12: "l",
		13: "m",
		14: "n",
		15: "o",
		16: "p",
		17: "q",
		18: "r",
		19: "s",
		20: "t",
		21: "u",
		22: "v",
		23: "w",
		24: "x",
		25: "y",
		26: "z",
		27: " ",
		28: ".",
		29: ",",
		30: "!",
		31: "\"",
		32: "'",
		33: "?",
		34: "\n",
		35: "\t",
	}
	return charMap, intMap
}

func encodeImage(message string, sourceImagePath string, destPath string) {
	//convert characters in message into numbers
	charMap, _ := initCharacterMap()
	downCaseMessage := strings.ToLower(message)
	slMessage := strings.Split(downCaseMessage, "")
	fmt.Println(slMessage)
	var slEncodedChars []int
	for _, char := range slMessage {
		slEncodedChars = append(slEncodedChars, charMap[char])
	}
	//open source image
	reader, err := os.Open(sourceImagePath)
	if err != nil {
		fmt.Println("Error in File Read:", err)
	}
	defer reader.Close()
	//decode source image
	img, _, err := image.Decode(reader)
	if err != nil {
		fmt.Println("Error in Decode", err)
	}
	//get source image "pixel" dimensions
	imgBounds := img.Bounds()
	minX := imgBounds.Min.X
	minY := imgBounds.Min.Y
	maxX := imgBounds.Max.X
	maxY := imgBounds.Max.Y

	//create a new RGBA
	encodedRGBA := image.NewRGBA(imgBounds)
	var curColor color.Color
	//iterate over each "pixel"
	for i := minX; i < maxX; i++ {
		for j := minY; j < maxY; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			if len(slEncodedChars) > i+j {
				r += uint32(slEncodedChars[i+j])
			}

			curColor = color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			}
			encodedRGBA.Set(i, j, curColor)
		}
	}
	outfile, err := os.Create("encoded.png")
	if err != nil {
		fmt.Println(err)
	}
	png.Encode(outfile, encodedRGBA)
}

func decodeImage(encodedImagePath string, masterImagePath string) string {
	_, intMap := initCharacterMap()

	encodedFile, err := os.Open(encodedImagePath)
	if err != nil {
		fmt.Println(err)
	}
	encodedImage, _, err := image.Decode(encodedFile)
	if err != nil {
		fmt.Println(err)
	}

	masterFile, err := os.Open(masterImagePath)
	if err != nil {
		fmt.Println(err)
	}
	masterImage, _, err := image.Decode(masterFile)
	if err != nil {
		fmt.Println(err)
	}

	bounds := encodedImage.Bounds()
	minX := bounds.Min.X
	maxX := bounds.Max.X
	minY := bounds.Min.Y
	maxY := bounds.Max.Y

	var encodedColor color.Color
	var masterColor color.Color

	colorDeltas := make(map[string][]uint8, 0)

	for i := minX; i < maxX; i++ {
		for j := minY; j < maxY; j++ {
			encodedColor = encodedImage.At(i, j)
			masterColor = masterImage.At(i, j)

			encR, encG, encB, encA := encodedColor.RGBA()
			mastR, mastG, mastB, mastA := masterColor.RGBA()

			curEncColor := color.RGBA{
				R: uint8(encR),
				G: uint8(encG),
				B: uint8(encB),
				A: uint8(encA),
			}
			curMastColor := color.RGBA{
				R: uint8(mastR),
				G: uint8(mastG),
				B: uint8(mastB),
				A: uint8(mastA),
			}

			colorDeltas["Red"] = append(colorDeltas["Red"], curEncColor.R-curMastColor.R)
			colorDeltas["Green"] = append(colorDeltas["Green"], curEncColor.G-curMastColor.G)
			colorDeltas["Blue"] = append(colorDeltas["Blue"], curEncColor.B-curMastColor.B)
			colorDeltas["Alpha"] = append(colorDeltas["Alpha"], curEncColor.A-curMastColor.A)
		}
	}
	// just for red for now
	message := ""
	for _, code := range colorDeltas["Red"] {
		message = fmt.Sprintf("%v%v", message, intMap[int(code)])
	}

	return message
}

func parseArgs() (process string, encodedPath string, masterPath string, message string) {
	pProcess := flag.String("process", "", "set this flag to 'encode' to encode an image, and decode to decode an image")
	pEncodedPath := flag.String("encodedPath", "", "set this flag to either write destination for a new encoded image or the path of an existing encoded image when decoding")
	pMasterPath := flag.String("masterPath", "", "set this flag when encoding and decoding as the path to the source image")
	pMessage := flag.String("message", "", "set the message you wish to encode in the image")
	flag.Parse()
	return *pProcess, *pEncodedPath, *pMasterPath, *pMessage
}

func main() {
	//process, encodedPath, masterPath, message := parseArgs()
	//temporary arguments to make iterating easier
	process := "encode"
	encodedPath := "/Users/tjudd/Documents/work/sandbox/stegano/encoded.png"
	masterPath := "/Users/tjudd/Documents/work/sandbox/stegano/squirrel.png"
	message := "Hello World"
	//message, err := ioutil.ReadFile("source.txt")
	//if err != nil {
	//	fmt.Println(err)
	//}
	if process == "encode" {
		encodeImage(string(message), masterPath, encodedPath)
	}
	if process == "decode" {
		receivedMessage := decodeImage(encodedPath, masterPath)
		fmt.Println("received:", receivedMessage)
	}
	if process == "" {
		fmt.Println("Invalid process, requre encode or decode")
	}
	return
}
