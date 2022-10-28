package util

import (
	"fmt"
	"time"
)

// CheckTime true:didn't expire;false:expired
func CheckTime(t *time.Time, expire int) (b bool) {
	fmt.Println("checkTime...")
	defer fmt.Println("checkTime finish")
	now := time.Now()
	newT := t.Add(time.Second * time.Duration(expire))
	if newT.After(now) {
		b = false
	} else {
		b = true
	}
	return
}
