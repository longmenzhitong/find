// Package order gathers all supported orders and implements concerned methods.
package order

import "strings"

const (
	Find   = "find"
	Add    = "add"
	Delete = "del"
	Modify = "mod"
	Exit   = "exit"
)

// orders is a string slice persist all of order.
var orders = []string{
	Find,
	Add,
	Delete,
	Modify,
	Exit,
}

// Order is used to parse order from user's input,
// returning the order which user want to execute.
func Order(input string) string {
	for _, order := range orders {
		if strings.HasPrefix(input, order) {
			return order
		}
	}
	return ""
}

// Param is used to parse param from user's input,
// returning the param which order execution need.
func Param(input string) string {
	return strings.TrimSpace(strings.TrimPrefix(input, Order(input)))
}

// Fast is used to check if user want to execute the order rapidly,
// returning check result and handled param.
func Fast(param string) (bool, string) {
	if strings.HasPrefix(param, "-f ") || strings.Contains(param, " -f ") {
		return true, strings.TrimSpace(strings.ReplaceAll(param, "-f", ""))
	}
	return false, param
}

// All is used to check if user want to influence all concerned notes,
// sometimes it means a fuzzy match of the key (e.g. del -a),
// returning check result and handled param.
func All(param string) (bool, string) {
	if strings.HasPrefix(param, "-a ") || strings.Contains(param, " -a ") {
		return true, strings.TrimSpace(strings.ReplaceAll(param, "-a", ""))
	}
	return false, param
}
