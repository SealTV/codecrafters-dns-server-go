package types

import (
	"bytes"
	"fmt"
	"strings"
)

type DBSQuestions []DNSQuestion

func (questions *DBSQuestions) Serialize(buf *bytes.Buffer) error {
	hash := map[string]uint16{}

	for _, q := range *questions {
		for _, label := range strings.Split(q.QName, ".") {
			position := buf.Len()

			if p, ok := hash[label]; !ok {
				hash[label] = uint16(position)

				buf.WriteByte(byte(len(label)))
				buf.WriteString(label)
			} else {
				buf.WriteByte(byte(3)<<6 | byte(p>>8))
				buf.WriteByte(byte(p))
			}
		}

		buf.WriteByte(0)

		buf.WriteByte(byte(q.QType >> 8))
		buf.WriteByte(byte(q.QType))

		buf.WriteByte(byte(q.QClass >> 8))
		buf.WriteByte(byte(q.QClass))
	}

	return nil
}

func (questions *DBSQuestions) Parse(count uint16, data []byte, offset int) (int, error) {
	hash := map[uint16]string{}

	for i := uint16(0); i < count; i++ {
		q := DNSQuestion{}

		labels := []string{}
		for {
			position := offset

			b1 := data[offset]
			offset++

			if b1>>6 == 3 {
				targetPosition := uint16(b1&0x3F)<<8 | uint16(data[offset])
				offset++

				if label, ok := hash[targetPosition]; ok {
					labels = append(labels, label)
				} else {
					return offset, fmt.Errorf("Unexpected index for compression string, for index: %v, known labels: %v", targetPosition, hash)
				}

				continue
			}

			length := int(b1)

			if length == 0 {
				break
			}

			label := data[offset : offset+length]
			offset += length

			labels = append(labels, string(label))
			hash[uint16(position)] = string(label)
		}

		q.QName = strings.Join(labels, ".")

		q.QType = QTYPE(uint16(data[offset])<<8 | uint16(data[offset+1]))
		offset += 2

		q.QClass = QCLASS(uint16(data[offset])<<8 | uint16(data[offset+1]))
		offset += 2

		*questions = append(*questions, q)
	}

	return offset, nil
}

type DNSQuestion struct {
	QName  string // The domain name to query.
	QType  QTYPE  // The type of the query.
	QClass QCLASS // The class of the query.
}

func (q *DNSQuestion) Parse(data []byte, offset int) (int, error) {
	// check is compression
	if data[offset]>>6 == 3 {
		offset += 2
		return offset, nil
	}

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

	q.QName = strings.Join(labels, ".")

	q.QType = QTYPE(uint16(data[offset+length])<<8 | uint16(data[offset+length+1]))
	offset += 2

	q.QClass = QCLASS(uint16(data[offset+length])<<8 | uint16(data[offset+length+1]))
	offset += 2

	return offset, nil
}
