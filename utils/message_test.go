package utils

import (
	"testing"
)

func TestSendDingMsg(t *testing.T) {
	webHook := "https://oapi.dingtalk.com/robot/send?access_token=80d19d9084c2ac0231a223afe49255d5f18a9cdc144bacbf764b9f292c7c0f7a"
	SendDingMsg(webHook,"test")
}
