package type1

import (
	"path/filepath"
	"testing"
)

func TestTrimSuffix(t *testing.T) {
	fn := "hello.txt"
	expected := "hello"
	if res := trimSuffix(fn); res != expected {
		t.Errorf("trimSuffix(%s) = %s, want %s", fn, res, expected)
	}
	fn = "hello"
	expected = "hello"
	if res := trimSuffix(fn); res != expected {
		t.Errorf("trimSuffix(%s) = %s, want %s", fn, res, expected)
	}
}

func TestLoadFont(t *testing.T) {
	t1font, err := LoadFont(filepath.Join("_testdata", "cmr10.pfb"), "")
	if err != nil {
		t.Error(err)
	}
	expected := "Computer"
	if t1font.FamilyName != expected {
		t.Errorf("t1font.FamilyName got %s, want %s", t1font.FamilyName, expected)
	}
}
