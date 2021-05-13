package app

import (
	"fmt"
	"log"

	"github.com/ddo/go-fast"
	"github.com/showwin/speedtest-go/speedtest"
)

type Speedy struct {
	log *log.Logger
}

// New constructs a User for api access.
func New(log *log.Logger) Speedy {
	return Speedy{
		log: log,
	}
}

type SpeednetKbps struct {
	Name  string  `json:"name"`
	DMbps float64 `json:"Download-mbps"`
	UMbps float64 `json:"Upload-mbps"`
}

func (s Speedy) Fasttest() ([]*SpeednetKbps, error) {
	fmt.Println("Started")

	fastCom := fast.New()

	// init
	err := fastCom.Init()
	if err != nil {
		panic(err)
	}

	// get urls
	urls, err := fastCom.GetUrls()
	if err != nil {
		panic(err)
	}

	// measure
	KbpsChan := make(chan float64)

	var speed []*SpeednetKbps
	var v []float64
	go func() {
		for Kbps := range KbpsChan {
			fmt.Printf("%.2f Kbps %.2f Mbps\n", Kbps, Kbps/1000)
			v = append(v, Kbps)
		}
		b := v[len(v)-1] / 1000
		f := SpeednetKbps{
			Name:  "fast",
			DMbps: b,
		}
		speed = append(speed, &f)
		fmt.Println("done")
	}()
	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		panic(err)
	}

	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)
		s := SpeednetKbps{
			Name:  "speedtest",
			DMbps: s.DLSpeed,
			UMbps: s.ULSpeed,
		}
		speed = append(speed, &s)
		fmt.Println("Finished")
	}
	return speed, nil
}
 