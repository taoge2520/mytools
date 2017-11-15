package file

import (
	"bufio"
	//	"database/sql"
	"fmt"
	"io"
	"os"
	"regexp"

	//	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//打开并逐行读取文件内容,根据tab或者空格分割该行内容
func GetContent(fileName string) (err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		inputstring, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		req := regexp.MustCompile(`[^\s]+`)
		str := req.FindAllString(inputstring, -1)
		fmt.Println(str)
		//根据需要自行处理
	}
	return
}

//打开并逐行读取文件内容，原样获取文件内容
func GetLine(fileName string) (err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file) //带缓冲区的读写
	for {
		line, err := reader.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		fmt.Print(line)

	}
	return
}

//往文件末尾追加内容
func AppendToFile(fileName string, content string) (err error) {
	// 以只写的模式，打开文件
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := file.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = file.WriteAt([]byte(content), n)
	}
	defer file.Close()
	return
}
