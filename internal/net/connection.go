package net

import "net"

type Connection struct {
	ID          int64
	tcp         *net.Conn
	isConnected bool
	/* ID of character set in HEX
	| latin1_swedish_ci   |  8  |  0x08  |
	| utf8_general_ci     |  33 |  0x21  |
	| binary     		  |  35 |  0x3f  |
	*/
	CharSet       string
	MaxPacketSize uint32
	ConnectWithDb string
}

var connIdInc int64

func NewConnection(tcp *net.Conn, isConnected bool, charSet string) *Connection {
	connIdInc++
	return &Connection{ID: connIdInc, tcp: tcp, isConnected: isConnected, CharSet: charSet}
}

//type Connector interface {
//	Connect() (tcp *net.TCPConn, err error)
//}
//
//func (conn *Connection) Connect() (tcp *net.TCPConn, err error) {
//
//	return tcp, err
//}
