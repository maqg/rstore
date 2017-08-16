package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpGet(urls string) int {
	resp, err := http.Get(urls)
	if err != nil {
		fmt.Printf("get url %s error\n", urls)
		return -1
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read body from url %s error\n", urls)
		return -1
	}

	fmt.Println(string(body))

	return 0
}

func httpPost(urls string) int {

	resp, err := http.PostForm(urls,
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		fmt.Printf("Post data to url %s error\n", urls)
		return -1
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read body from url %s error\n", urls)
		return -1
	}

	fmt.Println(string(body))

	return 0
}

func httpDo() {

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func main() {
	httpGet("http://www.baidu.com")
	httpPost("http://localhost:8000/v1/xxxx/manafiest/")
}
