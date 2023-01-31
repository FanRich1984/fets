package idmanage

import (
	"fmt"
	"testing"
	"time"

	"github.com/FanRich1984/fets/vatools"
)

func TestInviationCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Microsecond * time.Duration(vatools.CRnd(100, 300)))
		id, _ := ObFlakId.GetId()
		str := CreateInvitationCode(id)
		fmt.Println(id, " => ", str)
	}
}
