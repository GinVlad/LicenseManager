package handlers

import (
	"net/http"

	"github.com/cheo/licensemanager/internal/models"
	"github.com/cheo/licensemanager/internal/services"
	"github.com/gin-gonic/gin"
)

type LicenseHandler struct {
	svc *services.LicenseService
}

func NewLicenseHandler(svc *services.LicenseService) *LicenseHandler {
	return &LicenseHandler{svc: svc}
}

// Validate godoc
// @Summary      Validate a license key
// @Tags         license
// @Accept       json
// @Produce      json
// @Security     AppApiKey
// @Param        body body models.ValidateRequest true "Validation payload"
// @Success      200 {object} models.APIResponse{data=models.ValidateResponse}
// @Failure      400,401 {object} models.APIResponse
// @Router       /license/validate [post]
func (h *LicenseHandler) Validate(c *gin.Context) {
	var req models.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(err.Error()))
		return
	}

	result := h.svc.Validate(c.Request.Context(), req)
	c.JSON(http.StatusOK, models.OK(result))
}

// RemoveHWID godoc
// @Summary      Remove a device binding (client-side)
// @Tags         license
// @Produce      json
// @Security     AppApiKey
// @Param        licenseId query string true "License ID"
// @Param        hwidId    query string true "HWID record ID"
// @Success      200 {object} models.APIResponse
// @Failure      400,401 {object} models.APIResponse
// @Router       /license/hwid [delete]
func (h *LicenseHandler) RemoveHWID(c *gin.Context) {
	licenseID := c.Query("licenseId")
	hwidID := c.Query("hwidId")

	if licenseID == "" || hwidID == "" {
		c.JSON(http.StatusBadRequest, models.Fail("licenseId and hwidId are required"))
		return
	}

	if err := h.svc.RemoveHWID(c.Request.Context(), licenseID, hwidID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("failed to remove device"))
		return
	}

	c.JSON(http.StatusOK, models.OK(nil))
}
