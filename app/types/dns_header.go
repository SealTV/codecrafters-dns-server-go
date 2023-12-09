package types

import "fmt"

type DNSHeader struct {
	ID      uint16 // A random ID assigned to query packets. Response packets must reply with the same ID.
	QR      byte   // 1 bit. 1 for a reply packet, 0 for a question packet.
	OPCODE  OPCODE // 4 bits. 0 for a standard query, 1 for an inverse query, 2 for a server status request, 3-15 reserved.
	AA      byte   // 1 bit. 1 for an authoritative response, 0 otherwise.
	TC      byte   // 1 bit. 1 if the response was truncated, 0 otherwise.
	RD      byte   // 1 bit. 1 if recursion is desired, 0 otherwise.
	RA      byte   // 1 bit. 1 if recursion is available, 0 otherwise.
	Z       byte   // 3 bits. Reserved for future use. Must be 0.
	RCODE   RCODE  // 4 bits. 0 for no error, 1 for a format error, 2 for a server failure, 3 for a name error, 4 for a not implemented error, 5 for a refused error, 6-15 reserved.
	QDCount uint16 // The number of questions in the question section.
	ANCount uint16 // The number of resource records in the answer section.
	NSCount uint16 // The number of name server resource records in the authority records section.
	ARCount uint16 // The number of resource records in the additional records section.
}

func (h *DNSHeader) Parse(data []byte) (int, error) {
	if len(data) < 12 {
		return 0, fmt.Errorf("Not enough data to parse header")
	}

	h.ID = uint16(data[0])<<8 | uint16(data[1])

	h.QR = data[2] >> 7
	h.OPCODE = OPCODE((data[2] >> 3) & 0x0F)
	h.AA = (data[2] >> 2) & 0x01
	h.TC = (data[2] >> 1) & 0x01
	h.RD = data[2] & 0x01

	h.RA = data[3] >> 7
	h.Z = (data[3] >> 4) & 0x07
	h.RCODE = RCODE(data[3] & 0x0F)

	h.QDCount = uint16(data[4])<<8 | uint16(data[5])
	h.ANCount = uint16(data[6])<<8 | uint16(data[7])
	h.NSCount = uint16(data[8])<<8 | uint16(data[9])
	h.ARCount = uint16(data[10])<<8 | uint16(data[11])

	return 12, nil
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
