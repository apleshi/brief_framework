## 基于gin实现的简单webserver框架
### 使用前准备工作

首先获取需要导入的包

```sh
$ go get github.com/gin-gonic/gin
$ go get github.com/go-redis/redis
$ go get github.com/Unknwon/goconfig
```
注：如果导入失败，可以尝试执行如下命令 :
```sh
$ yum update nss curl libcurl
```
### 使用

1. import

```go
import "brief_framework"
```

2. 实现一个函数

Get方法，需要实现GetHandle

```go
type GetHandler func(map[string]string)(interface{})
```

参数为Query中的参数, 需要返回一个JSON的结构体. 如: 

```go
var get_handle plugin.GetHandler = func(m map[string]string) interface{} {
	type JsonHolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	return JsonHolder{Id: 77, Name: "Get_handle"}
}
```

3. init中调用AddProcess

第一个参数是相对路径，如：

```go
func init() {
	plugin.AddProcess("/get", sampleGet)
	plugin.AddProcess("/post", samplePost)
}
```

4. 启动
```go
func main() {
	brief_framework.Serve()
} 
```
