// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"fmt"
	"net/url"
	"strings"
)

// param defines a parameter from HTTP POST.
type param struct {
	Name     string
	Desc     string
	Required bool
	Options  []string // possible options for this param, optional
	Value    *string  // where the value[name] will be stored, optional
}

// ValidOption checks whether the given value match the defined
// options for the parameter, if options are provided.
func (p param) ValidOption(v string) bool {
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

// Validate validates each parameter in this form, and store the
// parameter's value in the Value field, if provided.
func (f form) Validate(values url.Values) error {
	for _, p := range f {
		v := values.Get(p.Name)
		if v == "" && p.Required {
			return fmt.Errorf("missing parameter %q: %s",
				p.Name, p.Desc)
		}
		if v != "" && !p.ValidOption(v) {
			return fmt.Errorf("invalid parameter %q=%q; supported: %s",
				p.Name, v, strings.Join(p.Options, ", "))
		}
		if p.Value != nil {
			*p.Value = v
		}
	}
	return nil
}
