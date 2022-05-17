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

func main() {
	fmt.Println("myprocspy starts")

	c := http.Client{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(time.Millisecond*500, func() {
		wg.Done()
	})

	go reqs(c, wg)
	go spies(wg)

	wg.Wait()
	fmt.Println("myprocspy ends")
}

func reqs(c http.Client, wg *sync.WaitGroup) {
	wg.Add(1)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		req(c)
	}
	wg.Done()
}

func spies(wg *sync.WaitGroup) {
	wg.Add(1)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Printf("%d", spy())
	}
	wg.Done()
}

func req(c http.Client) {
	res, _ := c.Get("http://jsonplaceholder.typicode.com/todos/1")
	_, _ = ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Print(".")
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
