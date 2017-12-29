package announce

import (
	"github.com/vvampirius/retracker/bittorrent/tracker"
	Response "github.com/vvampirius/retracker/bittorrent/response"
	"github.com/vvampirius/retracker/bittorrent/common"
	"fmt"
	CoreCommon "github.com/vvampirius/retracker/core/common"
	"net/http"
	"net/url"
	"io/ioutil"
)

func (self *Announce) ProcessAnnounce(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
	event string) *Response.Response {
		if request, err := tracker.MakeRequest(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
			event, self.Logger); err==nil {
			if self.Logger != nil {	self.Logger.Println(request.String()) }

			response := Response.Response{
				Interval: 30,
			}

			if request.Event != `stopped` {
				self.Storage.Update(*request)
				response.Peers = self.Storage.GetPeers(request.InfoHash)
				response.Peers = append(response.Peers, self.makeForwards(*request)...)
			} else { self.Storage.Delete(*request) }

			return &response
		} else {
			if self.Logger != nil {	self.Logger.Println(err.Error()) }
		}

		return nil
}

func (self *Announce) makeForwards(request tracker.Request) []common.Peer {
	peers := make([]common.Peer, 0)
	forwardsCount := len(self.Config.Forwards)
	if forwardsCount > 0 {
		ch := make(chan []common.Peer, forwardsCount)
		for _, v := range self.Config.Forwards {
			self.makeForward(v, request, ch)
		}
		for i := 0; i < forwardsCount; i++ {
			peers = append(peers, <-ch...)
		}
		fmt.Println("makeForwards exit")
	}
	return peers
}

func (self *Announce) makeForward(forward CoreCommon.Forward, request tracker.Request, ch chan<- []common.Peer) {
	peers := make([]common.Peer, 0)
	fmt.Println(forward, request)
	uri := fmt.Sprintf("%s?info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d", forward.Uri, url.QueryEscape(string(request.InfoHash)),
		url.QueryEscape(string(request.PeerID)), request.Port, request.Uploaded, request.Downloaded, request.Left)
	if forward.Ip != `` {
		uri = fmt.Sprintf("%s&ip=%s", uri, forward.Ip)
	}
	fmt.Println(uri)
	if resp, err := http.Get(uri); err==nil {
		if resp.StatusCode == http.StatusOK {
			if b, err := ioutil.ReadAll(resp.Body); err==nil {
				//if f, err := ioutil.TempFile("/tmp", "bencode_"); err==nil {
				//	if _, err := f.Write(b); err!=nil {
				//		fmt.Println(forward.Uri, err.Error())
				//	}
				//	f.Close()
				//} else { fmt.Println(forward.Uri, err.Error()) }
				if response, err := Response.Load(b); err==nil {
					peers = append(peers, response.Peers...)
				}  else { fmt.Println(forward.Uri, string(b), err.Error()) }
				if err := resp.Body.Close(); err!=nil { fmt.Println(forward.Uri, err.Error()) }
			} else { fmt.Println(forward.Uri, err.Error()) }
		} else { fmt.Println(forward.Uri, resp.Status) }
	} else { fmt.Println(forward.Uri, err.Error()) }
	ch <- peers
}