package main

import (
	"fmt"
	"log"
	"net"
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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		msg := DNSMessage{
			Header: DNSHeader{
				ID: 1234,
				QR: 1,
			},
		}

		response := msg.Serialize()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

type DNSMessage struct {
	Header   DNSHeader
	Question DNSQuestion
	Answer   DNSAnswer
	Space    []byte
}

func (m *DNSMessage) Serialize() []byte {
	result := m.Header.Serialize()

	//TODO: Add question and answer sections

	return result
}

type DNSQuestion struct{}

type DNSAnswer struct{}

type DNSHeader struct {
	ID      uint16 // A random ID assigned to query packets. Response packets must reply with the same ID.
	QR      byte   // 1 bit. 1 for a reply packet, 0 for a question packet.
	OPCODE  byte   // 4 bits. 0 for a standard query, 1 for an inverse query, 2 for a server status request, 3-15 reserved.
	AA      byte   // 1 bit. 1 for an authoritative response, 0 otherwise.
	TC      byte   // 1 bit. 1 if the response was truncated, 0 otherwise.
	RD      byte   // 1 bit. 1 if recursion is desired, 0 otherwise.
	RA      byte   // 1 bit. 1 if recursion is available, 0 otherwise.
	Z       byte   // 3 bits. Reserved for future use. Must be 0.
	RCODE   byte   // 4 bits. 0 for no error, 1 for a format error, 2 for a server failure, 3 for a name error, 4 for a not implemented error, 5 for a refused error, 6-15 reserved.
	QDCount uint16 // The number of questions in the question section.
	ANCount uint16 // The number of resource records in the answer section.
	NSCount uint16 // The number of name server resource records in the authority records section.
	ARCount uint16 // The number of resource records in the additional records section.
}

func (h *DNSHeader) Parse(data []byte) error {
	if len(data) != 12 {
		return fmt.Errorf("Not enough data to parse header")
	}

	h.ID = uint16(data[0])<<8 | uint16(data[1])
	h.QDCount = uint16(data[4])<<8 | uint16(data[5])
	h.ANCount = uint16(data[6])<<8 | uint16(data[7])
	h.NSCount = uint16(data[8])<<8 | uint16(data[9])
	h.ARCount = uint16(data[10])<<8 | uint16(data[11])

	return nil
}

func (h *DNSHeader) Serialize() []byte {
	data := make([]byte, 12)

	data[0] = byte(h.ID >> 8)
	data[1] = byte(h.ID)

	data[2] = byte(h.QR<<7) | byte(h.OPCODE<<3) | byte(h.AA<<2) | byte(h.TC<<1) | byte(h.RD)
	data[3] = byte(h.RA<<7) | byte(h.Z<<4) | byte(h.RCODE)

	data[4] = byte(h.QDCount >> 8)
	data[5] = byte(h.QDCount)

	data[6] = byte(h.ANCount >> 8)
	data[7] = byte(h.ANCount)

	data[8] = byte(h.NSCount >> 8)
	data[9] = byte(h.NSCount)

	data[10] = byte(h.ARCount >> 8)
	data[11] = byte(h.ARCount)

	return data
}
