package fakegit

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func Run(dir string, args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitCode int, err error) {
	fmt.Fprintf(os.Stderr, "FAKEGIT: git %s", strings.Join(args, " "))
	return bytes.Buffer{}, bytes.Buffer{}, 0, nil
}
