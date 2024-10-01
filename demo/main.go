package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var textToDisplay string
var adr string

func led() {
	matrix, width := textToLedMatrix(textToDisplay)
	destWidth := 47
	pos := 0
	for {
		ledMatrix := repeatMatrix(matrix, width, pos, destWidth)
		if pos < width {
			pos = pos + 1
		} else {
			pos = 0
		}

		buf, err := json.Marshal(ledMatrix)
		if err != nil {
			fmt.Errorf("%v", err)
		}

		http.Post(adr+"/event/led", "application/json", bytes.NewReader(buf))
		time.Sleep(time.Duration(50) * time.Millisecond)
	}
}

func streamtime() {
	for {
		currentTime := time.Now().UnixMilli()
		time.Sleep(time.Duration(16) * time.Millisecond)
		body := fmt.Sprintf("{\"time\":%d}", currentTime)
		http.Post(adr+"/event/streamtime", "application/json", strings.NewReader(body))
	}
}

var gliphs = makeGlyphs()

// and gaps = letters - 1 (1 column)
func textToLedMatrix(text string) ([]int, int) {
	runes := []rune(text)
	width := len(runes)*5 + (len(runes) - 1)
	buf := make([]int, width*6)
	widthPos := 0

	for ri, r := range runes {
		isLastRune := ri == (len(runes) - 1)
		g := gliphs[r]
		for rowNum := range 6 {
			destOffset := rowNum*width + widthPos
			srcOffset := rowNum * 5
			copy(buf[destOffset:destOffset+5], g[srcOffset:srcOffset+5])
			if !isLastRune {
				buf[destOffset+6] = 0
			}
		}
		if isLastRune {
			widthPos = widthPos + 5
		} else {
			widthPos = widthPos + 6
		}
	}

	return buf, width
}

// width * 6 rows
func repeatMatrix(srcMatrix []int, srcWidth, srcOffset, destWidth int) []int {
	buf := make([]int, destWidth*6)

	for rowNum := range 6 {
		isFirstInRow := true
		destPos := 0

		for destPos < destWidth {
			if isFirstInRow {
				copy(buf[rowNum*destWidth:], srcMatrix[rowNum*srcWidth+srcOffset:rowNum*srcWidth+srcWidth])
				destPos = srcWidth - srcOffset
				isFirstInRow = false
			} else {
				takeWidth := 0
				destFreeSpace := destWidth - destPos
				if destFreeSpace < srcWidth {
					takeWidth = destFreeSpace
				} else {
					takeWidth = srcWidth
				}

				copy(buf[rowNum*destWidth+destPos:], srcMatrix[rowNum*srcWidth:rowNum*srcWidth+takeWidth])
				destPos = destPos + takeWidth
			}
		}
	}

	return buf
}

func main() {
	flag.StringVar(&textToDisplay, "text", "", "text to display")
	flag.StringVar(&adr, "address", "", "SyncStreamer inbound address http[s]://[host]:[port]")
	flag.Parse()
	textToDisplay = strings.ToUpper(textToDisplay)

	if textToDisplay == "" || adr == "" {
		flag.PrintDefaults()
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go led()

	wg.Add(1)
	go streamtime()

	wg.Wait()
}
