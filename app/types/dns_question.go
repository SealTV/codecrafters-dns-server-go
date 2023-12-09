package types

import (
	"bytes"
	"strings"
)

type DBSQuestions []DNSQuestion

func (questions *DBSQuestions) Serialize(buf *bytes.Buffer) error {
	hash := map[string]uint16{}

	for _, q := range *questions {
		labels := strings.Split(q.QName, ".")
		name := q.QName

		commpessed := false

		for i := range labels {
			label := labels[i]
			position := buf.Len()

			if p, ok := hash[name]; !ok {
				buf.WriteByte(byte(len(label)))
				buf.WriteString(label)

				hash[name] = uint16(position)
				name = strings.Join(labels[i+1:], ".")
			} else {
				buf.WriteByte(byte(3)<<6 | byte(p>>8))
				buf.WriteByte(byte(p))

				commpessed = true
				break
			}
		}

		if !commpessed {
			buf.WriteByte(0)
		}

		buf.WriteByte(byte(q.QType >> 8))
		buf.WriteByte(byte(q.QType))

		buf.WriteByte(byte(q.QClass >> 8))
		buf.WriteByte(byte(q.QClass))
	}

	// fmt.Println("Hash:", hash)

	return nil
}

func (questions *DBSQuestions) Parse(count uint16, data []byte, offset int) (int, error) {
	for i := uint16(0); i < count; i++ {
		q := DNSQuestion{}

		labels := []string{}
		for {
			b1 := data[offset]
			offset++

			if b1>>6 == 3 {
				targetPosition := uint16(b1&0x3F)<<8 | uint16(data[offset])
				offset++

				// fmt.Println("try to read compressed string from position:", targetPosition)
				str, err := readCompressedString(data, int(targetPosition))
				if err != nil {
					return 0, err
				}

				// fmt.Println("read compressed string:", str)

				labels = append(labels, str)

				break
			}

			length := int(b1)

			if length == 0 {
				break
			}

			label := data[offset : offset+length]
			offset += length

			labels = append(labels, string(label))
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

func readCompressedString(data []byte, offset int) (string, error) {
	labels := []string{}
	for {

		b1 := data[offset]
		offset++

		length := int(b1)

		if length == 0 {
			break
		}

		label := data[offset : offset+length]
		offset += length

		labels = append(labels, string(label))
	}

	return strings.Join(labels, "."), nil
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
