// 中国联通组网终端与智能网关自动连接接口要求v1.0(20180612）
// iptables -t nat -A OUTPUT -p tcp -m tcp --dport 32768 -j DNAT --to-destination {elink_server_ip}
package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net"
	"runtime"
	"sync"
)

var (
	GlobalFlag uint32 = 0x3F721FB5
	GlobalMac         = "12:34:56:78:9a:bc"
	SessionMap        = sync.Map{}
	HeaderLen         = 8
	MaxDataLen uint32 = 1024 * 1024 * 4 // 4M
)

func init() {
	disableColor := false
	if runtime.GOOS == "windows" {
		disableColor = true
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: disableColor,
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	go WebStart()
}

func main() {
	listener, err := net.Listen("tcp4", ":32768")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	session := &ElinkSession{
		UID:    uuid.New().String(),
		Conn:   conn,
		Mac:    GlobalMac,
		PerMac: "",
		Seq: &Seq{
			RecvSeq: 0,
			SendSeq: rand.Intn(100000),
		},
		key: nil,
	}

	defer func() {
		log.Info(conn.Close())
		SessionMap.Delete(session.UID)
	}()

	SessionMap.Store(session.UID, session)

	for {
		headerBuf := make([]byte, HeaderLen)
		_, err := io.ReadFull(conn, headerBuf)
		if err != nil {
			log.Println("read error", conn.RemoteAddr(), conn.LocalAddr())
			return
		}

		header, err := parseHeader(headerBuf)
		if err != nil {
			log.Println("parse header error:", err)
			return
		}
		if header.Len > MaxDataLen {
			log.Println("header len is greater than MaxDataLen:", MaxDataLen)
			return
		}
		dataBuf := make([]byte, header.Len)

		n, err := io.ReadFull(conn, dataBuf)
		if err != nil {
			log.Printf("not read enough data: read %d Bytes, but need %d Bytes\n", n, header.Len)
			return
		}

		packet := Elink{
			header,
			dataBuf,
			nil,
		}
		log.Println(session.handlePacket(conn, packet))
	}
}
