package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"

    "github.com/joho/godotenv"

    "safeexec/pkg/cache"
    "safeexec/pkg/ipc"
    "safeexec/pkg/parser"
    "safeexec/pkg/tier0"
    "safeexec/pkg/tier1"
    "safeexec/pkg/tier2"
)

func main() {
    fmt.Println("Starting SafeExec Daemon (Tier 0 -> Tier 1 -> Cache -> Tier 2)...")

    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, relying on system environment variables.")
    }

    t0Engine := tier0.NewEngine()
    t1Scorer := tier1.NewScorer()

    dbPath := filepath.Join(".", "safeexec_cache.db")
    dbCache, err := cache.NewSQLiteCache(dbPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[FATAL] Failed to initialize SQLite cache at %s: %v\n", dbPath, err)
        os.Exit(1)
    }

    // Graceful Linux Shutdown Handling
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        fmt.Println("\n[Daemon] Caught termination signal. Shutting down cleanly...")
        dbCache.Close()
        os.Remove(ipc.SocketPath)
        os.Exit(0)
    }()

    handler := func(req ipc.Request) ipc.Response {
        fmt.Printf("\n[Daemon] Session: %s | Cwd: %s\n", req.SessionID, req.Cwd)
        fmt.Printf("[Daemon] Received: %s\n", req.Command)

        parsed, err := parser.Parse(req.Command)
        if err != nil {
            fmt.Printf("[Parser] Error: %v. Failing open.\n", err)
            return ipc.Response{Action: "ALLOW", Message: "Parse error, failing open"}
        }

        t0Action, t0Message := t0Engine.Evaluate(parsed)
        if t0Action == tier0.ActionAllow {
            fmt.Println("[Tier 0] 🟢 FAST ALLOW")
            return ipc.Response{Action: "ALLOW", Message: ""}
        } else if t0Action == tier0.ActionBlock {
            fmt.Printf("[Tier 0] 🔴 BLOCKED: %s\n", t0Message)
            return ipc.Response{Action: "BLOCK", Message: t0Message}
        }

        fmt.Println("[Tier 0] ⚪ INCONCLUSIVE. Escalating to Tier 1 ML...")

        isConclusive, t1Action, t1Message := t1Scorer.Evaluate(parsed)
        if isConclusive {
            if t1Action == "ALLOW" {
                fmt.Println("[Tier 1] 🟢 ML ALLOW")
                return ipc.Response{Action: "ALLOW", Message: ""}
            } else if t1Action == "BLOCK" {
                fmt.Printf("[Tier 1] 🔴 ML BLOCKED: %s\n", t1Message)
                return ipc.Response{Action: "BLOCK", Message: t1Message}
            }
        }

        fmt.Println("[Tier 1] ⚪ INCONCLUSIVE (The Gray Zone). Escalating to Tier 2 (LLM)...")

        hash := cache.HashCanonical(parsed.Canonical)
        cachedVerdict, err := dbCache.Get(hash)
        if err == nil && cachedVerdict != nil {
            fmt.Printf("[Cache] HIT for hash %s (Canonical: %s)\n", hash[:8], parsed.Canonical)
            finalMessage := cache.Interpolate(cachedVerdict.Message, parsed.Original, parsed.Canonical)
            return ipc.Response{
                Action:      cachedVerdict.Action,
                Message:     finalMessage,
                MitreID:     cachedVerdict.MitreID,
                BlastRadius: cachedVerdict.BlastRadius,
            }
        }

        fmt.Printf("[Cache] MISS for hash %s. Escalating to Tier 2 AI...\n", hash[:8])

        apiKey := os.Getenv("API_KEY")
        if apiKey == "" {
            fmt.Println("[Tier 2] WARNING: API_KEY not set. Failing open for ambiguity.")
            return ipc.Response{Action: "ALLOW", Message: "Ambiguous command allowed (LLM disabled)."}
        }

        llmResp, err := tier2.Evaluate(parsed.Original, apiKey)
        if err != nil {
            fmt.Printf("[Tier 2] AI Error: %v\n", err)
            return ipc.Response{Action: "ALLOW", Message: "Ambiguous command allowed (LLM error)."}
        }

        llmVerdict := cache.Verdict{
            Action:      llmResp.Action,
            Message:     llmResp.Message,
            MitreID:     llmResp.MitreID,     
            BlastRadius: llmResp.BlastRadius, 
        }
        
        dbCache.Put(hash, llmVerdict)
        finalMessage := cache.Interpolate(llmVerdict.Message, parsed.Original, parsed.Canonical)

        return ipc.Response{
            Action:      llmVerdict.Action,
            Message:     finalMessage,
            MitreID:     llmResp.MitreID,
            BlastRadius: llmResp.BlastRadius,
        }
    }
    
    if err := ipc.StartServer(handler); err != nil {
        fmt.Fprintf(os.Stderr, "[FATAL] Daemon IPC server crashed: %v\n", err)
        os.Exit(1)
    }
}