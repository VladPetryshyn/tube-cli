package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ktr0731/go-fuzzyfinder"
)

type Video struct {
  Title, Author, Url, Date, Views string
}

func main() {
  if len(os.Args) < 2 {
    fmt.Println("Please enter your request")
    os.Exit(1)
  }

  query := strings.Join(os.Args[1:], "+")
  uri := fmt.Sprintf("https://invidious-us.kavin.rocks/search?q=%s", query)

  response, err := http.Get(uri)
  if err != nil {
    log.Printf("Error on request: %v\n", err)
  }

  document, err := goquery.NewDocumentFromReader(response.Body)

  if err != nil {
    log.Printf("Error when parsing html: %v\n", err)
  }
  var list []Video

  document.Find(".pure-u-1 .pure-u-md-1-4 .h-box").Each(func(i int, s *goquery.Selection) {
    item := Video{}
    item.Title = s.Find("a [dir=auto]").Text()
    url, ok := s.Find("a").Attr("href")
    if ok {
      item.Url = fmt.Sprintf("https://youtube.com%s", url)
    }
    item.Author = s.Find("p.channel-name").Text()
    videInfo := s.Find(".video-card-row.flexible:last-child")
    item.Date = videInfo.Find(".flex-left p").Text()
    item.Views = videInfo.Find(".flex-right p").Text()
    if item.Author != "" {
      list = append(list, item)
    }
  })

  result, err := fuzzyfinder.Find(list, func(i int) string {
    return list[i].Title
  }, fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
        if i == -1 {
            return ""
        }
        return fmt.Sprintf("Author: %s\nViews: %s\n%s",
                list[i].Author,
                list[i].Views,
                list[i].Date,
                )
    }))

  if err == nil {
    exec.Command("mpv", list[result].Url).Run()
  }
}
