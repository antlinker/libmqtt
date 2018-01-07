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
	"fmt"
	"io"
	"math"
	"sync"
)

type BufferWriter interface {
	io.Writer
	io.ByteWriter
}

func boolToByte(flag bool) byte {
	if flag {
		return 1
	}
	return 0
}

func recvKey(packetID uint16) string {
	return fmt.Sprintf("%s%d", "R", packetID)
}

func sendKey(packetID uint16) string {
	return fmt.Sprintf("%s%d", "S", packetID)
}

type idGenerator struct {
	usedIds *sync.Map
}

func newIDGenerator() *idGenerator {
	return &idGenerator{
		usedIds: &sync.Map{},
	}
}

func (g *idGenerator) next(extra interface{}) uint16 {
	var i uint16
	for i = 1; i < math.MaxUint16; i++ {
		if _, ok := g.usedIds.Load(i); !ok {
			g.usedIds.Store(i, extra)
			return i
		}
	}
	return 1
}

func (g *idGenerator) free(id uint16) {
	g.usedIds.Delete(id)
}

func (g *idGenerator) getExtra(id uint16) (interface{}, bool) {
	return g.usedIds.Load(id)
}
