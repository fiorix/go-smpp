// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"fmt"
	"net/http"
	"strings"
)

// param defines a parameter from HTTP POST.
type param struct {
	Name     string
	Desc     string
	Required bool
	Options  []string
	Value    *string
}

// Valid checks whether the given value match the defined options for
// for the parameter, if options are defined.
func (p param) Valid(v string) bool {
	if p.Options == nil {
		return true
	}
	for _, opt := range p.Options {
		if v == opt {
			return true
		}
	}
	return false
}

// form is a list of parameters.
type form []param

// Validate validates each Param in this form, and store the parameter's
// value in the Value field.
func (f form) Validate(r *http.Request) error {
	for _, p := range f {
		v := r.FormValue(p.Name)
		if v == "" && p.Required {
			return fmt.Errorf("missing parameter %q: %s",
				p.Name, p.Desc)
		}
		if v != "" && !p.Valid(v) {
			return fmt.Errorf("invalid parameter %q=%q; supported: %s",
				p.Name, v, strings.Join(p.Options, ", "))
		}
		*p.Value = v
	}
	return nil
}
