// Copyright (c) CattleCloud LLC
// SPDX-License-Identifier: BSD-3-Clause

package scope

import (
	"testing"
)

func Test_New(t *testing.T) {
	t.Parallel()

	if c := New(); c == nil {
		t.Fail()
	}
}
