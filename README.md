# hacker-support
The tools you may need to make the game easier.

### Usage

* NewRequest 建立新的 http 请求

  ```go
  // method 为 http 方法
  // url 为请求地址
  // body 为请求 body
  // bodyType 为 body 的类型，这里支持一般格式和文件格式两种
  // 通过 httptool.DEFAULT 和 httptool.FILE 引用
  // 返回 ...
  func NewRequest(method, url, body string, bodyType int) (*HttpRequest, error)
  
  // 调用例
  	request, err := httptool.NewRequest(
  		"",
  		"http://127.0.0.1:8080/api/v1/organization/code",
  		"",
  		httptool.DEFAULT)
  	// handle err
  ```

  
