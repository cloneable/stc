package git

import "testing"

func TestParseShowRefLine(t *testing.T) {
	ref, err := ParseRef("0ee9e9fa90a6d36494576c1e750ddad5e176e0be refs/heads/master")
	if err != nil {
		t.Errorf("parseShowRefLine: %v", err)
	}
	if got, want := ref.name, RefName("refs/heads/master"); got != want {
		t.Errorf("ref.Name = %s, want %s", got, want)
	}
	if got, want := ref.objectName, ObjectName("0ee9e9fa90a6d36494576c1e750ddad5e176e0be"); got != want {
		t.Errorf("ref.Commit = %s, want %s", got, want)
	}

	// TODO: test edge/failure cases
	// TODO: fuzz
}