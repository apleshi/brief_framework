package brief_framework

import (
	"testing"
	"server_plugin"
)


func init() {
	server_plugin.AddProcessor("/test/abc", postpHandle)
	server_plugin.AddProcessor("/test/abc", getHandle)
	server_plugin.AddProcessor("/test/aaa", getHandle)
}


// 处理多层实验的函数
var postpHandle server_plugin.PostHandler = func(ba []byte) interface{} {
	type JsonHolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		kkk  []string
	}
	return JsonHolder{Id: 1, Name: "post"}
}

var getHandle server_plugin.GetHandle = func(reqMap map[string]string) interface{} {
	type JsonHolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		kkk  []string
	}
	return JsonHolder{Id: 2, Name: "get"}
}

func TestConfig(t *testing.T) {

	Serve()

}

