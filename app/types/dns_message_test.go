package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDNSMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     DNSMessage
		wantErr bool
	}{
		{
			name: "Parse",
			msg: DNSMessage{
				Header: DNSHeader{
					ID:      1,
					QR:      1,
					OPCODE:  3,
					AA:      1,
					TC:      1,
					RD:      1,
					RA:      1,
					Z:       1,
					RCODE:   1,
					QDCount: 1,
					ANCount: 1,
					NSCount: 0,
					ARCount: 0,
				},
				Questions: []DNSQuestion{{
					QName:  "wwe.google.com",
					QType:  1,
					QClass: 1,
				}},
				Answers: []DNSAnswer{
					{
						Name:     "wwe.google.com",
						Type:     1,
						Class:    1,
						TTL:      60,
						RDLenght: 4,
						RDATA:    []byte{8, 8, 8, 8},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parse multiple questions",
			msg: DNSMessage{
				Header: DNSHeader{
					ID:      1,
					QR:      1,
					OPCODE:  3,
					AA:      1,
					TC:      1,
					RD:      1,
					RA:      1,
					Z:       1,
					RCODE:   1,
					QDCount: 2,
					ANCount: 1,
					NSCount: 0,
					ARCount: 0,
				},
				Questions: []DNSQuestion{
					{
						QName:  "wwe.google.com",
						QType:  1,
						QClass: 1,
					},
					{
						QName:  "wwe.google.com",
						QType:  1,
						QClass: 1,
					},
				},
				Answers: []DNSAnswer{
					{
						Name:     "wwe.google.com",
						Type:     1,
						Class:    1,
						TTL:      60,
						RDLenght: 4,
						RDATA:    []byte{8, 8, 8, 8},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parse multiple questions with compression",
			msg: DNSMessage{
				Header: DNSHeader{
					ID:      1,
					QR:      1,
					OPCODE:  3,
					AA:      1,
					TC:      1,
					RD:      1,
					RA:      1,
					Z:       1,
					RCODE:   1,
					QDCount: 2,
					ANCount: 1,
					NSCount: 0,
					ARCount: 0,
				},
				Questions: []DNSQuestion{
					{
						QName:  "wwe.google.com",
						QType:  1,
						QClass: 1,
					},
					{
						QName:  "wwe.google.com",
						QType:  1,
						QClass: 1,
					},
				},
				Answers: []DNSAnswer{
					{
						Name:     "wwe.google.com",
						Type:     1,
						Class:    1,
						TTL:      60,
						RDLenght: 4,
						RDATA:    []byte{8, 8, 8, 8},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parse multiple questions with compression",
			msg: DNSMessage{
				Header: DNSHeader{
					ID:      20152,
					RD:      1,
					QDCount: 2,
					OPCODE:  QUERY,
				},
				Questions: []DNSQuestion{
					{
						QName:  "abc.longassdomainname.com",
						QType:  QTYPE(IN),
						QClass: QCLASS(A),
					},
					{
						QName:  "def.longassdomainname.com",
						QType:  QTYPE(IN),
						QClass: QCLASS(A),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.msg.Serialize()

			t.Log("data:", data)
			newQ := DNSMessage{}
			err := newQ.Parse(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DNSMessage.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.msg, newQ); diff != "" {
				t.Errorf("DNSMessage.Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
