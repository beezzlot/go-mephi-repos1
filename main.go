package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	errors := 0
	addr := "http://srv.msk01.gigacorp.local/_stats"

	for {
		resp, err := http.Get(addr)
		if err != nil {
			errors++
			if errors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			errors++
			resp.Body.Close()
			if errors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		buf := make([]byte, 256)
		n, err := resp.Body.Read(buf)
		resp.Body.Close()

		if n == 0 || err != nil {
			errors++
			if errors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		text := strings.TrimSpace(string(buf[:n]))
		items := strings.Split(text, ",")

		if len(items) != 6 {
			errors++
			if errors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		var vals [6]float64
		ok := true
		for i, s := range items {
			v, e := strconv.ParseFloat(strings.TrimSpace(s), 64)
			if e != nil {
				ok = false
				break
			}
			vals[i] = v
		}

		if !ok {
			errors++
			if errors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		errors = 0

		if vals[0] > 30 {
			fmt.Printf("Load Average is too high: %v\n", vals[0])
		}

		if vals[1] > 0 {
			memPercent := (vals[2] / vals[1]) * 100
			if memPercent > 80 {
				fmt.Printf("Memory usage too high: %v%%\n", memPercent)
			}
		}

		if vals[3] > 0 {
			diskPercent := (vals[4] / vals[3]) * 100
			if diskPercent > 90 {
				freeMb := (vals[3] - vals[4]) / (1024 * 1024)
				fmt.Printf("Free disk space is too low: %v Mb left\n", freeMb)
			}
		}

		bandwidth := 6551603348.0
		if bandwidth > 0 {
			netPercent := (vals[5] / bandwidth) * 100
			if netPercent > 90 {
				freeMbit := (bandwidth - vals[5]) * 8 / (1024 * 1024)
				fmt.Printf("Network bandwidth usage high: %v Mbit/s available\n", freeMbit)
			}
		}

		time.Sleep(10 * time.Second)
	}
}