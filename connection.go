package main

import (
	"fmt"
	"log"
	"net"
)

type Connection struct {
	conn   *net.UDPConn
	remote *net.UDPAddr
	key    []byte
}

func newConnection(remotePort int, key []byte) (*Connection, error) {
	lstnAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", remotePort))
	if nil != err {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", lstnAddr)
	if nil != err {
		return nil, err
	}
	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", *remoteIP, *port))
	if nil != err {
		log.Fatalln("Unable to resolve remote:", err)
	}
	return &Connection{
		conn:   conn,
		remote: remoteAddr,
		key:    key,
	}, nil
}

func (c *Connection) Write(data []byte) (int, error) {
	cipher, err := Encrypt(data, c.key)
	if err != nil {
		return 0, err
	}
	return c.conn.WriteToUDP(cipher, c.remote)
}

func (c *Connection) Read(b []byte) (int, error) {
	buf := make([]byte, BUFFERSIZE)
	n, _, err := c.conn.ReadFromUDP(buf)
	if err != nil {
		return n, err
	}
	plaintext, err := Decrypt(buf[:n], c.key)
	if err != nil {
		return 0, err
	}
	copy(b, plaintext)
	return len(plaintext), nil
}

func (c *Connection) Close() {
	c.conn.Close()
}
