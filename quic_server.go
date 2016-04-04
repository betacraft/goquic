package goquic

import (
	"errors"
	"net"

	"github.com/vanillahsu/go_reuseport"
)

type QuicServer struct {
	Q        *QuicSpdyServer
	Writer   *ServerWriter
	ReadChan chan UdpData
	StatChan chan statCallback
	conn     net.PacketConn
}

func NewQuicServer(addr string) (*QuicServer, error) {
	_, err := parsePort(addr)
	if err != nil {
		return nil, err
	}
	qs := &QuicServer{Q: &QuicSpdyServer{Addr: addr, numOfServers: 1, isSecure: false}}
	conn, err := reuseport.NewReusablePortPacketConn("udp4", qs.Q.Addr)
	if err != nil {
		return nil, err
	}
	qs.conn = conn
	return qs, nil
}

func (qs *QuicServer) Serve() error {
	udp_conn, ok := qs.conn.(*net.UDPConn)
	if !ok {
		return errors.New("ListenPacket did not return net.UDPConn")
	}

	listen_addr, err := net.ResolveUDPAddr("udp", udp_conn.LocalAddr().String())
	if err != nil {
		return err
	}
	wch := make(chan UdpData, 500) // TODO(serialx, hodduc):
	qs.Writer = NewServerWriter(wch)
	qs.ReadChan = make(chan UdpData, 500)
	qs.StatChan = make(chan statCallback)

	go qs.Q.Serve(listen_addr, qs.Writer, qs.ReadChan, qs.StatChan)

	readFunc := func(conn *net.UDPConn, srv *QuicServer) {
		buf := make([]byte, 65535)

		for {
			n, peer_addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				panic(err)
			}
			buf_new := make([]byte, n)
			copy(buf_new, buf[:n])
			srv.ReadChan <- UdpData{Addr: peer_addr, Buf: buf_new}
		}
	}

	writeFunc := func(conn *net.UDPConn, writer *ServerWriter) {
		for dat := range writer.Ch {
			conn.WriteToUDP(dat.Buf, dat.Addr)
		}
	}

	go writeFunc(udp_conn, qs.Writer)
	go readFunc(udp_conn, qs)

	return nil
}

func (qs *QuicServer) Close() {
	qs.conn.Close()
}
