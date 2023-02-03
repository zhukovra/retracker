package tracker

import (
	"github.com/zeebo/bencode"
	"github.com/zhukovra/retracker/bittorrent/common"
)

type Response struct {
	Interval int           `bencode:"interval"`
	Peers    []common.Peer `bencode:"peers"`
}

func (self *Response) Bencode() (string, error) {
	return bencode.EncodeString(self)
}
