package httptool

import (
	"encoding/base64"
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
	fmt.Println("response body:")
	fmt.Println("1.Message:")
	fmt.Println(r.Body.Message)
	fmt.Println("2.Text:")
	fmt.Println(r.Body.Data.Text)
	fmt.Println("3.ExtraInfo:")
	fmt.Println(r.Body.Data.ExtraInfo)
}

// ShowHeader ... 输出 response header
func (r *HttpResponse) ShowHeader() {
	fmt.Println("response header:")
	for key, value := range r.raw.Header {
		if key != "code" {
			for _, v := range value {
				fmt.Println(key + " : " + v)
			}
		}
	}
}

// GetHeader ... 获取请求头
func (r *HttpResponse) GetHeader(key string) ([]string, error) {
	if key == "" {
		return nil, errors.New("GetHeader Err : key is invalid")
	}
	return r.raw.Header.Values(key), nil
}

// DownloadFile ... 下载文件
func (r *HttpResponse) Save(path string) (err error) {
	if path == "" {
		return errors.New("path is invalid")
	}

	imageBytes, err := base64.StdEncoding.DecodeString(r.Body.Data.ExtraInfo)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, imageBytes, 0777)
	if err != nil {
		return err
	}

	// fmt.Println(fmt.Sprintf("download file successed at %s", path))

	return nil
}

func resolveResponse(response *HttpResponse) (err error) {
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
