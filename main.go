package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jodios/samplemovement/game"
	"github.com/jodios/samplemovement/game/characters"
)

const (
	screenWidth  = 800
	screenHeight = 700
)
const (
	// insideWidth  = 160
	// insideHeight = 140
	insideWidth  = 320
	insideHeight = 280
)

var (
	assets map[string]*ebiten.Image = make(map[string]*ebiten.Image)
)

func init() {
	filesystem := os.DirFS("assets/images")
	files, err := ioutil.ReadDir("assets/images")
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range files {
		file, err := filesystem.Open(fileInfo.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		assets[fileInfo.Name()] = ebiten.NewImageFromImage(img)
	}
}

//go:generate go run github.com/unitoftime/packer/cmd/packer@latest --input images --stats
func main() {
	game := game.Game{
		InnerWidth:  insideWidth,
		InnerHeight: insideHeight,
		Character:   characters.NewOrange(assets),
	}
	ebiten.SetWindowTitle("Super Awesome Game!!!")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	ebiten.RunGame(&game)
}
