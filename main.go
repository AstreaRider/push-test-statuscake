package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

type Pinger struct {
	Timeout time.Duration
	Count   int
}

func getAvgRtt(host string, pc, timeout Pinger) (time.Duration, error) {
	pinger, err := ping.NewPinger(host)

	if err != nil {
		return 0, err
	}

	pinger.Timeout = timeout.Timeout
	pinger.Count = pc.Count
	err = pinger.Run()
	if err != nil {
		return 0, err
	}
	avgRtt := pinger.Statistics().AvgRtt

	return avgRtt, nil
}

func getFlags(hosts, pk, testID string) ([]string, string, string, error) {
	flag.StringVar(&hosts, "host", "", "List of Server Hosts")
	flag.StringVar(&pk, "pk", "", "Statuscake PK")
	flag.StringVar(&testID, "test-id", "", "Statuscake Test ID")
	flag.Parse()

	if hosts == "" || pk == "" || testID == "" {
		return nil, "", "", errors.New("Error making request, usage example push-test-statuscake --host=<host-1,host2,..> --pk=<pk-statuscake> --test-id=<test-id-statuscake>\n")
	}

	hostList := strings.Split(hosts, ",")
	return hostList, pk, testID, nil
}

func getAvgRtts(rtts []time.Duration) time.Duration {
	var avg time.Duration
	for _, v := range rtts {
		avg += v
	}
	return avg / time.Duration(len(rtts))
}

func timeToIntConverter(v time.Duration) int {
	rounded := v / time.Millisecond * time.Millisecond
	return int(rounded.Milliseconds())
}

func createPushUrl(pk, testId string, timeLoad int) string {
	return fmt.Sprintf("https://push.statuscake.com/?PK=%v&TestID=%v&time=%v", pk, testId, timeLoad)
}

func createPushTest(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	var hosts string
	var primaryKey string
	var testID string
	timeout := 5 * time.Second

	listOfHost, pk, tID, errFlags := getFlags(hosts, primaryKey, testID)
	if errFlags != nil {
		panic(errFlags)
	}

	avgRtts := make([]time.Duration, len(listOfHost))
	for i, host := range listOfHost {
		avgRtt, err := getAvgRtt(host, Pinger{Count: 4}, Pinger{Timeout: timeout})
		if err != nil {
			fmt.Printf("Error pinging host: %v\n", err)
			return
		}
		avgRtts[i] = avgRtt
	}
	avg := getAvgRtts(avgRtts)
	timeLoad := timeToIntConverter(avg)
	pushUrl := createPushUrl(pk, tID, timeLoad)
	pushTest, errCreatePushTest := createPushTest(pushUrl)
	if errCreatePushTest != nil {
		fmt.Printf("Error making request: %v", errCreatePushTest)
		return
	}
	fmt.Printf("Push test result : %v\n", pushTest)
}
