package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type TextInfo struct {
	Text      string `json:"text"`       // 服务端在这里返回每个检查点的提示
	ExtraInfo string `json:"extra_info"` // 服务端在这里返回额外，这里放图片 []byte
}

type Response struct {
	Code    int      `json:"code"`
	Data    TextInfo `json:"data"`
	Message string   `json:"message"`
}

type HttpResponse struct {
	Body Response `json:"body"`
	raw  *http.Response
}

func (r *HttpResponse) ShowBody() {
	body := io.Reader(r.raw.Body)
	b, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Err:", err.Error())
	}

	fmt.Println(string(b))
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

func DownloadFile(r *HttpResponse, path string) (err error) {
	imageBytes := r.Body.Data.ExtraInfo

	err = ioutil.WriteFile(path, []byte(imageBytes), 0777)

	if err != nil {
		fmt.Println(fmt.Sprintf("write file failed! cause %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("download file successed at %s", path))

	return nil
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
