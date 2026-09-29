package management

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/registry"
)

// RefreshModelCatalog fetches the remote model catalogs now instead of waiting for the
// periodic refresh, so a newly released model routes without restarting the proxy.
func (h *Handler) RefreshModelCatalog(c *gin.Context) {
	ctx := c.Request.Context()
	source, changed, errModels := registry.RefreshModelsNow(ctx)
	if errModels != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": errModels.Error()})
		return
	}
	if changed == nil {
		changed = []string{}
	}

	codex := gin.H{}
	codexSource, codexChanged, errCodex := registry.RefreshCodexClientModelsNow(ctx)
	if errCodex != nil {
		codex["error"] = errCodex.Error()
	} else {
		codex["source"] = codexSource
		codex["changed"] = codexChanged
	}

	c.JSON(http.StatusOK, gin.H{
		"source":              source,
		"changed_providers":   changed,
		"codex_client_models": codex,
	})
}
