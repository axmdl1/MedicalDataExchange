package blockchain

import (
	"fmt"
	"time"
)

func logBc(format string, args ...interface{}) {
	ts := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("🟦 [BLOCKCHAIN %s] %s\n", ts, msg)
}

func logBcErr(format string, args ...interface{}) {
	ts := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("🟥 [BLOCKCHAIN ERROR %s] %s\n", ts, msg)
}
