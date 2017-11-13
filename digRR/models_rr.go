package digRR

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type RRs struct {
	Domain string
	A      []string
	NS     []string
	CNAME  []string
	MX     []string
	TXT    []string
}

var (
	Mapns_rr map[string][]string //问询对象的map
)

func init() {
	Mapns_rr = make(map[string][]string)
	err := initfirst()
	if err != nil {
		log.Println("init err:", err)
	}
	log.Println("init success!")
}

func GetDigRR(domain string) (re RRs, err error) {

	servers := Get_ask_server(domain)
	re.A = DigA(domain, servers)
	re.NS = DigNS(domain, servers)
	re.CNAME = DigCNAME("www."+domain, re.NS) //默认填充www.
	re.MX = DigMX(domain, re.NS)
	re.TXT = DigTXT(domain, re.NS)
	return
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
			fmt.Println(strs)
			continue
		}
		AddNS(strs[0], strs[1])
		//转为小写
		/*strs[0] = strings.ToLower(strs[0])
		strs[4] = strings.ToLower(strs[4])
		if strs[3] == "NS" {
			AddNS(strs[0], strs[4])
		}*/

	}

	return

}
func AddNS(qname string, name string) {

	if qname == name {
		fmt.Println("AddNS qname == name !!!! %v %v", qname, name)
		//panic(fmt.Sprintf("qname %v", qname))
		return
	}
	//Dsmap._lock.Lock()
	//defer Dsmap._lock.Unlock()
	for _, v := range Mapns_rr[name] {
		if v == qname {
			fmt.Println("CRI AddNS qname == name !!!! %v %v", qname, name)
			return
		}
	}

	ish := false
	for _, v := range Mapns_rr[qname] {
		if v == name {
			ish = true
			break
		}
	}
	if !ish {
		Mapns_rr[qname] = append(Mapns_rr[qname], name)
	}

}

func Get_ask_server(name string) (com_ns []string) {
	var key string
	h := strings.Split(name, ".")
	if len(h) == 3 {
		key = h[1] + "." + h[2]
	} else {
		t := len(h) - 1
		key = h[t]
	}

	//key := h[t] + "." //使用root.z时需要这句
	/*for _, v := range Mapns_rr[key] {
		//for _, v := range Mapns_rr[key] {
		com_ns = append(com_ns, v)
	}*/
	if v, ok := Mapns_rr[key]; ok {
		//存在
		return v
	}

	return

}

/*
 		 A               1 a host address
         NS              2 an authoritative name server
         MD              3 a mail destination (Obsolete - use MX)
         MF              4 a mail forwarder (Obsolete - use MX)
         CNAME           5 the canonical name for an alias
         SOA             6 marks the start of a zone of authority
         MB              7 a mailbox domain name (EXPERIMENTAL)
         MG              8 a mail group member (EXPERIMENTAL)
         MR              9 a mail rename domain name (EXPERIMENTAL)
         NULL            10 a null RRs (EXPERIMENTAL)
         WKS             11 a well known service description
         PTR             12 a domain name pointer
         HINFO           13 host information
         MINFO           14 mailbox or mail list information
         MX              15 mail exchange
         TXT             16 text strings
 		 AXFR            252 A request for a transfer of an entire zone
         MAILB           253 A request for mailbox-related records (MB, MG or MR)
         MAILA           254 A request for mail agent RRs (Obsolete - see MX)
         *               255 A request for all records
*/
func DigA(domain string, servers []string) (rr []string) {

	for _, v := range servers {
		remsg, err := send(v, domain, 1, 300*time.Millisecond) //seecd.net
		if err != nil {
			continue
		}

		for _, v := range remsg.Extra {
			if a, ok := v.(*A); ok {
				//fmt.Println("A:", a.A.String())
				rr = append(rr, a.A.String())

			}
		}
		break
	}
	return
}
func DigNS(domain string, servers []string) (rr []string) {

	for _, v := range servers {
		remsg, err := send(v, domain, 2, 300*time.Millisecond) //seecd.net
		if err != nil {
			continue
		}

		for _, v := range remsg.Ns {
			if ns, ok := v.(*NS); ok {
				rr = append(rr, ns.Ns)
			}
		}
		break
	}
	return
}
func DigCNAME(domain string, servers []string) (rr []string) {

	for _, v := range servers {
		fmt.Println(v, domain)
		remsg, err := send(v, domain, 5, 300*time.Millisecond) //seecd.net
		if err != nil {
			continue
		}
		for _, v := range remsg.Answer {
			if cn, ok := v.(*CNAME); ok {
				fmt.Println("CNAME:", cn.Target)
				rr = append(rr, cn.Target)
			}
		}

		break
	}
	return
}

func DigMX(domain string, servers []string) (rr []string) {

	for _, v := range servers {
		remsg, err := send(v, domain, 15, 300*time.Millisecond) //seecd.net
		if err != nil {
			continue
		}

		for _, v := range remsg.Answer {
			if m, ok := v.(*MX); ok {
				rr = append(rr, m.Mx)
			}
		}
		break
	}
	return
}

func DigTXT(domain string, servers []string) (rr []string) {

	for _, v := range servers {
		remsg, err := send(v, domain, 16, 300*time.Millisecond) //seecd.net
		if err != nil {
			//fmt.Println("debug:", err)
			continue
		}

		for _, v := range remsg.Answer {
			if t, ok := v.(*TXT); ok {

				rr = append(rr, t.TXT)
			}
		}

		break
	}
	return
}
