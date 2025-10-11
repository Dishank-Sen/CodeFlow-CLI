package history

import "exp1/internal/types"

// history --> ICreateHistory --> Create functionality

type ICreateHistory interface {
	Create(path string, data types.FileRecord) error
}