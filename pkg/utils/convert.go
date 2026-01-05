package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ToString 将任意类型转换为字符串
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return strconv.FormatBool(v)
	default:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}
}

// ToInt64 将字符串转换为int64
func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ToInt64OrDefault 将字符串转换为int64，失败返回默认值
func ToInt64OrDefault(s string, defaultValue int64) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToInt 将字符串转换为int
func ToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ToIntOrDefault 将字符串转换为int，失败返回默认值
func ToIntOrDefault(s string, defaultValue int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToFloat64 将字符串转换为float64
func ToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ToFloat64OrDefault 将字符串转换为float64，失败返回默认值
func ToFloat64OrDefault(s string, defaultValue float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToBool 将字符串转换为bool
func ToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ToBoolOrDefault 将字符串转换为bool，失败返回默认值
func ToBoolOrDefault(s string, defaultValue bool) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return v
}

// Int64ToString 将int64转换为字符串
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// IntToString 将int转换为字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// SplitToInt64Slice 将逗号分隔的字符串转换为int64切片
func SplitToInt64Slice(s string, sep string) ([]int64, error) {
	if s == "" {
		return []int64{}, nil
	}

	parts := strings.Split(s, sep)
	result := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		num, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}

// SplitToIntSlice 将逗号分隔的字符串转换为int切片
func SplitToIntSlice(s string, sep string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}

	parts := strings.Split(s, sep)
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}

// Int64SliceToString 将int64切片转换为逗号分隔的字符串
func Int64SliceToString(slice []int64, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	strs := make([]string, len(slice))
	for i, v := range slice {
		strs[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(strs, sep)
}

// InterfaceToInt64 将interface{}转换为int64
func InterfaceToInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int64:
		return val, nil
	case int:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

// StructToMap 将结构体转换为map
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// MapToStruct 将map转换为结构体
func MapToStruct(data map[string]interface{}, result interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, result)
}
