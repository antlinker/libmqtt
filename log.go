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

package libmqtt

import (
	"log"
	"os"
)

type LogLevel int

const (
	Silent LogLevel = iota
	Verbose
	Debug
	Info
	Warning
	Error
)

type logger struct {
	verbose *log.Logger
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

const (
	logFlag = log.Ltime | log.Ldate
)

func newStdLogger() *log.Logger {
	l := &log.Logger{}
	l.SetFlags(logFlag)
	l.SetOutput(os.Stderr)
	return l
}

func newLogger(l LogLevel) *logger {
	lo := &logger{}

	if l <= Error {
		lo.error = newStdLogger()
		lo.error.SetPrefix("[LIBMQTT] Error: ")
	}

	if l <= Warning {
		lo.warning = newStdLogger()
		lo.warning.SetPrefix("[LIBMQTT] Warning: ")
	}

	if l <= Info {
		lo.info = newStdLogger()
		lo.info.SetPrefix("[LIBMQTT] Info: ")
	}

	if l <= Debug {
		lo.debug = newStdLogger()
		lo.debug.SetPrefix("[LIBMQTT] Debug: ")
	}

	if l <= Verbose {
		lo.verbose = newStdLogger()
		lo.verbose.SetPrefix("[LIBMQTT] Verbose: ")
	}

	if l <= Silent {
		lo = nil
	}

	return lo
}

func (l *logger) v(data ...interface{}) {
	if l == nil || l.verbose == nil {
		return
	}
	l.verbose.Println(data...)
}

func (l *logger) d(data ...interface{}) {
	if l == nil || l.debug == nil {
		return
	}
	l.debug.Println(data...)
}

func (l *logger) i(data ...interface{}) {
	if l == nil || l.info == nil {
		return
	}
	l.info.Println(data...)
}

func (l *logger) w(data ...interface{}) {
	if l == nil || l.warning == nil {
		return
	}
	l.warning.Println(data...)
}

func (l *logger) e(data ...interface{}) {
	if l == nil || l.error == nil {
		return
	}
	l.error.Println(data...)
}
