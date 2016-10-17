package ini

import (
	"flag"
	"reflect"
	"testing"

	"github.com/conveyer/xflag/cflag/slices"
)

func TestConfigParse_NonExistentFile(t *testing.T) {
	if err := (&Config{}).New().AddFile("file_that_does_not_exist"); err == nil {
		t.Errorf("File does not exist, error expected.")
	}
}

func TestConfigParse_InvalidConfig(t *testing.T) {
	if err := (&Config{}).New().AddFile("./testdata/invalid.ini"); err == nil {
		t.Errorf("INI file is invalid, error expected.")
	}
}

func TestConfigAddFile(t *testing.T) {
	c := &Config{
		body: map[string]map[string]interface{}{},
	}
	err := c.AddFile("./testdata/config1.ini")
	assertNil(t, err)

	err = c.AddFile("./testdata/config2.ini")
	assertNil(t, err)

	err = c.AddFile("./testdata/config3.ini")
	assertNil(t, err)

	exp := map[string]map[string]interface{}{
		"": {
			"key1": "value1",
			"key2": "value3",
		},
		"some_section1": {
			"key1": "value1",
			"key2": "value2",
		},
		"some_section2": {
			"key1": "value3",
		},
		"some_section3": {
			"key1": "value2",
		},
	}
	if !reflect.DeepEqual(exp, c.body) {
		t.Errorf("Incorrectly parsed / joined config. Expected:\n%v.\nGot:\n%v.", exp, c.body)
	}
}

func TestConfigPrepare(t *testing.T) {
	c := &Config{
		body: map[string]map[string]interface{}{
			"": {
				"key1":  "value1",
				"arr[]": []string{"1", "2", "3"},
			},
			"paths": {
				"xxx": "${GOPATH} - ${GOPATH}",
			},
		},
	}
	for _, v := range []struct {
		Flag *flag.Flag
		Exp  string
	}{
		{&flag.Flag{Name: "key1", Value: &stringFlag{}}, "value1"},
		{&flag.Flag{Name: "paths:xxx", Value: &stringFlag{}}, "${GOPATH} - ${GOPATH}"},
		{&flag.Flag{Name: "non-existent-key", Value: &stringFlag{}}, ""},
		{&flag.Flag{Name: "non-existent-section:key", Value: &stringFlag{}}, ""},
		{&flag.Flag{Name: "arr[]", Value: &slices.Strings{}}, "[1; 2; 3]"},
	} {
		c.Prepare(v.Flag)
		if res := v.Flag.Value.String(); res != v.Exp {
			t.Errorf(
				`"%s": Expected "%s", got "%s".`,
				v.Flag.Name, v.Exp, res,
			)
		}
	}
}

func TestParseArgName(t *testing.T) {
	for k, vs := range map[string][]string{
		"key":               {"", "key"},
		"section:key":       {"section", "key"},
		"section:":          {"section", ""},
		":key":              {"", "key"},
		":":                 {"", ""},
		"":                  {"", ""},
		"::":                {"", ":"},
		"section:some:key:": {"section", "some:key:"},
	} {
		if sec, key := parseArgName(k); sec != vs[0] || key != vs[1] {
			t.Errorf(
				`Input "%s": Expected "%s", "%s"; got "%s", "%s".`,
				k, vs[0], vs[1], sec, key,
			)
		}
	}
}

func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`No error expected, got "%v".`, err)
	}
}

type stringFlag struct {
	d string
}

func (sf *stringFlag) String() string { return string(sf.d) }
func (sf *stringFlag) Set(v string) error {
	sf.d = v
	return nil
}