package services

import (
	"path/filepath"
	"testing"
)

func TestLogFilter(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "log",
			args: args{
				name: "test_2023-01-06.log",
			},
			want:  "2023-01-06",
			want1: true,
		},
		{
			name: "ziplog",
			args: args{
				name: "test_2023-01-06.log.zip",
			},
			want:  "2023-01-06",
			want1: true,
		},
		{
			name: "nolog",
			args: args{
				name: "test.log",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LogFilter(tt.args.name)
			if got != tt.want {
				t.Errorf("logFilter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("logFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestArchiveLogByRecent(t *testing.T) {
	p, _ := filepath.Abs("../logs")
	type args struct {
		dir    string
		recent int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "archive_log",
			args: args{
				dir:    p,
				recent: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ArchiveLogByRecent(tt.args.dir, tt.args.recent); (err != nil) != tt.wantErr {
				t.Errorf("ArchiveLogByRecent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
