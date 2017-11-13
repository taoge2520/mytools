package gotool

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIPv4toUint32(t *testing.T) {
	Convey("IP转为uint32", t, func() {

		Convey("正确的IP格式", func() {
			u, err := IPv4toUint32("0.0.0.1")
			So(u, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})
		Convey("不正确的IP格式", func() {
			u, err := IPv4toUint32("abc")
			So(u, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}
func TestUint32toIPv4(t *testing.T) {
	Convey("uint32转为IP格式", t, func() {
		Convey("uint32转IP格式(string)", func() {
			So(Uint32toIPv4(uint32(0)), ShouldEqual, "0.0.0.0")
		})
	})

}
