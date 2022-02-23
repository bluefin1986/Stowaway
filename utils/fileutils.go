package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadFile(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("\r\n[*] file %s doesn't exist!", filepath)
		filepath = "./cfg.properties"
		file, err = os.Open(filepath)
		if err != nil {
			panic(err)
		}
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(content)
}
