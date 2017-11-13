//需要在linux下运行
package models

import (
	//	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	TIMEOUT  = int64(1000)
	COUNT    = 4
	SIZE_BUF = 32
	TO       = false
)

type Pingdata struct {
	Msg      string
	SendN    int
	RecvN    int
	LostN    int
	ShortT   int
	LongT    int
	LostRate int
	SumT     int
}

func Ping(host string) (data Pingdata, err error) {

	var count int
	var size int
	var timeout int64
	var neverstop bool
	count = COUNT
	size = SIZE_BUF
	timeout = TIMEOUT
	neverstop = TO

	cname, err := net.LookupCNAME(host)
	if err != nil {
		return
	}
	starttime := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", host, time.Duration(timeout*1000*1000))
	if err != nil {
		return
	}
	ip := conn.RemoteAddr()
	data.Msg += "Ping " + cname + " [" + ip.String() + "] 具有 32 字节的数据:\n"
	//fmt.Println("正在 Ping " + cname + " [" + ip.String() + "] 具有 32 字节的数据:")

	var seq int16 = 1
	id0, id1 := genidentifier(host)
	const ECHO_REQUEST_HEAD_LEN = 8

	data.SendN = 0
	data.RecvN = 0
	data.LostN = 0
	data.ShortT = -1
	data.LongT = -1
	data.SumT = 0

	for count > 0 || neverstop {
		data.SendN++
		var msg []byte = make([]byte, size+ECHO_REQUEST_HEAD_LEN)
		msg[0] = 8                        // echo
		msg[1] = 0                        // code 0
		msg[2] = 0                        // checksum
		msg[3] = 0                        // checksum
		msg[4], msg[5] = id0, id1         //identifier[0] identifier[1]
		msg[6], msg[7] = gensequence(seq) //sequence[0], sequence[1]

		length := size + ECHO_REQUEST_HEAD_LEN

		check := checkSum(msg[0:length])
		msg[2] = byte(check >> 8)
		msg[3] = byte(check & 255)

		conn, err = net.DialTimeout("ip:icmp", host, time.Duration(timeout*1000*1000))

		checkError(err)

		starttime = time.Now()
		conn.SetDeadline(starttime.Add(time.Duration(timeout * 1000 * 1000)))
		_, err = conn.Write(msg[0:length])
		if err != nil {
			return
		}

		const ECHO_REPLY_HEAD_LEN = 20

		var receive []byte = make([]byte, ECHO_REPLY_HEAD_LEN+length)
		_, err = conn.Read(receive)

		if err != nil {
			return
		}
		var endduration int = int(int64(time.Since(starttime)) / (1000 * 1000))

		data.SumT += endduration

		time.Sleep(1000 * 1000 * 1000)

		if err != nil || receive[ECHO_REPLY_HEAD_LEN+4] != msg[4] || receive[ECHO_REPLY_HEAD_LEN+5] != msg[5] || receive[ECHO_REPLY_HEAD_LEN+6] != msg[6] || receive[ECHO_REPLY_HEAD_LEN+7] != msg[7] || endduration >= int(timeout) || receive[ECHO_REPLY_HEAD_LEN] == 11 {
			data.LostN++
			//fmt.Println("对 " + cname + "[" + ip.String() + "]" + " 的请求超时。")
			data.Msg += "对 " + cname + "[" + ip.String() + "]" + " 的请求超时\n"
		} else {
			if data.ShortT == -1 {
				data.ShortT = endduration
			} else if data.ShortT > endduration {
				data.ShortT = endduration
			}
			if data.LongT == -1 {
				data.LongT = endduration
			} else if data.LongT < endduration {
				data.LongT = endduration
			}
			data.RecvN++
			ttl := int(receive[8])
			//			fmt.Println(ttl)
			//fmt.Println("来自 " + cname + "[" + ip.String() + "]" + " 的回复: 字节=32 时间=" + strconv.Itoa(endduration) + "ms TTL=" + strconv.Itoa(ttl))
			data.Msg += "来自 " + cname + "[" + ip.String() + "]" + " 的回复: 字节=32 时间=" + strconv.Itoa(endduration) + "ms TTL=" + strconv.Itoa(ttl) + "\n"
		}

		seq++
		count--
	}

	//remsg = stat(ip.String(), sendN, lostN, recvN, shortT, longT, sumT)
	return

}

func checkSum(msg []byte) uint16 {
	sum := 0

	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256 // notice here, why *256?
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func gensequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genidentifier(host string) (byte, byte) {
	return host[0], host[1]
}

func stat(ip string, sendN int, lostN int, recvN int, shortT int, longT int, sumT int) (str string) {

	str += ip + " 的 Ping 统计信息:\n"
	str += "数据包:已发送=" + strconv.Itoa(sendN) + ",已接收 = " + strconv.Itoa(recvN) + ",丢失 =" + strconv.Itoa(lostN) + "(" + strconv.Itoa(int(lostN*100/sendN)) + "丢失)\n"
	str += "往返行程的估计时间(以毫秒为单位):\n"
	if recvN != 0 {
		str += `最短 = ` + strconv.Itoa(shortT) + `ms，最长 = ` + strconv.Itoa(longT) + `ms，平均 = ` + strconv.Itoa(sumT/sendN) + `ms`
	}
	return
}
