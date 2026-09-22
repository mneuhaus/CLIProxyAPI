package registry

import "strings"

const (
	// claudeProviderKey is the provider identifier used for Anthropic credentials.
	claudeProviderKey = "claude"

	// claudeModelPrefix marks an Anthropic model ID.
	claudeModelPrefix = "claude-"

	// claudeCloakedModelPrefix mirrors claudeDDModelPrefix in
	// internal/client/claude/models. It is duplicated rather than imported
	// because that package builds catalogs on top of this one.
	claudeCloakedModelPrefix = "claude-fable-5-dd-"
)

// IsClaudePassthroughModel reports whether modelID may be routed to a Claude
// credential even though no catalog entry registered it.
//
// The remote catalog trails Anthropic launches by days, so a model released
// today is unroutable here while the upstream account can already serve it.
// Accepting unknown claude-* IDs closes that window and leaves the upstream API
// as the authority: an ID Anthropic does not know still fails, just one hop
// later and with Anthropic's own error.
//
// Cloaked IDs are excluded. Those are rewritten before routing, so an unknown
// one signals a bug in the rewrite rather than a new release.
func IsClaudePassthroughModel(modelID string) bool {
	id := strings.ToLower(strings.TrimSpace(modelID))
	if !strings.HasPrefix(id, claudeModelPrefix) {
		return false
	}
	return !strings.HasPrefix(id, claudeCloakedModelPrefix)
}

// hasProviderClientLocked reports whether any registered client serves provider.
// Callers must hold at least a read lock on r.mutex.
func (r *ModelRegistry) hasProviderClientLocked(provider string) bool {
	for _, clientProvider := range r.clientProviders {
		if strings.EqualFold(strings.TrimSpace(clientProvider), provider) {
			return true
		}
	}
	return false
}
