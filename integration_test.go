package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/kurth4cker/httpserver/server"
)

func TestRecordingAndRetrievingWins(t *testing.T) {
	store := NewInMemoryPlayerStore()
	svr := server.NewPlayerServer(store)
	player := "Pepper"

	for range 3 {
		request, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)
	}

	t.Run("get score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/"+player, nil)
		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)

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
	})

	t.Run("get league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		want := []server.Player{
			{Name: "Pepper", Wins: 3},
		}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func getLeagueFromResponse(t testing.TB, body io.Reader) []server.Player {
	t.Helper()
	league := make([]server.Player, 0)
	if err := json.NewDecoder(body).Decode(&league); err != nil {
		t.Fatalf("failed to parse %q: %s", body, err)
	}
	return league
}
