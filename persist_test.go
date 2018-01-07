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
	"os"
	"testing"
	"time"
)

var (
	// drop packets
	testPersistStrategy = &PersistStrategy{
		Interval:         time.Second,
		MaxCount:         1,
		DropOnExceed:     true,
		DuplicateReplace: false,
	}

	testPersistKeys    = []string{"foo", "foo", "bar"}
	testPersistPackets = []Packet{
		&SubscribePacket{
			Topics: []*Topic{
				{Name: "test"},
			},
		},
		&PublishPacket{},
		&ConnPacket{},
	}
)

func testPersist(p PersistMethod, t *testing.T) {
	if _, ok := p.Load(testPersistKeys[len(testPersistKeys)-1]); ok {
		t.Log("persist strategy failed")
		t.Fail()
	}

	if v, ok := p.Load(testPersistKeys[0]); !ok {
		t.Log("load persisted packet fail, packet =", v)
		t.Fail()
	} else {
		if v.Type() == CtrlSubscribe {
			if v.(*SubscribePacket).Topics[0].Name !=
				testPersistPackets[0].(*SubscribePacket).Topics[0].Name {
				t.Log("source topic name =", v.(*SubscribePacket).Topics[0].Name)
				t.Log("target topic name =", testPersistPackets[0].(*SubscribePacket).Topics[0].Name)
				t.Fail()
			}
		}
	}
}

func TestMemPersist(t *testing.T) {
	p := NewMemPersist(testPersistStrategy)

	for i, k := range testPersistKeys {
		if err := p.Store(k, testPersistPackets[i]); err != nil {
			if err != PacketDroppedByStrategy {
				t.Log(err)
				t.Fail()
			}
		}
	}

	if p.n != 1 {
		t.Log("persist strategy failed, count =", p.n)
		t.Fail()
	}

	testPersist(p, t)
}

func TestFilePersist(t *testing.T) {
	dirPath := "test-file-persist"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	p := NewFilePersist(dirPath, testPersistStrategy)

	for i, k := range testPersistKeys {
		if err := p.Store(k, testPersistPackets[i]); err != nil {
			if err != PacketDroppedByStrategy {
				t.Log(err)
				t.Fail()
			}
		}
	}

	if p.n == 1 {
		t.Fail()
	}

	<-time.After(750 * time.Millisecond)

	if p.n == 1 {
		t.Fail()
	}

	<-time.After(750 * time.Millisecond)

	if p.n != 1 {
		t.Log("persist strategy failed, count =", p.n)
		t.Fail()
	}

	testPersist(p, t)
	err = os.RemoveAll(dirPath)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
