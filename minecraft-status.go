package minecraftstatus

import (
	"encoding/json"
	"github.com/geNAZt/minecraft-status/data"
	"github.com/geNAZt/minecraft-status/protocol"
	"time"
)

func GetStatus(host string) (*data.Status, error) {
	// Create Client
	conn, err := protocol.NewNetClient(host)
	if err != nil {
		return nil, err
	}

	// Send handshake and switch protocol state
	conn.SendHandshake()
	conn.State = protocol.Status

	// Send the request so the server sends us a status
	conn.SendStatusRequest()
	statusPacket, errPacket := conn.ReadPacket()
	if errPacket != nil {
		return nil, errPacket
	}

	// Now we need to send a Ping
	pingTime, errPing := getPing(conn)
	if errPing != nil {
		return nil, errPing
	}

	// Parse the status
	status := &data.Status{}
	errJson := json.Unmarshal([]byte(statusPacket.(protocol.StatusResponse).Data), status)
	if errJson != nil {
		return nil, errJson
	}

	status.Ping = pingTime

	return status, nil
}

func getPing(conn *protocol.Conn) (time.Duration, error) {
	starttime := time.Now()

	conn.SendClientStatusPing()
	pingPacket, errPingPacket := conn.ReadPacket()
	if errPingPacket != nil && pingPacket != nil {
		return 0, errPingPacket
	}

	return time.Now().Sub(starttime), nil
}
