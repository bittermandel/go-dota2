package dota2

import (
	"github.com/Philipp15b/go-steam"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestDota2(t *testing.T) {
	// This test just runs forever. Made this instead of a main.go (I'm lazy)
	client := steam.NewClient()

	sentryFileHash, err := ioutil.ReadFile("./.sentry")
	if err != nil {
		log.Println("Could not find sentry file.")
	}

	auth := &steam.LogOnDetails{
		Username:       os.Getenv("STEAM_USERNAME"),
		Password:       os.Getenv("STEAM_PASSWORD"),
		AuthCode:		os.Getenv("STEAM_AUTHCODE"),
		SentryFileHash: sentryFileHash,
	}
	New(client, auth)
	for {
		// Keep running and just listen
		continue
	}
}
