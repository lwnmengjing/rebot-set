package lark

import (
	"os"
	"testing"
)

func TestSendLarkMessage(t *testing.T) {
	type args struct {
		webhook string
		secret  string
		color   string
		title   string
		foot    string
		lines   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case0",
			args: args{
				webhook: os.Getenv("LARK_WEBHOOK"),
				secret:  os.Getenv("LARK_SECRET"),
				color:   "red",
				title:   "test",
				foot:    "此消息由系统自动生成，请勿回复。",
				lines:   []string{"test0", "test1", "test2", "test3"},
			},
			wantErr: false,
		},
		{
			name: "case1",
			args: args{
				webhook: os.Getenv("LARK_WEBHOOK"),
				secret:  os.Getenv("LARK_SECRET"),
				color:   "green",
				title:   "title",
				foot:    "此消息由系统自动生成，请勿回复。",
				lines:   []string{"content"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendLarkMessage(tt.args.webhook, tt.args.secret, tt.args.color, tt.args.title, tt.args.foot, tt.args.lines); (err != nil) != tt.wantErr {
				t.Errorf("SendLarkMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
