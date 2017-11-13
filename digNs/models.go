package digNs

import (
	"bufio"
	"io"
	"log"

	"os"
	"regexp"

	"strings"
	"time"
	"unicode"
)

var (
	root_ns   [13]string
	Mapns     map[string][]string //问询对象的map
	Check_ch  chan string
	Answer_ch chan Datas
)

type Datas struct {
	Domain string
	Ns     []string
}

const (
	NORMAL_TIME  = 500
	TIMEOUT_TIME = 1000
)

func init() {
	root_ns = [13]string{ //根
		"a.root-servers.net.",
		"b.root-servers.net.",
		"c.root-servers.net.",
		"d.root-servers.net.",
		"e.root-servers.net.",
		"f.root-servers.net.",
		"g.root-servers.net.",
		"h.root-servers.net.",
		"i.root-servers.net.",
		"j.root-servers.net.",
		"k.root-servers.net.",
		"l.root-servers.net.",
		"m.root-servers.net."}
	Mapns = make(map[string][]string)
	Check_ch = make(chan string, 100)
	Answer_ch = make(chan Datas, 100)
	err := initfirst()
	if err != nil {
		log.Println(err)

	}
	log.Println("len of mapns :", len(Mapns))

}
func initfirst() (err error) {
	inputFile, err := os.Open("./root1.txt")
	if err != nil {
		return
	}
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			return
		}

		req := regexp.MustCompile(`[^\s]+`)
		strs := req.FindAllString(inputString, -1)
		if len(strs) != 2 {
			log.Println(strs)
			continue
		}
		AddNS(strs[0], strs[1])

	}

	return

}
func AddNS(qname string, name string) {

	if qname == name {
		log.Println("AddNS qname == name !!!! %v %v", qname, name)
		//panic(log.Sprintf("qname %v", qname))
		return
	}
	//Dsmap._lock.Lock()
	//defer Dsmap._lock.Unlock()
	for _, v := range Mapns[name] {
		if v == qname {
			log.Println("CRI AddNS qname == name !!!! %v %v", qname, name)
			return
		}
	}

	ish := false
	for _, v := range Mapns[qname] {
		if v == name {
			ish = true
			break
		}
	}
	if !ish {
		Mapns[qname] = append(Mapns[qname], name)
	}

}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
func Checkns(value string) (ns []string) {

	Ns_name := get_ask_server(value)
	is_ch := IsChineseChar(value)
	if len(Ns_name) == 0 && !is_ch { //xn--4gq62f39mfm3c.xn--fiqs8s

		strs := strings.Split(value, ".")
		com_ns, err := dig_root(strs[len(strs)-1], NORMAL_TIME*time.Millisecond)
		if err != nil {
			log.Println(err)
		}
		ns = dig(value, com_ns, NORMAL_TIME*time.Millisecond)
		return
	}
	//和原数据类型有差异
	if is_ch && len(Ns_name) == 0 { //时时彩.中国

		changed, err := ToASCII(value)
		if err != nil {
			log.Println(err)
		}
		strs := strings.Split(changed, ".")
		com_ns, err := dig_root(strs[len(strs)-1], NORMAL_TIME*time.Millisecond)
		if err != nil {
			log.Println(err)
		}
		ns = dig(changed, com_ns, NORMAL_TIME*time.Millisecond)
		return
	}
	if is_ch { //时时彩.com

		changed, err := ToASCII(value)
		if err != nil {
			log.Println(err)
		}
		ns = dig(changed, Ns_name, NORMAL_TIME*time.Millisecond)

		return
	} //baidu.com
	if len(Ns_name) == 0 {
		strs := strings.Split(value, ".")
		com_ns, err := dig_root(strs[len(strs)-1], NORMAL_TIME*time.Millisecond)
		if err != nil {
			log.Println(err)
		}
		ns = dig(value, com_ns, NORMAL_TIME*time.Millisecond)
		return
	}

	ns = dig(value, Ns_name, NORMAL_TIME*time.Millisecond)

	return

}

func dig_root(domain string, sc time.Duration) (gtld_ns []string, err error) {

	for _, v := range root_ns {

		remsg, err := send(v, domain, 2, sc)
		if err != nil {
			continue
		}

		for _, v := range remsg.Ns {
			if ns, ok := v.(*NS); ok {
				gtld_ns = append(gtld_ns, ns.Ns)
			}
		}
	}

	return gtld_ns, nil
}
func dig(d string, ns []string, sc time.Duration) (answer_ns []string) {

	for _, v := range ns {

		remsg, err := send(v, d, 2, sc)

		if err != nil {
			log.Println(err)
			continue
		}

		if remsg.MsgHdr.Authoritative {
			ns = append(ns, "SOA")
			return
		}

		for _, v := range remsg.Ns {
			if ns, ok := v.(*NS); ok {
				answer_ns = append(answer_ns, ns.Ns)
			}
		}
		//		for _, v := range remsg.Answer {
		//			if ns, ok := v.(*NS); ok {
		//				answer_ns = append(answer_ns, ns.Ns)
		//			}
		//		}

		return

	}
	return
}
func get_ask_server(name string) (com_ns []string) {
	var key string
	h := strings.Split(name, ".")
	if len(h) == 3 {
		key = h[1] + "." + h[2]
	} else {
		t := len(h) - 1
		key = h[t]
	}

	if v, ok := Mapns[key]; ok {
		//存在
		return v
	}
	//fmt.Println("can't find ", name, "of ns server")
	return

}
