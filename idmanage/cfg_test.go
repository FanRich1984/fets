package idmanage

import (
	"fmt"
	"testing"
)

func TestCfg(t *testing.T) {
	c := NewCfg()
	fmt.Println("MachineId:", c.MachineId)
	fmt.Println("Quantity:", c.Quantity)
}
