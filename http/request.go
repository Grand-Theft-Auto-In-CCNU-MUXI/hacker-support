package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//bodyType 常量
const (
	Default = 0
	File    = 1
)

type Request struct {
	Content   string `json:"content"`    // 目前只有 checkpoint3 需要这个请求，都放这
	ExtraInfo string `json:"extra_info"` // 此字段暂时不用
}

type HttpRequest struct {
	req      *http.Request
	BodyJSON Request // 上述的 request 格式
}

func (r *HttpRequest) ShowBody() {
	body := io.Reader(r.req.Body)
	b, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Err:", err.Error())
	}

	fmt.Println(string(b))
}

func (r *HttpRequest) ShowHeader() {
	header := make(map[string][]string)
	for key, value := range r.req.Header {
		if key == "code" {
		} else {
			header[key] = value
		}
	}
	fmt.Println(header)
}

func (req *HttpRequest) AddHeader(key string, value string) {
	req.req.Header.Add(key, value)
}

func (req *HttpRequest) SetHeader(key string, value string) {
	req.req.Header.Set(key, value)
}

// bodyType has two choice `file` or `default`
func NewRequest(method, url, body string, bodyType int) (req *HttpRequest, err error) {
	var r io.Reader

	if bodyType == File {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		file, errFile := os.Open(body)
		if errFile != nil {
			fmt.Println(errFile)
			return
		}
		defer file.Close()
		part1, errFile := writer.CreateFormFile("file", filepath.Base(body))
		if errFile != nil {
			fmt.Println(errFile)
			return
		}

		_, errFile = io.Copy(part1, file)
		if errFile != nil {
			fmt.Println(errFile)
			return
		}
		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		req.req, err = http.NewRequest("POST", url, payload)
		if err != nil {
			fmt.Println("create a new request failed")
			return
		}
		req.req.Header.Set("Content-Type", writer.FormDataContentType())
	} else if bodyType == Default {
		r = strings.NewReader(body)

		if method == "" {
			method = "GET"
		}
		if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "CATCH" {
			return nil, errors.New("your method is wrong")
		}

		req.req, err = http.NewRequest(method, url, r)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("your bodytype is wrong")
	}

	user_name := os.Getenv("USERNAME")
	req.AddHeader("user_name", user_name)

	return
}

func SendRequest(request *HttpRequest) (response *HttpResponse, err error) {
	client := http.Client{}
	response.raw, err = client.Do(request.req)
	if err != nil {
		return nil, err
	}

	r, err := GetResponse(response)

	if err != nil {
		fmt.Println(err, r)
		return nil, errors.New("get response failed")
	}

	if response.raw.StatusCode == 200 {
		fmt.Println("request success the data is: ", r.Data.Text)
	} else if response.raw.StatusCode == 401 {
		fmt.Println("failed! the wrong data is: ", r.Message)
	} else if response.raw.StatusCode == 404 {
		fmt.Println("we can not find the path, did you input the right information? the wrong message is: ", r.Message)
	}

	return

}
