package minecraftstatus

import (
	"encoding/json"
	"github.com/geNAZt/minecraft-status/data"
	"github.com/geNAZt/minecraft-status/protocol"
	"github.com/tatsushid/go-fastping"
	"net"
	"reflect"
	"sync"
	"time"
)

var (
	pinger *fastping.Pinger
	lock   sync.Mutex
)

func init() {
	pinger = fastping.NewPinger()
}

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
	conn.SendClientStatusPing()
	_, errPingPacket := conn.ReadPacket()
	if errPingPacket != nil {
		return nil, errPingPacket
	}

	// Get the ping
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
	status.Favicons = []data.Favicon{}

	if animatedFavicon || status.Players.Online == 0 {
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
			} else if status.Players.Online != 0 {
				break
			}
		}
	}

	conn.Close()

	return status, nil
}

func getPing(conn *protocol.Conn) (time.Duration, error) {
	lock.Lock()
	defer lock.Unlock()

	// Parse the IPAddr
	ipAddr, errIp := net.ResolveIPAddr("ip4", conn.IP)
	if errIp != nil {
		return 0, errIp
	}

	// Only ping on IP at a time
	pinger.AddIPAddr(ipAddr)

	// When a ping response got back
	ch := make(chan time.Duration, 1)
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		ch <- rtt
	}

	// When timeout
	pinger.OnIdle = func() {
		ch <- time.Second
	}

	// Run
	err := pinger.Run()
	if err != nil {
		return 0, err
	}

	chOut := <-ch
	return chOut, nil
}
