package loader

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"strconv"
	"strings"
)

type ShuBa69 struct {
}

func NewShuBa69() *ShuBa69 {
	return &ShuBa69{}
}

func (s *ShuBa69) Category(link string) (title string, category []string, err error) {
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
	var minNum, maxNum int
	booksLinks := make(map[int]string)
	doc.Find(".catalog ul li").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a").Attr("href")
		if !ok {
			return
		}
		dataNum, ok := s.Attr("data-num")
		if !ok {
			return
		}
		num, e := strconv.Atoi(dataNum)
		if e != nil {
			return
		}
		if num < minNum {
			minNum = num
		}
		if num > maxNum {
			maxNum = num
		}
		booksLinks[num] = href
	})
	category = make([]string, 0, maxNum-minNum+1)
	for i := minNum; i <= maxNum; i++ {
		category = append(category, booksLinks[i])
	}
	title = doc.Find(".bread a").Last().Text()
	return
}

func (s *ShuBa69) Book(link string) (title, content string, err error) {
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
	title = doc.Find("h1").Text()
	doc.Find(".txtnav p").Each(func(i int, s *goquery.Selection) {
		content += strings.TrimSpace(s.Text()) + "\n"
	})
	return
}
