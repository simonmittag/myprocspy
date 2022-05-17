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

const concurrency = 10

func main() {
	fmt.Println("myprocspy starts")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(time.Millisecond*500, func() {
		wg.Done()
	})

	go spies()
	mreqs(wg)

	wg.Wait()
	fmt.Println("myprocspy ends")
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
	for i := 0; i < 1000; i++ {
		res, _ := c.Get("http://jsonplaceholder.typicode.com/todos/1")
		_, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
	}
	wg.Done()
}

func spies() {
	for {
		time.Sleep(time.Millisecond * 100)
		p := fmt.Sprintf("%02d-", spy())
		fmt.Print(p)
	}
}

func spy() int {
	pid := os.Getpid()
	cs, _ := procspy.Connections(true)
	d := 0
	for c := cs.Next(); c != nil; c = cs.Next() {
		if c.PID == uint(pid) && c.RemotePort == 80 {
			d++
		}
	}
	return d
}
