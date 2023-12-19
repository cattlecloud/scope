// Copyright (c) NOXIDE.LOL
// SPDX-License-Identifier: BSD-3-Clause

// Package xtc implements a simplified and more convenient API for creating context.Context values with often used
// functionality like timeouts.
package xtc

import (
	"context"
	"time"
)

// C is an alias for context.Context
type C = context.Context

// Cancel is an alias for context.CancelFunc
type Cancel = context.CancelFunc

// New will create a fresh context not part of any preceding chain of values.
func New() C {
	return context.Background()
}

// TTL will create a fresh context not part of any preceding chain of values,
// and will expire after the given duration.
func TTL(duration time.Duration) (C, Cancel) {
	return context.WithTimeout(New(), duration)
}

// Cancelable will create a fresh context not part of any preceding chain of
// values, and includes a Cancel function.
func Cancelable() (C, Cancel) {
	return context.WithCancel(New())
}

// WithCancel wraps an existing Context with a Cancel function.
func WithCancel(c C) (C, Cancel) {
	return context.WithCancel(c)
}
