package worker

import (
	"github.com/lobsterk/otus-abf/internal/models"
	"github.com/lobsterk/otus-abf/internal/repositories/mock"
	"github.com/lobsterk/otus-abf/internal/services"
	"github.com/lobsterk/otus-abf/internal/services/ip_guard"
	"testing"
)

func TestReloaderMasks_Exec(t *testing.T) {
	errChan := make(chan error, 10)

	masksRepository1 := &mock.MasksRepository{
		Rows: []models.Mask{
			{Id: 1, Mask: "123.23.44.55/8", ListId: 1},
		},
	}

	masksRepository2 := &mock.MasksRepository{
		Rows: []models.Mask{
			{Id: 2, Mask: "122.27.44.55/8", ListId: 1},
		},
	}

	lists := make([]services.IpGuard, 0, 2)
	lists = append(lists, ip_guard.NewMemoryIpGuard(1, masksRepository1))
	lists = append(lists, ip_guard.NewMemoryIpGuard(2, masksRepository2))

	w := NewReloaderMasks(lists, errChan)
	w.Exec()

	if len(errChan) != 0 {
		t.Error(<-errChan)
	}

	if ok, err := lists[0].Contains("123.23.44.55"); !ok {
		t.Error("!ok")
	} else if err != nil {
		t.Error(err.Error())
	}

	if ok, err := lists[1].Contains("122.27.44.55"); !ok {
		t.Error("!ok")
	} else if err != nil {
		t.Error(err.Error())
	}
}
