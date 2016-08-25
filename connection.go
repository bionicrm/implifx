package implifx

import (
	"net"
	"strconv"
	"github.com/bionicrm/controlifx"
)

type Connection struct {
	conn *net.UDPConn
}

func Listen(host string) (Connection, error) {
	return ListenOnPort(host, controlifx.DefaultPort)
}

func ListenOnPort(host string, port uint16) (o Connection, err error) {
	portStr := strconv.Itoa(int(port))

	laddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, portStr))
	if err != nil {
		return
	}

	o.conn, err = net.ListenUDP("udp", laddr)

	return
}

func (o Connection) LocalAddr() net.Addr {
	return o.conn.LocalAddr()
}

func (o Connection) Send(addr *net.UDPAddr, b []byte) error {
	_, err := o.conn.WriteTo(b, addr)

	return err
}

func (o Connection) Receive() (n int, msg ReceivableLanMessage, raddr *net.UDPAddr, err error) {
	for {
		b := make([]byte, controlifx.MaxReadSize)
		n, raddr, err = o.conn.ReadFromUDP(b)
		if err != nil {
			return
		}
		b = b[:n]

		msg = ReceivableLanMessage{}
		if err = msg.UnmarshalBinary(b); err == nil {
			break
		}
	}

	return
}

func (o Connection) Close() error {
	if o.conn != nil {
		return o.conn.Close()
	}

	return nil
}
