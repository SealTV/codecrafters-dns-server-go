package main

import (
	"fmt"
	"log"
	"net"
	"strings"
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

		in := DNSMessage{}
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

func MakerResponse(in DNSMessage) DNSMessage {
	rcode := byte(0)
	if in.Header.OPCODE != 0 {
		rcode = 4
	}

	return DNSMessage{
		Header: DNSHeader{
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
		Question: DNSQuestion{
			QNAME:  "codecrafters.io",
			QTYPE:  1,
			QCLASS: 1,
		},
		Answer: DNSAnswer{
			Name:     "codecrafters.io",
			Type:     1,
			Class:    1,
			TTL:      60,
			RDLenght: 4,
			RDATA:    []byte{8, 8, 8, 8},
		},
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

	result = append(result, m.Question.Serialize()...)

	result = append(result, m.Answer.Serialize()...)

	return result
}

func (m *DNSMessage) Parse(data []byte) error {
	err := m.Header.Parse(data)
	if err != nil {
		return err
	}

	err = m.Question.Parse(data)
	if err != nil {
		return err
	}

	err = m.Answer.Parse(data)
	if err != nil {
		return err
	}

	return nil
}

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
	if len(data) < 12 {
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

type DNSQuestion struct {
	QNAME  string // The domain name to query.
	QTYPE  uint16 // The type of the query.
	QCLASS uint16 // The class of the query.
}

func (q *DNSQuestion) Parse(data []byte) error {
	return nil
}

func (q *DNSQuestion) Serialize() []byte {
	labels := strings.Split(q.QNAME, ".")

	result := make([]byte, 0, len(q.QNAME)+4)

	for _, label := range labels {
		result = append(result, byte(len(label)))
		result = append(result, []byte(label)...)
	}

	result = append(result, 0)

	result = append(result, byte(q.QTYPE>>8), byte(q.QTYPE))
	result = append(result, byte(q.QCLASS>>8), byte(q.QCLASS))

	return result
}

type DNSAnswer struct {
	Name     string // The domain name encoded as a sequence of labels.
	Type     uint16 // The type of the resource record. 1 for an A record, 5 for a CNAME record etc., full list here
	Class    uint16 // The class of the resource record. 1 for an internet address, full list here. Usually set to 1 (full list here)
	TTL      uint32 // The duration in seconds a record can be cached before requerying.
	RDLenght uint16 // The length of the RDATA field in bytes.
	RDATA    []byte // The data specific to the record type.
}

func (a *DNSAnswer) Parse(data []byte) error {
	return nil
}

func (a *DNSAnswer) Serialize() []byte {
	labels := strings.Split(a.Name, ".")

	result := make([]byte, 0, len(a.Name)+10+int(a.RDLenght))

	for _, label := range labels {
		result = append(result, byte(len(label)))
		result = append(result, []byte(label)...)
	}

	result = append(result, 0)

	result = append(result, byte(a.Type>>8), byte(a.Type))
	result = append(result, byte(a.Class>>8), byte(a.Class))
	result = append(result, byte(a.TTL>>24), byte(a.TTL>>16), byte(a.TTL>>8), byte(a.TTL))
	result = append(result, byte(a.RDLenght>>8), byte(a.RDLenght))

	result = append(result, a.RDATA...)

	return result
}
