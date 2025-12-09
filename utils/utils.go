package utils

import (
	"errors"
	"fmt"
	"os"
)

func CheckDirExist(path string) bool{
	info, err := os.Stat(path)
	if info.IsDir(){
		return true
	}
	if err != nil || os.IsNotExist(err){
		return false
	}
	return false
}

func CheckFileExist(path string) bool{
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func CreateFile(path string){
	f, err := os.Create(path)
	defer f.Close()
	if err != nil{
		Error(err.Error())
	}
}