package networking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jodios/samplemovement/game"
	"github.com/jodios/samplemovement/game/characters"
)

var lock = sync.Mutex{}

type Networking struct {
	Connection       *websocket.Conn
	CharacterChannel chan characters.Orange
	CurrentGameState *game.GameState
	Context          context.Context
}

func NewNetworking(ctx context.Context, conn *websocket.Conn) *Networking {
	n := &Networking{
		Connection:       conn,
		CharacterChannel: make(chan characters.Orange),
		Context:          ctx,
	}
	return n
}

func (n *Networking) Start(name string, sprites map[string]*ebiten.Image) {
	go stateListener(n.CurrentGameState, name, sprites, n.Connection)
	go publishCharUpdate(n.CharacterChannel, name)
}

func publishCharUpdate(charChan chan characters.Orange, name string) {
	for gs := range charChan {
		data, _ := json.Marshal(gs)
		http.Post(
			"http://localhost:8080/publish?name="+name,
			"application/text",
			bytes.NewBuffer(data),
		)
	}
}

func stateListener(cgs *game.GameState, name string, sprites map[string]*ebiten.Image, conn *websocket.Conn) {
	for {
		newState := &game.GameState{}
		conn.ReadJSON(newState)
		if newState.AddedCharacter != nil {
			fmt.Println(len(newState.Characters))
			fmt.Println(len(cgs.Characters))
			if len(newState.Characters) >= len(cgs.Characters) {
				fmt.Printf("%v has joined the game\n", *newState.AddedCharacter)
			} else {
				fmt.Printf("%v has left the game\n", *newState.AddedCharacter)
				delete(cgs.Characters, *newState.AddedCharacter)
			}
		}
		for k, v := range newState.Characters {
			if c, ok := cgs.Characters[k]; ok {
				c.PosX = v.PosX
				c.PosY = v.PosY
				c.Speed = v.Speed
			} else {
				cgs.Characters[k] = characters.NewOrange(sprites)
				cgs.Characters[k].PosX = v.PosX
				cgs.Characters[k].PosY = v.PosY
				cgs.Characters[k].Speed = v.Speed
			}
		}
	}
}
