// Copyright (c) CattleCloud LLC
// SPDX-License-Identifier: BSD-3-Clause

// Package scope implements a simplified and more convenient API for creating
// context.Context values with often used functionality like timeouts.
package scope

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

// Deadline will create a fresh context not part of any preceding chain of
// values, and will expire at the given expiration time.
func Deadline(expiration time.Time) (C, Cancel) {
	return context.WithDeadline(New(), expiration)
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

// WithTTL wraps an existing Context that will expire after the given duration.
func WithTTL(c C, duration time.Duration) (C, Cancel) {
	return context.WithTimeout(c, duration)
}

// WithValue wraps an existing Context with value set for key.
func WithValue[K, V any](c C, key K, value V) C {
	return context.WithValue(c, key, value)
}

// Value retrieves the value associated with the given key.
func Value[K, V any](c C, key K) V {
	return c.Value(key).(V)
}
