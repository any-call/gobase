package mysql

import "testing"

func TestSum(t *testing.T) {
	t.Log(Sum("case when money_type = ? Then money else 0.0 end", "0.0", "total_flow"))
}
