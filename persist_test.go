package libmqtt

import "testing"

func TestMemPersist(t *testing.T) {
	p := NewMemPersist(&PersistStrategy{
		MaxCount:         1,
		DropOnExceed:     true,
		DuplicateReplace: false,
	})

	p.Store("foo", nil)
	p.Store("foo", &SubscribePacket{})
	p.Store("bar", nil)

	if p.count != 1 {
		t.Log("count =", p.count)
		t.Fail()
	}

	if v, ok := p.Load("foo"); !ok || v != nil {
		t.Log("pkt =", v)
		t.Fail()
	}
}

func TestFilePersist(t *testing.T) {

}
