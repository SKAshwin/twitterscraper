package twitterscraper

import (
	"encoding/json"
	"log"
	"net/http"
	urlutil "net/url"
	"strings"

	"errors"

	"github.com/PuerkitoBio/goquery"
)

// GetTweets returns an array of tweets (as strings) by a twitter user
// handle is the twitter handle of the user
// pages is the pages of tweets from that user you want
func GetTweets(handle string, pages uint8) ([]string, error) {
	url := "https://twitter.com/i/profiles/show/" + handle + "/timeline/tweets?include_available_features=1&include_entities=1&include_new_items_bar=true"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Hello world", handle, pages)

	req.Header = map[string][]string{
		"Accept":                []string{"application/json", "text/javascript", "*/*; q=0.01"},
		"Referer":               []string{"https://twitter.com/" + handle},
		"User-Agent":            []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8"},
		"X-Twitter-Active-User": []string{"yes"},
		"X-Requested-With":      []string{"XMLHttpRequest"},
		"Accept-Language":       []string{"en-US"},
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	var tweets []string

	for i := uint8(0); i < pages; i++ {
		var returnval struct {
			Items string `json:"items_html"`
		}
		err = json.NewDecoder(resp.Body).Decode(&returnval)
		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(returnval.Items))
		if err != nil {
			return nil, err
		}

		tweetsHTML := doc.Find(".tweet-text")
		curPageTweets := tweetsHTML.Map(func(index int, tweetHTML *goquery.Selection) string {
			return tweetHTML.Text()
		})
		tweets = append(tweets, curPageTweets...)

		finalIDonPage, exists := doc.Find(".stream-item").Last().Attr("data-item-id")
		if !exists {
			return nil, errors.New("Page fetched formatted incorrectly")
		}

		form, _ := urlutil.ParseQuery(req.URL.RawQuery)
		form.Add("max_position", finalIDonPage)
		req.URL.RawQuery = form.Encode()

		resp, err = (&http.Client{}).Do(req)
		if err != nil {
			return nil, err
		}
	}

	return tweets, nil

}