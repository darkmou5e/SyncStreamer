package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"github.com/darkmou5e/syncstreamer/processor"
	"github.com/darkmou5e/syncstreamer/timeframe/eventframe"
	"github.com/darkmou5e/syncstreamer/timestamp"
	"github.com/darkmou5e/syncstreamer/types"
)

func inServer(proc *processor.Processor) {
	muxIn := http.NewServeMux()

	muxIn.HandleFunc("/event/{channel}", func(resp http.ResponseWriter, req *http.Request) {
		// if req.Method == http.MethodPost {
		// }

		buf := make([]byte, req.ContentLength)
		req.Body.Read(buf)
		req.Body.Close()

		contentType := req.Header.Get("Content-Type")
		channelId := req.PathValue("channel")

		proc.AddEvent(&eventframe.Event{
			ChannelId: types.Id(channelId),
			EventType: types.ChannelType(contentType),
			EventData: buf,
		})

		resp.WriteHeader(http.StatusOK)
	})

	serverIn := http.Server{
		Addr:    ":4242",
		Handler: muxIn,
	}

	serverIn.ListenAndServe()
}

func outServer(proc *processor.Processor) {
	muxIn := http.NewServeMux()

	muxIn.Handle("/", http.FileServer(http.Dir("./client")))

	muxIn.HandleFunc("/frame", func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			fmt.Println("!!!")
			resp.WriteHeader(http.StatusNotFound)
			return
		} else {
			fmt.Println("OK!")
		}

		type ResponseItem struct {
			StartAt int
			EndAt   int
			Id      string
		}

		items := proc.GetTimeframes()
		its := make([]*ResponseItem, len(items))

		for i, x := range items {
			its[i] = &ResponseItem{
				StartAt: int(x.StartAt),
				EndAt:   int(x.EndAt),
				Id:      strconv.Itoa(int(x.StartAt)),
			}
		}

		b, err := json.Marshal(its)

		if err != nil {
			panic("out server panic")
		}

		resp.Header().Add("Content-Type", "application/json")
		resp.Write(b)
	})

	muxIn.HandleFunc("/frame/{frameId}", func(resp http.ResponseWriter, req *http.Request) {
		// if req.Method == http.MethodGet {
		// 	fmt.Println("!!!")
		// }

		frameId, err := strconv.ParseInt(req.PathValue("frameId"), 0, 0)
		fmt.Println(frameId)
		if err != nil {
			panic("frameId isn't int")
		}

		frames := proc.GetTimeframes()

		i := slices.IndexFunc(frames, func(frm *processor.TimeframeItem) bool {
			return frm.StartAt == timestamp.Timestamp(frameId)
		})

		if i < 0 {
			resp.WriteHeader(http.StatusNotFound)
		} else {
			resp.Header().Add("Content-Type", "application/octet-stream")
			resp.WriteHeader(http.StatusOK)
			dt := frames[i].Data
			resp.Write(dt)
		}
	})

	serverIn := http.Server{
		Addr:    ":8080",
		Handler: muxIn,
	}

	serverIn.ListenAndServe()
}

func main() {
	wg := sync.WaitGroup{}
	proc := processor.StartNewProcessor()

	wg.Add(1)
	go func() {
		inServer(proc)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		outServer(proc)
		wg.Done()
	}()

	wg.Wait()
}
