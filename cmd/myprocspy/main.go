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

	c := http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: concurrency,
		},
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(time.Millisecond*500, func() {
		wg.Done()
	})

	go spies()
	mreqs(c, wg)

	wg.Wait()
	fmt.Println("myprocspy ends")
}

func mreqs(c http.Client, wg *sync.WaitGroup) {
	for i := 0; i < concurrency; i++ {
		time.Sleep(time.Millisecond * 100)
		go reqs(c, wg)
	}
}

func reqs(c http.Client, wg *sync.WaitGroup) {
	wg.Add(1)
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
