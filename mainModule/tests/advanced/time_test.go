package advanced

import (
	"fmt"
	"mainModule/collects"
	"mainModule/tests"
	"testing"
	"time"
)

// NOTE GO语言中的访问控制
// NOTE 所有语法可见性均定义在package这个级别
func TestPrivateOrPublic(t *testing.T) {
	collects.Cao()
	// NOTE 只能访问到结构体的公有的属性
	se := tests.Second{}
	fmt.Println(se.A)
	m := More{A: 12, b: false}
	println(m.A, m.b)
}

func TestTimeZone(t *testing.T) {
	// 可以在 $GOROOT/lib/time/zoneinfo.zip 文件下看到所有时区。
	// 默认UTC
	loc, _ := time.LoadLocation("")
	fmt.Println(loc)
	// 服务器设定的时区，一般为CST
	loc, _ = time.LoadLocation("Local")
	fmt.Println(loc)
	// 美国洛杉矶PDT
	loc, _ = time.LoadLocation("America/Los_Angeles")
	fmt.Println(loc)
	// 获取指定时区的时间点
	local, _ := time.LoadLocation("America/Los_Angeles")
	fmt.Println(time.Date(2018, 1, 1, 12, 0, 0, 0, local))
}

func TestTime(t *testing.T) {
	// time.Time
	// time.Duration
	// time.Location
	// time.Timer
	// time.Ticker

	/*
		世界协调时间(UTC) 中国标准时间(CST)
		galang 中指定的特定时间格式为 "2006-01-02 15:04:05 -0700 MST"，
		为了记忆方便，按照美式时间格式 月日时分秒年 外加时区 排列起来依次是 01/02 03:04:05PM ‘06 -0700
	*/
	println("-------------------------------------------------------------------------------------")
	asTimezone, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(asTimezone)

	layout := "2006-01-02 15:04:05"
	var tm = time.Now()
	fmt.Println("Local：", time.Local, "Now：", tm,
		"Location：", tm.Location(), "Location：", tm.Local(), time.Local == tm.Location())

	// NOTE Time Parse
	tm, err := time.Parse(layout, "2011-04-23 12:24:51") // use UTC
	if err == nil {
		fmt.Println("Default Parse：", tm)
	} else {
		fmt.Println(err)
	}
	tm, _ = time.ParseInLocation(layout, "2021-08-27 16:30:05", asTimezone)
	fmt.Println("Asia/Shanghai Parse：", tm)
	tm, _ = time.ParseInLocation(layout, "2021-08-27 16:30:05", time.Local)
	fmt.Println("Local Parse：", tm)

	// NOTE timestamp to datetime
	tm = time.Unix(1630053284, 276208) // Use time.Local
	fmt.Println("Unix From Local", tm)

	// NOTE Date
	date := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(),
		tm.Minute(), tm.Second(), tm.Nanosecond(), tm.Location())
	fmt.Println(date)

	loc, _ := time.LoadLocation("America/Los_Angeles")
	amerTime := date.In(loc)
	fmt.Println(amerTime)

	// NOTE format
	formatStr := time.Now().Format("2006/01/02 03:04:05.000000000000 MST")
	fmt.Println(formatStr, " | ", time.Now().Format(time.UnixDate))

	// NOTE datetimne to timestamp
	fmt.Println(tm.Unix(), " | ", tm.UnixNano()) // 纳秒
	println("-------------------------------------------------------------------------------------")
	y, m, d := tm.Date()
	fmt.Println("日期", y, m.String(), d)
	h, M, s := tm.Clock()
	fmt.Println("时钟", h, M, s)
	fmt.Println("年", tm.Year())
	fmt.Println("月", tm.Month())
	fmt.Println("日", tm.Day())
	fmt.Println("时", tm.Hour())
	fmt.Println("分", tm.Minute())
	fmt.Println("秒", tm.Second())
	fmt.Println("纳秒", tm.Nanosecond())
	fmt.Println("周", tm.Weekday().String())
	y, w := tm.ISOWeek()
	fmt.Println(y, "年中的第", w, "周")
	fmt.Println("一年中的第", tm.YearDay(), "天")
	fmt.Println("转换为当前设置的时区时间", tm.Local())
	fmt.Println("当前设置的时区", tm.Location())
	name, offset := tm.Zone()
	fmt.Println("时区", name, offset)
	fmt.Println("零", tm.IsZero())
	fmt.Println("无纳秒时间戳", tm.Unix())
	println("-------------------------------------------------------------------------------------")
}

func TestDuration(t *testing.T) {
	// yd, _ := time.ParseDuration("1.5y")
	// md, _ := time.ParseDuration("1.5m")
	// dd, _ := time.ParseDuration("1.5d")

	hd, _ := time.ParseDuration("1.5h")
	md, _ := time.ParseDuration("1.5m")
	sd, _ := time.ParseDuration("1.5s")
	mid, _ := time.ParseDuration("1.5ms")
	vd, _ := time.ParseDuration("1.5us")
	nd, _ := time.ParseDuration("1.5ns")

	fmt.Println(hd, md, sd, mid, vd, nd)
	fmt.Println(hd.Hours(), hd.Minutes(), hd.Seconds(),
		hd.Microseconds(), hd.Milliseconds(), hd.Nanoseconds(), hd.String())
	var dur time.Duration = 1 * 60 * 60 * 1000 * 1000 * 1000
	fmt.Println(dur.Hours())
	fmt.Println("四舍五入：", hd.Round(dur))
	fmt.Println("向下取整：", hd.Truncate(dur))
}

// 不走缓存 -count=1
func TestTimeOperate(t *testing.T) {
	// type()  type-conversion
	// NOTE Sleep
	d := time.Duration(2) * time.Second
	// fmt.Println("第一个 sleep 2秒")
	// time.Sleep(d)
	// d = 3 * time.Second
	// fmt.Println("第二个 sleep 2秒")
	// time.Sleep(d)

	// TODO After
	fmt.Println("第一个 after 2秒")
	time.After(d)
	fmt.Println("----------------------------------")

	// NOTE Since | Until
	start := time.Now()
	fmt.Println(time.Since(start))
	fmt.Println(time.Until(start))
	fmt.Println(time.Until(time.Now()))

	// NOTE Add | Sub
	dayDur := 24 * time.Hour
	now := time.Now()
	yesterday := now.Add(-dayDur)
	tomorrow := now.Add(dayDur)
	fmt.Println(dayDur, "\n", now, "\n", yesterday, "\n", tomorrow)

	// NOTE Go 中有很少的操作符重载
	dayDur = tomorrow.Sub(yesterday)
	fmt.Println(dayDur)
	now = yesterday.AddDate(1, 1, 0)
	fmt.Println(now)
	fmt.Println(tomorrow.After(yesterday), yesterday.Before(tomorrow), !yesterday.Equal(tomorrow))

	// NOTE 时区转换
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	fmt.Println(time.Date(2018, 1, 2, 0, 0, 0, 0, beijing))
	// 当前时间转为指定时区时间
	fmt.Println(time.Now().In(beijing))

	// 指定时间转换成指定时区对应的时间
	dt, err := time.ParseInLocation("2006-01-02 15:04:05", "2017-05-11 14:06:06", time.Local)
	fmt.Println(dt, err)
	// 当前时间在零时区
	now = time.Now()
	utcNow := now.UTC()
	zoneName, offset := now.Zone()
	fmt.Println(now, utcNow, zoneName, time.Duration(offset)*time.Second)

}

// TODO
func TestSelectAfter(t *testing.T) {
	// LATER select 和 switch ??

	// select {
	// case m := <-c:
	// 	do something
	// case <-time.After(time.Duration(1) * time.Second):
	// 	fmt.Println("time out")
	// }
}

// TODO
func TestTicker(t *testing.T) {
	// Ticker 类型包含一个 channel，每隔一段时间执行(比如设置心跳时间等)
	// 无法取消
	tick := time.Tick(1 * time.Minute)
	for range tick {
		// do something
	}

	// 可通过调用ticker.Stop取消
	ticker := time.NewTicker(1 * time.Minute)
	// for _ = range tick {
	// do something
	// }
	ticker.Stop()
}

// TODO
func TestTimer(t *testing.T) {
	// Timer 类型用来代表一个单独的事件，当设置的时间过期后，
	// 发送当前的时间到 channel, 我们可以通过以下两种方式来创建
	// func AfterFunc(d Duration, f func()) *Timer
	// func NewTimer(d Duration) *Timer

	// 以上两函数都可以使用 Reset, 需要注意的地方是使用 Reset 时需要确保 t.C 通道被释放时才能调用，
	// 以防止发生资源竞争的问题，可通过以下方式解决

	// if !t.Stop() {
	// 	<-t.C
	// }
	// t.Reset(d)

}
