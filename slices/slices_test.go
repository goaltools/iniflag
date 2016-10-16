package slices

import (
	"reflect"
	"testing"
)

func TestBaseInit(t *testing.T) {
	msg := `Initialized must return "%v", got "%v".`

	e := false
	b := &base{}
	ok := b.Initialized()
	if ok {
		t.Errorf(msg, e, ok)
	}

	e = false
	b.RequireInit(e)
	ok = b.Initialized()
	if !ok {
		t.Errorf(msg, !e, !ok)
	}

	e = true
	b.RequireInit(e)
	ok = b.Initialized()
	if ok {
		t.Errorf(msg, e, ok)
	}
}

func TestStr(t *testing.T) {
	for exp, inp := range map[string]*test{
		"[]":        &test{},
		"[a]":       &test{d: []string{"a"}},
		"[a; b; c]": &test{d: []string{"a", "b", "c"}},
	} {
		if res := str(inp); res != exp {
			t.Errorf(`"%v": Expected "%s", got "%s".`, inp, exp, res)
		}
	}
}

func TestSet(t *testing.T) {
	for _, v := range []struct {
		fn  func(*testing.T, *test)
		exp []string
	}{
		{
			func(t *testing.T, st *test) {
			},
			nil,
		},
		{
			func(t *testing.T, st *test) {
				set(st, "a")
			},
			[]string{"a"},
		},
		{
			func(t *testing.T, st *test) {
				set(st, "a")
				set(st, "b")
				set(st, "c")
			},
			[]string{"a", "b", "c"},
		},
		{
			func(t *testing.T, st *test) {
				set(st, "a")
				set(st, "b")
				set(st, "c")

				set(st, EOI)

				set(st, "x")
				set(st, "y")
				set(st, "z")
			},
			[]string{"x", "y", "z"},
		},
	} {
		st := &test{}
		v.fn(t, st)
		if !reflect.DeepEqual(st.d, v.exp) {
			t.Errorf("Incorrect slice values. Expected:\n`%#v`.\nGot:\n`%#v`.", v.exp, st.d)
		}
	}
}

//
// Test object that implements a slice interface is below.
//

type test struct {
	base
	d []string
}

func (t *test) Len() int {
	return len(t.d)
}

func (t *test) Get(i int) string {
	return t.d[i]
}

func (t *test) Alloc() {
	t.d = []string{}
}

func (t *test) Add(v string) {
	t.d = append(t.d, v)
}