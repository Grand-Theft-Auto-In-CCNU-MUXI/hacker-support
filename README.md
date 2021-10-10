# hacker-support
The tools you may need to make the game easier.

### 1. import 导入

```go
// 分开导入比较方便
import "github.com/Grand-Theft-Auto-In-CCNU-MUXI/hacker-support/httptool" // http 模块
import "github.com/Grand-Theft-Auto-In-CCNU-MUXI/hacker-support/encrypt" // 加密解密模块
```

### 2. hacker-support 提供的类型

#### 常量

```go
// request body 的类型，支持通常字符串和文件
DEFAULT = 0
FILE    = 1

// http 方法
GETMETHOD    = "GET"
POSTMETHOD   = "POST"
PUTMETHOD    = "PUT"
DELETEMETHOD = "DELETE"
PATCHMETHOD = "PATCH"
```

#### 封装类型

```go
type Request struct {
	Content   string `json:"content"`    
	ExtraInfo string `json:"extra_info"` 
} //请求游戏服务端的底层格式

type HttpRequest struct {
	Req      *http.Request
	Body     *Request // 上述的 request 格式
	BodyType int
} // 用户的请求类型, 用户用这个结构体就可以了

type TextInfo struct {
	Text      string `json:"text"`       
	ExtraInfo string `json:"extra_info"`
} // 服务端响应的底层格式

type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    TextInfo `json:"data"`
} // 响应内容

type HttpResponse struct {
	Body Response `json:"body"`
	Raw  string
	raw  *http.Response
} // 用户的响应类型，用户用这个结构体就可以了
```

### 3. 函数

* NewRequest 建立新的 http 请求

  ```go
  // method 为 http 方法
  // url 为请求地址
  // body 为请求 body
  // bodyType 为 body 的类型，这里支持一般格式和文件格式两种
  // 通过 httptool.DEFAULT 和 httptool.FILE 引用
  // 返回 响应类型指针与error类型
  func NewRequest(method, url, body string, bodyType int) (*HttpRequest, error)
  
  // 调用例1
  	request, err := httptool.NewRequest(
  		httptool.GETMETHOD,
  		"http://127.0.0.1:8080/api/v1/organization/code",
  		"",
  		httptool.DEFAULT)
  	// handle err
  
  // 调用例2
  	request, err := httptool.NewRequest(
  		httptool.POSTMETHOD,
  		"http://127.0.0.1:8080/api/v1/organization/code",
  		"C:\Users\username\Desktop\test.jpg",
  		httptool.FILE)
  	// handle err
  ```

+ aes256加密

  ```go
  // plaintext 为待加密的字符切片
  // key 为密钥
  // 返回加密后的字符切片与error类型
  func AESEncryptOutInBase64(plaintext []byte, key []byte) ([]byte, error)
  
  // 调用例
  	data, err := encrypt.AESEncryptOutInBase64([]byte(errorCode), []byte(secretKey))
  	// handle err
  ```

+ aes256解密

  ```go
  // plaintext为待解密的字符切片
  // key为密钥
  // 返回解密后的字符切片与error类型
  func AESDecodeAfterBase64(data []byte, key []byte) ([]byte, error)
  
  // 调用例
  	errorCode, err := encrypt.AESDecodeAfterBase64(data, []byte(secretKey))
  	// handle err
  ```

+ base64解密

  ```go
  // str为待解密的字符串
  // 返回解密后的字符与error类型
  func Base64Decode(str string) (string, error)
  
  // 调用例
  	errorCode, err := encrypt.Base64Decode(data）
  	// handle err
  ```

+ base64加密

  ```go
  // plaintext为待加密的字符切片
  // 返回加密后的字符串
  func Base64Encode(plaintext []byte) string
  
  // 调用例
  data := encrypt.Base64Encode([]byte(errorCode))
  	// handle err
  ```

### 4. 方法

#### HttpRequest

+ SendRequest 发送**请求**

  ```go
  // r是用户请求类型结构体指针
  // 返回响应类型结构体指针与error类型
  func (r *HttpRequest) SendRequest() (*HttpResponse, error)
  
  // 调用例 (request是*HttpRequest类型变量)
  	resp, err := request.SendRequest()
  	// handle err
  ```

+ ShowBody打印**请求**的body部分

  ```go
  // r是用户请求类型结构体指针
  func (r *HttpRequest) ShowBody()
  
  // 调用例 (request是*HttpRequest类型变量)
  	request.ShowBody()
  ```

+ ShowHeader打印**请求**的header部分

  ```go
  // r是用户请求类型结构体指针
  func (r *HttpRequest) ShowHeader()
  
  // 调用例 (request是*HttpRequest类型变量)
  	request.ShowHeader()
  ```

+ AddHeader添加**请求头**

  ```go
  // r是用户请求类型结构体指针
  // key为header中需要新添加的成员的名字，value为它的值
  func (r *HttpRequest) AddHeader(key string, value string)
  
  // 调用例 (request是*HttpRequest类型变量)
  request.AddHeader("passport", "1234abcd")
  ```

+ SetHeader更换**请求头**

  ```go
  // r是用户请求类型结构体指针
  // key为header中需要更改的成员的名字，value为它的新值
  func (r *HttpRequest) SetHeader(key string, value string)
  
  // 调用例 (request是*HttpRequest类型变量)
  request.AddHeader("passport", "56789efgh")
  ```


#### HttpResponse

+ ShowBody 打印响应的body部分

  ```go
  // r是用户请求类型结构体指针
  func (r *HttpResponse) ShowBody()
  
  // 调用例 (resp是*HttpResponse类型变量)
  resp.ShowBody()
  ```

+ ShowHeader打印响应的header部分

  ```go
  // r是用户请求类型结构体指针
  func (r *HttpResponse) ShowHeader()
  
  // 调用例 (resp是*HttpResponse类型变量)
  resp.ShowHeader()
  ```

+ GetHeader获取响应头中某一成员的值

  ```go
  // r是响应类型结构体指针
  // key为想要获取的成员名称
  // 返回字符串数组（key的值）与error类型
  func (r *HttpResponse) GetHeader(key string) ([]string, error)
  
  // 调用例 (resp是*HttpResponse类型变量)
  	value, err := resp.GetHeader("Passport")
  	// handle err
  ```

+ Save下载文件

  ```go
  // r是响应类型结构体指针
  // path为下载的文件路径
  // 返回error类型
  func (r *HttpResponse) Save(path string) (err error)
  
  // 调用例 (resp是*HttpResponse类型变量)
  	err := resp.Save("C:\Users\username\Desktop\test.jpg")
  	// handle err
  ```


  
