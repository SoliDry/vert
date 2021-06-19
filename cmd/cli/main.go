package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/SoliDry/vert/internal/db"
	vnet "github.com/SoliDry/vert/internal/net"
	"github.com/arthurkushman/pgo"
	"github.com/magiconair/properties"
)

const configFile = "/etc/vert/vert.conf.d/vert.cnf"

const defaultPort = "9909"

func main() {
	var l net.Listener
	var conn net.Conn
	var err error

	db.CreateSysFilesAndTables()

	if pgo.FileExists(configFile) == false {
		log.Fatal(fmt.Sprintf("error: there is no main config file - %s", configFile))
	}

	// read props from config
	p := properties.MustLoadFile(configFile, properties.UTF8)
	bindAddr, ok := p.Get("bind.address")
	if !ok {
		bindAddr = "0.0.0.0"
	}

	l, err = net.Listen("tcp", bindAddr+":"+defaultPort)
	if err != nil {
		log.Println(fmt.Errorf("error connecting TCP: %w", err))
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on " + bindAddr + ":" + defaultPort)

	// connections accepting worker
	for {
		conn, err = l.Accept()
		if err != nil {
			log.Println(fmt.Errorf("error accepting: %w", err))
		}

		go handleReq(conn)
	}
}

func handleReq(conn net.Conn) {
	buf := make([]byte, 4096)

	_, err := conn.Read(buf)
	if err != nil {
		log.Println(fmt.Errorf("could not read from tcp: %w", err))
	}

	vnetConn := vnet.NewConnection(&conn, true, "utf8_unicode_ci")
	dataOut, err := vnet.GetHandshakeSeq(*vnetConn, buf)
	if err != nil {
		log.Println(fmt.Errorf("handshake error: %w", err))
	}

	_, err = conn.Write(dataOut)
	if err != nil {
		log.Println(fmt.Errorf("could not write to tcp: %w", err))
	}

	conn.Close()
}
