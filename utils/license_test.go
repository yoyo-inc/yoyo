package utils

import "testing"

func TestSNLicense_Activate(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "activate",
			args: args{code: "isYk24mulZgMsLl1CKFHLmOCCD7wFw5PpqKQww1YOxwLFHOc1FYOicoTYUNpl3KQJocC3d8OzQh4OnsofflPlueA5xthbB7hNGqjjZzXyyDQxaPrlx3OjHaCKkwcMCVvIL5kuMJybwwdcwredeq1wda2oevaw4LP3ey9XIPvUTxHKb6mOSLo2VBNrZOeb2VwHfeYh9V6ByF9uUgPbByElzoIwzDuQi2ibiCk4Yu84wgOZWse0ifAjseJ0uMUDMcHGOA"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := &SNLicense{}
			if got := sn.Activate(tt.args.code); got != tt.want {
				t.Errorf("Activate() = %v, want %v", got, tt.want)
			}
		})
	}
}
