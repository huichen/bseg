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
	"unicode"
	"unicode/utf8"
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

	tokens := []string{}
	segments := []uint8{}

	prevSeg := bseg.FIXSEG
	inStopToken := true
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		ws := splitTextToWords([]byte(line))
		for _, word := range ws {
			if !IsBoundary(string(word)) {
				segments = append(segments, bseg.SEG)
				tokens = append(tokens, string(word))
				inStopToken = false
			} else {
				if !inStopToken {
					inStopToken = true
					if prevSeg == bseg.SEG {
						segments[len(segments)-1] = bseg.FIXSEG
					}
					prevSeg = bseg.FIXSEG
				}
			}
		}
	}

	if segments[len(segments)-1] != bseg.NOSEG {
		segments = segments[0 : len(segments)-1]
	}

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

// 将文本划分成字元
func splitTextToWords(text []byte) [][]byte {
	output := make([][]byte, len(text))
	current := 0
	currentWord := 0
	inAlphanumeric := true
	alphanumericStart := 0
	for current < len(text) {
		r, size := utf8.DecodeRune(text[current:])
		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			// 当前是拉丁字母或数字（非中日韩文字）
			if !inAlphanumeric {
				alphanumericStart = current
				inAlphanumeric = true
			}
		} else {
			if inAlphanumeric {
				inAlphanumeric = false
				if current != 0 {
					output[currentWord] = toLower(text[alphanumericStart:current])
					currentWord++
				}
			}
			output[currentWord] = text[current : current+size]
			currentWord++
		}
		current += size
	}

	// 处理最后一个字元是英文的情况
	if inAlphanumeric {
		if current != 0 {
			output[currentWord] = toLower(text[alphanumericStart:current])
			currentWord++
		}
	}

	return output[:currentWord]
}

// 将英文词转化为小写
func toLower(text []byte) []byte {
	output := make([]byte, len(text))
	for i, t := range text {
		if t >= 'A' && t <= 'Z' {
			output[i] = t - 'A' + 'a'
		} else {
			output[i] = t
		}
	}
	return output
}

func IsBoundary(word string) bool {
	stopTokens := map[string]string{
		"。":  ".",
		"，":  ".",
		",":  ".",
		":":  ".",
		"：":  ".",
		"“":  ".",
		"”":  ".",
		"\"": ".",
		"'":  ".",
		"《":  ".",
		"》":  ".",
		"!":  ".",
		"！":  ".",
		";":  ".",
		"；":  ".",
		"…":  ".",
		"—":  ".",
		" ":  ".",
		"?":  ".",
		"？":  ".",
		"<":  ".",
		">":  ".",
		"·":  ".",
		"~":  ".",
		"|":  ".",
		"、":  ".",
		"「":  ".",
		"」":  ".",
		"（":  ".",
		"）":  ".",
		"(":  ".",
		")":  ".",
		"．":  ".",
		"‘":  ".",
		"’":  ".",
		"　":  ".",
	}
	_, found := stopTokens[word]
	return found
}
