## 基于gin实现的简单webserver框架
### 使用前准备工作

首先获取需要导入的包

```sh
$ go get github.com/gin-gonic/gin
$ go get github.com/go-redis/redis
```
### 使用

1. import

```go
import "brief_framework"
```

2. 实现一个函数

Get方法，需要实现GetHandle

```go
type GetHandle func(map[string]string)(interface{})
```

参数为Query中的参数, 需要返回一个JSON的结构体. 如: 

```go
var get_handle server_plugin.GetHandle = func(m map[string]string) interface{} {
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
	server_plugin.AddProcess("/get", SampleGet)
	server_plugin.AddProcess("/post", SamplePost)
}
```
4. 启动
```go
func main() {
	brief_framework.Serve()
} 
```
