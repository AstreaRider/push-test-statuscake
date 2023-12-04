package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

func pingServer(host string, pingCount int) (*ping.Statistics, error) {
	// Membuat instance pinger
	pinger, err := ping.NewPinger(host)

	if err != nil {
		return nil, err
	}

	// Konfigurasi jumlah paket ping yang akan dikirim
	pinger.Count = pingCount

	// Mulai ping
	err = pinger.Run()
	if err != nil {
		return nil, err
	}

	// Mendapatkan hasil ping setelah selesai
	return pinger.Statistics(), nil
}

func getServers(hosts string) []string {
	flag.StringVar(&hosts, "host", "", "List of Server Hosts")
	flag.Parse()

	hostList := strings.Split(hosts, " ")
	return hostList
}

func printResult(host string) {
	stats, err := pingServer(host, 4)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("===============================================")
		fmt.Printf("Host: %v\n", host)
		fmt.Printf("Ping Result - Sent: %v, Received: %v, Lost: %v\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("Min/Max/Avg latency: %v / %v / %v\n", stats.MinRtt, stats.MaxRtt, stats.AvgRtt)
	}
}

func getAvgRtt(host string) time.Duration {
	stats, _ := pingServer(host, 4)
	return stats.AvgRtt
}

func getAverage(values []time.Duration) time.Duration {
	var avg time.Duration
	for _, v := range values {
		avg += v
	}
	return avg / time.Duration(len(values))
}

func main() {
	var hosts string

	hostList := getServers(hosts)

	avgRtt := make([]time.Duration, len(hostList))

	for i, host := range hostList {
		printResult(host)
		avgRtt[i] = getAvgRtt(host)
	}
	avg := getAverage(avgRtt)
	fmt.Println("===============================================")
	fmt.Printf("Rata-Rata: %v\n", avg)
}
