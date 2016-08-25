package implifx

import (
	"net"
	"strconv"
	"github.com/bionicrm/controlifx"
	"encoding"
)

type Connection struct {
	Mac uint64

	conn *net.UDPConn
}

func Listen(host string) (Connection, error) {
	return ListenOnOtherPort(host, strconv.Itoa(controlifx.DefaultPort))
}

func ListenOnOtherPort(host, port string) (o Connection, err error) {
	laddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))
	if err != nil {
		return
	}

	o.conn, err = net.ListenUDP("udp", laddr)

	return
}

func (o Connection) LocalAddr() net.Addr {
	return o.conn.LocalAddr()
}

func (o Connection) Port() uint16 {
	// An invalid port will never be returned, so error checking is not
	// necessary.
	_, portStr, _ := net.SplitHostPort(o.LocalAddr().String())
	port, _ := strconv.Atoi(portStr)

	return uint16(port)
}

func (o Connection) Send(addr *net.UDPAddr, b []byte) error {
	_, err := o.conn.WriteTo(b, addr)

	return err
}

func (o Connection) Receive() (n int, raddr *net.UDPAddr, msg ReceivableLanMessage, err error) {
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

func (o Connection) Respond(always bool, triggeringAddr *net.UDPAddr, triggeringMsg ReceivableLanMessage, t uint16, payload encoding.BinaryMarshaler) (tx int, _ error) {
	msg := controlifx.SendableLanMessage{
		Header:controlifx.LanHeader{
			Frame:controlifx.LanHeaderFrame{
				Size:controlifx.LanHeaderSize,
				Source:triggeringMsg.Header.Frame.Source,
			},
			FrameAddress:controlifx.LanHeaderFrameAddress{
				Target:o.Mac,
				Sequence:triggeringMsg.Header.FrameAddress.Sequence,
			},
		},
	}

	send := func() (int, error) {
		b, err := msg.MarshalBinary()
		if err != nil {
			return 0, err
		}

		return len(b), o.Send(triggeringAddr, b)
	}

	if triggeringMsg.Header.FrameAddress.AckRequired {
		msg.Header.ProtocolHeader.Type = controlifx.AcknowledgementType

		if n, err := send(); err != nil {
			return 0, err
		} else {
			tx = n
		}
	}

	if always || triggeringMsg.Header.FrameAddress.ResRequired {
		msg.Header.ProtocolHeader.Type = t
		msg.Payload = payload

		b, err := payload.MarshalBinary()
		if err != nil {
			return 0, err
		}
		msg.Header.Frame.Size += uint16(len(b))

		if n, err := send(); err != nil {
			return tx, err
		} else {
			tx += n
		}
	}

	return
}
