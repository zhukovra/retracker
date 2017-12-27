package announce

import (
	"github.com/vvampirius/retracker/bittorrent/tracker"
	"github.com/vvampirius/retracker/bittorrent/common"
	"fmt"
	CoreCommon "github.com/vvampirius/retracker/core/common"
)

func (self *Announce) ProcessAnnounce(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
	event string) *tracker.Response {
		if request, err := tracker.MakeRequest(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
			event, self.Logger); err==nil {
			if self.Logger != nil {	self.Logger.Println(request.String()) }

			response := tracker.Response{
				Interval: 30,
			}

			if request.Event != `stopped` {
				self.Storage.Update(*request)
				response.Peers = self.Storage.GetPeers(request.InfoHash)
				response.Peers = append(response.Peers, self.makeForwards(*request)...)
			} else {
				self.Storage.Delete(*request)
				//TODO: make another response ?
			}

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
			fmt.Println(<-ch)
		}
		fmt.Println("makeForwards exit")
	}
	return peers
}

func (self *Announce) makeForward(forward CoreCommon.Forward, request tracker.Request, ch chan<- []common.Peer) {
	peers := make([]common.Peer, 0)
	fmt.Println(forward, request)
	ch <- peers
}