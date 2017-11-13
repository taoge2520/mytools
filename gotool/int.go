package gotool

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

//将数字转成逗号分隔形式便于阅读,如 1000 转成 1,000
func Comma(v int64) string {
	sign := ""
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}

	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

const (
	Byte = 1.0
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
	TB   = 1024 * GB
)

//将字节转成可读性好的格式 如 1024 转成 1.00 K
func ByteFormat(bytes uint64) string {
	value := float64(bytes)
	var unit string
	switch {
	case value == 0:
		return "0"
	case value < KB:
		value = value
		unit = "B"
	case value < MB:
		value = value / KB
		unit = "K"
	case value < GB:
		value = value / MB
		unit = "M"
	case value < TB:
		value = value / GB
		unit = "G"
	}
	return fmt.Sprintf("%.2f %v", value, unit)
}
