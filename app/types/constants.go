package types

type TYPE uint16

const (
	A     TYPE = 1
	NS    TYPE = 2
	MD    TYPE = 3
	MF    TYPE = 4
	CNAME TYPE = 5
	SOA   TYPE = 6
	MB    TYPE = 7
	MG    TYPE = 8
	MR    TYPE = 9
	NULL  TYPE = 10
	WKS   TYPE = 11
	PTR   TYPE = 12
	HINFO TYPE = 13
	MINFO TYPE = 14
	MX    TYPE = 15
	TXT   TYPE = 16
)

type QTYPE TYPE

const (
	AXFR  QTYPE = 252
	MAILB QTYPE = 253
	MAILA QTYPE = 254
	ALL   QTYPE = 255
)

type CLASS uint16

const (
	IN CLASS = 1
	CS CLASS = 2
	CH CLASS = 3
	HS CLASS = 4
)

type QCLASS CLASS

const (
	ANY QCLASS = 255
)

type OPCODE byte

const (
	QUERY  OPCODE = 0
	IQUERY OPCODE = 1
	STATUS OPCODE = 2
)

type RCODE byte

const (
	NOERROR  RCODE = 0
	FORMERR  RCODE = 1
	SERVFAIL RCODE = 2
	NXDOMAIN RCODE = 3
	NOTIMP   RCODE = 4
	REFUSED  RCODE = 5
)
