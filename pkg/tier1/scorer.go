package tier1

import (
	"fmt"
	"safeexec/pkg/parser"
)


const Threshold float32 = 0.75

type Scorer struct {
	extractor *FeatureExtractor
	model     *Model
}

func NewScorer() *Scorer {
	return &Scorer{
		extractor: &FeatureExtractor{},
		model:     &Model{},
	}
}


func (s *Scorer) Evaluate(cmd *parser.ParsedCommand) (bool, string, string) {
	features := s.extractor.Extract(cmd)
	probability := s.model.Score(features)

	fmt.Printf("[Tier 1] ML Features: %v | Score: %.2f\n", features, probability)

	if probability >= Threshold {
		
		msg := fmt.Sprintf("ML Risk Scorer flagged command as highly dangerous (Confidence: %d%%)", int(probability*100))
		return true, "BLOCK", msg
	} else if probability < 0.40 {
		
		return true, "ALLOW", ""
	}

	
	return false, "SKIP", ""
}
