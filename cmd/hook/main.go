package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net"
    "os"
    "strings"
    "time"

    "safeexec/pkg/ipc"
)

func main() {
    if len(os.Args) < 2 {
        os.Exit(0) 
    }

    command := strings.Join(os.Args[1:], " ")

    // 10ms fail-open dial timeout over Unix socket
    conn, err := net.DialTimeout("unix", ipc.SocketPath, 10*time.Millisecond)
    if err != nil {
        os.Exit(0)
    }
    defer conn.Close()

    conn.SetDeadline(time.Now().Add(60 * time.Second))

    cwd, _ := os.Getwd()
    req := ipc.Request{
        Command:   command,
        SessionID: os.Getenv("TERM_SESSION_ID"),
        Cwd:       cwd,
    }

    out, err := json.Marshal(req)
    if err != nil {
        os.Exit(0)
    }

    _, err = conn.Write(append(out, '\n'))
    if err != nil {
        os.Exit(0)
    }

    scanner := bufio.NewScanner(conn)
    if scanner.Scan() {
        var res ipc.Response
        if err := json.Unmarshal(scanner.Bytes(), &res); err == nil {
            if res.Action == "BLOCK" {
                fmt.Fprintf(os.Stderr, "\033[91m\033[1m[SafeExec Blocked]\033[0m \033[91m%s\033[0m\n", res.Message)
                if res.MitreID != "" && res.MitreID != "None" {
                    fmt.Fprintf(os.Stderr, "\033[93m► MITRE ATT&CK: %s\033[0m\n", res.MitreID)
                    fmt.Fprintf(os.Stderr, "\033[93m► Blast Radius: %s\033[0m\n", res.BlastRadius)
                }
                os.Exit(1) 
            }
            if res.Action == "WARN" {
                fmt.Fprintf(os.Stderr, "\033[93m[SafeExec Warning]\033[0m %s\n", res.Message)
            }
        }
    }
    os.Exit(0)
}