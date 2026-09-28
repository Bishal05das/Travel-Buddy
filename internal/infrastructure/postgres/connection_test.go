package db

import (
	"testing"

	"github.com/bishal05das/travelbuddy/config"
)

func TestGetConnectionStringQuotesValues(t *testing.T) {
	got := GetConnectionString(&config.DBConfig{
		User: "app", Password: `p a'ss\word`, Host: "localhost", Port: 5432, Name: "tour",
	})
	want := `user='app' password='p a\'ss\\word' host='localhost' port=5432 dbname='tour' sslmode=disable`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}
