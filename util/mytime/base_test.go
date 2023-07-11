package mytime

import (
	"fmt"
	"testing"
	"time"
)

func Test_Truncate(t *testing.T) {
	now := time.Now()

	fmt.Println("origin time:", now)
	clearMillSec := TruncateMillSec(now)
	clearSec := TruncateSec(now)
	fmt.Println("origin time:", now)
	fmt.Println("clear MillSec time:", clearMillSec)
	fmt.Println("clear sec time:", clearSec)
}
