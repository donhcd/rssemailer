package main

import (
	"bytes"
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
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

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	for _, newItem := range newitems {
		ciTit := strings.ToLower(newItem.Title)
		isGood := false
		for _, wantedSubString := range strings.Split("sfo,san francisco,sjc,san jose,oak,oakland", ",") {
			if strings.Contains(ciTit, wantedSubString) {
				isGood = true
			}
		}
		if !isGood {
			continue
		}
		cmd := exec.Command("mail", "-s", "Flight deal: "+newItem.Title, "don.hcd@gmail.com")
		cmd.Stdin = bytes.NewReader([]byte(fmt.Sprintf("Link: %s\n\nDescription: %s\n", newItem.Comments, newItem.Description)))
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "[e] mail error: %s; output: %q\n", err, string(output))
		}
	}
}
