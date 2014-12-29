package protocol

type Handshake struct {
	ProtocolVersion VarInt
	Address         string
	Port            uint16
	State           VarInt
}

func (c *Conn) SendHandshake() {
	c.WritePacket(Handshake{
		ProtocolVersion: 4,
		Address:         c.Host,
		Port:            c.Port,
		State:           1,
	})
}
