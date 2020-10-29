package repositories

import (
	"github.com/lobsterk/otus-abf/internal/constants"
	"github.com/lobsterk/otus-abf/internal/models"
)

type Masks interface {
	Get(listId constants.ListId) ([]models.Mask, error)
	Add(mask *models.Mask) error
	Drop(id int) error
}
