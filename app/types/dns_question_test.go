package types

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDNSQuestion(t *testing.T) {
	tests := []struct {
		name    string
		q       DBSQuestions
		wantErr bool
	}{
		{
			name: "Parse",
			q: DBSQuestions{{
				QName:  "wwe.google.com",
				QType:  QTYPE(A),
				QClass: QCLASS(IN),
			}},
			wantErr: false,
		},
		{
			name: "Parse 2",
			q: DBSQuestions{
				{
					QName:  "wwe.google.com",
					QType:  QTYPE(A),
					QClass: QCLASS(IN),
				},
				{
					QName:  "qwe.wwe.google.com",
					QType:  QTYPE(A),
					QClass: QCLASS(IN),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}

			if err := tt.q.Serialize(&buf); err != nil {
				t.Errorf("DNSQuestion.Serialize() error = %v, bytes: %v", err, buf.Bytes())
				return
			}

			data := buf.Bytes()

			t.Logf("serialized questions: %v", string(data))
			t.Logf("serialized questions: %+v", data)

			newQ := DBSQuestions{}
			got, err := newQ.Parse(uint16(len(tt.q)), data, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("DNSQuestion.Parse() error = %v, wantErr %v, bytes: %v", err, tt.wantErr, string(data))
				return
			}

			if got != len(data) {
				t.Errorf("DNSQuestion.Parse() = %v, want %v", got, len(data))
			}

			if diff := cmp.Diff(tt.q, newQ); diff != "" {
				t.Errorf("DNSQuestion.Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
