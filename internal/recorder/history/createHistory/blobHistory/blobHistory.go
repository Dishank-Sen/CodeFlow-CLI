package blobhistory

import (
	"exp1/internal/recorder/history/createHistory/blobHistory/delta"
	"exp1/internal/recorder/history/createHistory/blobHistory/snapshot"
	"exp1/internal/types"
	"fmt"
	"log"
)

type BlobHistory struct{
	snapshot *snapshot.Snapshot
	delta *delta.Delta
}

func NewBlobHistory() *BlobHistory{
	return &BlobHistory{
		snapshot: snapshot.NewSnapshot(),
		delta: delta.NewDelta(),
	}
}

func (b *BlobHistory) Create(path string, data any) error {
	if b.isSnapshot(data){
		return b.snapshot.Create(path, data)
	}else{
		return b.delta.Create(path, data)
	}
}

func (b *BlobHistory) isSnapshot(data any) bool{
	dataType := b.getFileType(data)
	if dataType == ""{
		log.Fatal("cannot destructure into blob type")
	}
	if dataType == "snapshot"{
		return true
	}
	return false
}

func (b *BlobHistory) getFileType(data any) string{
	fmt.Println(data)
	if fileData, ok := data.(types.FileRecord); ok {
		return fileData.Type
	}
	return ""
}