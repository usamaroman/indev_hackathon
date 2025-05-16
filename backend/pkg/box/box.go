package box

import (
	"fmt"
	"net"

	api "github.com/usamaroman/demo_indev_hackathon/backend/api/proto"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(host, port string) (*Client, error) {
	deviceAddr := fmt.Sprintf("%s:%s", host, port) // "192.168.1.100:7000"

	// Connect to the device
	conn, err := net.Dial("tcp", deviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v\n", err)
	}

	fmt.Println("Connected to device")

	b := &Client{
		addr: deviceAddr,
		conn: conn,
	}

	return b, nil
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

func (b *Client) Close() {
	b.conn.Close()
}
