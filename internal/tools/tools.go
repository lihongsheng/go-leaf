package tools

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"html/template"
	"io"
	apiError "message-center/api/error"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// isTimeActive 判断时间是否在一个范围内
func isTimeActive(now, start, stop time.Time) bool {
	return (start.Before(now) || start.Equal(now)) && now.Before(stop)
}

func FileWithLineNum(skip int) (file string, line int) {
	_, f, l, _ := runtime.Caller(skip)
	return f, l
}

func FileWithLineNumToStr(skip int) string {
	f, l := FileWithLineNum(skip + 1)
	return fmt.Sprintf("File:%s;Line:%d", f, l)
}

// Equal 对比两个数据是否相等
func Equal(x, y any) bool {
	return reflect.DeepEqual(x, y)
}

func Unix() int64 {
	return time.Now().Unix()
}

func UnixToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

func DateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func Date(t time.Time) string {
	return t.Format("2006-01-02")
}

// Lock 加锁
func Lock(ctx context.Context, key string, lockTime time.Duration, redis *redis.Client) (success bool, unlock func(), err error) {
	value := time.Now().UnixMicro()
	val := strconv.FormatInt(value, 10)
	stm, err := redis.SetNX(ctx, key, val, lockTime).Result()
	if err != nil {
		return false, nil, apiError.ErrorClientRedisSetError("set redis lock error %v", err)
	}
	if stm == false {
		return false, nil, nil
	}
	return true, func() {
		c := context.Background()
		v := redis.Get(c, key)
		if v.Err() != nil {
			return
		}
		if v.String() != val {
			return
		}
		redis.Del(c, key)
	}, nil
}

func Md5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func TemplateParse(v interface{}, content string) (string, error) {
	t := template.New("tem")
	tpl, err := t.Parse(content)
	if err != nil {
		return "", apiError.ErrorSystemTemplateParseError("模板解析失败 %s, Template%s", err, "content")
	}
	byteStr := bytes.NewBufferString("")
	err = tpl.Execute(byteStr, v)
	if err != nil {
		return "", apiError.ErrorSystemTemplateParseError("模板解析失败 %s, var%s , Template%s", err, v, "content")
	}
	return byteStr.String(), nil
}

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func IsURI(fl string) bool {
	if !strings.Contains(fl, "http://") && !strings.Contains(fl, "https://") {
		fl = "http://" + fl
	}
	s := fl
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func IsEmail(email string) bool {
	emailRegexString := "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	emailRegex := regexp.MustCompile(emailRegexString)
	return emailRegex.MatchString(email)
}

func GenerateID() string {
	return uuid.New().String()
}

func GetReason(err error) string {
	reason := ""
	if err != nil {
		if e, ok := err.(*errors.Error); ok {
			reason = e.Reason
		} else {
			reason = apiError.ErrorReason_SYSTEM_UNKNOWN_ERROR.String()
		}
	}
	return reason
}
