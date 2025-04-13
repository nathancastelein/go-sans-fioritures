package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Stub struct{}

func (s *Stub) ListStones(ctx context.Context) []*InfinityStone {
	return []*InfinityStone{
		{
			Name:   "space",
			Color:  "blue",
			Power:  "Teleportation",
			Status: "secured",
		},
		{
			Name:   "mind",
			Color:  "yellow",
			Power:  "Mind Control",
			Status: "secured",
		},
	}
}

func (s *Stub) GetStone(ctx context.Context, stoneName string) (*InfinityStone, error) {
	return &InfinityStone{stoneName, "white", stoneName, "secured"}, nil
}

func (s *Stub) ReportSuspiciousActivity(ctx context.Context, report Report) {

}

func (s *Stub) LogUser(ctx context.Context, username string, password string) (*User, error) {
	return &User{
		ID:        1,
		FirstName: "Peter",
		LastName:  "Parker",
		HeroName:  "Spider-Man",
	}, nil
}

func TestServer(t *testing.T) {
	t.Parallel()
	stub := &Stub{}
	server := httptest.NewServer(NewServer(stub, stub, stub))
	t.Cleanup(func() {
		server.Close()
	})

	t.Run("GET /stones", func(t *testing.T) {
		t.Parallel()
		// Arrange
		request, err := http.NewRequest(http.MethodGet, server.URL+"/stones", nil)
		if err != nil {
			t.Fatalf("expected no error while creating request, got: %s", err)
		}
		request.SetBasicAuth("Peter", "Parker")
		want := stub.ListStones(context.Background())

		// Act
		response, err := http.DefaultClient.Do(request)

		// Assert
		if err != nil {
			t.Fatalf("expected no error while performing request, got: %s", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("expected status code 200, got: %d", response.StatusCode)
		}

		defer response.Body.Close()

		var got []*InfinityStone
		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Fatalf("could not decode response body: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected JSON response.\nGot:  %#v\nWant: %#v", got, want)
		}
	})

	t.Run("GET /stones/space", func(t *testing.T) {
		t.Parallel()
		// Arrange
		request, err := http.NewRequest(http.MethodGet, server.URL+"/stones/space", nil)
		if err != nil {
			t.Fatalf("expected no error while creating request, got: %s", err)
		}
		request.SetBasicAuth("Peter", "Parker")
		want, err := stub.GetStone(context.Background(), "space")
		if err != nil {
			t.Fatalf("expected no error getting expected stone, got: %s", err)
		}

		// Act
		response, err := http.DefaultClient.Do(request)

		// Assert
		if err != nil {
			t.Fatalf("expected no error while performing request, got: %s", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("expected status code 200, got: %d", response.StatusCode)
		}

		defer response.Body.Close()

		got := &InfinityStone{}
		if err := json.NewDecoder(response.Body).Decode(got); err != nil {
			t.Fatalf("could not decode response body: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected JSON response.\nGot:  %#v\nWant: %#v", got, want)
		}
	})

	t.Run("POST /stones/report", func(t *testing.T) {
		t.Parallel()
		// Arrange
		request, err := http.NewRequest(http.MethodPost, server.URL+"/stones/report", bytes.NewBufferString(`{"stone": "space", "report": "lost"}`))
		if err != nil {
			t.Fatalf("expected no error while creating request, got: %s", err)
		}
		request.SetBasicAuth("Peter", "Parker")

		// Act
		response, err := http.DefaultClient.Do(request)

		// Assert
		if err != nil {
			t.Fatalf("expected no error while performing request, got: %s", err)
		}

		if response.StatusCode != http.StatusAccepted {
			t.Fatalf("expected status code 200, got: %d", response.StatusCode)
		}
	})
}
