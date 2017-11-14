package formattime

import (
	"fmt"
	"strings"
	//	"time"
)

var (
	month map[string]string
)

func init() {
	month = make(map[string]string)
	month["jan"] = "01"
	month["feb"] = "02"
	month["mar"] = "03"
	month["apr"] = "04"
	month["may"] = "05"
	month["jun"] = "06"
	month["jul"] = "07"
	month["aug"] = "08"
	month["sep"] = "09"
	month["oct"] = "10"
	month["nov"] = "11"
	month["dec"] = "12"

	month["january"] = "01"
	month["february"] = "02"
	month["march"] = "03"
	month["april"] = "04"
	month["may"] = "05"
	month["june"] = "06"
	month["july"] = "07"
	month["august"] = "08"
	month["september"] = "09"
	month["october"] = "10"
	month["november"] = "11"
	month["december"] = "12"
}
func Usefunc() {
	tt, err := GetTime("june 18 2017")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tt)
}

func GetTime(t string) (times string, err error) {
	//一类
	if strings.Contains(t, "-") {
		//2018-02-02t11:01:44z
		//2017-11-18t14:08:09.219z
		if strings.HasSuffix(t, "z") {
			if strings.Contains(t, "t") {
				h := strings.Split(t, "t")
				if len(h) == 2 {
					h[1] = h[1][:8]
					times = strings.Join(h, " ")
					return
				}

			}
		}
		//2018-07-11 09:49:03 clt
		if strings.HasSuffix(t, "clt") {
			t = strings.TrimRight(t, " clt")
			times = t
			return
		}

		//2017-01-22
		if len(t) == 10 {
			times = t
			times += " 00:00:00"
			return

		}
		//2018-03-30t13:06:18+13:00
		if strings.Index(t, "t") == 10 {
			h := strings.Split(t, "t")
			if len(h) == 2 {
				h[1] = h[1][:8]
				times = strings.Join(h, " ")
				return
			}

		}

		//15-aug-2018 23:59:59 utc
		if strings.HasSuffix(t, "utc") {
			t = strings.TrimRight(t, " utc")
			h := strings.Split(t, " ")
			if len(h) == 2 {
				date := strings.Split(h[0], "-")
				if len(date) == 3 {
					date[1] = month[date[1]]
					date[0], date[2] = date[2], date[0]
					times = strings.Join(date, "-")
					times = times + " " + h[1]
					return
				}
			}

		}

		//2018-04-28 00:00:00
		h := strings.Split(t, " ")
		if len(h) == 2 {
			times = t
			return
		}

		//19-april-2018
		date := strings.Split(t, "-")
		if len(date) == 3 {
			date[1] = month[date[1]]
			date[0], date[2] = date[2], date[0]
			times = strings.Join(date, "-")
			times = times + " " + "00:00:00"
			return
		}

	}

	//二类
	//06/07/2016
	if strings.Contains(t, "/") {
		s := strings.Split(t, "/")
		if len(s) == 3 {
			s[0], s[2] = s[2], s[0]
			times = strings.Join(s, "-")
			times += " 00:00:00"
			return
		}
	}

	//三类
	if strings.Contains(t, ".") {
		//2020. 02. 29.
		h := strings.Split(t, ".")
		if len(h) == 4 {
			h[1] = strings.Trim(h[1], " ")
			h[2] = strings.Trim(h[2], " ")

			h = append(h[:3], h[4:]...)
			times = strings.Join(h, "-")
			times += " 00:00:00"
			return
		}
		if len(h) == 3 {
			//24.04.2018 08:06:04
			if len(h[2]) > 4 {
				ytime := strings.Split(h[2], " ")
				if len(ytime) == 2 {
					times = ytime[0] + "-" + h[1] + "-" + h[0] + " " + ytime[1]
					return
				}
			}
			h[0], h[2] = h[2], h[0]
			times = strings.Join(h, "-") + " 00:00:00"
			return
		}
	}
	//四类
	if strings.Contains(t, "gmt") {
		h := strings.Split(t, " ")
		if len(h) == 6 {
			h[1] = month[h[1]]
			times = h[5] + "-" + h[1] + "-" + h[2] + " " + h[3]
			return
		}
	}
	h := strings.Split(t, " ")
	if len(h) == 3 {
		if _, ok := month[h[0]]; ok {
			h[0] = month[h[0]]
			times = h[2] + "-" + h[0] + "-" + h[1] + " 00:00:00"
			return
		}
		if _, ok := month[h[1]]; ok {
			h[1] = month[h[1]]
			times = h[2] + "-" + h[1] + "-" + h[0] + " 00:00:00"
			return
		}
	}
	return
}
