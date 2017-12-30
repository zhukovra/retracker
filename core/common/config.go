package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Listen string
	Debug bool
	Age float64
	XRealIP bool
	Forwards []Forward
	ForwardTimeout int
}

func (self *Config) ReloadForwards(fileName string) error {
	if b, err := ioutil.ReadFile(fileName); err==nil {
		self.Forwards = make([]Forward, 0)
		return yaml.Unmarshal(b, &self.Forwards)
	} else { return err }
	return nil
}