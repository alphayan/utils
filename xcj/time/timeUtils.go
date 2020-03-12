package xcjTime

import (
	"log"
	"time"
)

//将时间转换为字符串格式
func FormatTimeToString(date time.Time) string {
	//获取时间戳
	timestamp := date.Unix()
	//格式化为字符串,tm为Time类型
	tm := time.Unix(timestamp, 0)
	//dateStr := tm.Format("2006-01-02 15:04:05 T")
	dateStr := tm.Format("2006-01-02T15:04:05Z07:00")
	return dateStr
}

func FormatTimeToLocalString(date time.Time) string {
	//获取时间戳
	timestamp := date.Unix()
	//格式化为字符串,tm为Time类型
	tm := time.Unix(timestamp, 0)
	//dateStr := tm.Format("2006-01-02 15:04:05 T")
	dateStr := tm.Format("2006-01-02 15:04:05")
	return dateStr
}

//将时间字符串转化为时间
func FormatTimeStrToTime(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02 03:04:05 PM", dateStr)
	if err != nil {
		log.Println("字符串转换为时间出错：", err)
	}
	return date
}
/**
   * @Description 算出两点时间差天数
   * @Author Acemon
   * @Date  2018/9/14 10:35 
   * @return 
*/
func DayDiffer(startTime,endTime time.Time) int64 {
	var day  int64
	if  startTime.Before(endTime) {
		differ := endTime.Unix() - startTime.Unix()
		day = differ/86400
		if (differ%86400) > 0 {
			day = day + 1
			return day
		}
	}
	return day
}