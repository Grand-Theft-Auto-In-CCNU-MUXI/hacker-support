package hackersupport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

type HttpResponse struct {
	body Response
	raw  *http.Response
}

type Request struct {
	req *http.Request
}

func (r *Request) ShowBody() {
	body := io.Reader(r.req.Body)
	b, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Err:", err.Error())
	}

	fmt.Println(string(b))
}

func (r *HttpResponse) ShowBody() {
	body := io.Reader(r.raw.Body)
	b, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Err:", err.Error())
	}

	fmt.Println(string(b))
}

func (r *Request) ShowHeader() {
	header := make(map[string][]string)
	for key, value := range r.req.Header {
		if key == "code" {
		} else {
			header[key] = value
		}
	}
	fmt.Println(header)
}

func (h *HttpResponse) ShowHeader() {
	header := make(map[string][]string)
	for key, value := range h.raw.Header {
		if key == "code" {
		} else {
			header[key] = value
		}
	}
	fmt.Println(header)
}

func (h *HttpResponse) GetHeader(key string) []string {
	return h.raw.Header.Values(key)
}

func (req *Request) AddHeader(key string, value string) {
	req.req.Header.Add(key, value)
}

func (req *Request) SetHeader(key string, value string) {
	req.req.Header.Set(key, value)
}

// bodyType has two choice `file` or `default`
func NewRequest(method, url, body, bodyType string) (req *Request, err error) {
	var r io.Reader

	if bodyType == "file" {
		r, err = os.Open(body)

		if err != nil {
			return nil, errors.New("open file failed!")
		}
	} else {
		r = strings.NewReader(body)
	}

	if method == "" {
		method = "GET"
	}
	if method != "GET" || method != "POST" || method != "PUT" || method != "DELETE" || method != "CATCH" {
		return nil, errors.New("your method is wrong")
	}

	req.req, err = http.NewRequest(method, url, r)

	if err != nil {
		return nil, err
	}

	user_name := os.Getenv("USERNAME")
	req.AddHeader("user_name", user_name)

	return
}

func SendRequest(request http.Request) (response *HttpResponse, err error) {
	client := http.Client{}
	response.raw, err = client.Do(&request)
	if err != nil {
		return nil, err
	}
	return
}

func GetResponse(response *HttpResponse) (resp Response, err error) {
	body, err := ioutil.ReadAll(response.raw.Body)

	defer response.raw.Body.Close()

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}

	return

}
