package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	filePath := "E:/data/spider/movies.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for i := 0; i < 5; i++ {
		write.WriteString("test \r\n")
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}
