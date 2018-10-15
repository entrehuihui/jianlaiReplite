package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var accept = "ext/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
var acceptEncoding = "gzip, deflate, br"
var acceptLanguage = "zh-CN,zh;q=0.9"
var cacheControl = "max-age=0"
var connection = "keep-alive"
var cookie = "_ga=GA1.2.640308922.1517208694; _gid=GA1.2.1969781450.1517208694; Hm_lvt_9734b08ecbd8cf54011e088b00686939=1517208694; LXB_REFER=www.baidu.com; Hm_lvt_50b0b45724f4f39e2a94cb8af0e9b547=1517208709; SMM_auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjZWxscGhvbmUiOiIxMzAyNzk0MTg4MCIsImNvbXBhbnlfaWQiOjAsImNvbXBhbnlfc3RhdHVzIjowLCJjcmVhdGVfYXQiOjE1MTcyMDkxMDksImVtYWlsIjoiIiwiZW5fZW5kX3RpbWUiOjAsImVuX3JlZ2lzdGVyX3N0ZXAiOjEsImVuX3JlZ2lzdGVyX3RpbWUiOjAsImVuX3N0YXJ0X3RpbWUiOjAsImVuX3VzZXJfdHlwZSI6MCwiZW5kX3RpbWUiOjAsImlzX21haWwiOjAsImlzX3Bob25lIjoxLCJsYW5ndWFnZSI6ImhxIiwibHlfZW5kX3RpbWUiOjAsImx5X3N0YXJ0X3RpbWUiOjAsImx5X3VzZXJfdHlwZSI6MCwicmVnaXN0ZXJfdGltZSI6MTUxNzIwOTEwOSwicyI6MjAsInN0YXJ0X3RpbWUiOjAsInVzZXJfaWQiOjEyNDM5ODUsInVzZXJfbmFtZSI6IlNNTTE1MTcyMDkxMDlaQyIsInVzZXJfdHlwZSI6MCwienhfZW5kX3RpbWUiOjAsInp4X3N0YXJ0X3RpbWUiOjAsInp4X3VzZXJfdHlwZSI6MH0.bPvDyGFCPxIu4lYIh0RrLjUozu-vtKN_GcePDo20lm4; _gat=1; _gat_UA-102039857-2=1; Hm_lpvt_9734b08ecbd8cf54011e088b00686939=1517210593; Hm_lpvt_50b0b45724f4f39e2a94cb8af0e9b547=1517210593"
var host = "hq.smm.cn"
var upgradeInsecureRequests = "1"

var url = "https://www.baidu.com/s?wd=%E5%89%91%E6%9D%A5"

//限制进程通道
//var dateChan = make(chan []string, 5)

var userAgent = [...]string{
	"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)"}

//GetRandomUserAgent 设置header user-Agent
func GetRandomUserAgent() string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return userAgent[r.Intn(len(userAgent))]
}

//保存最新章节
var chapter string

func main() {
	oldchapter, err := ioutil.ReadFile("1.txt")
	if err == nil {
		chapter = string(oldchapter)
	}
	fmt.Println(chapter, "#########")
	ticker := time.NewTicker(1 * time.Minute)
	for _ = range ticker.C {
		check()
	}
}

func check() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", accept)
	//req.Header.Add("Accept-Encoding", acceptEncoding)	//加上会乱码
	req.Header.Add("Accept-Language", acceptLanguage)
	req.Header.Add("Cache-Control", cacheControl)
	req.Header.Add("Connection", connection)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Host", host)
	req.Header.Add("Upgrade-Insecure-Requests", upgradeInsecureRequests)
	req.Header.Add("User-Agent", GetRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	query, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}
	sq := query.Find("div.container_s").Find("#content_left").Find(".result-op").Find(".op_tb_more")
	s := sq.Find("p").First().Find("span").First().Text()
	s1 := sq.Find("p").First().Find("a").Text()
	s2, _ := sq.Find("p").First().Find("a").Attr("href")
	fmt.Println(s)
	fmt.Println(s1)
	fmt.Println(s2)
	if s1 == "" {
		return
	}
	if strings.Compare(s1, chapter) == 0 {
		return
	}
	chapter = s1
	var buffer bytes.Buffer
	buffer.WriteString(s)
	buffer.WriteString("\n")
	buffer.WriteString(s1)
	buffer.WriteString("\n")
	buffer.WriteString(s2)
	buffer.WriteString("\n")
	ioutil.WriteFile("1.txt", []byte(chapter), 0777)
	SendMail(s1, buffer.String())
}
