package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kurth4cker/httpserver/server"
)

func TestRecordingAndRetrievingWins(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := server.PlayerServer{
		Store: store,
	}
	player := "Pepper"

	for range 3 {
		request, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
	}

	request, _ := http.NewRequest(http.MethodGet, "/players/"+player, nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	{
		got := response.Code
		want := http.StatusOK
		if got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	}

	{
		got := response.Body.String()
		want := "3"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}
