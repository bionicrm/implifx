package implifx

import (
	"github.com/bionicrm/controlifx"
	"fmt"
	"encoding/binary"
	"bytes"
	"math"
	"encoding"
)

type (
	SendableLanMessage   encoding.BinaryMarshaler
	ReceivableLanMessage controlifx.ReceivableLanMessage

	SetPowerLanMessage      controlifx.SetPowerLanMessage
	SetLabelLanMessage      controlifx.SetLabelLanMessage
	EchoRequestLanMessage   controlifx.EchoRequestLanMessage
	LightSetColorLanMessage controlifx.LightSetColorLanMessage
	LightSetPowerLanMessage controlifx.LightSetPowerLanMessage

	StateServiceLanMessage      controlifx.StateServiceLanMessage
	StateHostInfoLanMessage     controlifx.StateHostInfoLanMessage
	StateHostFirmwareLanMessage controlifx.StateHostFirmwareLanMessage
	StateWifiInfoLanMessage     controlifx.StateWifiInfoLanMessage
	StateWifiFirmwareLanMessage controlifx.StateWifiFirmwareLanMessage
	StatePowerLanMessage        controlifx.StatePowerLanMessage
	StateLabelLanMessage        controlifx.StateLabelLanMessage
	StateVersionLanMessage      controlifx.StateVersionLanMessage
	StateInfoLanMessage         controlifx.StateInfoLanMessage
	StateLocationLanMessage     controlifx.StateLocationLanMessage
	StateGroupLanMessage        controlifx.StateGroupLanMessage
	EchoResponseLanMessage      controlifx.EchoResponseLanMessage
	LightStateLanMessage        controlifx.LightStateLanMessage
	LightStatePowerLanMessage   controlifx.LightStatePowerLanMessage
)

func (o *ReceivableLanMessage) UnmarshalBinary(data []byte) error {
	// Header.
	o.Header = controlifx.LanHeader{}
	if err := o.Header.UnmarshalBinary(data[:controlifx.LanHeaderSize]); err != nil {
		return err
	}

	// Payload.
	payload, err := getReceivablePayloadOfType(o.Header.ProtocolHeader.Type)
	if err != nil {
		return err
	}
	if payload == nil {
		return nil
	}

	o.Payload = payload

	return o.Payload.UnmarshalBinary(data[controlifx.LanHeaderSize:])
}

func getReceivablePayloadOfType(t uint16) (encoding.BinaryUnmarshaler, error) {
	var payload encoding.BinaryUnmarshaler

	switch t {
	case controlifx.SetPowerType:
		payload = &SetPowerLanMessage{}
	case controlifx.SetLabelType:
		payload = &SetLabelLanMessage{}
	case controlifx.EchoRequestType:
		payload = &EchoRequestLanMessage{}
	case controlifx.LightSetColorType:
		payload = &LightSetColorLanMessage{}
	case controlifx.LightSetPowerType:
		payload = &LightSetPowerLanMessage{}
	case controlifx.GetServiceType:
		fallthrough
	case controlifx.GetHostInfoType:
		fallthrough
	case controlifx.GetHostFirmwareType:
		fallthrough
	case controlifx.GetWifiInfoType:
		fallthrough
	case controlifx.GetWifiFirmwareType:
		fallthrough
	case controlifx.GetPowerType:
		fallthrough
	case controlifx.GetLabelType:
		fallthrough
	case controlifx.GetVersionType:
		fallthrough
	case controlifx.GetInfoType:
		fallthrough
	case controlifx.GetLocationType:
		fallthrough
	case controlifx.GetGroupType:
		fallthrough
	case controlifx.LightGetType:
		fallthrough
	case controlifx.LightGetPowerType:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot create new payload of type %d; is it binary decodable?", t)
	}

	return payload, nil
}

func (o *SetPowerLanMessage) UnmarshalBinary(data []byte) error {
	o.Level = controlifx.PowerLevel(binary.LittleEndian.Uint16(data[:2]))

	return nil
}

func (o *SetLabelLanMessage) UnmarshalBinary(data []byte) error {
	o.Label = controlifx.Label(bytes.TrimRight(data[:32], "\x00"))

	return nil
}

func (o *EchoRequestLanMessage) UnmarshalBinary(data []byte) error {
	copy(o.Payload[:], data[:64])

	return nil
}

func (o *LightSetColorLanMessage) UnmarshalBinary(data []byte) error {
	if err := o.Color.UnmarshalBinary(data[1:9]); err != nil {
		return err
	}

	o.Duration = binary.LittleEndian.Uint32(data[9:13])

	return nil
}

func (o *LightSetPowerLanMessage) UnmarshalBinary(data []byte) error {
	o.Level = controlifx.PowerLevel(binary.LittleEndian.Uint16(data[:2]))
	o.Duration = binary.LittleEndian.Uint32(data[2:6])

	return nil
}

func (o StateServiceLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 5)

	// Service.
	data[0] = byte(o.Service)

	// Port.
	binary.LittleEndian.PutUint32(data[1:], o.Port)

	return
}

func (o StateHostInfoLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 12)

	// Signal.
	binary.LittleEndian.PutUint32(data[:4], math.Float32bits(o.Signal))

	// Tx.
	binary.LittleEndian.PutUint32(data[4:8], o.Tx)

	// Rx.
	binary.LittleEndian.PutUint32(data[8:12], o.Rx)

	return
}

func (o StateHostFirmwareLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 12)

	// Build.
	binary.LittleEndian.PutUint64(data[:8], o.Build)

	// Version.
	binary.LittleEndian.PutUint32(data[8:], o.Version)

	return
}

func (o StateWifiInfoLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 12)

	// Signal.
	binary.LittleEndian.PutUint32(data[:4], math.Float32bits(o.Signal))

	// Tx.
	binary.LittleEndian.PutUint32(data[4:8], o.Tx)

	// Rx.
	binary.LittleEndian.PutUint32(data[8:12], o.Rx)

	return
}

func (o StateWifiFirmwareLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 12)

	// Build.
	binary.LittleEndian.PutUint64(data[:8], o.Build)

	// Version.
	binary.LittleEndian.PutUint32(data[8:], o.Version)

	return
}

func (o StatePowerLanMessage) MarshalBinary() ([]byte, error) {
	return o.Level.MarshalBinary()
}

func (o StateLabelLanMessage) MarshalBinary() ([]byte, error) {
	return o.Label.MarshalBinary()
}

func (o StateVersionLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 12)

	// Vendor.
	binary.LittleEndian.PutUint32(data[:4], o.Vendor)

	// Product.
	binary.LittleEndian.PutUint32(data[4:8], o.Product)

	// Version.
	binary.LittleEndian.PutUint32(data[8:], o.Version)

	return
}

func (o StateInfoLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 24)

	// Time.
	binary.LittleEndian.PutUint64(data[:8], uint64(o.Time))

	// Uptime.
	binary.LittleEndian.PutUint64(data[8:16], o.Uptime)

	// Downtime.
	binary.LittleEndian.PutUint64(data[16:24], o.Downtime)

	return
}

func (o StateLocationLanMessage) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 56)

	// Location.
	copy(data[:16], o.Location[:])

	// Label.
	b, err := o.Label.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[16:48], b)

	// Updated at.
	b, err = o.UpdatedAt.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[48:], b)

	return
}

func (o StateGroupLanMessage) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 56)

	// Location.
	copy(data[:16], o.Group[:])

	// Label.
	b, err := o.Label.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[16:48], b)

	// Updated at.
	b, err = o.UpdatedAt.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[48:], b)

	return
}

func (o EchoResponseLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 64)

	copy(data, o.Payload[:])

	return
}

func (o LightStateLanMessage) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 44)

	// Color.
	b, err := o.Color.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[:8], b)

	// Power.
	b, err = o.Power.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[10:12], b)

	// Label.
	b, err = o.Label.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[12:], b)

	return
}

func (o LightStatePowerLanMessage) MarshalBinary() ([]byte, error) {
	return o.Level.MarshalBinary()
}
