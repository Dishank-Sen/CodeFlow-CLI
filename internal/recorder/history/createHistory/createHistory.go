package createhistory

import (
	blobhistory "exp1/internal/recorder/history/createHistory/blobHistory"
	"exp1/internal/types"
)

type CreateHistory struct{
	blobHistory *blobhistory.BlobHistory
}

func NewCreateHistory() *CreateHistory{
	return &CreateHistory{
		blobHistory: blobhistory.NewBlobHistory(),
	}
}

func (c *CreateHistory) Create(path string, data types.FileRecord) error {
	return c.blobHistory.Create(path, data)
}