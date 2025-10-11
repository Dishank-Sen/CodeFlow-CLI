package createhistory

import (
	blobhistory "exp1/internal/recorder/history/createHistory/blobHistory"
	nonblobhistory "exp1/internal/recorder/history/createHistory/nonBlobHistory"
	"exp1/internal/types"
)

type CreateHistory struct{
	blobHistory *blobhistory.BlobHistory
	nonBlobHistory *nonblobhistory.NonBlobHistory
}

func NewCreateHistory() *CreateHistory{
	return &CreateHistory{
		blobHistory: blobhistory.NewBlobHistory(),
		nonBlobHistory: nonblobhistory.NewNonBlobHistory(),
	}
}

func (c *CreateHistory) Create(path string, data types.FileRecord) error {
	if isBlob(data) {
		return c.blobHistory.Create(path, data)
	}
	return c.nonBlobHistory.Create(path, data)
}

func isBlob(data types.FileRecord) bool {
	return data.IsBlobType
}