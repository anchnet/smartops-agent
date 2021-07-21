package executor

import "testing"

func Test_formatOutputGather(t *testing.T) {
	type args struct {
		resourceName string
		elems        []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				resourceName: "prefix",
				elems:        []byte("a\nb"),
			},
			want: "prefix: a\nprefix: b",
		},
		{
			args: args{
				resourceName: "prefix",
				elems:        []byte("addd\ncccb"),
			},
			want: "prefix: addd\nprefix: cccb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatOutputGather(tt.args.resourceName, tt.args.elems); got != tt.want {
				t.Errorf("formatOutputGather() = %v, want %v", got, tt.want)
			}
		})
	}
}
