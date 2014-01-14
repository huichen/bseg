package bseg

import (
	"unicode"
	"unicode/utf8"
)

func GetSegmentsFromText(text []string) (tokens []string, segments []uint8) {
	prevSeg := FIXSEG
	inStopToken := true
	for _, line := range text {
		ws := splitTextToWords([]byte(line))
		for _, word := range ws {
			if !IsBoundary(string(word)) {
				segments = append(segments, SEG)
				tokens = append(tokens, string(word))
				inStopToken = false
			} else {
				if !inStopToken {
					inStopToken = true
					if prevSeg == SEG {
						segments[len(segments)-1] = FIXSEG
					}
					prevSeg = FIXSEG
				}
			}
		}
	}

	if segments[len(segments)-1] != NOSEG {
		segments = segments[0 : len(segments)-1]
	}

	return
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
