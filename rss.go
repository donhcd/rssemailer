package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	feed := rss.New(5, true, chanHandler, itemHandler)
	for {
		if err := feed.Fetch("http://feeds.feedburner.com/TheFlightDeal", nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s\n", err)
			return
		}
		<-time.After(1 * time.Hour)
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
}

type OredInterests []string

type Interests struct {
	Email          string
	AndedInterests []OredInterests
}

var interests []Interests

func getRecipientsForItem(item *rss.Item) []string {
	var recipients []string
	for _, interest := range interests {
		ciTit := strings.ToLower(item.Title)
		matchesAnd := true
		for _, oredInterest := range interest.AndedInterests {
			matchesOr := false
			for _, wantedSubstring := range oredInterest {
				matchesOr = matchesOr || strings.Contains(ciTit, wantedSubstring)
				if matchesOr {
					break
				}
			}
			matchesAnd = matchesAnd && matchesOr
			if !matchesAnd {
				break
			}
		}
		if matchesAnd {
			recipients = append(recipients, interest.Email)
		}
	}
	return recipients
}

func sendEmailForItem(item *rss.Item, email string) {
	// println(email, item.Title)
	// _ = exec.Command
	// _ = bytes.NewReader
		cmd := exec.Command("mail", "-s", "Flight deal: "+item.Title, email)
		cmd.Stdin = bytes.NewReader([]byte(fmt.Sprintf("Link: %s\n\nDescription: %s\n", item.Comments, item.Description)))
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "[e] mail error: %s; output: %q\n", err, string(output))
		}
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	for _, newItem := range newitems {
		recipients := getRecipientsForItem(newItem)
		for _, recipient := range recipients {
			sendEmailForItem(newItem, recipient)
		}
	}
}

func init() {
	if bytes, err := ioutil.ReadFile("config.json"); err != nil {
		panic(err)
	} else if err := json.Unmarshal(bytes, &interests); err != nil {
		panic(err)
	}
}
