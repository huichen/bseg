/*

bseg关键词提取服务器同时提供了两种模式：

	"/"	演示网页
	"/json"	JSON格式的RPC服务
		输入：
			POST或GET模式输入text参数
		输出JSON格式：
			{
				phrases:[
					{"text":"服务器", "count":"10"},
					{"text":"指令", "count":"8"},
					...
				]
			}


测试服务器见 http://bseg.weiboglass.com

*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/huichen/bseg"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
)

var (
	host         = flag.String("host", "", "HTTP服务器主机名")
	port         = flag.Int("port", 8080, "HTTP服务器端口")
	staticFolder = flag.String("static_folder", "static", "静态页面存放的目录")
)

type JsonResponse struct {
	Phrases []*Phrase `json:"phrases"`
}

type Phrase struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
}

func JsonRpcServer(w http.ResponseWriter, req *http.Request) {
	// 得到要分词的文本
	text := req.URL.Query().Get("text")
	if text == "" {
		text = req.PostFormValue("text")
	}

	lines := strings.Split(text, "\n")
	tokens, segments := bseg.GetSegmentsFromText(lines)

	seg := bseg.NewBSeg()
	seg.ProcessText(tokens, segments)
	ts := seg.GetDict()

	// 整理为输出格式
	ps := []*Phrase{}
	for _, token := range ts {
		ps = append(ps, &Phrase{Text: token.Name, Count: token.Count})
	}
	response, _ := json.Marshal(&JsonResponse{Phrases: ps})
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(response))
}

func main() {
	flag.Parse()

	// 将线程数设置为CPU数
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/json", JsonRpcServer)
	http.Handle("/", http.FileServer(http.Dir(*staticFolder)))
	log.Print("服务器启动")
	http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
}
