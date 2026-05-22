// Copyright (c) CattleCloud LLC
// SPDX-License-Identifier: BSD-3-Clause

package scope

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestJoin_CancelA(t *testing.T) {
	t.Parallel()

	ctxA, cancelA := Cancelable()
	ctxB := New()

	j, _ := Join(ctxA, ctxB)
	cancelA()

	<-j.Done()
	if !errors.Is(j.Err(), context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", j.Err())
	}
}

func TestJoin_CancelB(t *testing.T) {
	t.Parallel()

	ctxA := New()
	ctxB, cancelB := Cancelable()

	j, _ := Join(ctxA, ctxB)
	cancelB()

	<-j.Done()
	if !errors.Is(j.Err(), context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", j.Err())
	}
}

func TestJoin_CancelSelf(t *testing.T) {
	t.Parallel()

	ctxA := New()
	ctxB := New()

	j, cancel := Join(ctxA, ctxB)
	cancel()

	<-j.Done()
	if !errors.Is(j.Err(), context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", j.Err())
	}
}

func TestJoin_AlreadyDone(t *testing.T) {
	t.Parallel()

	t.Run("a already done", func(t *testing.T) {
		ctxA, cancelA := Cancelable()
		cancelA()
		ctxB := New()

		j1, _ := Join(ctxA, ctxB)
		select {
		case <-j1.Done():
		default:
			t.Fatal("expected joined context to be done immediately when A is already done")
		}
		if !errors.Is(j1.Err(), context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", j1.Err())
		}
	})

	t.Run("b already done", func(t *testing.T) {
		ctxC := New()
		ctxD, cancelD := Cancelable()
		cancelD()

		j2, _ := Join(ctxC, ctxD)
		select {
		case <-j2.Done():
		default:
			t.Fatal("expected joined context to be done immediately when B is already done")
		}
		if !errors.Is(j2.Err(), context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", j2.Err())
		}
	})
}

func TestJoin_Value(t *testing.T) {
	t.Parallel()

	type key string

	ctxA := WithValue(New(), key("k1"), "v1")
	ctxB := WithValue(New(), key("k2"), "v2")
	ctxBConflicting := WithValue(ctxB, key("k1"), "v1-b")

	j1, cancel1 := Join(ctxA, ctxB)
	defer cancel1()

	if v := j1.Value(key("k1")); v != "v1" {
		t.Errorf("expected v1, got %v", v)
	}
	if v := j1.Value(key("k2")); v != "v2" {
		t.Errorf("expected v2, got %v", v)
	}
	if v := j1.Value(key("missing")); v != nil {
		t.Errorf("expected nil, got %v", v)
	}

	j2, cancel2 := Join(ctxA, ctxBConflicting)
	defer cancel2()

	// A must take precedence over B
	if v := j2.Value(key("k1")); v != "v1" {
		t.Errorf("expected v1 (from ctxA), got %v", v)
	}
}

func TestJoin_Deadline(t *testing.T) {
	t.Parallel()

	now := time.Now()

	ctxNoDeadline := New()

	ctxEarly, cancelEarly := Deadline(now.Add(1 * time.Hour))
	defer cancelEarly()

	ctxLate, cancelLate := Deadline(now.Add(2 * time.Hour))
	defer cancelLate()

	tests := []struct {
		name         string
		a            context.Context
		b            context.Context
		wantDeadline time.Time
		wantOk       bool
	}{
		{"no deadlines", ctxNoDeadline, ctxNoDeadline, time.Time{}, false},
		{"only a has deadline", ctxEarly, ctxNoDeadline, now.Add(1 * time.Hour), true},
		{"only b has deadline", ctxNoDeadline, ctxEarly, now.Add(1 * time.Hour), true},
		{"both have deadlines, a is earlier", ctxEarly, ctxLate, now.Add(1 * time.Hour), true},
		{"both have deadlines, b is earlier", ctxLate, ctxEarly, now.Add(1 * time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, cancel := Join(tt.a, tt.b)
			defer cancel()

			gotDeadline, gotOk := j.Deadline()

			if gotOk != tt.wantOk {
				t.Errorf("Deadline() ok = %v, want %v", gotOk, tt.wantOk)
			}

			if gotOk && !gotDeadline.Equal(tt.wantDeadline) {
				t.Errorf("Deadline() = %v, want %v", gotDeadline, tt.wantDeadline)
			}
		})
	}
}
