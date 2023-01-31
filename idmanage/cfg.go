package idmanage

import (
	"github.com/FanRich1984/fets/config"
	"github.com/FanRich1984/fets/vatools"
)

type Cfg struct {
	MachineId int16 // 机器ID，如果为零则使用本机IP地址
	Quantity  int16 // 要开启的数量
}

func NewCfg() *Cfg {
	cfg := &Cfg{
		MachineId: 0,
		Quantity:  1,
	}
	c := config.NewConfig("./cfg.ini")
	if mp, ok := c.GetNode("idmanage"); ok {
		for k, v := range mp {
			switch k {
			case "machineid":
				cfg.MachineId = vatools.SInt16(v)
			case "quantity":
				cfg.Quantity = vatools.SInt16(v)
			}
		}
	}
	return cfg
}
