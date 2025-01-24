package main

import (
	"flag"
	"fmt"
	"github.com/tbxark-arc/novelspider/loader"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Retry = 3
)

type Config struct {
	link       string
	dir        string
	startIndex int
	interval   int
}

func main() {
	var conf Config
	flag.StringVar(&conf.link, "link", "", "link of the books")
	flag.StringVar(&conf.dir, "dir", "./", "directory to save the books")
	flag.IntVar(&conf.startIndex, "start", -1, "start index of the books")
	flag.IntVar(&conf.interval, "interval", 1, "seconds between each download")
	parser := flag.String("parser", "", "parser of the books, support: shuba69, sksw")
	help := flag.Bool("help", false, "show help")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	if conf.link == "" {
		log.Panic("Please provide the link of the books")
	}
	load, err := createParser(*parser)
	if err != nil {
		log.Panic(err)
	}
	err = startDownload(load, &conf)
	if err != nil {
		log.Panic(err)
	}
}

func createParser(name string) (loader.Parser, error) {
	switch name {
	case "shuba69":
		return loader.NewShuBa69(), nil
	case "sksw":
		return loader.NewSksw(), nil
	default:
		return nil, fmt.Errorf("unknown parser %s", name)
	}
}

func startDownload(parser loader.Parser, conf *Config) error {
	title, category, err := parser.Category(conf.link)
	if err != nil {
		return err
	}
	fileName, err := filepath.Abs(filepath.Join(conf.dir, title+".txt"))
	if err != nil {
		return err
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	log.Printf("Create file %s", fileName)
	log.Printf("Start download %s", title)
	defer file.Close()
DOWNLOAD:
	for idx, book := range category {
		if idx < conf.startIndex {
			continue
		}
		if !strings.HasPrefix(book, "http") {
			continue
		}
		for i := 0; i < Retry; i++ {
			if loadAndWriteCategory(parser, book, file) == nil {
				time.Sleep(time.Duration(conf.interval) * time.Second)
				continue DOWNLOAD
			}
			log.Printf("Retry %d times for %s", i+1, book)
			time.Sleep(time.Second)
		}
		log.Panicf("Failed to download %s", book)
	}
	log.Println("Downloaded all books")
	return nil
}

func loadAndWriteCategory(parser loader.Parser, book string, save *os.File) error {
	log.Printf("Downloading %s", book)
	cateTitle, content, err := parser.Book(book)
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
	return save.Sync()
}
