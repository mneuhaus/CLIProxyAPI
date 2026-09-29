package registry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestDetectChangedProviders_KimiAliases(t *testing.T) {
	oldData := &staticModelsJSON{
		Kimi: []*ModelInfo{{ID: "kimi-k2"}},
	}
	newData := &staticModelsJSON{
		Kimi: []*ModelInfo{{ID: "kimi-k2"}, {ID: "kimi-k3"}},
	}

	changed := detectChangedProviders(oldData, newData)
	expected := map[string]bool{
		"kimi":     false,
		"kimi-ai":  false,
		"kimi.ai":  false,
		"kimi.com": false,
	}

	for _, p := range changed {
		if _, ok := expected[p]; ok {
			expected[p] = true
		}
	}

	for p, found := range expected {
		if !found {
			t.Errorf("expected changed provider %q to be reported, got %v", p, changed)
		}
	}
}

func TestRefreshModelsNow_RegistersNewModelWithoutWaitingForTick(t *testing.T) {
	var catalog map[string]any
	if err := json.Unmarshal(embeddedModelsJSON, &catalog); err != nil {
		t.Fatalf("decode embedded catalog: %v", err)
	}
	codexPro, _ := catalog["codex-pro"].([]any)
	catalog["codex-pro"] = append(codexPro, map[string]any{"id": "gpt-refresh-test", "object": "model", "owned_by": "openai"})
	body, err := json.Marshal(catalog)
	if err != nil {
		t.Fatalf("encode catalog: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}))
	defer server.Close()

	oldURLs, oldData := modelsURLs, getModels()
	modelsURLs = []string{server.URL}
	var notified []string
	SetModelRefreshCallback(func(changed []string) { notified = changed })
	defer func() {
		modelsURLs = oldURLs
		SetModelRefreshCallback(nil)
		modelsCatalogStore.mu.Lock()
		modelsCatalogStore.data = oldData
		modelsCatalogStore.mu.Unlock()
	}()
	if err := loadModelsFromBytes(embeddedModelsJSON, "embed"); err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}

	source, changed, err := RefreshModelsNow(context.Background())
	if err != nil {
		t.Fatalf("RefreshModelsNow: %v", err)
	}
	if source != server.URL {
		t.Errorf("source = %q, want %q", source, server.URL)
	}
	if !slices.Equal(changed, []string{"codex"}) || !slices.Equal(notified, []string{"codex"}) {
		t.Errorf("changed = %v, notified = %v, want [codex] for both", changed, notified)
	}
	found := false
	for _, m := range getModels().CodexPro {
		found = found || (m != nil && m.ID == "gpt-refresh-test")
	}
	if !found {
		t.Error("refreshed catalog does not contain gpt-refresh-test")
	}

	modelsURLs = []string{server.URL + "/missing"}
	server.Config.Handler = http.NotFoundHandler()
	if _, _, err := RefreshModelsNow(context.Background()); err == nil {
		t.Error("RefreshModelsNow succeeded although every URL failed")
	}
}
