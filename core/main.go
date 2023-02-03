package core

import (
	"fmt"
	"github.com/zhukovra/retracker/core/common"
	Receiver "github.com/zhukovra/retracker/core/receiver"
	Storage "github.com/zhukovra/retracker/core/storage"
	"net/http"
)

type Core struct {
	Config   *common.Config
	Storage  *Storage.Storage
	Receiver *Receiver.Receiver
}

func New(config *common.Config) *Core {
	storage := Storage.New(config)
	core := Core{
		Config:   config,
		Storage:  storage,
		Receiver: Receiver.New(config, storage),
	}
	http.HandleFunc("/announce", core.Receiver.Announce.HttpHandler)
	if err := http.ListenAndServe(config.Listen, nil); err != nil { // set listen port
		fmt.Println(err)
	}
	//TODO: do it with context
	return &core
}
