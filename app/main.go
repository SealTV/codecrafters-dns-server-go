package main

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/types"
	// Uncomment this block to pass the first stage
	// "net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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

		fmt.Printf("Parces Message: %+v\n", in)

		msg := MakerResponse(in)
		_, err = udpConn.WriteToUDP(msg.Serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func MakerResponse(in types.DNSMessage) types.DNSMessage {
	rcode := types.NOERROR

	if in.Header.OPCODE != types.IQUERY {
		rcode = types.NOTIMP
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
			QDCount: 1,
			ANCount: 1,
			NSCount: 0,
			ARCount: 0,
		},
		Questions: []types.DNSQuestion{{
			QName:  in.Questions[0].QName,
			QType:  1,
			QClass: 1,
		}},
		Answers: []types.DNSAnswer{{
			Name:     in.Questions[0].QName,
			Type:     1,
			Class:    1,
			TTL:      60,
			RDLenght: 4,
			RDATA:    []byte{8, 8, 8, 8},
		}},
	}
}
