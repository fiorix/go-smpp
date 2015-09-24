// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import "net/http"

// auth provides an authentication handler. Hook up your own stuff here.
func auth(f http.HandlerFunc) http.HandlerFunc {
	return f
}
