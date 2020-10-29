package mock

import (
	"github.com/lobsterk/otus-abf/internal/constants"
	"github.com/lobsterk/otus-abf/internal/models"
)

type MasksRepository struct {
	Rows []models.Mask
}

func (r *MasksRepository) Get(listId constants.ListId) ([]models.Mask, error) {
	return r.Rows, nil
}

func (r *MasksRepository) Add(mask *models.Mask) error {
	mask.Id = len(r.Rows)
	r.Rows = append(r.Rows, *mask)
	return nil
}

func (r *MasksRepository) Drop(id int) error {
	for i := range r.Rows {
		if r.Rows[i].Id != id {
			continue
		}
		r.Rows = append(r.Rows[:i], r.Rows[i+1:]...)
		return nil
	}
	return nil
}
