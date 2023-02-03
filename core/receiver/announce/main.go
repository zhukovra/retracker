package announce

import (
	CoreCommon "github.com/zhukovra/retracker/core/common"
	Storage "github.com/zhukovra/retracker/core/storage"
	"log"
	"os"
)

type Announce struct {
	Config  *CoreCommon.Config
	Logger  *log.Logger
	Storage *Storage.Storage
}

func New(config *CoreCommon.Config, storage *Storage.Storage) *Announce {
	announce := Announce{
		Config:  config,
		Logger:  log.New(os.Stdout, `announce `, log.Flags()),
		Storage: storage,
	}
	return &announce
}
