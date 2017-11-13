// -------------------------
//
// Copyright 2015, undiabler
//
// git: github.com/undiabler/golang-whois
//
// http://undiabler.com
//
// Released under the Apache License, Version 2.0
//
//--------------------------

package whois

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"time"
)

//Simple connection to whois servers with default timeout 5 sec
func GetWhois(domain string) (string, error) {
	return GetWhoisTimeout(domain, time.Second*5)

}

//Connection to whois servers with various time.Duration
func GetWhoisTimeout(domain string, timeout time.Duration) (result string, err error) {

	var (
		parts []string

		buffer     []byte
		connection net.Conn
	)

	parts = strings.Split(domain, ".")
	if len(parts) < 2 {
		err = fmt.Errorf("Domain(%s) name is wrong!", domain)
		return
	}
	//last part of domain is zome
	zone := parts[len(parts)-1]

	server, ok := servers[zone]
	//fmt.Println("this domain server is:", server)
	if !ok {
		err = fmt.Errorf("No such server for zone %s. Domain %s.", zone, domain)
		return
	}
	//connection, err = Dial("tcp", "10.10.100.188", server)
	connection, err = net.DialTimeout("tcp", net.JoinHostPort(server, "43"), timeout)

	if err != nil {
		//return net.Conn error
		return
	}

	defer connection.Close()

	connection.Write([]byte(domain + "\r\n"))

	buffer, err = ioutil.ReadAll(connection)

	if err != nil {
		return
	}

	result = string(buffer[:])

	return
}
func GetWhois2(domain string, remote string, timeout time.Duration) (result string, err error) {

	var (
		parts []string

		buffer     []byte
		connection net.Conn
	)

	parts = strings.Split(domain, ".")
	if len(parts) < 2 {
		err = fmt.Errorf("Domain(%s) name is wrong!", domain)
		return
	}
	//connection, err = Dial("tcp", "10.10.100.18", remote)
	connection, err = net.DialTimeout("tcp", net.JoinHostPort(remote, "43"), timeout)

	if err != nil {
		//return net.Conn error
		return
	}

	defer connection.Close()

	connection.Write([]byte(domain + "\r\n"))

	buffer, err = ioutil.ReadAll(connection)

	if err != nil {
		return
	}

	result = string(buffer[:])

	return
}
func Dial(network string, local string, remote string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   500 * time.Millisecond, //超时设置
		KeepAlive: 1 * time.Second,
	}
	local = local + ":0" //端口0,系统会自动分配本机端口
	switch network {
	case "udp":
		addr, err := net.ResolveUDPAddr(network, local)
		if err != nil {
			return nil, err
		}
		dialer.LocalAddr = addr
	case "tcp":
		addr, err := net.ResolveTCPAddr(network, local)
		if err != nil {
			return nil, err
		}
		dialer.LocalAddr = addr
	}
	return dialer.Dial(network, remote+":43")
}
