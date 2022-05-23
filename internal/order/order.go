// Package order gathers all supported orders and implements concerned methods.
package order

import "strings"

const (
	Find       = "find"
	Add        = "add"
	Delete     = "del"
	FastDelete = "fdel"
	Modify     = "mod"
	Exit       = "exit"
)

// orders is a string slice persist all of order.
var orders = []string{
	Find,
	Add,
	Delete,
	FastDelete,
	Modify,
	Exit,
}

// Parse is used to parse order from user's input.
func Parse(input string) string {
	for _, order := range orders {
		if strings.HasPrefix(input, order) {
			return order
		}
	}
	return ""
}
