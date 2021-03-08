package check

import (
	"net"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var numericOnlyRegexp = regexp.MustCompile(`^\d+$`)

func IsNumeric(str string) bool {
	if len(str) == 0 {
		return false
	}
	return numericOnlyRegexp.MatchString(str)
}

func IsStartWithWx(str string) bool {
	return strings.ToLower(str[0:2]) == "wx"
}

func IsStartWithQQ(str string) bool {
	return strings.ToLower(str[0:2]) == "qq"
}

func IsStartWithAli(str string) bool {
	return strings.ToLower(str[0:3]) == "ali"
}

// GoContainStrArr — Checks if a value exists in an array
func GoContainStrArr(element string, arrs []string, caseSensitive bool) bool {
	if len(arrs) == 0 {
		return false
	}
	for _, el := range arrs {
		if !caseSensitive {
			el = strings.ToLower(el)
			element = strings.ToLower(element)
		}
		if strings.Contains(element, el) {
			return true
		}
	}
	return false
}

func GoInIntArr(element int, arrs []int) bool {
	if len(arrs) == 0 {
		return false
	}
	for _, el := range arrs {
		if element == el {
			return true
		}
	}
	return false
}

func GoInInt32Arr(element int32, arrs []int32) bool {
	if len(arrs) == 0 {
		return false
	}
	for _, el := range arrs {
		if element == el {
			return true
		}
	}
	return false
}

func GoInInt64Arr(element int64, arrs []int64) bool {
	if len(arrs) == 0 {
		return false
	}
	for _, el := range arrs {
		if element == el {
			return true
		}
	}
	return false
}

// CheckIPType 0:invalid,1:ipv4,2:ipv6
func CheckIPType(ip string) int {
	if net.ParseIP(ip) == nil {
		return 0
	}
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return 1
		case ':':
			return 2
		}
	}
	return 0
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// Reflect if an interface is either a struct or a pointer to a struct
// and has the defined member method. If error is nil, it means
// the MethodName is accessible with reflect.
func IsStructMethodExist(structObject interface{}, MethodName string) bool {
	ValueInterface := reflect.ValueOf(structObject)
	// Check if the passed interface is a pointer
	if ValueInterface.Type().Kind() != reflect.Ptr {
		// Create a new type of structObject, so we have a pointer to work with
		ValueInterface = reflect.New(reflect.TypeOf(structObject))
	}
	// Get the method by name
	Method := ValueInterface.MethodByName(MethodName)
	if !Method.IsValid() {
		return false
	}
	return true
}

// Reflect if an interface is either a struct or a pointer to a struct
// and has the defined member field, if error is nil, the given
// FieldName exists and is accessible with reflect.
func IsStructFieldExist(structObject interface{}, FieldName string) bool {
	ValueInterface := reflect.ValueOf(structObject)
	// Check if the passed interface is a pointer
	if ValueInterface.Type().Kind() != reflect.Ptr {
		// Create a new type of structObject's Type, so we have a pointer to work with
		ValueInterface = reflect.New(reflect.TypeOf(structObject))
	}
	// 'dereference' with Elem() and get the field by name
	Field := ValueInterface.Elem().FieldByName(FieldName)
	if !Field.IsValid() {
		return false
	}
	return true
}
