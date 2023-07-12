package mytime

import (
	"fmt"
	"testing"
	"time"
)

func Test_Truncate(t *testing.T) {
	now := time.Now()

	clearMillSec := TruncateMillSec(now)
	clearSec := TruncateSec(now)
	clearMin := TruncateMinute(now)
	fmt.Println("origin time:", now)
	fmt.Println("clear MillSec time:", clearMillSec)
	fmt.Println("clear sec time:", clearSec)
	fmt.Println("clear minute time:", clearMin)
}
