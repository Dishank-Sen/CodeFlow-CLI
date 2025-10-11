package blob

import (
	"log"
	"os"
	"path/filepath"
)

type SaveBlob struct{

}

func NewSaveBlob() *SaveBlob{
	return &SaveBlob{}
}

// save blob and return its path
func (s *SaveBlob) saveBlobFromPath(hash string, path string) string{
	blobPath, err := s.getBlobPath(&hash)
	if err != nil{
		log.Fatal(err)
	}

	content, err := s.getFileContent(path)
	if err != nil{
		log.Fatal(err)
	}

	err = s.writeContentToBlob(blobPath, content)
	if err != nil{
		log.Fatal(err)
	}
	return blobPath
}

func (s *SaveBlob) saveBlobFromContent(hash string, content string) string{
	blobPath, err := s.getBlobPath(&hash)
	if err != nil{
		log.Fatal(err)
	}

	err = s.writeContentToBlob(blobPath, []byte(content))
	if err != nil{
		log.Fatal(err)
	}
	return blobPath
}

func (s *SaveBlob) getBlobPath(hash *string) (string, error) {
    // Dereference hash pointer
    h := *hash

    // Construct blob path like .rec/blob/ab/abcdef123...
    blobPath := filepath.Join(".rec", "blob", h[:2], h)

    // Ensure parent directory exists
    if err := os.MkdirAll(filepath.Dir(blobPath), 0755); err != nil {
        return "", err
    }

    return blobPath, nil
}

func (s *SaveBlob) getFileContent(path string) ([]byte, error){
	return os.ReadFile(path)
}

func (s *SaveBlob) writeContentToBlob(blobPath string, content []byte) error{
	if err := os.WriteFile(blobPath, content, 0644); err != nil {
        return err
    }
	return nil
}