package term

import (
	"fmt"
	"os"
	"image"
	"math"
	"image/color"
	_ "image/png"
	"image/draw"
)

type font struct {
	brightness map[rune][]float32
	height, width int
}

var (
	// This is by far the most annoying string I have ever hand-typed.
	// (even assembly byte-code is better)
	// Also, the ☺ is a placeholder for non-existant characters.
	runeOrder = []rune(
		"☺0123456789ABCDEFG" +
        "HIJKLMNOPQRSTUVWXY" +
        "Zabcdefghijklmnopq" +
        "rstuvwxyzæÆαβγΔδξε" +
        "θλμπΣστφψΩáâäàåêëè" +
        "íïîìóôöòúûüùÿħñāēō" +
        "`~!@#$%^&*()-_=+[{" +
        "]}\\|;:'\",<.>/?÷≈∞¿" +
        "¡«»¶•₹ℓ₴₮£¥€Ж¢§‡∂☺" +
        "ЯЉЊΞ☺†Ю☺≤≥∫√…ÇçШЂ ")
)


const (
	blockRows = 10
	blockCols = 18
	nullRune = '☺'
)

func runeDimensions(img image.Image) (height, width int) {
	// Returns the height and width of a single rune based off of the
	// input image.

	b := img.Bounds()
	if b.Dx() % blockCols != 0 || b.Dy() % blockRows != 0 { 
		panic(fmt.Sprintf("Oops! Your font file is %d pixels wide and %d" +
			" pixels tall, but its width needs to be divisible by %d and its" +
			" height needs to be divisible by %d.", b.Dx(), b.Dy(), 
			blockCols, blockRows))
	}
	height, width = b.Dy() / blockRows, b.Dx() / blockCols

	return

}

func addPixel(f *font, r rune, c color.NRGBA) {
	// Puts a boolean into the correct place in the correct array. Creates
	// an array if it doesn't already exist.
	if r == nullRune { return }		

	if _, ok := f.brightness[r]; !ok {
		f.brightness[r] = make([]float32, 0, f.height * f.width)
	}

	total := c.R / 3 + c.G / 3 + c.B / 3

	f.brightness[r] = append(f.brightness[r], 
		1 - float32(total) /  float32(math.MaxUint8))

	if len(f.brightness[r]) > f.height * f.width {
		panic(fmt.Sprintf("Rune '%c' accessed more than the maximum" + 
			" number of times", r))
	}
}

func buildFont(fileName string) *font {
	// Creates a font based off of the image file stored at fileName
	//
	// First character in sheet is ignored. Any color found in that cell is
	// ignored. Any other color is replaced by the color of the cell. I may
	// change this to support shading at some point.

	file, err := os.Open(fileName)
	if err != nil { panic(err) }
	defer file.Close()

	if len(runeOrder) != blockRows * blockCols {
		panic(fmt.Sprintf("|runeOrder| = %d, was expecting %d * %d = %d",
			len(runeOrder), blockRows, blockCols, blockRows * blockCols))
	}

	decodedImg, _, err := image.Decode(file)
	if err != nil { panic(err) }
	b := decodedImg.Bounds()

	blockImg := image.NewNRGBA(b)
	draw.Draw(blockImg, b, decodedImg, 
		image.Pt(b.Min.X, b.Min.Y), draw.Src)

	runeHeight, runeWidth := runeDimensions(blockImg)
	
	f := font{
		brightness: make(map[rune][]float32), 
		height: runeHeight, 
		width: runeWidth,
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			runeX := (x - b.Min.X) / runeWidth
			runeY := (y - b.Min.Y) / runeHeight
			i := runeX + blockCols * runeY

			c := blockImg.At(x, y)

			if nrgba, ok := c.(color.NRGBA); ok {
				addPixel(&f, runeOrder[i], nrgba)
			} else {
				panic("Can't happen (I hope).")
			}
		}
	}

	return &f
}