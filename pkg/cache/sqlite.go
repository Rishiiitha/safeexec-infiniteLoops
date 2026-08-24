package cache

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	
	_ "modernc.org/sqlite"
)


type Verdict struct {
	Action      string
	Message     string 
	MitreID     string 
	BlastRadius string 
}

type cacheEntry struct {
	Hash    string
	Verdict Verdict
}

type SQLiteCache struct {
	db         *sql.DB
	writeQueue chan cacheEntry
}

func NewSQLiteCache(dbPath string) (*SQLiteCache, error) {
	
	
	db, err := sql.Open("sqlite", "debug_cache.db")
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")

	query := `
	CREATE TABLE IF NOT EXISTS tier2_cache (
		hash TEXT PRIMARY KEY,
		action TEXT,
		message TEXT,
		mitre_id TEXT,
		blast_radius TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	c := &SQLiteCache{
		db:         db,
		writeQueue: make(chan cacheEntry, 1000),
	}

	go c.writeBehindWorker()

	return c, nil
}

func HashCanonical(canonical string) string {
	hasher := sha256.New()
	hasher.Write([]byte(canonical))
	return hex.EncodeToString(hasher.Sum(nil))
}


func (c *SQLiteCache) Get(canonicalHash string) (*Verdict, error) {
	var v Verdict
	row := c.db.QueryRow("SELECT action, message, mitre_id, blast_radius FROM tier2_cache WHERE hash = ?", canonicalHash)
	err := row.Scan(&v.Action, &v.Message, &v.MitreID, &v.BlastRadius)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}
	return &v, nil 
}


func (c *SQLiteCache) Put(canonicalHash string, verdict Verdict) {
	select {
	case c.writeQueue <- cacheEntry{Hash: canonicalHash, Verdict: verdict}:
	default:
		
		log.Printf("[Cache] Write queue full, dropping cache write for %s", canonicalHash)
	}
}

func (c *SQLiteCache) writeBehindWorker() {
	for entry := range c.writeQueue {
		_, err := c.db.Exec(
			"INSERT OR REPLACE INTO tier2_cache (hash, action, message, mitre_id, blast_radius) VALUES (?, ?, ?, ?, ?)",
			entry.Hash, entry.Verdict.Action, entry.Verdict.Message, entry.Verdict.MitreID, entry.Verdict.BlastRadius,
		)
		if err != nil {
			log.Printf("[Cache] Async write failed: %v", err)
		}
	}
}


func Interpolate(templateMsg string, originalCmd string, canonicalCmd string) string {
	
	if strings.Contains(templateMsg, "<ARG>") {
		
		return strings.ReplaceAll(templateMsg, "<ARG>", "[target]") + fmt.Sprintf(" (Context: %s)", originalCmd)
	}
	return templateMsg
}

func (c *SQLiteCache) Close() error {
	close(c.writeQueue)
	return c.db.Close()
}