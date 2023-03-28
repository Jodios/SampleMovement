package main

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jodios/samplemovement/constants"
	"github.com/jodios/samplemovement/game"
	"github.com/jodios/samplemovement/game/characters"
	"github.com/jodios/samplemovement/networking"
)

var (
	assets map[string]*ebiten.Image = make(map[string]*ebiten.Image)
	ctx                             = context.Background()
)

const (
	charName string = "a"
	subUrl   string = "ws://localhost:8080/subscribe"
	pubUrl   string = "http://localhost:8080/publish"
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

func establishConnection() *networking.Networking {
	url := subUrl + "?name=" + charName
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error whilst connecting to server %v\n", err)
	}
	return networking.NewNetworking(ctx, conn)
}

//go:generate go run github.com/unitoftime/packer/cmd/packer@latest --input images --stats
func main() {
	networking := establishConnection()
	gs := &game.GameState{
		Characters: make(map[string]*characters.Orange),
	}
	g := game.Game{
		InnerWidth:       constants.InsideWidth,
		InnerHeight:      constants.InsideHeight,
		Character:        characters.NewOrange(assets),
		CharacterChannel: networking.CharacterChannel,
		GameState:        gs,
	}
	networking.CurrentGameState = gs
	defer networking.Connection.Close()
	defer close(networking.CharacterChannel)

	state := &characters.Orange{}
	err := networking.Connection.ReadJSON(state)
	if err != nil {
		fmt.Printf("Error reading user data %v\n", err)
	}
	g.Character.PosX = state.PosX
	g.Character.PosY = state.PosY
	g.Character.CharName = charName
	g.Character.Speed = state.Speed
	g.GameState.Characters[charName] = g.Character
	networking.Start(charName, assets)

	ebiten.SetWindowTitle("Sample Multiplayer Game")
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	ebiten.RunGame(&g)
}

type GameClient struct {
}

func (gc *GameClient) ServeHTTP() {

}
