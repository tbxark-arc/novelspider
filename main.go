package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	id := flag.Int("id", 0, "id of the books")
	file := flag.String("file", "book.txt", "path to save the books")
	flag.Parse()
	if *id == 0 {
		log.Panic("Please provide the id of the books")
	}
	category, err := parseCategory(*id)
	if err != nil {
		log.Panic(err)
	}
	save, err := os.Create(*file)
	if err != nil {
		log.Panic(err)
	}
	defer save.Close()
	for _, book := range category {
		log.Printf("Downloading %s", book)
		title, content, e := parseBook(book)
		if e != nil {
			log.Panic(e)
		}
		_, e = save.WriteString(title + "\n")
		if e != nil {
			log.Panic(e)
		}
		_, e = save.WriteString(content + "\n")
		if e != nil {
			log.Panic(e)
		}
	}
	log.Println("Downloaded all books")
}

func parseCategory(id int) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://69shuba.cx/book/"+strconv.Itoa(id)+"/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0.1 Safari/605.1.15")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	var booksLinks []string
	doc.Find(".catalog ul li").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")
		if strings.HasPrefix(link, "http") {
			booksLinks = append(booksLinks, link)
		}
	})
	sort.Sort(sort.StringSlice(booksLinks))
	return booksLinks, nil
}

func parseBook(link string) (title, content string, err error) {
	title, content = "", ""
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0.1 Safari/605.1.15")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	// 将 GBK 或 GB18030 转换为 UTF-8
	utf8Body, err := simplifiedchinese.GB18030.NewDecoder().Bytes(body) // 或使用 simplifiedchinese.GBK
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(utf8Body))
	if err != nil {
		return
	}
	title = doc.Find("h1").Text()
	doc.Find(".txtnav p").Each(func(i int, s *goquery.Selection) {
		content += strings.TrimSpace(s.Text()) + "\n"
	})
	return
}
