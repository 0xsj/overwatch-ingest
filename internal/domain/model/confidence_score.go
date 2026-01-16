package model

type ConfidenceScore struct {
	overall           float64
	sourceReliability float64
	dataCompleteness  float64
	temporalFreshness float64
	signatureTrust    float64
	factors           []ConfidenceFactor
}

type ConfidenceFactor struct {
	name   string
	score  float64
	weight float64
	reason string
}

func NewConfidenceScore(
	sourceReliability float64,
	dataCompleteness float64,
	temporalFreshness float64,
	signatureTrust float64,
	factors []ConfidenceFactor,
) ConfidenceScore {
	overall := calculateOverallScore(
		sourceReliability,
		dataCompleteness,
		temporalFreshness,
		signatureTrust,
	)

	return ConfidenceScore{
		overall:           overall,
		sourceReliability: sourceReliability,
		dataCompleteness:  dataCompleteness,
		temporalFreshness: temporalFreshness,
		signatureTrust:    signatureTrust,
		factors:           factors,
	}
}

func ReconstructConfidenceScore(
	overall float64,
	sourceReliability float64,
	dataCompleteness float64,
	temporalFreshness float64,
	signatureTrust float64,
	factors []ConfidenceFactor,
) ConfidenceScore {
	return ConfidenceScore{
		overall:           overall,
		sourceReliability: sourceReliability,
		dataCompleteness:  dataCompleteness,
		temporalFreshness: temporalFreshness,
		signatureTrust:    signatureTrust,
		factors:           factors,
	}
}

func (c ConfidenceScore) Overall() float64            { return c.overall }
func (c ConfidenceScore) SourceReliability() float64  { return c.sourceReliability }
func (c ConfidenceScore) DataCompleteness() float64   { return c.dataCompleteness }
func (c ConfidenceScore) TemporalFreshness() float64  { return c.temporalFreshness }
func (c ConfidenceScore) SignatureTrust() float64     { return c.signatureTrust }
func (c ConfidenceScore) Factors() []ConfidenceFactor { return c.factors }

func (c ConfidenceScore) IsAboveThreshold(threshold float64) bool {
	return c.overall >= threshold
}

func (c ConfidenceScore) IsBelowThreshold(threshold float64) bool {
	return c.overall < threshold
}

func (c ConfidenceScore) IsAcceptable(acceptThreshold float64) bool {
	return c.overall >= acceptThreshold
}

func (c ConfidenceScore) ShouldQuarantine(acceptThreshold, rejectThreshold float64) bool {
	return c.overall < acceptThreshold && c.overall >= rejectThreshold
}

func (c ConfidenceScore) ShouldReject(rejectThreshold float64) bool {
	return c.overall < rejectThreshold
}

func (c ConfidenceScore) WithFactor(factor ConfidenceFactor) ConfidenceScore {
	newFactors := make([]ConfidenceFactor, len(c.factors)+1)
	copy(newFactors, c.factors)
	newFactors[len(c.factors)] = factor

	return ConfidenceScore{
		overall:           c.overall,
		sourceReliability: c.sourceReliability,
		dataCompleteness:  c.dataCompleteness,
		temporalFreshness: c.temporalFreshness,
		signatureTrust:    c.signatureTrust,
		factors:           newFactors,
	}
}

func NewConfidenceFactor(name string, score float64, weight float64, reason string) ConfidenceFactor {
	return ConfidenceFactor{
		name:   name,
		score:  clampScore(score),
		weight: weight,
		reason: reason,
	}
}

func (f ConfidenceFactor) Name() string    { return f.name }
func (f ConfidenceFactor) Score() float64  { return f.score }
func (f ConfidenceFactor) Weight() float64 { return f.weight }
func (f ConfidenceFactor) Reason() string  { return f.reason }

func (f ConfidenceFactor) WeightedScore() float64 {
	return f.score * f.weight
}

const (
	weightSourceReliability = 0.30
	weightDataCompleteness  = 0.25
	weightTemporalFreshness = 0.20
	weightSignatureTrust    = 0.25
)

func calculateOverallScore(
	sourceReliability float64,
	dataCompleteness float64,
	temporalFreshness float64,
	signatureTrust float64,
) float64 {
	sourceReliability = clampScore(sourceReliability)
	dataCompleteness = clampScore(dataCompleteness)
	temporalFreshness = clampScore(temporalFreshness)
	signatureTrust = clampScore(signatureTrust)

	overall := (sourceReliability * weightSourceReliability) +
		(dataCompleteness * weightDataCompleteness) +
		(temporalFreshness * weightTemporalFreshness) +
		(signatureTrust * weightSignatureTrust)

	return clampScore(overall)
}

func clampScore(score float64) float64 {
	if score < 0.0 {
		return 0.0
	}
	if score > 1.0 {
		return 1.0
	}
	return score
}

func ZeroConfidenceScore(reason string) ConfidenceScore {
	return ConfidenceScore{
		overall:           0.0,
		sourceReliability: 0.0,
		dataCompleteness:  0.0,
		temporalFreshness: 0.0,
		signatureTrust:    0.0,
		factors: []ConfidenceFactor{
			NewConfidenceFactor("zero_score", 0.0, 1.0, reason),
		},
	}
}

func DefaultConfidenceScore() ConfidenceScore {
	return NewConfidenceScore(
		0.5,
		1.0,
		1.0,
		1.0,
		nil,
	)
}
