package rule

import (
	"errors"
	"testing"
)

func TestRegex(t *testing.T) {
	tests := []*struct {
		params   []string
		value    string
		path     string
		expected bool
		err      error
	}{
		{params: []string{".+"}, value: "regex", path: "", expected: true, err: nil},
		{params: []string{"[0-9]+"}, value: "123", path: "", expected: true, err: nil},
		{params: []string{"![0-9]+"}, value: "123", path: "", expected: false, err: nil},
		{params: []string{"[a-z]+"}, value: "123", path: "", expected: false, err: nil},
		{params: []string{"[a-z\\-]+"}, value: "google-test", path: "", expected: true, err: nil},
		{params: []string{"[a-z\\-]+"}, value: "google.test", path: "", expected: false, err: nil},
		{params: []string{"${1}_${0}"}, value: "google_test", path: "google/test", expected: true, err: nil},
		{params: []string{"${1}"}, value: "swu1", path: "gen/swu1/data", expected: true, err: nil}, // github.com/loeffel-io/ls-lint/issues/307
		{params: []string{"${1}_${0}"}, value: "test", path: "google/test", expected: false, err: nil},
	}

	i := 0
	for _, test := range tests {
		rule := new(Regex).Init()

		// parameters
		err := rule.SetParameters(test.params)

		if !errors.Is(err, test.err) {
			t.Errorf("Test %d failed with unmatched error - %e", i, err)
			return
		}

		// validate
		res, err := rule.Validate(test.value, test.path, true)

		if err != nil && err != test.err {
			t.Errorf("Test %d failed with unmatched error - %s", i, err.Error())
			return
		}

		if res != test.expected {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}
