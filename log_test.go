package lock

import (
	"testing"
)

func TestLog_PrintLockUsageTime(t *testing.T) {
	type fields struct {
		Logger Logger
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "print1",
			fields: fields{
				Logger: NewLog(),
			},
			args: args{
				format: "name:%s use %d ",
				args:   []interface{}{"lock1", 10},
			},
		},
		{
			name: "print2",
			fields: fields{
				Logger: NewLog(),
			},
			args: args{
				format: "name:%s use %d ",
				args:   []interface{}{"lock2", 20},
			},
		},
		{
			name: "print3",
			fields: fields{
				Logger: NewLog(),
			},
			args: args{
				format: "name:%s use %d ",
				args:   []interface{}{"lock3", 30},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.fields.Logger
			l.PrintLockUsageTime(tt.args.format, tt.args.args...)
		})
	}
}
