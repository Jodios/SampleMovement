package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error establishing TCP connection: %v", err)
	}
	s := &http.Server{
		Handler:      NewGameServer(),
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
	return s.Shutdown(ctx)
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

type GameServer struct {
	subscribeMessageBuffer int
	publishLimiter         *rate.Limiter
	serveMux               http.ServeMux
	subscribersMu          sync.Mutex
	subscribers            map[*subscriber]struct{}
}

func NewGameServer() *GameServer {
	gs := &GameServer{
		subscribeMessageBuffer: 16,
		subscribers:            make(map[*subscriber]struct{}),
		publishLimiter:         rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	gs.serveMux.HandleFunc("/subscribe", gs.subscribeHandler)
	gs.serveMux.HandleFunc("/publish", gs.publishHandler)
	return gs
}

func (gs *GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}

func (gs *GameServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = gs.subscribe(r.Context(), c)
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
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
	}
	gs.publish(msg)
	w.WriteHeader(http.StatusAccepted)
}

func (gs *GameServer) subscribe(ctx context.Context, c *websocket.Conn) error {
	ctx = c.CloseRead(ctx)
	s := &subscriber{
		msgs: make(chan []byte, gs.subscribeMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
	gs.addSubscriber(s)
	defer gs.deleteSubscriber(s)
	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
func (gs *GameServer) publish(msg []byte) {
	gs.subscribersMu.Lock()
	defer gs.subscribersMu.Unlock()

	gs.publishLimiter.Wait(context.Background())
	for s := range gs.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

func (gs *GameServer) addSubscriber(s *subscriber) {
	gs.subscribersMu.Lock()
	gs.subscribers[s] = struct{}{}
	gs.subscribersMu.Unlock()
}
func (gs *GameServer) deleteSubscriber(s *subscriber) {
	gs.subscribersMu.Lock()
	delete(gs.subscribers, s)
	gs.subscribersMu.Unlock()
}
func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, msg)
}
