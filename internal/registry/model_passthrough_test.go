package registry

import "testing"

func TestGetModelProvidersPassesThroughUnknownClaudeModel(t *testing.T) {
	r := newTestModelRegistry()
	r.RegisterClient("claude-1", "claude", []*ModelInfo{{ID: "claude-opus-5"}})

	providers := r.GetModelProviders("claude-opus-5-5")
	if len(providers) != 1 || providers[0] != "claude" {
		t.Fatalf("expected claude passthrough, got %v", providers)
	}
	if !r.ClientSupportsModel("claude-1", "claude-opus-5-5") {
		t.Fatal("expected claude client to accept the unregistered model")
	}
}

func TestGetModelProvidersRejectsUnknownClaudeModelWithoutClaudeClient(t *testing.T) {
	r := newTestModelRegistry()
	r.RegisterClient("gemini-1", "gemini", []*ModelInfo{{ID: "gemini-2.5-pro"}})

	if providers := r.GetModelProviders("claude-opus-5-5"); len(providers) != 0 {
		t.Fatalf("expected no providers without a claude credential, got %v", providers)
	}
	if r.ClientSupportsModel("gemini-1", "claude-opus-5-5") {
		t.Fatal("expected a non-claude client to reject the model")
	}
}

func TestGetModelProvidersDoesNotPassThroughNonClaudeModels(t *testing.T) {
	r := newTestModelRegistry()
	r.RegisterClient("claude-1", "claude", []*ModelInfo{{ID: "claude-opus-5"}})

	for _, modelID := range []string{
		"gpt-6-astra",
		"claude-fable-5-dd-reversed",
	} {
		if providers := r.GetModelProviders(modelID); len(providers) != 0 {
			t.Fatalf("expected no passthrough for %q, got %v", modelID, providers)
		}
		if r.ClientSupportsModel("claude-1", modelID) {
			t.Fatalf("expected claude client to reject %q", modelID)
		}
	}
}

func TestGetModelProvidersStillPrefersRegisteredProviders(t *testing.T) {
	r := newTestModelRegistry()
	r.RegisterClient("claude-1", "claude", []*ModelInfo{{ID: "claude-opus-5"}})

	providers := r.GetModelProviders("claude-opus-5")
	if len(providers) != 1 || providers[0] != "claude" {
		t.Fatalf("expected registered lookup to win, got %v", providers)
	}
}
