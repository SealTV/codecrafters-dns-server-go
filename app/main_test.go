package main

import "testing"

func TestDNSQuestion_Parse(t *testing.T) {
	tests := []struct {
		name    string
		q       DNSQuestion
		wantErr bool
	}{
		{
			name: "Parse",
			q: DNSQuestion{
				QNAME:  "www.google.com",
				QTYPE:  1,
				QCLASS: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.q.Serialize()

			newQ := DNSQuestion{}
			got, err := newQ.Parse(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DNSQuestion.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != len(data) {
				t.Errorf("DNSQuestion.Parse() = %v, want %v", got, len(data))
			}

			if newQ.QNAME != tt.q.QNAME {
				t.Errorf("DNSQuestion.Parse() QNAME = %v, want %v", newQ.QNAME, tt.q.QNAME)
			}

			if newQ.QTYPE != tt.q.QTYPE {
				t.Errorf("DNSQuestion.Parse() QTYPE = %v, want %v", newQ.QTYPE, tt.q.QTYPE)
			}

			if newQ.QCLASS != tt.q.QCLASS {
				t.Errorf("DNSQuestion.Parse() QCLASS = %v, want %v", newQ.QCLASS, tt.q.QCLASS)
			}
		})
	}
}
