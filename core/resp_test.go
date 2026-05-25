package core

import (
	"testing"
	"reflect"
)

func TestReadSimpleString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantPos int
		wantErr bool
	}{
		{
			name:    "basic string",
			input:   []byte("+OK\r\n"),
			want:    "OK",
			wantPos: 5,
		},
		{
			name:    "empty string",
			input:   []byte("+\r\n"),
			want:    "",
			wantPos: 3,
		},
		{
			name:    "string with spaces",
			input:   []byte("+hello world\r\n"),
			want:    "hello world",
			wantPos: 14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, pos, err := readSimpeString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("value = %q, want %q", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("pos = %d, want %d", pos, tt.wantPos)
			}
		})
	}
}

func TestReadError(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantPos int
		wantErr bool
	}{
		{
			name:    "generic error",
			input:   []byte("-ERR unknown command\r\n"),
			want:    "ERR unknown command",
			wantPos: 22,
		},
		{
			name:    "wrong type error",
			input:   []byte("-WRONGTYPE value is not a string\r\n"),
			want:    "WRONGTYPE value is not a string",
			wantPos: 34,
		},
		{
			name:    "empty error message",
			input:   []byte("-\r\n"),
			want:    "",
			wantPos: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, pos, err := readError(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("value = %q, want %q", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("pos = %d, want %d", pos, tt.wantPos)
			}
		})
	}
}

func TestReadInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int64
		wantPos int
		wantErr bool
	}{
		{
			name:    "positive integer",
			input:   []byte(":42\r\n"),
			want:    42,
			wantPos: 5,
		},
		{
			name:    "zero",
			input:   []byte(":0\r\n"),
			want:    0,
			wantPos: 4,
		},
		{
			name:    "negative integer",
			input:   []byte(":-1\r\n"),
			want:    -1,
			wantPos: 5,
		},
		{
			name:    "large number",
			input:   []byte(":1000000\r\n"),
			want:    1000000,
			wantPos: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, pos, err := readInt64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("value = %d, want %d", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("pos = %d, want %d", pos, tt.wantPos)
			}
		})
	}
}

func TestReadBulkString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantPos int
		wantErr bool
	}{
		{
			name:    "basic string",
			input:   []byte("$5\r\nhello\r\n"),
			want:    "hello",
			wantPos: 11,
		},
		{
			name:    "single character",
			input:   []byte("$1\r\na\r\n"),
			want:    "a",
			wantPos: 7,
		},
		{
			name:    "empty string",
			input:   []byte("$0\r\n\r\n"),
			want:    "",
			wantPos: 6,
		},
		{
			name:    "string with spaces",
			input:   []byte("$11\r\nhello world\r\n"),
			want:    "hello world",
			wantPos: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, pos, err := readBulkString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("value = %q, want %q", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("pos = %d, want %d", pos, tt.wantPos)
			}
		})
	}
}

func TestReadArray(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []interface{}
		wantPos int
		wantErr bool
	}{
		{
			name:    "array of integers",
			input:   []byte("*3\r\n:1\r\n:2\r\n:3\r\n"),
			want:    []interface{}{int64(1), int64(2), int64(3)},
			wantPos: 16,
		},
		{
			name:    "array of bulk strings",
			input:   []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
			want:    []interface{}{"foo", "bar"},
			wantPos: 22,
		},
		{
			name:    "single element",
			input:   []byte("*1\r\n+OK\r\n"),
			want:    []interface{}{"OK"},
			wantPos: 9,
		},
		{
			name:    "mixed types",
			input:   []byte("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n"),
			want:    []interface{}{int64(1), int64(2), int64(3), int64(4), "hello"},
			wantPos: 31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, pos, err := readArray(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("value = %v, want %v", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("pos = %d, want %d", pos, tt.wantPos)
			}
		})
	}
}
