package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/syncstreamer/server/processor"
	"github.com/syncstreamer/server/timeframe/eventframe"
	"github.com/syncstreamer/server/timestamp"
	"github.com/syncstreamer/server/types"
)

var inAddr string
var outAddr string
var serveStatic bool
var timeframeDuration int
var timeframeHistoryItems int

func readParams() {
	const (
		defaultTimeframeDuration     = 10_000
		defaultTimeframeHistoryItems = 5
	)
	inAddrEnv, _ := os.LookupEnv("SYNCSTREAMER_IN_ADDRESS")
	flag.StringVar(&inAddr, "in_addr", inAddrEnv, "Inbound address \"[host]:[port]\"")
	outAddrEnv, _ := os.LookupEnv("SYNCSTREAMER_OUT_ADDRESS")
	flag.StringVar(&outAddr, "out_addr", outAddrEnv, "Outbound address \"[host]:[port]\"")
	serveStaticEnv, _ := os.LookupEnv("SYNCSTREAMER_SERVE_STATIC")
	serveStaticEnvBool := serveStaticEnv == "true"
	flag.BoolVar(&serveStatic, "serve_static", serveStaticEnvBool,
		"set to true if the server should serve client static too, default: false")
	timeframeDurationEnv, _ := os.LookupEnv("SYNCSTREAME_TIMEFRAME_DURATION")
	timeframeDurationEnvInt, _ := strconv.ParseInt(timeframeDurationEnv, 10, 64)
	flag.IntVar(&timeframeDuration, "timeframe_duration", int(timeframeDurationEnvInt),
		fmt.Sprintf("timeframe duration in ms, default: %d", defaultTimeframeDuration))
	timeframeHistoryItemsEnv, _ := os.LookupEnv("SYNCSTREAM_TIMEFRAME_HISTORY_ITEMS")
	timeframeHistoryItemsEnvInt, _ := strconv.ParseInt(timeframeHistoryItemsEnv, 10, 64)
	flag.IntVar(&timeframeHistoryItems, "timeframe_history_items", int(timeframeHistoryItemsEnvInt),
		fmt.Sprintf("timeframe history items number, default: %d", defaultTimeframeHistoryItems))
	flag.Parse()

	if inAddr == "" || outAddr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if timeframeDuration == 0 {
		timeframeDuration = defaultTimeframeDuration
	}

	if timeframeHistoryItems == 0 {
		timeframeDuration = defaultTimeframeHistoryItems
	}
}

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
		Addr:    inAddr,
		Handler: muxIn,
	}

	serverIn.ListenAndServe()
}

func outServer(proc *processor.Processor) {
	muxIn := http.NewServeMux()

	if serveStatic {
		muxIn.Handle("/", http.FileServer(http.Dir("./static")))
	}

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
		Addr:    outAddr,
		Handler: muxIn,
	}

	serverIn.ListenAndServe()
}

func main() {
	readParams()

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
