package main

import (
	"sync"
	"testing"
)

func TestMutex_IsLocked(t *testing.T) {
	type fields struct {
		mu sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test",
			fields: fields{},
			want:   false,
		},
		{
			name:   "test1",
			fields: fields{},
			want:   false,
		},
		{
			name:   "test2",
			fields: fields{},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mutex{
				mu: tt.fields.mu,
			}
			if got := m.IsLocked(); got != tt.want {
				t.Errorf("IsLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMutex_Lock(t *testing.T) {
	type fields struct {
		mu sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "test",
			fields: fields{},
		},
		{
			name:   "test1",
			fields: fields{},
		},
		{
			name:   "test2",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mutex{
				mu: tt.fields.mu,
			}
			m.Lock()
			m.Unlock()
		})
	}
}

func TestMutex_TryLock(t *testing.T) {
	type fields struct {
		mu sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test",
			fields: fields{},
			want:   true,
		},
		{
			name:   "test1",
			fields: fields{},
			want:   true,
		},
		{
			name:   "test2",
			fields: fields{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mutex{
				mu: tt.fields.mu,
			}
			if got := m.TryLock(); got != tt.want {
				t.Errorf("TryLock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMutex_Unlock(t *testing.T) {
	type fields struct {
		mu sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "test",
			fields: fields{},
		},
		{
			name:   "test1",
			fields: fields{},
		},
		{
			name:   "test2",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mutex{
				mu: tt.fields.mu,
			}
			m.Lock()
			m.Unlock()
		})
	}
}
