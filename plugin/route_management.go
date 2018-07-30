package plugin

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"brief_framework/logger"
	"encoding/json"
)

type GetHandler func(map[string]string) interface{}
type PostHandler func([]byte) interface{}
type PostHandlerMap func(map[string]string) interface{}

type PathHandler struct {
	Path string
	Handler interface{}
}

var ResSlice []*PathHandler


//Add handler to map
func RegisterHandler(relativePath string, p interface{}) {
	ResSlice = append(ResSlice, &PathHandler{relativePath, p})
}

func getQueryMap(c *gin.Context) map[string]string {
	m := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		m[k] = v[0]
	}

	return m
}

func getPostByteSlice(c *gin.Context) []byte {
	buf := make([]byte, 4096)
	n, err := c.Request.Body.Read(buf)
	if err == io.EOF {
		//read all
	} else {
		log.Printf("Read Error: %s\n", err)
	}
	c.Set("postByte", string(buf[0:n]))
	return buf[0:n]
}

func getPostMap(c *gin.Context) map[string]string {
	m := make(map[string]string)
	for k, v := range c.Request.PostForm {
		m[k] = v[0]
	}
	return m
}

func DoGetProcess(c *gin.Context, p GetHandler) {
	obj := p(getQueryMap(c))
	doFinalProcess(c, obj)
}

func DoPostProcess(c *gin.Context, p PostHandler) {
	obj := p(getPostByteSlice(c))
	doFinalProcess(c, obj)
}

func DoPostProcessMap(c *gin.Context, p PostHandlerMap) {
	obj := p(getPostMap(c))
	doFinalProcess(c, obj)
}

func doFinalProcess(c *gin.Context, obj interface{}) {
	r, err := json.Marshal(obj)
	var s string
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusOK, string(r))
	}
	c.Set("returnString", string(s)) //for logging
}

func InitHandler(router *gin.Engine)  {
	for _, v := range ResSlice {
		if p, ok := v.Handler.(GetHandler); ok {
			router.GET(v.Path, func(c *gin.Context) {
				DoGetProcess(c, p)
			})
		} else if p, ok := v.Handler.(PostHandler); ok {
			router.POST(v.Path, func(c *gin.Context) {
				DoPostProcess(c, p)
			})
		} else if p, ok := v.Handler.(PostHandlerMap); ok {
			router.POST(v.Path, func(c *gin.Context) {
				DoPostProcessMap(c, p)
			})
		} else {
			logger.Instance().Warn("Handler can not recognize")
		}
	}
}
