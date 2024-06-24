package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SendHttpGet(address string) (string, error) {
	timeout := time.Duration(time.Second * 1)

	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(address)
	// resp, err := http.Get(address)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func SendHttpPost(address string, contenttype, content string) (string, error) {
	//"application/json"
	resp, err := http.Post(address, contenttype, strings.NewReader(content))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func SendHttpPostForm(address string, values map[string]string) (string, error) {
	data := url.Values{}

	for k, v := range values {
		data[k] = []string{v}
	}

	resp, err := http.PostForm(address, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func SendHttpBody(server string, path string, body string, millisecond int64) (string, error) {

	timeout := time.Millisecond * time.Duration(millisecond)

	d := net.Dialer{Timeout: timeout, KeepAlive: timeout}

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	conn, err := d.DialContext(ctx, "tcp", server)

	if err != nil {

		return "", err
	}
	defer conn.Close()

	length := fmt.Sprintf("%d", len(body))
	data := "POST " + path + " HTTP/1.0\nUser-Agent: Cyeam\nContent-Length: " + length + "\n\n" + body + "\r\n\r\n"

	conn.Write([]byte(data))

	//接收服务端反馈
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(conn.RemoteAddr().String(), "read msg error: ", err)
		return "", err
	}
	result := string(buffer[:n])

	if len(result) > 0 {
		st := strings.Index(result, "\r\n\r\n")
		if st >= 0 {
			result = result[st:]
		} else {
			result = ""
		}
	}
	result = strings.Trim(result, "\r\n")

	return result, nil
}
