package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jodios/samplemovement/constants"
	"github.com/jodios/samplemovement/game"
	"github.com/jodios/samplemovement/game/characters"
	"nhooyr.io/websocket"
)

var (
	assets map[string]*ebiten.Image = make(map[string]*ebiten.Image)
)

const (
	charName string = "JohnA"
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

func establishConnection() (*websocket.Conn, *http.Response) {
	url := subUrl + "?name=" + charName
	conn, r, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{})
	if err != nil {
		log.Fatalf("Error whilst connecting to server %v\n", err)
	}
	return conn, r
}

//go:generate go run github.com/unitoftime/packer/cmd/packer@latest --input images --stats
func main() {
	conn, _ := establishConnection()
	g := game.Game{
		InnerWidth:  constants.InsideWidth,
		InnerHeight: constants.InsideHeight,
		Character:   characters.NewOrange(assets),
		Connection:  conn,
	}
	defer conn.Close(websocket.StatusNormalClosure, "Closing")

	_, msg, err := conn.Read(context.Background())
	if err != nil {
		fmt.Printf("Error reading user data %v\n", err)
	}
	state := &characters.Orange{}
	json.NewDecoder(bytes.NewBuffer(msg)).Decode(state)
	g.Character.PosX = state.PosX
	g.Character.PosY = state.PosY
	g.Character.CharName = charName
	g.Character.Speed = state.Speed

	ebiten.SetWindowTitle("Sample Multiplayer Game")
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	ebiten.RunGame(&g)
}

type GameClient struct {
}

func (gc *GameClient) ServeHTTP() {

}
