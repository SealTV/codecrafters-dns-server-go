package main

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/types"
)

func main() {
	resolver, err := GetDNSResolver()
	if err != nil {
		log.Fatal(err)
	}

	defer resolver.conn.Close()

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	log.Println("Listening on", udpAddr)

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		fmt.Printf("Received %d bytes from %s: %v\n", size, source, buf[:size])

		in := types.DNSMessage{}
		err = in.Parse(buf[:size])
		if err != nil {
			fmt.Println("Failed to parse message:", err)
			continue
		}

		fmt.Printf("Parsed Message: %+v, \n \n", in)

		msg := MakerResponse(resolver, in)
		_, err = udpConn.WriteToUDP(msg.Serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func MakerResponse(resolver *dnsResolver, in types.DNSMessage) types.DNSMessage {
	rcode := types.NOERROR

	if in.Header.OPCODE != types.QUERY {
		rcode = types.NOTIMP
	}

	questions := in.Questions
	answers := []types.DNSAnswer{}

	if rcode == types.NOERROR {
		var err error
		answers, err = resolver.ResolveAddress(questions)
		if err != nil {
			rcode = types.SERVFAIL
		}
	}

	return types.DNSMessage{
		Header: types.DNSHeader{
			ID:      in.Header.ID,
			QR:      1,
			OPCODE:  in.Header.OPCODE,
			AA:      0,
			TC:      0,
			RD:      in.Header.RD,
			RA:      0,
			Z:       0,
			RCODE:   rcode,
			QDCount: uint16(len(questions)),
			ANCount: uint16(len(answers)),
			NSCount: 0,
			ARCount: 0,
		},
		Questions: questions,
		Answers:   answers,
	}
}

type dnsResolver struct {
	conn net.Conn
}

func GetDNSResolver() (*dnsResolver, error) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:53")
	if err != nil {
		return nil, fmt.Errorf("Cannot resolve local DNS server address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("Cannot dial local DNS server: %v", err)
	}

	return &dnsResolver{
		conn: conn,
	}, nil
}

func (dnsr *dnsResolver) ResolveAddress(quesions types.DBSQuestions) ([]types.DNSAnswer, error) {
	msg := types.DNSMessage{
		Header: types.DNSHeader{
			ID:      1,
			OPCODE:  types.QUERY,
			RD:      1,
			QDCount: uint16(len(quesions)),
		},
		Questions: quesions,
	}

	_, err := dnsr.conn.Write(msg.Serialize())
	if err != nil {
		return nil, fmt.Errorf("Cannot send message to local DNS server: %v", err)
	}

	buf := make([]byte, 512)
	size, err := dnsr.conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("Cannot read response from local DNS server: %v", err)
	}

	res := types.DNSMessage{}
	err = res.Parse(buf[:size])
	if err != nil {
		return nil, fmt.Errorf("Cannot parse response from local DNS server: %v", err)
	}

	if res.Header.RCODE != types.NOERROR {
		return nil, fmt.Errorf("Local DNS server returned error: %v", res.Header.RCODE)
	}

	if len(res.Answers) == 0 {
		return nil, fmt.Errorf("Local DNS server returned no answers")
	}

	return res.Answers, nil
}
