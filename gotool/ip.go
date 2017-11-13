package gotool

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

//ip格式转uint32
func IPv4toUint32(ip string) (uint32, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return 0, errors.New("ParseIP error")
	}
	i = i.To4()
	return binary.BigEndian.Uint32(i), nil
}

////Uint32转ip格式 old
// func Uint32toIPv4(ipint uint32) string {
// 	a := ipint >> 24
// 	b := (ipint - (a << 24)) >> 16
// 	c := (ipint - (a << 24) - (b << 16)) >> 8
// 	d := ipint - (a << 24) - (b << 16) - (c << 8)
// 	return fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
// }

//Uint32转ip格式
func Uint32toIPv4(ipint uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipint)
	return ip.String()
}

//Uint32转成net.IP
func Uint32toIP(ipint uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipint)
	return ip
}

//获取本机网卡IP
func GetLocalIP() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, v := range addrs {
		if ipnet, ok := v.(*net.IPNet); ok {
			ip := ipnet.IP
			if ip.To4() != nil && !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips, nil
}

//将CIDR转成数字,如  1.0.0.0/24 转成 16777216 16777471
func CIDRToUint32(cidr string) (start uint32, end uint32, err error) {
	i, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return
	}
	mask, bit := n.Mask.Size()
	start, err = IPv4toUint32(i.String())
	if err != nil {
		return
	}
	network, err := IPv4toUint32(n.IP.String())
	if err != nil {
		return
	}
	end = 1<<(uint32(bit-mask)) + network - 1
	return
}

//将CIDR转成起始IP-结束IP,如  192.168.0.0/24 转成 192.168.0.0 192.168.0.255
func CIDRToIPRange(cidr string) (startip string, endip string, err error) {
	start, end, err := CIDRToUint32(cidr)
	if err != nil {
		return
	}
	startip = Uint32toIPv4(start)
	endip = Uint32toIPv4(end)
	return
}

//将起始IP-结束IP转成CIDR,如 192.168.0.0 192.168.0.255  转成 192.168.0.0/24
func IPRangeToCIDR(startip string, endip string) (cidr string, err error) {
	start, err := IPv4toUint32(startip)
	if err != nil {
		return
	}
	end, err := IPv4toUint32(endip)
	if err != nil {
		return
	}
	//取主机位长度
	bit := len(fmt.Sprintf("%b", start^end))
	//起始地址（网络号）
	ipint := (start >> uint32(bit)) << uint32(bit)
	cidr = fmt.Sprintf("%v/%v", Uint32toIPv4(ipint), 32-bit)
	return
}

// 用startip,endip计算子网掩码长度，如192.168.0.0 192.168.0.255 返回 24
func IPNetMaskBit(startip string, endip string) (bit int, err error) {
	start, err := IPv4toUint32(startip)
	if err != nil {
		return
	}
	end, err := IPv4toUint32(endip)
	if err != nil {
		return
	}
	bit = 32 - len(fmt.Sprintf("%b", start^end))
	return
}

//将 1.1.1.0/24  转成1.1.1.0 255.255.255.0
func CIDRToIPMask(cidr string) (string, string, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", err
	}
	bytesToIP := func(b []byte) string {
		i := binary.BigEndian.Uint32(b)
		return Uint32toIPv4(i)
	}
	ip := ipnet.IP.String()
	mask := bytesToIP([]byte(ipnet.Mask))
	return ip, mask, nil
}

//将 1.1.1.0 255.255.255.0   转成 1.1.1.0/24
func IPMaskToCIDR(ip string, mask string) string {
	IP := net.ParseIP(ip)
	if IP == nil {
		return ""
	}
	return fmt.Sprintf("%s/%v", IP, MaskLength(mask))
}

//计算反掩码  将255.255.255.0 转成0.0.0.255
func InverseMask(mask string) string {
	i, err := IPv4toUint32(mask)
	if err != nil {
		return ""
	}
	return Uint32toIPv4(^i)
}

//计算掩码长度  255.255.255.0 得出 24
func MaskLength(mask string) int {
	ip := net.ParseIP(mask).To4()
	if ip == nil {
		return 0
	}
	i := []byte(ip)
	ipmask := net.IPv4Mask(i[0], i[1], i[2], i[3])
	ones, _ := ipmask.Size()
	return ones
}

// Well-known IPv4 Private addresses
var (
	PrivateIPNet = []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
)

//是否为私网ip
func IsPrivateIP(ip string) bool {
	for _, ipnet := range PrivateIPNet {
		_, n, _ := net.ParseCIDR(ipnet)
		if n.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}
