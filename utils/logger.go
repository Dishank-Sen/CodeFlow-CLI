package utils

import (
	"fmt"
	"os"
)

func Error(out string){
	fmt.Println("Error[CLI]: ", out)
	os.Exit(0)
}