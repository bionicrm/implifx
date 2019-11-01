package implifx

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"github.com/yath/controlifx"
	"math"
)

type (
	SendableLanMessage   encoding.BinaryMarshaler
	ReceivableLanMessage controlifx.ReceivableLanMessage

	SetPowerLanMessage      controlifx.SetPowerLanMessage
	SetLabelLanMessage      controlifx.SetLabelLanMessage
	SetOwnerLanMessage      controlifx.SetOwnerLanMessage
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
	StateOwnerLanMessage        controlifx.StateOwnerLanMessage
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
	case controlifx.SetOwnerType:
		payload = &SetOwnerLanMessage{}
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
	case controlifx.GetOwnerType:
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
	// Level.
	o.Level = binary.LittleEndian.Uint16(data[:2])

	return nil
}

func (o *SetLabelLanMessage) UnmarshalBinary(data []byte) error {
	// Label.
	o.Label = controlifx.BToStr(data[:32])

	return nil
}

func (o *SetOwnerLanMessage) UnmarshalBinary(data []byte) error {
	// Owner.
	copy(o.Owner[:], data[:16])

	// Label.
	o.Label = controlifx.BToStr(data[16:48])

	// Updated at.
	o.UpdatedAt = binary.LittleEndian.Uint64(data[48:])

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
	o.Level = binary.LittleEndian.Uint16(data[:2])
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

func (o StatePowerLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 2)

	// Level.
	binary.LittleEndian.PutUint16(data, o.Level)

	return
}

func (o StateLabelLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 32)

	// Label.
	copy(data, o.Label)

	return
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

func (o StateLocationLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 56)

	// Location.
	copy(data[:16], o.Location[:])

	// Label.
	copy(data[16:48], o.Label)

	// Updated at.
	binary.LittleEndian.PutUint64(data[48:], o.UpdatedAt)

	return
}

func (o StateGroupLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 56)

	// Location.
	copy(data[:16], o.Group[:])

	// Label.
	copy(data[16:48], o.Label)

	// Updated at.
	binary.LittleEndian.PutUint64(data[48:], o.UpdatedAt)

	return
}

func (o StateOwnerLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 56)

	// Owner.
	copy(data[:16], o.Owner[:])

	// Label.
	copy(data[16:48], o.Label)

	// Updated at.
	binary.LittleEndian.PutUint64(data[48:], o.UpdatedAt)

	return
}

func (o EchoResponseLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 64)

	// Payload.
	copy(data, o.Payload[:])

	return
}

func (o LightStateLanMessage) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 52)

	// Color.
	b, err := o.Color.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[:8], b)

	// Power.
	binary.LittleEndian.PutUint16(data[10:12], o.Power)

	// Label.
	copy(data[12:44], o.Label)

	return
}

func (o LightStatePowerLanMessage) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 2)

	// Level.
	binary.LittleEndian.PutUint16(data, o.Level)

	return
}
