package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/0xsj/overwatch-pkg/provenance"
)

// Verifier is a mock implementation of provenance.Verifier.
type Verifier struct {
	mu sync.RWMutex

	Calls struct {
		Verify              int
		VerifyBase64        int
		VerifySignatureInfo int
		SupportsDID         int
	}

	Errors struct {
		Verify              error
		VerifyBase64        error
		VerifySignatureInfo error
	}

	// SupportsDIDResult controls the return value of SupportsDID.
	SupportsDIDResult bool
}

func NewVerifier() *Verifier {
	return &Verifier{
		SupportsDIDResult: true,
	}
}

func (v *Verifier) Verify(_ context.Context, _ string, _ []byte, _ []byte) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.Verify++
	return v.Errors.Verify
}

func (v *Verifier) VerifyBase64(_ context.Context, _ string, _ []byte, _ string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.VerifyBase64++
	return v.Errors.VerifyBase64
}

func (v *Verifier) VerifySignatureInfo(_ context.Context, _ []byte, _ *provenance.SignatureInfo) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.VerifySignatureInfo++
	return v.Errors.VerifySignatureInfo
}

func (v *Verifier) SupportsDID(_ string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.SupportsDID++
	return v.SupportsDIDResult
}

// SetVerifyError configures the verifier to return an error on all verify calls.
func (v *Verifier) SetVerifyError(err error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Errors.Verify = err
	v.Errors.VerifyBase64 = err
	v.Errors.VerifySignatureInfo = err
}

// SetVerifySuccess configures the verifier to succeed on all verify calls.
func (v *Verifier) SetVerifySuccess() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Errors.Verify = nil
	v.Errors.VerifyBase64 = nil
	v.Errors.VerifySignatureInfo = nil
}

// ─────────────────────────────────────────────────────────────────
// EnvelopeBuilder (mock via nil signer - handler handles nil gracefully)
// ─────────────────────────────────────────────────────────────────

// NilEnvelopeBuilder returns nil, which causes the handler to skip signing.
// The handler checks if signer == nil and returns early.
func NilEnvelopeBuilder() *provenance.EnvelopeBuilder {
	return nil
}

// MockSignatureInfo creates a mock SignatureInfo for testing.
func MockSignatureInfo(signerType provenance.SignerType) *provenance.SignatureInfo {
	return &provenance.SignatureInfo{
		DID:        "did:key:z6MkTestMockSigner",
		SignerType: signerType,
		Signature:  "bW9ja3NpZ25hdHVyZQ==",
		SignedAt:   time.Now().UTC(),
	}
}
