package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0.1 Safari/605.1.15"
	BaseURL   = "https://69shuba.cx"
	Retry     = 3
)

func main() {
	id := flag.Int("id", 29624, "id of the books")
	dir := flag.String("dir", "./", "directory to save the books")
	startIndex := flag.Int("start", -1, "start index of the books")
	interval := flag.Int("interval", 1, "seconds between each download")

	flag.Parse()
	if *id == 0 {
		log.Panic("Please provide the id of the books")
	}
	title, category, err := parseCategory(*id)
	if err != nil {
		log.Panic(err)
	}
	fileName := filepath.Join(*dir, title+".txt")
	file, err := openOrCreateFile(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
DOWNLOAD:
	for idx, book := range category {
		if idx < *startIndex {
			continue
		}
		if !strings.HasPrefix(book, "http") {
			continue
		}
		for i := 0; i < Retry; i++ {
			if downloadCategory(book, file) == nil {
				time.Sleep(time.Duration(*interval) * time.Second)
				continue DOWNLOAD
			}
			log.Printf("Retry %d times for %s", i+1, book)
			time.Sleep(time.Second)
		}
		log.Panicf("Failed to download %s", book)
	}
	log.Println("Downloaded all books")
}

func openOrCreateFile(fileName string) (*os.File, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return os.Create(fileName)
	} else {
		return os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	}
}

func downloadCategory(book string, save *os.File) error {
	log.Printf("Downloading %s", book)
	cateTitle, content, err := parseBook(book)
	log.Printf("Downloaded %s", cateTitle)
	if err != nil {
		return err
	}
	_, err = save.WriteString(cateTitle + "\n")
	if err != nil {
		return err
	}
	_, err = save.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return nil
}

func parseCategory(id int) (title string, category []string, err error) {
	req, err := http.NewRequest(http.MethodGet, BaseURL+"/book/"+strconv.Itoa(id)+"/", nil)
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
		link, ok := s.Find("a").Attr("href")
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
		booksLinks[num] = link
	})
	category = make([]string, 0, maxNum-minNum+1)
	for i := minNum; i <= maxNum; i++ {
		category = append(category, booksLinks[i])
	}
	title = doc.Find(".bread a").Last().Text()
	return
}

func parseBook(link string) (title, content string, err error) {
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
