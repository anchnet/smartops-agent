//Package filter 过滤用户不需要的数据， filter结构体存放用户需要的数据
package filter

import (
	"testing"
)

func TestSetFilter(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "testSetFilter",
			args: args{
				bytes: []byte(`{
					"cpu": ["system.cpu.idle"],
					"proc":[],
					"mem":[],
					"disk":["system.disk.free"],
					"mysql":[],
					"nginx":[]
				}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFilter(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("SetFilter() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(GetFilter().Map[metricFlied]) != 2 {
				t.Error("SetFilter() error: metric init fail ")
			}
		})
	}
}

func Test_find(t *testing.T) {
	type args struct {
		arr []string
		val string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "在数组里面",
			args: args{
				arr: []string{"a", "b"},
				val: "a",
			},
			want: true,
		},
		{
			name: "不在数组里面",
			args: args{
				arr: []string{"a", "b"},
				val: "c",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := find(tt.args.arr, tt.args.val); got != tt.want {
				t.Errorf("find() = %v, want %v", got, tt.want)
			}
		})
	}
}
