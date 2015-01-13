package minecraftstatus

import (
	"encoding/json"
	"github.com/geNAZt/minecraft-status/data"
	"github.com/geNAZt/minecraft-status/protocol"
	"reflect"
	"time"
)

func GetStatus(host string, animatedFavicon bool) (*data.Status, error) {
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
	var pingTime time.Duration
	starttime := time.Now()

	conn.SendClientStatusPing()
	_, errPingPacket := conn.ReadPacket()
	if errPingPacket != nil {
		return nil, errPingPacket
	}

	pingTime = time.Now().Sub(starttime)

	ch := make(chan bool, 1)
	go func() {
		// Get the ping
		tempPingTime, errPing := Ping(conn)
		if errPing == nil && tempPingTime > 0 {
			pingTime = tempPingTime
		}

		ch <- true
	}()

	// Parse the status
	status := &data.Status{}
	errJson := json.Unmarshal([]byte(statusPacket.(protocol.StatusResponse).Data), status)
	if errJson != nil {
		return nil, errJson
	}

	status.Ping = pingTime
	status.Favicons = []data.Favicon{}

	// Wait for additional Favicons (animated motd hack)
	for {
		starttime := time.Now()
		additionalStatusPacket, errPacket := conn.ReadPacket()
		if errPacket != nil {
			break
		}

		// Check if packet is correct
		if !reflect.TypeOf(additionalStatusPacket).AssignableTo(reflect.TypeOf(protocol.StatusResponse{})) {
			continue
		}

		// Parse the status
		additionalStatus := &data.Status{}
		errJson := json.Unmarshal([]byte(additionalStatusPacket.(protocol.StatusResponse).Data), additionalStatus)
		if errJson != nil {
			continue
		}

		status = additionalStatus

		if animatedFavicon {
			favicon := data.Favicon{
				Icon:        status.Favicon,
				DisplayTime: int32(time.Now().Sub(starttime) / time.Millisecond),
			}

			status.Favicons = append(status.Favicons, favicon)
			status.Favicon = additionalStatus.Favicon
		}
	}

	conn.Close()

	<-ch
	return status, nil
}
