package git

import (
	"strconv"
	"testing"
)

func TestShellQuote(t *testing.T) {
	t.Parallel()
	for i, testcase := range []struct {
		input string
		want  string
	}{{
		input: "-a-b-c-",
		want:  "-a-b-c-",
	}, {
		input: "",
		want:  "''",
	}, {
		input: " ",
		want:  "' '",
	}, {
		input: " space ",
		want:  "' space '",
	}, {
		input: "abc@123",
		want:  "abc@123",
	}, {
		input: "abc$123",
		want:  "'abc$123'",
	}, {
		input: "abc'123",
		want:  `"abc'123"`,
	}, {
		input: "\n'NL' $X \n",
		want:  `"\x0a'NL' \$X \x0a"`,
	}} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Helper()
			if got := shellQuote(testcase.input); got != testcase.want {
				t.Errorf("shellQuote(%q) = %q; want %q", testcase.input, got, testcase.want)
			}
		})
	}
}
