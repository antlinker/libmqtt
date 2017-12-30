/*
 * Copyright GoIIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
