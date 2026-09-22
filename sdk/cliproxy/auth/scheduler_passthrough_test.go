package auth

import "testing"

func TestScheduledAuthMetaSupportsUnregisteredClaudeModel(t *testing.T) {
	registered := map[string]struct{}{"claude-opus-5": {}}

	claudeMeta := buildScheduledAuthMetaWithModelSet(&Auth{ID: "claude-1", Provider: "claude"}, registered, 1)
	if !claudeMeta.supportsModel("claude-opus-5") {
		t.Fatal("expected registered model to be supported")
	}
	if !claudeMeta.supportsModel("claude-opus-5-5") {
		t.Fatal("expected unregistered claude model to pass through on a claude auth")
	}
	if claudeMeta.supportsModel("gpt-6-astra") {
		t.Fatal("expected non-claude model to stay rejected")
	}

	codexMeta := buildScheduledAuthMetaWithModelSet(&Auth{ID: "codex-1", Provider: "codex"}, map[string]struct{}{"gpt-6-astra": {}}, 1)
	if codexMeta.supportsModel("claude-opus-5-5") {
		t.Fatal("expected claude passthrough to be limited to claude auths")
	}
}
