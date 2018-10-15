package main

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
	"sync"
)

var (
	user = "3183015853@qq.com" //发送账户
	pwd  = "ygtfinsitxyxdeac"  //密码
	uto  = "273227543@qq.com"  //默认接收账户
	to   = []string{uto}       //接收账户
	mu   sync.Mutex            //互斥锁，防止多进程数据出错
)

//SendMail 发送提醒邮件
//默认收发账户为...
//参数：title 邮件名
//参数：date	邮件内容
func SendMail(title, date string) {
	fmt.Println("准备发送短信")
	nickname := "test" //邮件名
	subject := "数据报告"  //邮件标题
	body := "数据抓取出错了." //邮件内容
	if title != "" {
		subject = title
		nickname = title
	}
	if date != "" {
		body = date
	}

	auth := smtp.PlainAuth("", user, pwd, "smtp.qq.com")
	contentType := "Content-Type: text/plain; charset=UTF-8"
	mu.Lock()
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	mu.Unlock()
	err := smtp.SendMail("smtp.qq.com:587", auth, user, to, msg)
	if err != nil {
		//发送失败， 删除txt记录
		ioutil.WriteFile("1.txt", []byte(""), 0777)
		fmt.Println("短信发送失败，失败原因:", err)
		return
	}
	//发送成功
	fmt.Println("发送邮件成功!")
}

//Set 设置收发账户
//参数：u 发送账户
//参数：p 发送账户密码
func Set(u, p string) {
	if u != "" {
		mu.Lock()
		user = u
		pwd = p
		mu.Unlock()
	}
}

//AddTo 增加接收账户
//参数： ut 接收账户
func AddTo(ut string) {

	if ut != "" {
		mu.Lock()
		to = append(to, ut)
		fmt.Println(to)
		mu.Unlock()
	}
}

//SubTo 减少接收账户
//参数： ut 接收账户
func SubTo(ut string) {

	if ut == "" {
		return
	}
	mu.Lock()
	for i := 0; i < len(to); i++ {
		if to[i] == ut {
			to = append(to[:i], to[i+1:]...)
			fmt.Println(to)
		}
	}
	mu.Unlock()
}

//GetTo 获取接收账户
//返回字符串切片：[]string
func GetTo() []string {
	return to
}
