package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	Time      string `json:"time"`
	Title     string `json:"title"`
	Writer    string `json:"writer"`
	ReadCount string `json:"readCount"`
}

func spider(url string) (Page, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return Page{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return Page{}, err
	}
	defer resp.Body.Close()

	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		return Page{}, err
	}

	page := Page{}
	// body > section > section.n_container > div > div.n_right.fr > section > form > div > div.nav01
	docDetail.Find("div.nav01 ").
		Each(func(i int, s *goquery.Selection) {
			time := s.Find("h6 > span:nth-child(1)").Text()
			// body > section > section.n_container > div > div.n_right.fr > section > form > div > div.nav01 > h6 > span:nth-child(2)
			writer := strings.Trim(s.Find("h6 > span:nth-child(2)").Text(), " ")
			title := strings.Trim(s.Find("h3").Text(), " ")
			readCount := strings.TrimSpace(s.Find("span[name^=dynclicks_wbnews]").Text())
			page.ReadCount = readCount
			page.Title = title
			page.Time = time
			page.Writer = writer
		})

	return page, err
}

func main() {
	urls := []string{
		"https://news.fzu.edu.cn/info/1011/3616.htm",
		"https://news.fzu.edu.cn/info/1011/28280.htm",
		"https://news.fzu.edu.cn/info/1011/28240.htm",
		"https://news.fzu.edu.cn/info/1011/28224.htm",
		"https://news.fzu.edu.cn/info/1011/27992.htm",
		"https://news.fzu.edu.cn/info/1011/27968.htm",
		"https://news.fzu.edu.cn/info/1011/27775.htm",
		"https://news.fzu.edu.cn/info/1011/27741.htm",
		"https://news.fzu.edu.cn/info/1011/27548.htm",
		"https://news.fzu.edu.cn/info/1011/27500.htm",
	}
	ch := make(chan Page)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			page, err := spider(url)
			if err == nil {
				ch <- page
			}
		}(url)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for page := range ch {
		fmt.Printf("%s %s %s %s\n", page.Time, page.Title, page.Writer, page.ReadCount)
	}
}
