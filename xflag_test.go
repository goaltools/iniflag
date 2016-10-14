package xflag

import (
	"flag"
	"go/build"
	"os"
	"testing"
)

var (
	f1 = flag.String("key1", "value1_def", "flag from default section, file 1")
	f2 = flag.String("section:key1", "value2_def", "flag from `section`, file 2")
	f3 = flag.String("arg", "value_def", "flag from arguments")
	f4 = flag.String("doesNotExist", "default", "flag without input")
)

func TestParse(t *testing.T) {
	// Simulating "--arg value" input arguments.
	os.Args = []string{os.Args[0], "--arg", "value"}

	err := Parse("./testdata/file1.ini", "./testdata/file2.ini")
	if err != nil {
		t.Errorf(`No error expected, got "%v".`, err)
	}

	if *f1 != build.Default.GOPATH || *f2 != "value2" || *f3 != "value" || *f4 != "default" {
		t.Errorf("Incorrect value of one of the flags.")
	}
}

func TestParse_IncorrectFile(t *testing.T) {
	err := Parse("file_does_not_exist")
	if err == nil {
		t.Errorf("File does not exist, error expected.")
	}
}
