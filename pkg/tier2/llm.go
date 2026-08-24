package tier2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LLMResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MitreID     string `json:"mitre_id"`
	BlastRadius string `json:"blast_radius"`
}



type OpenAIPayload struct {
	Model          string         `json:"model"`
	ResponseFormat ResponseFormat `json:"response_format"`
	Messages       []Message      `json:"messages"`
	Temperature    float32        `json:"temperature"`
	MaxTokens      int            `json:"max_tokens"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func Evaluate(command string, apiKey string) (*LLMResponse, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is not set")
	}


		systemPrompt := `You are a strict Linux security EDR agent evaluating an ambiguous shell command.
	Analyze the command and determine if it is destructive, overly permissive, or malicious.
	Keep your explanation extremely concise (1 to 2 sentences maximum).

	If the command is malicious, you MUST map it to a specific MITRE ATT&CK Technique ID (e.g., T1105 for Ingress Tool Transfer, T1059 for Command and Scripting Interpreter). 
	Estimate the Blast Radius (e.g., "Local Directory", "User Privilege", "System-Wide", "Network Compromise").
	If the command is ALLOWED, set mitre_id and blast_radius to "None".

	Output valid JSON in this exact structure:
	{
	"action": "ALLOW" or "BLOCK",
	"message": "Explanation of the danger and a safer alternative.",
	"mitre_id": "TXXXX",
	"blast_radius": "Scope of potential damage"
	}`

	payload := OpenAIPayload{
		Model:       "openai/gpt-oss-20b", 
		Temperature: 0.1,
		MaxTokens:   500,
		ResponseFormat: ResponseFormat{
			Type: "json_object",
		},
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("Command to evaluate:\n%s", command)},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Groq request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Groq returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("empty response from Groq")
	}

	rawContent := result.Choices[0].Message.Content

	rawContent = strings.TrimSpace(rawContent)
	rawContent = strings.TrimPrefix(rawContent, "```json")
	rawContent = strings.TrimPrefix(rawContent, "```")
	rawContent = strings.TrimSuffix(rawContent, "```")
	rawContent = strings.TrimSpace(rawContent)

	var llmResp LLMResponse
	if err := json.Unmarshal([]byte(rawContent), &llmResp); err != nil {
		return nil, fmt.Errorf("failed to parse Groq JSON output: %v\nRaw output: %s", err, rawContent)
	}

	return &llmResp, nil
}