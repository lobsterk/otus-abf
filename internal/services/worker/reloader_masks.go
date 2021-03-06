package worker

import (
	"github.com/lobsterk/otus-abf/internal/services"
	"time"
)

func NewReloaderMasks(lists []services.IpGuard, errChan chan error) *ReloaderMasks {
	return &ReloaderMasks{
		lists:   lists,
		errChan: errChan,
	}
}

type ReloaderMasks struct {
	lists   []services.IpGuard
	errChan chan error
}

func (w *ReloaderMasks) Exec() {
	for i := range w.lists {
		if err := w.lists[i].Reload(); err != nil {
			w.errChan <- err
		}
	}
}

func (w *ReloaderMasks) GetInterval() time.Duration {
	return time.Second * 5
}
