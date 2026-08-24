package tier1

import (
	"strings"

	"safeexec/pkg/parser"
)


type FeatureExtractor struct{}


func (f *FeatureExtractor) Extract(cmd *parser.ParsedCommand) []float32 {
	features := make([]float32, 5)

	features[0] = float32(strings.Count(cmd.Original, "|"))


	features[1] = float32(strings.Count(cmd.Original, " -"))

	if strings.HasPrefix(strings.TrimSpace(cmd.Original), "sudo ") {
		features[2] = 1.0
	} else {
		features[2] = 0.0
	}

	networkTerms := []string{"curl", "wget", "nc", "netcat", "/dev/tcp"}
	for _, term := range networkTerms {
		if strings.Contains(cmd.Original, term) {
			features[3] += 1.0
		}
	}

	
	obfuscationScore := 0
	obfuscationScore += strings.Count(cmd.Original, "\\")
	obfuscationScore += strings.Count(cmd.Original, "base64")
	features[4] = float32(obfuscationScore)

	return features
}
