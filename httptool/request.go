package httptool

import (
	"bytes"
	"encoding/json"
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
	DEFAULT = 0
	FILE    = 1

	GETMETHOD    = "GET"
	POSTMETHOD   = "POST"
	PUTMETHOD    = "PUT"
	DELETEMETHOD = "DELETE"
	PATCHMETHOD  = "PATCH"
)

// Request ... 访问游戏服务端的底层格式
type Request struct {
	Content   string `json:"content"`    // 目前只有 checkpoint3 需要这个请求，都放这
	ExtraInfo string `json:"extra_info"` // 此字段暂时不用
}

// HttpRequest ... 提供给用户的请求类型
type HttpRequest struct {
	Req      *http.Request
	Body     *Request // 上述的 request 格式
	BodyType int
}

// ShowBody ... 打印请求 body
func (r *HttpRequest) ShowBody() {
	fmt.Println("request body:")
	if r.BodyType == FILE {
		body := io.Reader(r.Req.Body)
		b, err := ioutil.ReadAll(body)
		if err != nil {
			fmt.Println("ShowBody Err: ", err.Error())
		}
		fmt.Println(string(b))
		return
	}

	fmt.Println(r.Body.Content)
}

// ShowHeader ... 打印请求 header
func (r *HttpRequest) ShowHeader() {
	fmt.Println("request header:")
	for key, value := range r.Req.Header {
		if key != "Code" {
			for _, v := range value {
				fmt.Println(key + " : " + v)
			}
		}
	}
}

// AddHeader ... 增加请求头
func (r *HttpRequest) AddHeader(key string, value string) {
	if key == "" {
		fmt.Println("AddHeader Err: key is invalid")
		return
	}
	r.Req.Header.Add(key, value)
}

// SetHeader ... 更换请求头
func (r *HttpRequest) SetHeader(key string, value string) {
	if key == "" {
		fmt.Println("SetHeader Err: key is invalid")
		return
	}
	r.Req.Header.Set(key, value)
}

// NewRequest ... 创建请求
// bodyType has two choice `file` or `default`
func NewRequest(method, url, body string, bodyType int) (req *HttpRequest, err error) {
	if bodyType == FILE {
		req, err = handleFile(url, body)
		if err != nil {
			return nil, err
		}
	} else if bodyType == DEFAULT {
		req, err = handleDefault(method, url, body)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("your bodytype is wrong")
	}

	userName := os.Getenv("USERNAME")
	if userName == "" {
		userName = os.Getenv("USER")
	}
	req.AddHeader("code", userName)

	return
}

// SendRequest ... 发送消息，根据状态码输出提示
func (r *HttpRequest) SendRequest() (*HttpResponse, error) {
	response := new(HttpResponse)
	var err error
	client := &http.Client{}
	response.raw, err = client.Do(r.Req)
	if err != nil {
		return nil, err
	}

	if response.raw.StatusCode == 200 {
		err = resolveResponse(response)
		if err != nil {
			return nil, err
		}
		fmt.Println("Send request successfully! Please check your response body.")
		//fmt.Println("request success the data is: ")
		//fmt.Println(response.Body.Data.Text)
		//fmt.Println("the Extra info is: ")
		//fmt.Println(response.Body.Data.ExtraInfo)
	} else {
		body, err := ioutil.ReadAll(response.raw.Body)
		if err != nil {
			fmt.Println("read body error" + err.Error())
			return nil, err
		}
		if response.raw.StatusCode == 400 {
			fmt.Println("http 400 failed! the wrong data is: ")
			fmt.Println(string(body))
			return response, nil
		} else if response.raw.StatusCode == 404 {
			fmt.Println("http 404. we can not find the path, did you input the right information? the wrong message is: ")
			fmt.Println(string(body))
			return response, nil
		} else if response.raw.StatusCode == 500 {
			fmt.Println("http 500. server error, message: ")
			fmt.Println(string(body))
		}
	}

	return response, nil
}

func handleFile(url, path string) (*HttpRequest, error) {
	req := new(HttpRequest)
	req.BodyType = FILE

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile := os.Open(path)
	if errFile != nil {
		return nil, errFile
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	part1, errFile := writer.CreateFormFile("file", filepath.Base(path))
	if errFile != nil {
		return nil, errFile
	}

	_, errFile = io.Copy(part1, file)
	if errFile != nil {
		return nil, errFile
	}

	err := writer.Close()
	if err != nil {
		return nil, errFile
	}

	req.Req, err = http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.SetHeader("Content-Type", writer.FormDataContentType())

	return req, nil
}

func handleDefault(method, url, body string) (*HttpRequest, error) {
	req := new(HttpRequest)
	req.BodyType = DEFAULT
	myBody := Request{
		Content: body,
	}
	req.Body = &myBody

	bodyJSON, err := json.Marshal(myBody)
	if err != nil {
		return nil, err
	}

	payload := strings.NewReader(string(bodyJSON))

	if method == "" {
		method = "GET"
	}

	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {
		return nil, errors.New("your method is wrong")
	}

	req.Req, err = http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.AddHeader("Content-Type", "application/json")

	return req, nil
}
