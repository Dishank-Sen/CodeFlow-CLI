package blob

import (
	"exp1/internal/recorder/blob/hash"
	"log"
)

type Blob struct{
	saveBlob ISaveBlob
	retrieveBlob IRetrieveBlob
	hash *hash.Hash
}

func NewBlob() *Blob{
	return &Blob{
		saveBlob: NewSaveBlob(),
		retrieveBlob: NewRetrieveBlob(),
		hash: hash.NewHash(),
	}
}

func (b *Blob) CreateBlobFromPath(path string) string{
	hash, err := b.hash.HashFile(path)
	if err != nil{
		log.Fatal(err)
	}
	blobPath := b.saveBlob.saveBlobFromPath(hash, path)
	return blobPath
}

func (b *Blob) CreateBlobFromContent(content string) string{
	hash := b.hash.HashContent([]byte(content))
	blobPath := b.saveBlob.saveBlobFromContent(hash, content)
	return blobPath
}