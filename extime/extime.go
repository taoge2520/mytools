package extime

import (
	"time"
)

//时间戳转时间字符串
func StampTotimes(timestamp int64) (times string) {
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05 ")
}

//时间字符串转时间戳
func TimesTostamp(times string) (stamp int64, err error) {
	tm, err := time.Parse("2006-01-02 15:04:05", times)
	if err != nil {
		return
	}
	stamp = tm.Unix()
	return
}
