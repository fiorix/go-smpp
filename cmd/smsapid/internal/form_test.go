// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"net/http"
	"net/url"
	"testing"
)

func TestForm(t *testing.T) {
	var a, b, c, d string
	form := Form{
		{"a", "param a", true, nil, &a},
		{"b", "param b", false, nil, &b},
		{"c", "param c", true, []string{"1", "2"}, &c},
		{"d", "param d", false, []string{"3", "4"}, &d},
	}
	test := []struct {
		Desc string
		Form url.Values
		Good bool
	}{
		{"missing a, d", url.Values{"b": {"."}, "c": {"."}}, false},
		{"missing b, d", url.Values{"a": {"."}, "c": {"1"}}, true},
		{"missing c, d", url.Values{"a": {"."}, "b": {"3"}}, false},
		{"missing b", url.Values{"a": {"."}, "c": {"2"}, "d": {"4"}}, true},
		{"invalid c", url.Values{"a": {"."}, "c": {"."}}, false},
	}
	for _, el := range test {
		r := &http.Request{Form: el.Form}
		err := form.Validate(r)
		if el.Good && err != nil {
			t.Fatal(err)
		}
		if !el.Good && err == nil {
			t.Fatal("bad form parsed without error:", el.Desc)
		}
	}
}
