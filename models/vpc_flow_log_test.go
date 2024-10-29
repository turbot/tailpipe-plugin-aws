package models

import (
	"github.com/turbot/pipe-fittings/utils"
	"reflect"
	"testing"
)

func TestFlowLogFromString(t *testing.T) {
	type args struct {
		rowString string
		schema    []string
	}
	tests := []struct {
		name    string
		args    args
		want    *AwsVpcFlowLog
		wantErr bool
	}{
		{
			name: "Test FlowLogFromString",
			args: args{
				rowString: "2 123456789012 eni-1235b8ca",
				schema: []string{
					"version",
					"account-id",
					"interface-id",
				},
			},
			want: &AwsVpcFlowLog{
				Version:     utils.ToPointer[int32](2),
				AccountID:   utils.ToStringPointer("123456789012"),
				InterfaceID: utils.ToStringPointer("eni-1235b8ca"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlowLogFromString(tt.args.rowString, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlowLogFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlowLogFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
