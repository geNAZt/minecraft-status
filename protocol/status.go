package protocol

import "time"

type StatusResponse struct {
	Data string
}

type StatusPing struct {
	Time int64
}

type StatusGet struct {
}

type ClientStatusPing struct {
	Time int64
}

func (c *Conn) SendStatusRequest() {
	c.WritePacket(StatusGet{})
}

func (c *Conn) SendClientStatusPing() {
	c.WritePacket(ClientStatusPing{
		Time: time.Now().Unix(),
	})
}
