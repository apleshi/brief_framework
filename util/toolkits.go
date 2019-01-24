package util

import (
	"net"
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println("ip:", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func ZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	r, _ := zlib.NewReader(b)

	if r != nil {
		var out bytes.Buffer
		defer r.Close()
		io.Copy(&out, r)
		return out.Bytes()
	}

	return nil
}

func ZlibCompress(rawData string) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(rawData))
	w.Close()

	return b.Bytes()
}

func GetBaseDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		fmt.Print("get basepath failed.")
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}