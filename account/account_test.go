package account

import (
	"github.com/stlwtr/pan/conf"
	"testing"
)

func TestAccount_UserInfo(t *testing.T) {
	accountClient := NewAccountClient(conf.TestData.AccessToken)
	res, err := accountClient.UserInfo()
	if err != nil {
		t.Fail()
	}
	t.Logf("TestAccount_UserInfo res: %+v", res)
}

func TestAccount_Quota(t *testing.T) {
	accountClient := NewAccountClient(conf.TestData.AccessToken)
	res, err := accountClient.Quota()
	if err != nil {
		t.Fail()
	}
	t.Logf("TestAccount_Quota res: %+v", res)
}
