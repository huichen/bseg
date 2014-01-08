package main

import (
	"bufio"
	"flag"
	"github.com/huichen/bseg"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

var (
	input       = flag.String("input", "", "")
	output_dict = flag.String("output_dict", "dict.txt", "")
	cpuprofile  = flag.String("cpuprofile", "", "处理器profile文件")
)

func main() {
	flag.Parse()

	file, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

	tokens := []string{}
	segments := []uint8{}
	prevSeg := bseg.FIXSEG
	for _, t := range lines {
		if t == "." {
			if prevSeg == bseg.SEG {
				segments[len(segments)-1] = bseg.FIXSEG
			}
			prevSeg = bseg.FIXSEG
			continue
		}

		words := strings.Split(t, " ")
		for i := 0; i < len(words)-1; i++ {
			tokens = append(tokens, words[i])
			segments = append(segments, bseg.NOSEG)
		}
		tokens = append(tokens, words[len(words)-1])
		segments = append(segments, bseg.SEG)
		prevSeg = bseg.SEG
	}

	if segments[len(segments)-1] != bseg.NOSEG {
		segments = segments[0 : len(segments)-1]
	}

	if len(tokens) != len(segments)+1 {
		log.Fatal("not same")
	}

	seg := bseg.NewBSeg()

	// 打开处理器profile文件
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	seg.ProcessText(tokens, segments)

	seg.DumpDict(*output_dict)
}
