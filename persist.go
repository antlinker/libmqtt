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
	"sync"
	"sync/atomic"
	"time"
)

// PersistStrategy defines the details to be complied in persist methods
type PersistStrategy struct {
	// Interval applied to file/database persist
	// if this field is set to 0, means do persist per action
	// default value is 1s
	Interval time.Duration

	// MaxCount applied to all persist method
	// if this field set to 0, means no persist limit
	// for memory persist, means max in memory count
	// for file/database persist, means max entry in file/memory
	// default value is 0
	MaxCount uint32

	// DropOnExceed defines how to tackle with packets incoming
	// when max count is reached, default value is false
	DropOnExceed bool

	// DuplicateReplace defines whether duplicated key should
	// override previous one, default value is true
	DuplicateReplace bool
}

// DefaultPersistStrategy will create a default PersistStrategy
// Interval = 1s, MaxCount = 0, DropOnExceed = false, DuplicateReplace = true
func DefaultPersistStrategy() *PersistStrategy {
	return &PersistStrategy{
		Interval:         time.Second,
		MaxCount:         0,
		DropOnExceed:     false,
		DuplicateReplace: true,
	}
}

// PersistMethod defines the behavior of persist methods
type PersistMethod interface {
	// Name of what persist strategy used
	Name() string

	// Store a packet with key
	Store(key string, p Packet) error

	// Load a packet from stored data according to the key
	Load(key string) (Packet, bool)

	// Range over data stored, return false to break the range
	Range(func(key string, p Packet) bool)

	// Delete
	Delete(key string)

	// Destroy stored data
	Destroy() error
}

// NewMemPersist create a in memory persist method with provided strategy
// if no strategy provided (nil), then the default strategy will be used
func NewMemPersist(strategy *PersistStrategy) *MemPersist {
	p := &MemPersist{
		data:  &sync.Map{},
		count: 0,
	}
	if strategy == nil {
		p.strategy = DefaultPersistStrategy()
	} else {
		p.strategy = strategy
	}
	return p
}

// MemPersist is the in memory persist method
type MemPersist struct {
	data     *sync.Map
	count    uint32
	strategy *PersistStrategy
}

// Name of this persist method
func (m *MemPersist) Name() string {
	if m == nil {
		return "<nil>"
	}
	return "MemPersist"
}

// Store a key packet pair, in memory persist always return nil (no error)
func (m *MemPersist) Store(key string, p Packet) error {
	if m == nil {
		return nil
	}

	if m.strategy.MaxCount > 0 &&
		atomic.LoadUint32(&m.count) >= m.strategy.MaxCount &&
		m.strategy.DropOnExceed {
		return nil
	}

	if _, loaded := m.data.LoadOrStore(key, p); !loaded {
		atomic.AddUint32(&m.count, 1)
	} else if m.strategy.DuplicateReplace {
		m.data.Store(key, p)
	}
	return nil
}

// Load a packet with key, return nil, false when no packet found
func (m *MemPersist) Load(key string) (Packet, bool) {
	if m == nil {
		return nil, false
	}

	if p, ok := m.data.Load(key); ok {
		if p != nil {
			return p.(Packet), true
		}
	} else {
		return nil, false
	}

	return nil, true
}

// Range over all packet persisted
func (m *MemPersist) Range(f func(key string, p Packet) bool) {
	if m == nil || f == nil {
		return
	}

	m.data.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(Packet))
	})
}

// Delete a persisted packet with key
func (m *MemPersist) Delete(key string) {
	if m == nil {
		return
	}

	m.data.Delete(key)
}

// Destroy persist storage
func (m *MemPersist) Destroy() error {
	if m == nil {
		return nil
	}

	m.data = &sync.Map{}
	return nil
}

// NewFilePersist will create a file persist method with provided
// dirPath and strategy, if no strategy provided (nil), then the
// default strategy will be used
func NewFilePersist(dirPath string, strategy *PersistStrategy) *FilePersist {
	p := &FilePersist{
		dirPath: dirPath,
	}

	if strategy != nil {
		p.strategy = strategy
	}
	return p
}

// FilePersist is the file persist method
type FilePersist struct {
	dirPath  string
	strategy *PersistStrategy
}

// Name of this persist method
func (m *FilePersist) Name() string {
	if m == nil {
		return "<nil>"
	}
	return "FilePersist"
}

// Store a key packet pair, error happens when file access failed
func (m *FilePersist) Store(key string, p Packet) error {
	if m == nil {
		return nil
	}

	if m.strategy.MaxCount > 0 &&
		m.count() >= m.strategy.MaxCount &&
		m.strategy.DropOnExceed {
		return nil
	}

	if !m.exists(key) || m.strategy.DuplicateReplace {
		return m.store(key, p)
	}

	return nil
}

// Load a packet with key, return nil, false when no packet found
func (m *FilePersist) Load(key string) (Packet, bool) {
	if m == nil {
		return nil, false
	}

	return nil, false
}

// Range over all packet persisted
func (m *FilePersist) Range(f func(key string, p Packet) bool) {
	if m == nil || f == nil {
		return
	}
}

// Delete a persisted packet with key
func (m *FilePersist) Delete(key string) {
	if m == nil {
		return
	}
}

// Destroy persist storage
func (m *FilePersist) Destroy() error {
	if m == nil {
		return nil
	}
	return nil
}

func (m *FilePersist) store(key string, p Packet) error {
	return nil
}

func (m *FilePersist) exists(key string) bool {
	return false
}

func (m *FilePersist) count() uint32 {
	return 0
}
