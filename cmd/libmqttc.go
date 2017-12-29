package main

import (
	"bufio"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"

	mq "github.com/goiiot/libmqtt"
)

const (
	lineStart = "> "
)

var (
	client mq.Client
)

func main() {
	flag.Parse()
	osCh := make(chan os.Signal, 2)
	signal.Notify(osCh, os.Kill, os.Interrupt)
	wg := &sync.WaitGroup{}

	// handle system signals
	go func() {
		for range osCh {
			os.Exit(0)
		}
	}()
	// handle user input
	wg.Add(1)
	go func() {
		print(lineStart)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() != "" {
				args := strings.Split(scanner.Text(), " ")
				execCmd(args)
			}
			print(lineStart)
		}
		wg.Done()
	}()
	wg.Wait()
}

func execCmd(args []string) {
	cmd := strings.ToLower(args[0])
	args = args[1:]
	ok := false
	switch cmd {
	case "c", "conn":
		ok = execConn(args)
	case "p", "pub":
		ok = execPub(args)
	case "s", "sub":
		ok = execSub(args)
	case "u", "unsub":
		ok = execUnSub(args)
	case "q", "exit":
		ok = execDisConn(args)
	case "?", "h", "help":
		ok = usage()
	}

	if !ok {
		usage()
	}
}

func usage() bool {
	print("Usage\n\n")
	print("  ")
	connUsage()
	print("  ")
	pubUsage()
	print("  ")
	subUsage()
	print("  ")
	unSubUsage()
	println(`  q, exit [force] - disconnect and exit`)
	println(`  h, help - print this help message`)
	return true
}
