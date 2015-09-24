// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"net/http"
	"strings"
)

// cors is an HTTP handler for managing cross-origin resource sharing.
// Ref: https://en.wikipedia.org/wiki/Cross-origin_resource_sharing.
func cors(f http.HandlerFunc, methods ...string) http.HandlerFunc {
	ms := strings.Join(methods, ", ") + ", OPTIONS"
	md := make(map[string]struct{})
	for _, method := range methods {
		md[method] = struct{}{}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		origin := "*"
		if len(r.Header.Get("Origin")) > 0 {
			origin = r.Header.Get("Origin")
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", ms)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if _, exists := md[r.Method]; exists {
			f.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Allow", ms)
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}
}
