package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
}

func main() {
	flag.Parse()

	flagArgs := flag.Args()
	if len(flagArgs) != 2 {
		log.Fatal("Missing required params, usage: \"go-telnet host port\" or \"go-telnet --timeout=10s host port\"")
	}
	address := net.JoinHostPort(flagArgs[0], flagArgs[1])

	c := NewTelnetClient(address, timeout, ioutil.NopCloser(os.Stdin), os.Stdout)
	log.Printf("...Connected to %s\n", address)
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := c.Close()
		log.Fatal(err)
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			if err := c.Send(); err != nil {
				log.Fatal(err)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			if err := c.Receive(); err != nil {
				log.Fatal(err)
				return
			}
		}
	}()

	wg.Wait()
}
