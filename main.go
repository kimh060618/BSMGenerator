package main

import (
	"Autocrypt/BSMGenerator/bsm"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

type MessageType int

const (
	PCW MessageType = iota
)

const (
	FCW MessageType = 1 + iota
	EEBL
	CLW
)

const (
	BSW MessageType = 4 + iota
	LCW
	LTA
	RLVW
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func randomGenerator() int {
	rand.Seed(time.Now().UnixNano())
	p := rand.Float32()
	if p < 0.02 {
		return 0
	} else if p < 0.08 {
		return rand.Intn(3) + 1
	} else if p < 0.2 {
		return rand.Intn(4) + 4
	} else {
		return 8
	}
}

func levelCheck(t int) int {
	if t == int(PCW) {
		return 1
	} else if int(FCW) <= t && t < int(BSW) {
		return 2
	} else if t < 8 {
		return 3
	} else {
		return 4
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		m := &bsm.BSMMessage{}

		(*m).Type = randomGenerator()
		(*m).TimeStamp = time.Now().UnixNano()
		(*m).Level = levelCheck((*m).Type)

		err = conn.WriteJSON(m)
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1000 * time.Microsecond)
	}
}

func main() {
	mux := pat.New()
	mux.Get("/bsm", handler)

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
