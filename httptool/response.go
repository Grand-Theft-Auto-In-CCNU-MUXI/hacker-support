package httptool

import (
	"encoding/json"
	"errors"
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
	Message string   `json:"message"`
	Data    TextInfo `json:"data"`
}

type HttpResponse struct {
	Body Response `json:"body"`
	raw  *http.Response
}

// ShowBody ... 输出 response body
// TODO: 美化输出
func (r *HttpResponse) ShowBody() {
	fmt.Println("Message:")
	fmt.Println(r.Body.Message)
	fmt.Println("Text:")
	fmt.Println(r.Body.Data.Text)
	fmt.Println("ExtraInfo:")
	fmt.Println(r.Body.Data.ExtraInfo)
}

// ShowHeader ... 输出 response header
func (r *HttpResponse) ShowHeader() {
	header := make(map[string][]string)
	for key, value := range r.raw.Header {
		if key == "code" {
		} else {
			header[key] = value
		}
	}
	fmt.Println(header)
}

// GetHeader ... 获取请求头
func (r *HttpResponse) GetHeader(key string) ([]string,error) {
	if key == "" {
		return nil,errors.New("GetHeader Err : key is invalid")
	}
	return r.raw.Header.Values(key),nil
}

// DownloadFile ... 下载文件
func DownloadFile(r *HttpResponse, path string) (err error) {
	if path == "" {
		return errors.New("path is invalid")
	}

	imageBytes := r.Body.Data.ExtraInfo

	err = ioutil.WriteFile(path, []byte(imageBytes), 0777)

	if err != nil {
		return err
	}

	// fmt.Println(fmt.Sprintf("download file successed at %s", path))

	return nil
}

func ResolveResponse(response *HttpResponse) (err error) {
	body, err := ioutil.ReadAll(response.raw.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("GetResponse close body failed")
		}
	}(response.raw.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response.Body)
	if err != nil {
		return
	}

	return

}
