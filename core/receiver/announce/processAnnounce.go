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
	"context"
	"time"
	"os"
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
	logger := *self.Logger
	logger.SetPrefix(fmt.Sprintf("makeForwards(%x) ", request.InfoHash))
	if forwardsCount > 0 {
		if self.Config.Debug { logger.Printf("Making forwards to %d forwarders\n", forwardsCount)}
		ch := make(chan []common.Peer, forwardsCount)
		ctx, _  := context.WithTimeout(context.Background(), time.Second * time.Duration(self.Config.ForwardTimeout))
		for _, v := range self.Config.Forwards {
			go self.makeForward(v, request, ch)
		}
		for i := 0; i < forwardsCount; i++ {
			//peers = append(peers, <-ch...)
			select {
				case prs := <-ch:
					if self.Config.Debug { logger.Printf("Got %d peers\n", len(prs))}
					peers = append(peers, prs...)
				case <-ctx.Done():
					if self.Config.Debug { logger.Println(`Got timeout`) }
					i = forwardsCount
			}
		}
		if self.Config.Debug { logger.Printf("Finished with %d peers\n", len(peers)) }
	}
	return peers
}

func (self *Announce) makeForward(forward CoreCommon.Forward, request tracker.Request, ch chan<- []common.Peer) {
	peers := make([]common.Peer, 0)
	logger := *self.Logger
	logger.SetPrefix(fmt.Sprintf("makeForward(%x, %s) ", request.InfoHash, forward.Uri))
	uri := fmt.Sprintf("%s?info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d", forward.Uri, url.QueryEscape(string(request.InfoHash)),
		url.QueryEscape(string(request.PeerID)), request.Port, request.Uploaded, request.Downloaded, request.Left)
	if forward.Ip != `` {
		uri = fmt.Sprintf("%s&ip=%s&ipv4=%s", uri, forward.Ip, forward.Ip) //TODO: check for IPv4
	}
	if self.Config.Debug { logger.Println(uri) }
	if resp, err := http.Get(uri); err==nil {
		if resp.StatusCode == http.StatusOK {
			var tmpFileName string
			if b, err := ioutil.ReadAll(resp.Body); err==nil {
				if self.Config.Debug {
					if f, err := ioutil.TempFile(os.TempDir(), "bencode_"); err==nil {
						tmpFileName = f.Name()
						fmt.Fprintln(f, request.InfoHash)
						fmt.Fprintln(f, forward.Uri)
						if _, err := f.Write(b); err!=nil { logger.Println(err.Error())	}
						f.Close()
					} else { logger.Println(err.Error()) }
				}
				if response, err := Response.Load(b); err==nil {
					if self.Config.Debug { logger.Printf("Got %d peers\n", len(response.Peers)) }
					peers = append(peers, response.Peers...)
				}  else { logger.Printf("Can't load response (%s): %s\n", tmpFileName, err.Error()) }
				if err := resp.Body.Close(); err!=nil { logger.Printf("Can't Close() body: %s\n", err.Error()) }
			} else { logger.Printf("Can't read body: %s\n", err.Error()) }
		} else { logger.Printf("HTTP error: %d %s\n", resp.StatusCode, resp.Status) }
	} else { logger.Println(err.Error()) }
	ch <- peers
}