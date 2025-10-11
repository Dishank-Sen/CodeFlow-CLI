package recorder

import (
    "io/ioutil"
    "path/filepath"
)

type Reader struct {
    root string
}

func NewReader(root string) *Reader {
    return &Reader{root: root}
}

func (r *Reader) Load(path string, version string) ([]byte, error) {
    blobPath := filepath.Join(r.root, "blobs", version[:2], version)
    return ioutil.ReadFile(blobPath)
}
