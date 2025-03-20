package main

import (
	"net/http"
	"sync"

	"github.com/kurth4cker/httpserver/server"
)

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		Store: make(map[string]int),
	}
}

type InMemoryPlayerStore struct {
	Store      map[string]int
	storeMutex sync.Mutex
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	i.storeMutex.Lock()
	defer i.storeMutex.Unlock()
	return i.Store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.storeMutex.Lock()
	defer i.storeMutex.Unlock()
	i.Store[name]++
}

func main() {
	store := NewInMemoryPlayerStore()
	server := &server.PlayerServer{
		Store: store,
	}
	err := http.ListenAndServe("127.0.0.1:8080", server)
	if err != nil {
		panic(err)
	}
}
