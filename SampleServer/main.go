package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jodios/sampleserver/game"
	"nhooyr.io/websocket"
)

func main() {
	gameState, err := run()
	if err != nil {
		panic(err)
	}
	saveFile, _ := json.MarshalIndent(gameState, "", "\t")
	err = os.WriteFile("game.json", saveFile, 0666)
	if err != nil {
		fmt.Printf("Unable to save game: %v\n", err)
	}
}

func run() (*game.Game, error) {
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error establishing TCP connection: %v", err)
	}
	fmt.Print("Listening on port 8080\n")
	gs := NewGameServer()
	s := &http.Server{
		Handler:      gs,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("Failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("Terminating: %v", sig)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return gs.gameState, s.Shutdown(ctx)
}

type subscriber struct {
	msgs       chan []byte
	closeSlow  func()
	connection *websocket.Conn
}
type GameServer struct {
	subscribeMessageBuffer int
	serveMux               http.ServeMux
	subscribersMu          sync.Mutex
	subscribers            map[string]*subscriber
	gameState              *game.Game
	activeGameState        *game.Game
	count                  int
}

func NewGameServer() *GameServer {
	gameState := game.Game{}
	file, err := os.Open("game.json")
	if err != nil {
		log.Printf("Failed to open a game....Creating new one\n")
	}
	err = json.NewDecoder(file).Decode(&gameState)
	if err != nil {
		log.Printf("Failed to open a game....Creating new one\n")
		gameState = game.Game{
			Characters: make(map[string]*game.Character),
		}
	}
	gs := &GameServer{
		subscribeMessageBuffer: 1020,
		subscribers:            make(map[string]*subscriber),
		gameState:              &gameState,
		activeGameState: &game.Game{
			Characters: make(map[string]*game.Character),
		},
	}
	// go func() {
	// 	for {
	// 		msg, _ := json.Marshal(gs.activeGameState)
	// 		fmt.Println(string(msg))
	// 		time.Sleep(time.Second * 2)
	// 	}
	// }()
	gs.serveMux.HandleFunc("/subscribe", gs.subscribeHandler)
	gs.serveMux.HandleFunc("/publish", gs.publishHandler)
	return gs
}
func (gs *GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}
func (gs *GameServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("name")
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = gs.subscribe(r.Context(), c, user)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
func (gs *GameServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
	name := r.URL.Query().Get("name")
	err := json.NewDecoder(r.Body).Decode(gs.activeGameState.Characters[name])
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	msg, _ := json.Marshal(gs.activeGameState)
	gs.publish(r.Context(), msg)
	w.WriteHeader(http.StatusAccepted)
}
func (gs *GameServer) subscribe(ctx context.Context, c *websocket.Conn, user string) error {
	fmt.Printf("%v has joined....\n", user)
	ctx = c.CloseRead(ctx)
	s := &subscriber{
		connection: c,
		msgs:       make(chan []byte, gs.subscribeMessageBuffer),
		closeSlow: func() {
			c.Close(
				websocket.StatusPolicyViolation,
				"connection too slow to keep up with messages",
			)
		},
	}
	if userChar, ok := gs.gameState.Characters[user]; ok {
		gs.activeGameState.Characters[user] = userChar
		jsonChar, err := json.Marshal(userChar)
		if err != nil {
			fmt.Println(err)
		}
		s.msgs <- []byte(jsonChar)
	} else {
		newChar := &game.Character{
			PosX:  100,
			PosY:  100,
			Speed: 1,
		}
		gs.gameState.Characters[user] = newChar
		gs.activeGameState.Characters[user] = newChar
		jsonChar, err := json.Marshal(gs.gameState.Characters[user])
		if err != nil {
			fmt.Println(err)
		}
		s.msgs <- []byte(jsonChar)
	}
	gs.addSubscriber(ctx, s, user)
	defer gs.deleteSubscriber(user)
	for {
		select {
		case msg := <-s.msgs:
			gs.count++
			err := writeTimeout(ctx, time.Millisecond*1000, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
func (gs *GameServer) publish(ctx context.Context, msg []byte) {
	gs.subscribersMu.Lock()
	defer gs.subscribersMu.Unlock()
	for _, s := range gs.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}
func (gs *GameServer) addSubscriber(ctx context.Context, s *subscriber, user string) {
	gs.subscribersMu.Lock()
	gs.subscribers[user] = s
	gs.activeGameState.Added = &user
	msg, _ := json.Marshal(gs.activeGameState)
	for _, subs := range gs.subscribers {
		subs.msgs <- msg
	}
	gs.activeGameState.Added = nil
	gs.subscribersMu.Unlock()
}
func (gs *GameServer) deleteSubscriber(user string) {
	fmt.Printf("%v has left....\n", user)
	gs.subscribersMu.Lock()
	delete(gs.subscribers, user)
	delete(gs.activeGameState.Characters, user)
	gs.activeGameState.Added = &user
	msg, _ := json.Marshal(gs.activeGameState)
	for _, s := range gs.subscribers {
		s.msgs <- msg
	}
	gs.activeGameState.Added = nil
	gs.subscribersMu.Unlock()
}
func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, msg)
}
