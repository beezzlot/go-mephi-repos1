package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	errorCount := 0
	url := "http://srv.msk01.gigacorp.local/_stats"

	for {
		resp, err := http.Get(url)
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			errorCount++
			resp.Body.Close()
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		data := make([]byte, 128)
		n, _ := resp.Body.Read(data)
		resp.Body.Close()

		line := strings.TrimSpace(string(data[:n]))
		parts := strings.Split(line, ",")

		if len(parts) != 6 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		var nums [6]float64
		parseOk := true
		for i, p := range parts {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil {
				parseOk = false
				break
			}
			nums[i] = v
		}

		if !parseOk {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		errorCount = 0

		if nums[0] > 30 {
			fmt.Printf("Load Average is too high: %v\n", nums[0])
		}

		if nums[1] > 0 {
			memPercent := (nums[2] / nums[1]) * 100
			if memPercent > 80 {
				fmt.Printf("Memory usage too high: %v%%\n", memPercent)
			}
		}

		if nums[3] > 0 {
			diskPercent := (nums[4] / nums[3]) * 100
			if diskPercent > 90 {
				freeMB := (nums[3] - nums[4]) / (1024 * 1024)
				fmt.Printf("Free disk space is too low: %v Mb left\n", freeMB)
			}
		}

		if len(parts) > 5 {
			netUsage := nums[5]
			netCapacity := 1073741824.0
			if netCapacity > 0 {
				netPercent := (netUsage / netCapacity) * 100
				if netPercent > 90 {
					freeMbit := (netCapacity - netUsage) * 8 / (1024 * 1024)
					fmt.Printf("Network bandwidth usage high: %v Mbit/s available\n", freeMbit)
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}
