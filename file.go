package main

import (
	"fmt"
	"os"
)

func initFile() *os.File {
	filePath := "result.md"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("文件打开失败", err)
		return nil
	}

	return file

	//Flush将缓存的文件真正写入到文件中
	//write.Flush()
}
