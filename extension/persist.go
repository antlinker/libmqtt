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

package extension

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/go-redis/redis"
	lib "github.com/goiiot/libmqtt"
)

const defaultRedisKey = "libmqtt"

// NewRedisPersist will create a new RedisPersist for session persist
// with provided redis connection and key mainKey,
// if passed empty mainKey here, the default mainKey "libmqtt" will be used
// if no redis client (nil) provided, will return nil
func NewRedisPersist(conn *redis.Client, mainKey string) *RedisPersist {
	if conn == nil {
		return nil
	}

	if mainKey == "" {
		mainKey = defaultRedisKey
	}

	buf := &bytes.Buffer{}
	return &RedisPersist{
		conn:    conn,
		mainKey: mainKey,
		buf:     buf,
		w:       bufio.NewWriter(buf),
	}
}

// RedisPersist defines the persist method with redis
type RedisPersist struct {
	conn    *redis.Client
	w       *bufio.Writer
	buf     *bytes.Buffer
	mainKey string
}

// Name of RedisPersist is "RedisPersist"
func (r *RedisPersist) Name() string {
	if r == nil {
		return "<nil>"
	}

	return "RedisPersist"
}

// Store a packet with key
func (r *RedisPersist) Store(key string, p lib.Packet) error {
	if r == nil || r.conn == nil {
		return nil
	}

	if p.Bytes(r.w) != nil {
		if ok, err := r.conn.HSet(r.mainKey, key, r.buf.String()).Result(); !ok {
			return err
		}
	}
	r.buf.Reset()

	return nil
}

// Load a packet from stored data according to the key
func (r *RedisPersist) Load(key string) (lib.Packet, bool) {
	if r == nil || r.conn == nil {
		return nil, false
	}

	if rs, err := r.conn.HGet(r.mainKey, key).Result(); err != nil {
		if pkt, err := lib.DecodeOnePacket(strings.NewReader(rs)); err != nil {
			// delete wrong packet
			r.Delete(key)
		} else {
			return pkt, true
		}
	}

	return nil, false
}

// Range over data stored, return false to break the range
func (r *RedisPersist) Range(f func(string, lib.Packet) bool) {
	if r == nil || r.conn == nil {
		return
	}

	if set, err := r.conn.HGetAll(r.mainKey).Result(); err == nil {
		for k, v := range set {
			if pkt, err := lib.DecodeOnePacket(strings.NewReader(v)); err != nil {
				r.Delete(k)
				continue
			} else {
				if !f(k, pkt) {
					break
				}
			}
		}
	}
}

// Delete a persisted packet with key
func (r *RedisPersist) Delete(key string) error {
	if r == nil || r.conn == nil {
		return nil
	}
	_, err := r.conn.HDel(r.mainKey, key).Result()
	return err
}

// Destroy stored data
func (r *RedisPersist) Destroy() error {
	if r == nil || r.conn == nil {
		return nil
	}
	_, err := r.conn.Del(r.mainKey).Result()
	return err
}
