package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/kurth4cker/httpserver/server"
)

func TestGetPlayers(t *testing.T) {
	store := stupPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := server.NewPlayerServer(&store)

	t.Run("get Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		{
			got := response.Code
			want := http.StatusOK

			if got != want {
				t.Errorf("got status %d, want %d", got, want)
			}
		}

		got := response.Body.String()
		want := "20"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("get Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		{
			got := response.Code
			want := http.StatusOK

			if got != want {
				t.Errorf("got status %d, want %d", got, want)
			}
		}

		got := response.Body.String()
		want := "10"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("return 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound
		if got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	})
}

func TestStoreWins(t *testing.T) {
	store := stupPlayerStore{
		scores: make(map[string]int),
	}
	server := server.NewPlayerServer(&store)

	t.Run("return accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted
		if got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	})

	t.Run("record win when POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Floyd", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := len(store.winCalls)
		want := 2

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("record name on POST", func(t *testing.T) {
		player := "Edward"

		request, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := store.winCalls[2]
		want := "Edward"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestLeague(t *testing.T) {
	store := stupPlayerStore{}
	svr := server.NewPlayerServer(&store)

	t.Run("it returns 200 on league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	})

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []server.Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := stupPlayerStore{nil, nil, wantedLeague}
		svr := server.NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		got := newLeagueFromResponse(t, response.Body)

		if !slices.Equal(store.league, got) {
			t.Errorf("got %v, want %v", got, wantedLeague)
		}
	})

	t.Run("content-type should be JSON", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		got := response.Result().Header.Get("content-type")
		want := "application/json"

		if got != want {
			t.Errorf("got content-type %q, want %q", got, want)
		}
	})
}

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/players/%s", name),
		nil,
	)
	return request
}

func newLeagueFromResponse(t testing.TB, body io.Reader) []server.Player {
	t.Helper()
	league := make([]server.Player, 0)
	if err := json.NewDecoder(body).Decode(&league); err != nil {
		t.Fatalf("unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return league
}

type stupPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []server.Player
}

func (s *stupPlayerStore) GetLeague() []server.Player {
	return s.league
}

func (s *stupPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *stupPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
