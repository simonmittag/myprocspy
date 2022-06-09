package main

import (
	"fmt"
	"github.com/simonmittag/procspy"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const concurrency = 32

var url = "jsonplaceholder.typicode.com/todos/1"
var tlsmode = false

func main() {
	if len(os.Args) > 0 && os.Args[1] == "-s" {
		url = "https://" + url
		tlsmode = true
	} else {
		url = "http://" + url
	}
	fmt.Printf("\nmyprocspy-starts-open-conns-%v-c%v-", url, concurrency)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(time.Millisecond*500, func() {
		wg.Done()
	})

	go spies()
	mreqs(wg)

	wg.Wait()
	fmt.Println("myprocspy-ends\n")
}

func initHTTPClient() http.Client {
	c := http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: 1,
			IdleConnTimeout: time.Duration(1 * time.Second),
		},
	}
	return c
}

func mreqs(wg *sync.WaitGroup) {
	for i := 0; i < concurrency; i++ {
		time.Sleep(time.Millisecond * 100)
		go reqs(wg)
	}
}

func reqs(wg *sync.WaitGroup) {
	wg.Add(1)
	c := initHTTPClient()
	for i := 0; i < 50; i++ {
		res, _ := c.Get(url)
		_, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		time.Sleep(time.Millisecond * 20)
	}
	wg.Done()
}

func spies() {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		p := fmt.Sprintf("%02d-", spy())
		fmt.Print(p)
	}
}

func spy() int {
	pid := os.Getpid()
	cs, _ := procspy.Connections(true)
	d := 0
	var wantport uint16 = 80
	if tlsmode {
		wantport = 443
	}
	for c := cs.Next(); c != nil; c = cs.Next() {
		if c.PID == uint(pid) && c.RemotePort == wantport {
			d++
		}
	}
	return d
}
