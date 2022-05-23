// package stdin implements methods for handling standard input.
package stdin

import (
	"bufio"
	"os"
	"strings"
)

// ReadString is used to get user's input,
// returning space-trimmed string and error.
func ReadString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	} else {
		return strings.TrimSpace(input), nil
	}
}
