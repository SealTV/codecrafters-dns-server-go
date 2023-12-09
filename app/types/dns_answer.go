package types

import "strings"

type DNSAnswer struct {
	Name     string // The domain name encoded as a sequence of labels.
	Type     TYPE   // The type of the resource record. 1 for an A record, 5 for a CNAME record etc., full list here
	Class    CLASS  // The class of the resource record. 1 for an internet address, full list here. Usually set to 1 (full list here)
	TTL      uint32 // The duration in seconds a record can be cached before requerying.
	RDLenght uint16 // The length of the RDATA field in bytes.
	RDATA    []byte // The data specific to the record type.
}

func (a *DNSAnswer) Parse(data []byte, offset int) (int, error) {
	length := int(data[offset])
	offset++

	labels := []string{}
	for length > 0 {
		label := data[offset : offset+length]
		offset += len(label)

		length = int(data[offset])
		offset++

		labels = append(labels, string(label))
	}

	a.Name = strings.Join(labels, ".")

	a.Type = TYPE(uint16(data[offset])<<8 | uint16(data[offset+1]))
	offset += 2

	a.Class = CLASS(uint16(data[offset])<<8 | uint16(data[offset+1]))
	offset += 2

	a.TTL = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
	offset += 4

	a.RDLenght = uint16(data[offset])<<8 | uint16(data[offset+1])
	offset += 2

	a.RDATA = data[offset : offset+int(a.RDLenght)]
	offset += int(a.RDLenght)

	return offset, nil
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
