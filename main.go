package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	fails := 0
	host := "http://srv.msk01.gigacorp.local/_stats"

	for {
		resp, err := http.Get(host)
		if err != nil {
			fails++
			if fails >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			fails++
			resp.Body.Close()
			if fails >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		data := make([]byte, 512)
		n, err := resp.Body.Read(data)
		resp.Body.Close()

		if n == 0 {
			fails++
			if fails >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		line := strings.TrimSpace(string(data[:n]))
		parts := strings.Split(line, ",")

		if len(parts) != 7 {
			fails++
			if fails >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		var nums [7]float64
		good := true
		for i, p := range parts {
			v, e := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if e != nil {
				good = false
				break
			}
			nums[i] = v
		}

		if !good {
			fails++
			if fails >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(10 * time.Second)
			continue
		}

		fails = 0

		if nums[0] > 30 {
			fmt.Printf("Load Average is too high: %v\n", nums[0])
		}

		if nums[1] > 0 {
			memPct := (nums[2] / nums[1]) * 100
			if memPct > 80 {
				fmt.Printf("Memory usage too high: %v%%\n", memPct)
			}
		}

		if nums[3] > 0 {
			diskPct := (nums[4] / nums[3]) * 100
			if diskPct > 90 {
				freeMb := (nums[3] - nums[4]) / (1024 * 1024)
				fmt.Printf("Free disk space is too low: %v Mb left\n", freeMb)
			}
		}

		if nums[5] > 0 {
			netPct := (nums[6] / nums[5]) * 100
			if netPct > 90 {
				freeMbit := (nums[5] - nums[6]) * 8 / (1024 * 1024)
				fmt.Printf("Network bandwidth usage high: %v Mbit/s available\n", freeMbit)
			}
		}

		time.Sleep(10 * time.Second)
	}
}