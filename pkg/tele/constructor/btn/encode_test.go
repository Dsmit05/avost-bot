package btn

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantEvent string
		wantArgs  []string
	}{
		{
			name:      "1 case: event and 1 arg",
			data:      "one test",
			wantEvent: "one",
			wantArgs:  []string{"test"},
		},
		{
			name:      "2 case: only event",
			data:      "one",
			wantEvent: "one",
			wantArgs:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvent, gotArgs := DecodeStr(tt.data)
			if gotEvent != tt.wantEvent {
				t.Errorf("DecodeStr() gotEvent = %v, want %v", gotEvent, tt.wantEvent)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("DecodeStr() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestDecodeInt(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantEvent string
		wantArgs  []int
	}{
		{
			name:      "1 case: event and 1 arg",
			data:      "event 123",
			wantEvent: "event",
			wantArgs:  []int{123},
		},
		{
			name:      "2 case: only event",
			data:      "one",
			wantEvent: "one",
			wantArgs:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvent, gotArgs := DecodeInt(tt.data)
			assert.Equal(t, gotEvent, tt.wantEvent)
			assert.Equal(t, gotArgs, tt.wantArgs)
		})
	}
}
