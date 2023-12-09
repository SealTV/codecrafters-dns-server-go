package types

import (
	"bytes"
	"fmt"
)

type DNSMessage struct {
	Header    DNSHeader
	Questions DBSQuestions
	Answers   []DNSAnswer
	Space     []byte
}

func (m *DNSMessage) Serialize() []byte {
	buf := bytes.Buffer{}

	buf.Write(m.Header.Serialize())

	if err := m.Questions.Serialize(&buf); err != nil {
		fmt.Println("Error serializing questions:", err)
	}

	// result = append(result, m.Questions.Serialize()...)

	for _, a := range m.Answers {
		// result = append(result, a.Serialize()...)
		buf.Write(a.Serialize())
	}

	return buf.Bytes()
}

func (m *DNSMessage) Parse(data []byte) error {
	n, err := m.Header.Parse(data)
	if err != nil {
		return err
	}

	n, err = m.Questions.Parse(m.Header.QDCount, data, n)
	if err != nil {
		return fmt.Errorf("cannot parse questions: %w", err)
	}

	// fmt.Printf("QUSTIONS: %+v\n", m.Questions)

	for i := uint16(0); i < m.Header.ANCount; i++ {
		a := DNSAnswer{}
		n, err = a.Parse(data, n)
		if err != nil {
			return err
		}

		m.Answers = append(m.Answers, a)
	}

	return nil
}
