package history

import (
	createhistory "exp1/internal/recorder/history/createHistory"
	"exp1/internal/types"
)

type History struct {
	create ICreateHistory
}

func NewHistory() *History {
	return &History{
		create: createhistory.NewCreateHistory(),
	}
}

func (h *History) Create(path string, data types.FileRecord) error {
	return h.create.Create(path, data)
}