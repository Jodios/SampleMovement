package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	l, err := net.Listen("TCP", ":8080")
	if err != nil {
		fmt.Printf("Error establishing TCP connection: %v", err)
	}
	s := &http.Server{
		Handler:      &GameServer{},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	select {
	case err := <-errc:
		log.Printf("Failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("Terminating: %v", sig)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return s.Shutdown(ctx)
}

type GameServer struct {
}

func (s GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("Error accepting websocket: %v", err)
	}
	defer c.Close(websocket.StatusInternalError, "checkem")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		fmt.Printf("Error reading json: %v", err)
	}
	fmt.Printf("Received: %v", v)
	c.Close(websocket.CloseNormalClosure)
}
