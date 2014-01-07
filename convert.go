package main

import (
	"flag"
	"os"
	"log"
	"bufio"
	"unicode/utf8"
	"unicode"
	"fmt"
)

var (
	input = flag.String(
		"input",
		"",
		"")
	output = flag.String(
		"output",
		"output.txt",
		"")
)

func main() {
	flag.Parse()

        // 打开将要搜索的文件
        file, err := os.Open(*input)
        if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

        // 逐行读入
        log.Printf("读入文本 %s", *input)
        scanner := bufio.NewScanner(file)
        lines := []string{}
        for scanner.Scan() {
		text := scanner.Text()
                if text != "" {
                        lines = append(lines, text)
                }
        }
        log.Print("文件行数", len(lines))

	oFile, oErr := os.Create(*output)
        if oErr != nil {
		log.Fatal(oErr)
	}
	defer oFile.Close()

	stopTokens := map[string]string{
		"。":".",
		"，":".",
		",":".",
		":":".",
		"：":".",
		"“":".",
		"”":".",
		"\"":".",
		"'":".",
		"《":".",
		"》":".",
		"!":".",
		"！":".",
		";":".",
		"；":".",
		"…":".",
		"—":".",
		" ":".",
		"?":".",
		"？":".",
		"<":".",
		">":".",
		"·":".",
		"~":".",
		"|":".",
		"、":".",
		"「":".",
		"」":".",
	}

	w := bufio.NewWriter(oFile)
	inStopToken := false
	for _, t := range lines {
		ws := splitTextToWords([]byte(t))
		for _, word := range ws {
			_, has := stopTokens[string(word)]
			if !has {
				fmt.Fprintln(w, string(word))
				inStopToken = false
			} else {
				if !inStopToken {
					inStopToken = true
					fmt.Fprintln(w, ".")
				}
			}
		}
	}
	w.Flush()
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
