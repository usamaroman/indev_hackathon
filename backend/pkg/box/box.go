package box

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	api "github.com/usamaroman/demo_indev_hackathon/backend/api/proto"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/logger"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	log *slog.Logger

	addr string
	conn net.Conn

	bleName string
	mac     string
	token   string

	status string
}

func (b *Client) GetBleName() string {
	return b.bleName
}

func (b *Client) GetToken() string {
	return b.token
}

func New(log *slog.Logger, host, port string) (*Client, error) {
	log = log.With(slog.String("component", "box client"))

	deviceAddr := fmt.Sprintf("%s:%s", host, port) // "192.168.1.100:7000"

	b := &Client{
		log:  log,
		addr: deviceAddr,
	}

	conn, err := net.Dial("tcp", deviceAddr)
	if err != nil {
		b.log.Error("failed to connect", logger.Error(err))
		return nil, err
	}

	log.Info("connected to device", slog.String("addr", deviceAddr))

	b.conn = conn
	b.getInfo()

	return b, nil
}

func (b *Client) getInfo() error {
	msg := &api.ClientMessage{
		Message: &api.ClientMessage_GetInfo{
			GetInfo: &api.GetInfo{},
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		b.log.Error("failed to marshal", logger.Error(err))
		return err
	}

	_, err = b.conn.Write(data)
	if err != nil {
		b.log.Error("failed to send", logger.Error(err))
		return err
	}

	b.log.Debug("sent get info request")

	b.readResponse()

	return nil
}

func (b *Client) LightOn() error {
	msg := &api.ClientMessage{
		Message: &api.ClientMessage_SetState{
			SetState: &api.SetState{
				State: api.States_LightOn,
			},
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("failed to marshal: %v\n", err)
		return err
	}

	_, err = b.conn.Write(data)
	if err != nil {
		fmt.Printf("failed to send: %v\n", err)
		return err
	}

	return nil
}

func (b *Client) LightOff() error {
	msg := &api.ClientMessage{
		Message: &api.ClientMessage_SetState{
			SetState: &api.SetState{
				State: api.States_LightOff,
			},
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("failed to marshal: %v\n", err)
		return err
	}

	_, err = b.conn.Write(data)
	if err != nil {
		fmt.Printf("failed to send: %v\n", err)
		return err
	}

	return nil
}

func (b *Client) readResponse() {
	buf := make([]byte, 1024)

	err := b.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		b.log.Error("SetReadDeadline failed", logger.Error(err))
		return
	}

	n, err := b.conn.Read(buf)
	if err != nil {
		b.log.Error("failed to read response", logger.Error(err))
		return
	}

	response := &api.ControllerResponse{}
	err = proto.Unmarshal(buf[:n], response)
	if err != nil {
		b.log.Error("failed to unmarshal response", logger.Error(err))
		return
	}

	switch r := response.Response.(type) {
	case *api.ControllerResponse_Info:
		info := r.Info
		b.log.Debug("received Info response", slog.Any("info", info))

		fmt.Printf("IP: %s\n", info.Ip)
		fmt.Printf("MAC: %s\n", info.Mac)
		b.mac = info.Mac
		fmt.Printf("BLE Name: %s\n", info.BleName)
		b.bleName = info.BleName
		fmt.Printf("Token: %s\n", info.Token)
		b.token = info.Token
	case *api.ControllerResponse_State:
		state := r.State
		fmt.Printf("Received State response:\n")
		fmt.Printf("Light: %v\n", state.LightOn)
		fmt.Printf("Door Lock: %v\n", state.DoorLock)
		fmt.Printf("Channel 1: %v\n", state.Channel_1)
		fmt.Printf("Channel 2: %v\n", state.Channel_2)
		fmt.Printf("Temperature: %.2f\n", state.Temperature)
		fmt.Printf("Pressure: %.2f\n", state.Pressure)
		fmt.Printf("Humidity: %.2f\n", state.Humidity)
	case *api.ControllerResponse_Status:
		b.log.Debug("received Status response", slog.String("status", api.Statuses_name[int32(r.Status)]))
		b.status = api.Statuses_name[int32(r.Status)]

	default:
		b.log.Error("received unknown response type")
	}
}

func (b *Client) Close() {
	b.conn.Close()
}
