package mycert

import (
	"testing"

	"github.com/any-call/gobase/util/mytime"
)

func TestRenewCertificate(t *testing.T) {
	tt, err := CheckCertExpiryEx("proxy.stcoin.uk", "43.227.112.128", 15433)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("expired is :%s", tt.In(mytime.CstSh).Format(mytime.FormatYYYYMMDDHHMISS))
}
