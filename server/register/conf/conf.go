package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
	"github.com/oikomi/FishChatServer2/common/xtime"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	*commconf.CommConf
	configFile string
	RPCServer  *commconf.RPCServer
	Auth       *Auth
	Redis      *Redis
}

type Auth struct {
	Encryption string
	Salt       string
}

type Redis struct {
	*commconf.Redis
	Expire xtime.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "./register.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
