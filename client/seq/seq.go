package seq

import "sync"

// Seq 序列ID生成器（1~65535）
type Seq struct {
	id   uint16
	lock sync.Mutex
}

func (s *Seq) NextID() uint16 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.id += 1
	return s.id
}

var SeqIns = new(Seq)
