package whois

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

var servers = make(map[string]string)

type m struct {
	suffix string
	server string
}

func init() {
	datas, err := Get_servers()
	if err != nil {
		fmt.Println("init server fail:", err)
	}
	for _, v := range datas {
		servers[v.suffix] = v.server

	}

}

func Get_servers() (re []m, err error) {
	file, err := os.Open("./confserver.txt")
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
		var temp m
		temp.suffix = str[0]
		temp.server = str[1]
		re = append(re, temp)
	}
	return
}
