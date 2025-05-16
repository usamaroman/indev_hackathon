package main

import (
	"fmt"
	"net"
	"time"

	api "github.com/usamaroman/demo_indev_hackathon/backend/api/proto"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Device address
	deviceAddr := "192.168.1.100:7000"

	// Connect to the device
	conn, err := net.Dial("tcp", deviceAddr)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to device")

	// Example 1: Send GetInfo request
	sendGetInfo(conn)

	// Wait a bit for response
	time.Sleep(1 * time.Second)

	// Example 2: Send SetState request to turn light on
	sendSetState(conn, api.States_LightOn)

	// Example 3: Send GetState request
	sendGetState(conn)
}

func sendGetInfo(conn net.Conn) {
	// Create GetInfo message
	msg := &api.ClientMessage{
		Message: &api.ClientMessage_GetInfo{
			GetInfo: &api.GetInfo{},
		},
	}

	// Marshal to protobuf
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
		return
	}

	// Send to device
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Failed to send: %v\n", err)
		return
	}

	fmt.Println("Sent GetInfo request")
	readResponse(conn)
}

func sendSetState(conn net.Conn, state api.States) {
	// Create SetState message
	msg := &api.ClientMessage{
		Message: &api.ClientMessage_SetState{
			SetState: &api.SetState{
				State: state,
			},
		},
	}

	// Marshal to protobuf
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
		return
	}

	// Send to device
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Failed to send: %v\n", err)
		return
	}

	fmt.Printf("Sent SetState request (%v)\n", state)
	readResponse(conn)
}

func sendGetState(conn net.Conn) {
	// Create GetState message
	msg := &api.ClientMessage{
		Message: &api.ClientMessage_GetState{
			GetState: &api.GetState{},
		},
	}

	// Marshal to protobuf
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
		return
	}

	// Send to device
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Failed to send: %v\n", err)
		return
	}

	fmt.Println("Sent GetState request")
	readResponse(conn)
}

func readResponse(conn net.Conn) {
	// Create buffer for response
	buf := make([]byte, 1024)

	// Set read timeout
	err := conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		fmt.Printf("SetReadDeadline failed: %v\n", err)
		return
	}

	// Read response
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	// Unmarshal response
	response := &api.ControllerResponse{}
	err = proto.Unmarshal(buf[:n], response)
	if err != nil {
		fmt.Printf("Failed to unmarshal response: %v\n", err)
		return
	}

	// Handle different response types
	switch r := response.Response.(type) {
	case *api.ControllerResponse_Info:
		info := r.Info
		fmt.Printf("Received Info response:\n")
		fmt.Printf("IP: %s\n", info.Ip)
		fmt.Printf("MAC: %s\n", info.Mac)
		fmt.Printf("BLE Name: %s\n", info.BleName)
		fmt.Printf("Token: %s\n", info.Token)
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
		fmt.Printf("Received Status: %v\n", r.Status)
	default:
		fmt.Printf("Received unknown response type\n")
	}
}
