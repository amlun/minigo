package utils

import (
	"errors"
	"time"
)

// TimeRange 时间范围结构
type TimeRange struct {
	Start time.Time `form:"start" binding:"omitempty"` // 开始时间
	End   time.Time `form:"end" binding:"omitempty"`   // 结束时间
}

// Validate 验证时间范围的有效性
func (tr *TimeRange) Validate() error {
	if tr.Start.After(tr.End) {
		return errors.New("start time cannot be later than end time")
	}
	return nil
}

// IsEmpty 检查时间范围是否为空（零值）
func (tr *TimeRange) IsEmpty() bool {
	return tr.Start.IsZero() || tr.End.IsZero()
}

// ParseTimeRange 解析时间范围
func ParseTimeRange(startStr, endStr string) (TimeRange, error) {
	startTime, err := ParseDate(startStr)
	if err != nil {
		return TimeRange{}, err
	}
	endTime, err := ParseDate(endStr)
	if err != nil {
		return TimeRange{}, err
	}
	// 默认是23:59:59
	endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	return TimeRange{Start: startTime, End: endTime}, nil
}

func ParseDate(str string) (time.Time, error) {
	// Ignore null, like in the main JSON package.
	if str == "null" || str == "" {
		return time.Time{}, nil
	}
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return time.Time{}, err
	}
	return tt, nil
}
