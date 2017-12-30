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
	"os"
	"strconv"
	"strings"

	mq "github.com/goiiot/libmqtt"
)

func execConn(args []string) bool {
	if len(args) < 1 {
		return false
	}

	if len(strings.Split(args[0], ":")) != 2 {
		return false
	}

	options := make([]mq.Option, 0)
	if len(args) == 2 {
		opts := strings.Split(args[1], ",")
		var willQos mq.QosLevel
		var sslSkipVerify, ssl, will, willRetain bool
		var sslCert, sslKey, sslCA, sslServer, willTopic, willMsg string
		for _, v := range opts {
			kv := strings.Split(v, "=")
			if len(kv) != 2 {
				println(v, "option should be key=value")
				return true
			}

			switch kv[0] {
			case "clean":
				options = append(options, mq.WithCleanSession(kv[0] == "y"))
			case "ssl":
				ssl = kv[1] == "y"
			case "ssl_cert":
				sslCert = kv[1]
			case "ssl_key":
				sslKey = kv[1]
			case "ssl_ca":
				sslCA = kv[1]
			case "ssl_server":
				sslServer = kv[1]
			case "ssl_skip_verify":
				sslSkipVerify = kv[1] == "y"
			case "will":
				will = kv[1] == "y"
			case "will_topic":
				willTopic = kv[1]
			case "will_qos":
				qos, err := strconv.Atoi(kv[1])
				if err != nil || qos > 2 {
					invalidQos()
					return true
				}
				willQos = mq.QosLevel(qos)
			case "will_msg":
				willMsg = kv[1]
			case "will_retain":
				willRetain = kv[1] == "y"
			}
		}
		if ssl {
			options = append(options, mq.WithTLS(sslCert, sslKey, sslCA, sslServer, sslSkipVerify))
		}
		if will {
			options = append(options, mq.WithWill(willTopic, willQos, willRetain, []byte(willMsg)))
		}
	}
	options = append(options, mq.WithServer(args[0]))
	newClient(options)
	return true
}

func execDisConn(args []string) bool {
	if client != nil {
		client.Destroy(!(len(args) > 0 && args[1] != "force"))
	}

	os.Exit(0)

	return true
}

func newClient(options []mq.Option) {
	if client != nil {
		client.Destroy(true)
	}

	var err error
	client, err = mq.NewClient(options...)
	if err != nil {
		println(err.Error())
		return
	}
	client.Connect(connHandler)

	client.HandlePub(pubHandler)
	client.HandleSub(subHandler)
	client.HandleUnSub(unSubHandler)
	client.HandleNet(netHandler)
}

func connUsage() {
	println(`c, conn [server:port] [OPTIONS] - connect to server
    [OPTIONS] is a key=value config list separated by comma, valid keys:
      clean={y|n},ssl={y|n},ssl_cert={cert_path},
      ssl_key={key_path},ssl_ca={ca},
      ssl_server={server_name},ssl_skip_verify={y|n},
      will={y|n},will_topic={topic_name},will_qos={qos},
      will_msg={msg},will_retain={y|n}`)
}
