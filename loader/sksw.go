package loader

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"net/url"
)

type Sksw struct {
	userAgent string
}

func NewSksw() *Sksw {
	return &Sksw{}
}

func (s *Sksw) Category(link string) (title string, category []string, err error) {
	// example https://69shuba.cx/book/74678.htm
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(simplifiedchinese.GB18030.NewDecoder().Reader(res.Body))
	if err != nil {
		return
	}
	category = make([]string, 0, 100)
	uri, _ := url.Parse(link)
	doc.Find(".list-group.list-charts li").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a").Attr("href")
		if !ok {
			return
		}
		category = append(category, uri.Scheme+"://"+uri.Host+href)
	})
	title = doc.Find(".bread a").Last().Text()
	return
}

func (s *Sksw) Book(link string) (title, content string, err error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}
	doc, err := goquery.NewDocumentFromReader(simplifiedchinese.GB18030.NewDecoder().Reader(res.Body))
	if err != nil {
		return
	}
	title = doc.Find(".panel-heading").First().Text()
	content = doc.Find(".panel-body.content-body.content-ext").Text()
	return
}
