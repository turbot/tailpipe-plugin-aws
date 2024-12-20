package rows

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

func TestAwsElbAccessLog_InitialiseFromMap(t *testing.T) {

	type args struct {
		m map[string]string
	}
	tests := []struct {
		name    string
		want    *ElbAccessLog
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "All valid fields",
			args: args{
				m: map[string]string{
					"timestamp":                "2021-07-01T00:00:00Z",
					"type":                     "type",
					"elb":                      "elb",
					"client":                   "198.51.100.1:100",
					"target":                   "198.51.100.2:200",
					"request_processing_time":  "99.9",
					"target_processing_time":   "55.5",
					"response_processing_time": "33.3",
					"elb_status_code":          "200",
					"target_status_code":       "200",
					"received_bytes":           "100",
					"sent_bytes":               "100",
					"request":                  "request",
					"user_agent":               "agent",
					"ssl_cipher":               "sahr4098ewrjofdsjkdsf",
					"ssl_protocol":             "SSLv3",
					"target_group_arn":         "arn:partition:service:region:account-id:resource",
					"trace_id":                 "trace",
					"domain_name":              "domain",
					"chosen_cert_arn":          "arn:partition:service:region:account-id:resource",
					"matched_rule_priority":    "1",
					"request_creation_time":    "2021-07-01T00:00:00Z",
					"actions_executed":         "actions",
					"redirect_url":             "https://myapp.example.com/auth/callback?code=AUTH_CODE",
					"error_reason":             "reason",
					"target_list":              "list",
					"target_status_list":       "list",
					"classification":           "x",
					"classification_reason":    "reason",
					"conn_trace_id":            "abcd",
				},
			},

			want: &ElbAccessLog{

				Type:                   "type",
				Timestamp:              time.Date(2021, 07, 01, 00, 00, 00, 00, time.UTC),
				Elb:                    "elb",
				ClientIP:               "198.51.100.1",
				ClientPort:             100,
				TargetIP:               utils.ToStringPointer("198.51.100.2"),
				TargetPort:             200,
				RequestProcessingTime:  99.9,
				TargetProcessingTime:   55.5,
				ResponseProcessingTime: 33.3,
				ElbStatusCode:          utils.ToPointer(200),
				TargetStatusCode:       utils.ToPointer(200),
				ReceivedBytes:          utils.ToPointer(int64(100)),
				SentBytes:              utils.ToPointer(int64(100)),
				Request:                "request",
				UserAgent:              "agent",
				SslCipher:              "sahr4098ewrjofdsjkdsf",
				SslProtocol:            "SSLv3",
				TargetGroupArn:         "arn:partition:service:region:account-id:resource",
				TraceID:                "trace",
				DomainName:             "domain",
				ChosenCertArn:          "arn:partition:service:region:account-id:resource",
				MatchedRulePriority:    1,
				RequestCreationTime:    time.Date(2021, 07, 01, 00, 00, 00, 00, time.UTC),
				ActionsExecuted:        "actions",
				RedirectURL:            utils.ToStringPointer("https://myapp.example.com/auth/callback?code=AUTH_CODE"),
				ErrorReason:            utils.ToStringPointer("reason"),
				TargetList:             utils.ToStringPointer("list"),
				TargetStatusList:       utils.ToStringPointer("list"),
				Classification:         utils.ToStringPointer("x"),
				ClassificationReason:   utils.ToStringPointer("reason"),
				ConnTraceID:            utils.ToStringPointer("abcd"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid timestamp",
			args: args{
				m: map[string]string{
					"timestamp": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		//{
		//	name: "Non string map value",
		//	args: args{
		//		m: map[string]string{
		//			"timestamp": 1,
		//		},
		//	},
		//	wantErr: assert.Error,
		//},
		{
			name: "Invalid request_processing_time",
			args: args{
				m: map[string]string{
					"request_processing_time": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid target_processing_time",
			args: args{
				m: map[string]string{
					"target_processing_time": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid response_processing_time",
			args: args{
				m: map[string]string{
					"response_processing_time": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid elb_status_code",
			args: args{
				m: map[string]string{
					"elb_status_code": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid target_status_code",
			args: args{
				m: map[string]string{
					"target_status_code": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid received_bytes",
			args: args{

				m: map[string]string{
					"received_bytes": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid sent_bytes",
			args: args{
				m: map[string]string{
					"sent_bytes": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid matched_rule_priority",
			args: args{
				m: map[string]string{
					"matched_rule_priority": "invalid",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid request_creation_time",
			args: args{
				m: map[string]string{
					"request_creation_time": "invalid",
				},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ElbAccessLog{}
			err := l.InitialiseFromMap(tt.args.m)
			// Check if the error assertion passes or fails
			if !tt.wantErr(t, err, "InitialiseFromMap()") {
				// Stop immediately if error assertion fails
				t.FailNow()
			}
			if err != nil {
				return
			}
			// Tactical - ignore tp fields
			l.CommonFields = schema.CommonFields{}
			equal := cmp.Equal(l, tt.want, cmpopts.IgnoreUnexported(ElbAccessLog{}))
			if !equal {
				t.Errorf("InitialiseFromMap() = %v, want %v", l, tt.want)
			}
		})
	}
}
