package reflect

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

func GetBool(raw any) bool {
	input, isValid := GetValue(raw)
	if !isValid {
		return false
	}
	var boolValue = false
	switch input := input.(type) {
	case bool:
		boolValue = input
	case string:
		var err error
		boolValue, err = strconv.ParseBool(strings.ToUpper(input))
		if err != nil {
			logger.WarnErr("string parse bool failed", err)
		}
	case int, int8, int16, int32, int64:
		intValue := GetIntBase(input, 64)
		switch intValue {
		case 0:
			boolValue = false
		case 1:
			boolValue = true
		default:
			boolValue = false
		}
	case uint, uint8, uint16, uint32, uint64:
		uintValue := GetUintBase(input, 64)
		switch uintValue {
		case 0:
			boolValue = false
		case 1:
			boolValue = true
		default:
			boolValue = false
		}
	default:
		logger.Warn("conversion for bool type failed", zap.Any("type", reflect.TypeOf(input)), zap.Any("value", input))
	}
	return boolValue
}
