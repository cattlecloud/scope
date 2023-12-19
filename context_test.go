// Copyright (c) NOXIDE.LOL
// SPDX-License-Identifier: BSD-3-Clause

package xtc

import (
	"testing"
)

func Test_New(t *testing.T) {
	if c := New(); c == nil {
		t.Fail()
	}
}
