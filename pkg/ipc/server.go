package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Request struct {
	Command   string `json:"command"`
	SessionID string `json:"session_id"` 
	Cwd       string `json:"cwd"`
}

type Response struct {
	Action      string `json:"action"` 
	Message     string `json:"message"`
	MitreID     string `json:"mitre_id"`    
	BlastRadius string `json:"blast_radius"`
}

const Address = "127.0.0.1:9099"

func StartServer(handler func(req Request) Response) error {
	
	listener, err := net.Listen("tcp", Address)
	if err != nil {
		return fmt.Errorf("failed to start IPC server: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Daemon listening on tcp://%s\n", Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[IPC] Failed to accept connection: %v\n", err)
			continue
		}

		
		go handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler func(req Request) Response) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	scanner := bufio.NewScanner(conn)
	
	if scanner.Scan() {
		var req Request
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			fmt.Printf("[IPC] Invalid request payload: %v\n", err)
			return
		}

		res := handler(req)

		out, _ := json.Marshal(res)
		conn.Write(append(out, '\n'))
	}
}