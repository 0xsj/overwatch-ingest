package model

import (
	"github.com/0xsj/overwatch-pkg/types"
)

type SourceReliability struct {
	sourceID            types.ID
	tenantID            types.Optional[types.ID]
	reliabilityScore    float64
	totalRecords        int64
	acceptedRecords     int64
	rejectedRecords     int64
	quarantinedRecords  int64
	corroboratedRecords int64
	disputedRecords     int64
	calculatedAt        types.Timestamp
	windowStart         types.Timestamp
	windowEnd           types.Timestamp
}

func NewSourceReliability(
	sourceID types.ID,
	tenantID types.Optional[types.ID],
) *SourceReliability {
	now := types.Now()
	return &SourceReliability{
		sourceID:            sourceID,
		tenantID:            tenantID,
		reliabilityScore:    0.5,
		totalRecords:        0,
		acceptedRecords:     0,
		rejectedRecords:     0,
		quarantinedRecords:  0,
		corroboratedRecords: 0,
		disputedRecords:     0,
		calculatedAt:        now,
		windowStart:         now,
		windowEnd:           now,
	}
}

func ReconstructSourceReliability(
	sourceID types.ID,
	tenantID types.Optional[types.ID],
	reliabilityScore float64,
	totalRecords int64,
	acceptedRecords int64,
	rejectedRecords int64,
	quarantinedRecords int64,
	corroboratedRecords int64,
	disputedRecords int64,
	calculatedAt types.Timestamp,
	windowStart types.Timestamp,
	windowEnd types.Timestamp,
) *SourceReliability {
	return &SourceReliability{
		sourceID:            sourceID,
		tenantID:            tenantID,
		reliabilityScore:    reliabilityScore,
		totalRecords:        totalRecords,
		acceptedRecords:     acceptedRecords,
		rejectedRecords:     rejectedRecords,
		quarantinedRecords:  quarantinedRecords,
		corroboratedRecords: corroboratedRecords,
		disputedRecords:     disputedRecords,
		calculatedAt:        calculatedAt,
		windowStart:         windowStart,
		windowEnd:           windowEnd,
	}
}

func (s *SourceReliability) SourceID() types.ID                 { return s.sourceID }
func (s *SourceReliability) TenantID() types.Optional[types.ID] { return s.tenantID }
func (s *SourceReliability) ReliabilityScore() float64          { return s.reliabilityScore }
func (s *SourceReliability) TotalRecords() int64                { return s.totalRecords }
func (s *SourceReliability) AcceptedRecords() int64             { return s.acceptedRecords }
func (s *SourceReliability) RejectedRecords() int64             { return s.rejectedRecords }
func (s *SourceReliability) QuarantinedRecords() int64          { return s.quarantinedRecords }
func (s *SourceReliability) CorroboratedRecords() int64         { return s.corroboratedRecords }
func (s *SourceReliability) DisputedRecords() int64             { return s.disputedRecords }
func (s *SourceReliability) CalculatedAt() types.Timestamp      { return s.calculatedAt }
func (s *SourceReliability) WindowStart() types.Timestamp       { return s.windowStart }
func (s *SourceReliability) WindowEnd() types.Timestamp         { return s.windowEnd }

func (s *SourceReliability) RecordAccepted() {
	s.totalRecords++
	s.acceptedRecords++
	s.recalculateScore()
}

func (s *SourceReliability) RecordRejected() {
	s.totalRecords++
	s.rejectedRecords++
	s.recalculateScore()
}

func (s *SourceReliability) RecordQuarantined() {
	s.totalRecords++
	s.quarantinedRecords++
	s.recalculateScore()
}

func (s *SourceReliability) RecordCorroborated() {
	s.corroboratedRecords++
	s.recalculateScore()
}

func (s *SourceReliability) RecordDisputed() {
	s.disputedRecords++
	s.recalculateScore()
}

func (s *SourceReliability) recalculateScore() {
	s.reliabilityScore = s.calculateReliabilityScore()
	s.calculatedAt = types.Now()
	s.windowEnd = s.calculatedAt
}

func (s *SourceReliability) calculateReliabilityScore() float64 {
	if s.totalRecords == 0 {
		return 0.5
	}

	acceptanceRate := float64(s.acceptedRecords) / float64(s.totalRecords)

	corroborationRate := 0.0
	corroboratedTotal := s.corroboratedRecords + s.disputedRecords
	if corroboratedTotal > 0 {
		corroborationRate = float64(s.corroboratedRecords) / float64(corroboratedTotal)
	} else {
		corroborationRate = 0.5
	}

	const (
		weightAcceptance    = 0.7
		weightCorroboration = 0.3
	)

	score := (acceptanceRate * weightAcceptance) + (corroborationRate * weightCorroboration)

	return clampScore(score)
}

func (s *SourceReliability) AcceptanceRate() float64 {
	if s.totalRecords == 0 {
		return 0.0
	}
	return float64(s.acceptedRecords) / float64(s.totalRecords)
}

func (s *SourceReliability) RejectionRate() float64 {
	if s.totalRecords == 0 {
		return 0.0
	}
	return float64(s.rejectedRecords) / float64(s.totalRecords)
}

func (s *SourceReliability) QuarantineRate() float64 {
	if s.totalRecords == 0 {
		return 0.0
	}
	return float64(s.quarantinedRecords) / float64(s.totalRecords)
}

func (s *SourceReliability) IsReliable(threshold float64) bool {
	return s.reliabilityScore >= threshold
}

func (s *SourceReliability) IsUnreliable(threshold float64) bool {
	return s.reliabilityScore < threshold
}

func (s *SourceReliability) HasSufficientData(minRecords int64) bool {
	return s.totalRecords >= minRecords
}

func (s *SourceReliability) ResetWindow(windowStart types.Timestamp) {
	s.totalRecords = 0
	s.acceptedRecords = 0
	s.rejectedRecords = 0
	s.quarantinedRecords = 0
	s.corroboratedRecords = 0
	s.disputedRecords = 0
	s.reliabilityScore = 0.5
	s.windowStart = windowStart
	s.windowEnd = windowStart
	s.calculatedAt = types.Now()
}
