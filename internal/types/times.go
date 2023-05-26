package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 可以定义数据库里 date time类型字段为 LocalDateTime，因为实现了 gorm的val和scan接口会做自动转换
// 也因为实现了 JSON在序列化和反序列化的UnmarshalJSON和MarshalJSON接口，也会做自动转换为

const DateTimeFormat = "2006-01-02 15:04:05"
const DateFormat = "2006-01-02"

type LocalDateTime time.Time

func (t *LocalDateTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalDateTime(time.Time{})
		return
	}
	str := string(data)
	str = strings.Trim(str, "\"")
	now, err := time.Parse(DateTimeFormat, string(data))
	*t = LocalDateTime(now)
	return
}

func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.String())
	return []byte(output), nil
}

func (t LocalDateTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(t.String()), nil
}

func (t *LocalDateTime) Scan(v interface{}) error {
	tTime, _ := time.Parse(DateTimeFormat+" +0800 CST", v.(time.Time).String())
	*t = LocalDateTime(tTime)
	return nil
}

func (t LocalDateTime) String() string {
	return time.Time(t).Format(DateTimeFormat)
}

// LocalDate 0000-00-00 格式
type LocalDate time.Time

func (t *LocalDate) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalDate(time.Time{})
		return
	}
	str := string(data)
	str = strings.Trim(str, "\"")
	if len(str) == 10 {
		str = str + " 00:00:00"
	}
	now, err := time.Parse(DateTimeFormat, str)
	*t = LocalDate(now)
	return
}

func (t LocalDate) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.String())
	return []byte(output), nil
}

func (t LocalDate) Value() (driver.Value, error) {
	if t.String() == "0001-01-01" {
		return nil, nil
	}
	return []byte(t.String()), nil
}

func (t *LocalDate) Scan(v interface{}) error {
	tTime, _ := time.Parse(DateFormat+" +0800 CST", v.(time.Time).String())
	*t = LocalDate(tTime)
	return nil
}

func (t LocalDate) String() string {
	str := time.Time(t).Format(DateFormat)
	return str
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, "\"")
	duration, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	duration := int64(d)
	t := time.Duration(duration)
	output := strconv.FormatInt(int64(t.Seconds()), 10) + "s"
	return []byte(output), nil
}
