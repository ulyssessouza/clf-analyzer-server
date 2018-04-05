package data

import (
	"sync"
	"time"
)

type SynchBroadcastArray struct {
	sync.Mutex
	channels []chan struct{}
}

func (s *SynchBroadcastArray) Register(w chan struct{}) {
	s.Lock()
	defer s.Unlock()
	s.channels = append(s.channels, w)
}

func (s *SynchBroadcastArray) Deregister(w chan struct{}) {
	s.Lock()
	defer s.Unlock()
	// Delete by not including the channel in the new slice
	var newSlice []chan struct{}
	for _, v := range s.channels {
		if v == w {
			continue
		} else {
			newSlice = append(newSlice, v)
		}
	}
	s.channels = newSlice
}

func (s *SynchBroadcastArray) Broadcast(t *time.Ticker) {
	s.Lock()
	defer s.Unlock()

	<-t.C
	for _, c := range s.channels {
		c <- struct{}{}
	}
}