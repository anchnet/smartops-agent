package conf

import "testing"

func TestChangeConfSite(t *testing.T) {
	type args struct {
		site     string
		confPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "find site",
			args: args{
				site:     "test.org",
				confPath: "./conf_test.yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeConfSite(tt.args.site, tt.args.confPath); (err != nil) != tt.wantErr {
				t.Errorf("ChangeConfSite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
