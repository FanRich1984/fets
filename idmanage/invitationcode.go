// 邀请码生成
// 通过用户唯一数字ID生成唯一的邀请码

package idmanage

import "strings"

var invitationCodeSet = []rune{
	'1', '2', '3', '4', '5', '6', '7', '8', '9',
	'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q',
	// 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

var invitationCodeSetLen uint64 = uint64(len(invitationCodeSet))

// 通过用户ID生成邀请码
func CreateInvitationCode(id uint64) string {
	var code []rune
	for i := 0; i < 100; i++ {
		idx := id % invitationCodeSetLen
		code = append(code, invitationCodeSet[idx])
		id = id / invitationCodeSetLen
		if id <= invitationCodeSetLen {
			code = append(code, invitationCodeSet[id])
			break
		}
	}

	return string(code)
}

// 通过用户ID生成缴请码并转为大写
func CreateBigInvitationCode(id uint64) string {
	code := CreateInvitationCode(id)

	return strings.ToUpper(code)
}
