package main

import (
	"reflect"
	"testing"
)

func TestMutex_IsLocked(t *testing.T) {
	type fields struct {
		ch chan struct{}
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
				ch: tt.fields.ch,
			}
			if got := m.IsLocked(); got != tt.want {
				t.Errorf("IsLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMutex_Lock(t *testing.T) {
	type fields struct {
		ch chan struct{}
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
				ch: tt.fields.ch,
			}
			m.Lock()
			m.Unlock()
		})
	}
}

func TestMutex_TryLock(t *testing.T) {
	type fields struct {
		ch chan struct{}
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
				ch: tt.fields.ch,
			}
			if got := m.TryLock(); got != tt.want {
				t.Errorf("TryLock() = %v, want %v", got, tt.want)
			}
			m.Unlock()
			t.Logf("m is %v", m)
		})
	}
}

func TestMutex_Unlock(t *testing.T) {
	type fields struct {
		ch chan struct{}
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
				ch: tt.fields.ch,
			}
			m.Lock()
			m.Unlock()
			t.Logf("m is %v", m)
		})
	}
}

func TestNewMutex(t *testing.T) {
	tests := []struct {
		name string
		want *Mutex
	}{
		{
			name: "test",
			want: &Mutex{},
		},
		{
			name: "test1",
			want: &Mutex{},
		},
		{
			name: "test2",
			want: &Mutex{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMutex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMutex() = %v, want %v", got, tt.want)
			}
			t.Logf("data is %v", tt)
		})
	}
}
