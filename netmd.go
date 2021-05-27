package netmd

import (
	"bytes"
	"errors"
	"github.com/google/gousb"
	"log"
	"time"
)

type NetMD struct {
	debug bool
	index int
	devs  []*gousb.Device
	ctx   *gousb.Context
	out   *gousb.OutEndpoint
	ekb   *EKB
}

type Encoding byte

type Channels byte

const (
	EncSP  Encoding = 0x90
	EncLP2 Encoding = 0x92
	EncLP4 Encoding = 0x93

	ChanStereo Channels = 0x00
	ChanMono   Channels = 0x01
)

var (
	ByteArr16 = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func NewNetMD(index int, debug bool) (md *NetMD, err error) {
	md = &NetMD{
		index: index,
		debug: debug,
		ekb:   NewEKB(),
	}

	md.ctx = gousb.NewContext()
	md.devs, err = md.ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		for _, d := range Devices {
			if d.deviceId == desc.Product && d.vendorId == desc.Vendor {
				if md.debug {
					log.Printf("Found %s", d.name)
				}
				return true
			}
		}
		return false
	})

	if err != nil {
		return
	}

	if len(md.devs) == 0 || len(md.devs) <= md.index {
		err = errors.New("no compatible netmd device found or incorrect index")
		return
	}

	for num := range md.devs[md.index].Desc.Configs {
		config, _ := md.devs[md.index].Config(num)
		for _, desc := range config.Desc.Interfaces {
			intf, _ := config.Interface(desc.Number, 0)
			for _, endpointDesc := range intf.Setting.Endpoints {
				if endpointDesc.Direction == gousb.EndpointDirectionOut {
					if md.out, err = intf.OutEndpoint(endpointDesc.Number); err != nil {
						return
					}
					if md.debug {
						log.Printf("%s", endpointDesc)
					}
				}
			}
			config.Close()
		}
	}
	return
}

func (md *NetMD) Close() {
	for _, d := range md.devs {
		d.Close()
	}
	md.ctx.Close()
}

// Wait makes sure the device is truly finished, needed to prevent crashes on the SHARP IM-DR410/IM-DR420
// and the Sony MZ-N420D
func (md *NetMD) Wait() error {
	buf := make([]byte, 4)
	for i := 0; i < 10; i++ {
		c, err := md.devs[md.index].Control(gousb.ControlIn|gousb.ControlVendor|gousb.ControlInterface, 0x01, 0, 0, buf)
		if err != nil {
			return err
		}
		if c != 4 {
			if md.debug {
				log.Println("sync response != 4 bytes")
			}
		} else {
			if bytes.Equal(buf, []byte{0x00, 0x00, 0x00, 0x00}) {
				return nil
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	return errors.New("no sync response")
}

// RequestDiscCapacity returns the totals in seconds
func (md *NetMD) RequestDiscCapacity() (recorded uint64, total uint64, available uint64, err error) {
	md.release()
	r, err := md.rawCall([]byte{0x00, 0x18, 0x06, 0x02, 0x10, 0x10, 0x00}, []byte{0x30, 0x80, 0x03, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return
	}
	recorded = (hexToInt(r[29]) * 3600) + (hexToInt(r[30]) * 60) + hexToInt(r[31])
	total = (hexToInt(r[35]) * 3600) + (hexToInt(r[36]) * 60) + hexToInt(r[37])
	available = (hexToInt(r[42]) * 3600) + (hexToInt(r[43]) * 60) + hexToInt(r[44])
	return
}

// SetDiscHeader will write  a raw title to the disc
func (md *NetMD) SetDiscHeader(t string) error {
	md.poll()
	o, err := md.RequestDiscHeader()
	if err != nil {
		return err
	}
	j := len(o) // length of old title
	h := len(t) // length of new title
	c := []byte{0x00, 0x00, 0x30, 0x00, 0x0a, 0x00, 0x50, 0x00}
	c = append(c, intToHex16(int16(h))...)
	c = append(c, 0x00, 0x00)
	c = append(c, intToHex16(int16(j))...)
	c = append(c, []byte(t)...)
	_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x18, 0x01, 0x01}, []byte{0x00})
	_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x18, 0x01, 0x00}, []byte{0x00})
	_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x18, 0x01, 0x03}, []byte{0x00})
	_, err = md.rawCall([]byte{0x00, 0x18, 0x07, 0x02, 0x20, 0x18, 0x01}, c) // actual call
	_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x18, 0x01, 0x00}, []byte{0x00})
	if err != nil {
		return err
	}
	return nil
}

// RequestDiscHeader returns the raw title of the disc
func (md *NetMD) RequestDiscHeader() (string, error) {
	r, err := md.rawCall([]byte{0x00, 0x18, 0x06, 0x02, 0x20, 0x18, 0x01}, []byte{0x00, 0x00, 0x30, 0x00, 0x0a, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return "", err
	}
	return string(r[25:]), nil
}

// RecordingParameters current default recording parameters set on the NetMD
func (md *NetMD) RecordingParameters() (encoding Encoding, channels Channels, err error) {
	r, err := md.rawCall([]byte{0x00, 0x18, 0x09, 0x80, 0x01, 0x03, 0x30}, []byte{0x88, 0x01, 0x00, 0x30, 0x88, 0x05, 0x00, 0x30, 0x88, 0x07, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return
	}
	encoding = Encoding(r[34])
	channels = Channels(r[35])
	return
}

// RequestStatus returns known status flags
func (md *NetMD) RequestStatus() (disk bool, err error) {
	//_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x80, 0x00, 0x01}, []byte{0x00})
	r, err := md.rawCall([]byte{0x00, 0x18, 0x09, 0x80, 0x01, 0x02, 0x30}, []byte{0x88, 0x00, 0x00, 0x30, 0x88, 0x04, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return
	}
	disk = r[26] == 0x40 // 0x80 no disk
	return
}

func (md *NetMD) RequestTrackCount() (c int, err error) {
	_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x10, 0x01, 0x01}, []byte{0x00})
	r, err := md.rawCall([]byte{0x00, 0x18, 0x06, 0x02, 0x10, 0x10, 0x01}, []byte{0x30, 0x00, 0x10, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return
	}
	c = int(hexToInt16(r[23:]))
	return
}

// RequestTrackTitle returns the raw title of the trk number starting from 0
func (md *NetMD) RequestTrackTitle(trk int) (t string, err error) {
	r, err := md.rawCall([]byte{0x00, 0x18, 0x06, 0x02, 0x20, 0x18, byte(2) & 0xff}, []byte{0x00, byte(trk) & 0xff, 0x30, 0x00, 0x0a, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return
	}
	t = string(r[25:])
	return
}

// SetTrackTitle set the title of the trk number starting from 0, isNew can be be true if it's a newadded track
func (md *NetMD) SetTrackTitle(trk int, t string, isNew bool) (err error) {
	j := 0
	if !isNew {
		o, err := md.RequestTrackTitle(trk)
		if err != nil {
			return err
		}
		j = len(o) // length of old title
	}
	h := len(t) // length of new title
	s := []byte{0x00, byte(trk) & 0xff, 0x30, 0x00, 0x0a, 0x00, 0x50, 0x00}
	s = append(s, intToHex16(int16(h))...)
	s = append(s, 0x00, 0x00)
	s = append(s, intToHex16(int16(j))...)
	s = append(s, []byte(t)...)

	if !isNew {
		_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x18, 0x02, 0x00}, []byte{0x00})
		_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x18, 0x02, 0x03}, []byte{0x00})
	}

	_, err = md.rawCall([]byte{0x00, 0x18, 0x07, 0x02, 0x20, 0x18, byte(2) & 0xff}, s)

	if !isNew {
		_, err = md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x18, 0x02, 0x00}, []byte{0x00})
	}

	if err != nil {
		return
	}
	return
}

// EraseTrack will erase the trk number starting from 0
func (md *NetMD) EraseTrack(trk int) error {
	s := []byte{0x10, 0x01}
	s = append(s, intToHex16(int16(trk))...)
	_, err := md.rawCall([]byte{0x00, 0x18, 0x40, 0xff, 0x01, 0x00, 0x20}, s)
	if err != nil {
		return err
	}
	return nil
}

// MoveTrack will move the trk number to a new position
func (md *NetMD) MoveTrack(trk, to int) error {
	s := []byte{0x10, 0x01}
	s = append(s, intToHex16(int16(trk))...)
	s = append(s, 0x20, 0x10, 0x01)
	s = append(s, intToHex16(int16(to))...)
	_, err := md.rawCall([]byte{0x00, 0x18, 0x08, 0x10, 0x10, 0x01, 0x00}, []byte{0x00})
	_, err = md.rawCall([]byte{0x00, 0x18, 0x43, 0xff, 0x00, 0x00, 0x20}, s)
	if err != nil {
		return err
	}
	return nil
}

// RequestTrackLength returns the duration in seconds of the trk starting from 0
func (md *NetMD) RequestTrackLength(trk int) (duration uint64, err error) {
	s := []byte{0x01}
	s = append(s, intToHex16(int16(trk))...)
	s = append(s, 0x30, 0x00, 0x01, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00)
	r, err := md.rawCall([]byte{0x00, 0x18, 0x06, 0x02, 0x20, 0x10}, s)
	if err != nil {
		return
	}
	duration = (hexToInt(r[27]) * 3600) + (hexToInt(r[28]) * 60) + hexToInt(r[29])
	return
}

// RequestTrackEncoding returns the Encoding of the trk starting from 0
func (md *NetMD) RequestTrackEncoding(trk int) (encoding Encoding, err error) {
	s := append(intToHex16(int16(trk)), 0x30, 0x80, 0x07, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00)
	r, err := md.rawCall([]byte{0x00, 0x18, 0x06, 0x02, 0x20, 0x10, 0x01}, s)
	if err != nil {
		return
	}
	return Encoding(r[len(r)-2]), nil
}

func (md *NetMD) rawCall(chk []byte, payload []byte) ([]byte, error) {
	i := append(chk, payload...)
	if md.debug {
		log.Printf("md.rawCall send <- % x", i)
	}

	md.poll()
	if _, err := md.devs[md.index].Control(gousb.ControlOut|gousb.ControlVendor|gousb.ControlInterface, 0x80, 0, 0, i); err != nil {
		return nil, err
	}

	for tries := 0; tries < 10; tries++ {
		if h := md.poll(); h != -1 {
			b, err := md.receive(h)
			if err != nil {
				return nil, err
			}

			if bytes.Equal(b[1:len(chk)], chk[1:]) {
				return b, nil
			} else {
				if md.debug {
					log.Printf("Skipping mismatch: % x <-> % x", b[1:len(chk)], chk[1:])
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}

	return nil, errors.New("poll failed")
}

// acquire is part of SHARP NetMD protocols and probably do nothing on Sony devices
func (md *NetMD) acquire() error {
	_, err := md.rawCall([]byte{0x00, 0xff, 0x01}, []byte{0x0c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	if err != nil {
		return err
	}
	return nil
}

// release is part of the acquire lifecycle
func (md *NetMD) release() error {
	_, err := md.rawCall([]byte{0x00, 0xff, 0x01}, []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	if err != nil {
		return err
	}
	return nil
}

func (md *NetMD) receive(s int) ([]byte, error) {
	buf := make([]byte, s)
	if _, err := md.devs[md.index].Control(gousb.ControlIn|gousb.ControlVendor|gousb.ControlInterface, 0x81, 0, 0, buf); err != nil {
		return nil, err
	}
	if md.debug {
		if buf[0] == 0x0a {
			log.Printf(" -> Rejected -> % x", buf)
			return nil, errors.New("controlIn was rejected")
		} else if buf[0] == 0x09 {
			log.Printf(" -> Accepted -> % x", buf)
		} else if buf[0] == 0x0f {
			log.Printf(" -> Interim <-")
		} else if buf[0] == 0x08 {
			log.Printf(" -> notImplemented <-")
		} else {
			log.Printf(" -> Unknown  -> % x", buf)
		}
	}
	return buf, nil
}

func (md *NetMD) poll() int {
	buf := make([]byte, 4)
	md.devs[md.index].Control(gousb.ControlIn|gousb.ControlVendor|gousb.ControlInterface, 0x01, 0, 0, buf)
	if buf[0] == 0x01 { //&& buf[1] == 0x81
		return int(buf[2])
	}
	return -1
}
