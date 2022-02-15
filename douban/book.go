package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "root"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "douban_kgqa"
)

var DB *sql.DB

type BookData struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Publish string `json:"publish"`
	PubDate string `json:"pubDate"`
	Price   string `json:"price"`
	Quote   string `json:"quote"`
	Rating  string `json:"rating"`
}

func main() {
	InitDB()
	//for i := 0; i < 10; i++ {
	//	Spider(strconv.Itoa(i*25))
	//}
	ch := make(chan bool)
	for i := 0; i < 10; i++ {
		go Spider(strconv.Itoa(i*25), ch)
	}
	for i := 0; i < 10; i++ {
		<-ch
	}
}

func Spider(page string, ch chan bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://book.douban.com/top250?start="+page, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://book.douban.com/top250?icn=index-book250-all")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	docDetail.Find("#content > div > div.article > div > table").
		Each(func(i int, s *goquery.Selection) {
			var bookData BookData
			title, ok := s.Find("tbody > tr > td:nth-child(2) > div.pl2 > a").Attr("title")
			title2 := s.Find("tbody > tr > td:nth-child(2) > div.pl2 > a > span").Text()
			msg := s.Find("tbody > tr > td:nth-child(2) > p.pl").Text()
			msgSplit := strings.Split(msg, "/")
			author := msgSplit[0]
			publish := msgSplit[1]
			pubDate := msgSplit[2]
			price := msgSplit[3]
			quote := s.Find("tbody > tr > td:nth-child(2) > p.quote > span").Text()
			rating := s.Find("tbody > tr > td:nth-child(2) > div.star.clearfix > span.rating_nums").Text()

			if ok {
				bookData.Title = title + title2
				bookData.Author = author
				bookData.Publish = publish
				bookData.PubDate = pubDate
				bookData.Price = price
				bookData.Quote = quote
				bookData.Rating = rating
				fmt.Println(bookData)
				InsertSql(bookData)
			}

		})
	if ch != nil {
		ch <- true
	}
}

func InitDB() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail")
		return
	}
	fmt.Println("connnect success")
}

func InsertSql(bookData BookData) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO book_data (`Title`,`Author`,`Publish`,`PubDate`,`Price`,`Quote`,`Rating`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Prepare fail", err)
		return false
	}
	_, err = stmt.Exec(bookData.Title, bookData.Author, bookData.Publish, bookData.PubDate, bookData.Price, bookData.Quote, bookData.Rating)
	if err != nil {
		fmt.Println("Exec fail", err)
		return false
	}
	_ = tx.Commit()
	return true
}
