package twitterscraper_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/SKAshwin/twitterscraper"
)

func TestGetTweets(t *testing.T) {
	fmt.Println("Hello world")
	tweets, err := twitterscraper.GetTweets("realDonaldTrump", 25)
	if err != nil {
		t.Fatal(err)
	}

	for i, tweet := range tweets {
		log.Println(i, tweet)
	}
}
