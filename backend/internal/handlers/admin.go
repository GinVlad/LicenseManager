package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cheo/licensemanager/internal/models"
	"github.com/cheo/licensemanager/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	db        *pgxpool.Pool
	svc       *services.LicenseService
	jwtSecret string
}

func NewAdminHandler(db *pgxpool.Pool, svc *services.LicenseService, jwtSecret string) *AdminHandler {
	return &AdminHandler{db: db, svc: svc, jwtSecret: jwtSecret}
}

// Login godoc
// @Summary      Admin login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.AdminLoginRequest true "Credentials"
// @Success      200 {object} models.APIResponse{data=object{token=string}}
// @Failure      401 {object} models.APIResponse
// @Router       /admin/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("invalid credentials"))
		return
	}

	var id, hash string
	err := h.db.QueryRow(context.Background(),
		"SELECT id, password_hash FROM admins WHERE email = $1",
		req.Email,
	).Scan(&id, &hash)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("invalid credentials"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("token error"))
		return
	}

	c.JSON(http.StatusOK, models.OK(gin.H{"token": signed}))
}

// ListApps godoc
// @Summary      List applications
// @Tags         apps
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.APIResponse{data=[]models.Application}
// @Failure      401 {object} models.APIResponse
// @Router       /admin/apps [get]
func (h *AdminHandler) ListApps(c *gin.Context) {
	rows, err := h.db.Query(context.Background(),
		"SELECT id, name, slug, api_key, max_hwid_per_license, created_at FROM applications ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("query error"))
		return
	}
	defer rows.Close()

	var apps []models.Application
	for rows.Next() {
		var a models.Application
		rows.Scan(&a.ID, &a.Name, &a.Slug, &a.APIKey, &a.MaxHWIDPerLicense, &a.CreatedAt)
		apps = append(apps, a)
	}
	c.JSON(http.StatusOK, models.OK(apps))
}

// CreateApp godoc
// @Summary      Create application
// @Tags         apps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body models.CreateAppRequest true "App info"
// @Success      201 {object} models.APIResponse{data=models.Application}
// @Failure      400,401 {object} models.APIResponse
// @Router       /admin/apps [post]
func (h *AdminHandler) CreateApp(c *gin.Context) {
	var body struct {
		Name              string `json:"name" binding:"required"`
		Slug              string `json:"slug" binding:"required"`
		MaxHWIDPerLicense int    `json:"maxHwidPerLicense"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(err.Error()))
		return
	}
	if body.MaxHWIDPerLicense == 0 {
		body.MaxHWIDPerLicense = 2
	}

	var app models.Application
	err := h.db.QueryRow(context.Background(), `
		INSERT INTO applications (name, slug, max_hwid_per_license)
		VALUES ($1, $2, $3)
		RETURNING id, name, slug, api_key, max_hwid_per_license, created_at`,
		body.Name, body.Slug, body.MaxHWIDPerLicense,
	).Scan(&app.ID, &app.Name, &app.Slug, &app.APIKey, &app.MaxHWIDPerLicense, &app.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("failed to create app: "+err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.OK(app))
}

// UpdateApp godoc
// @Summary      Update an application
// @Tags         apps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                  true "App ID"
// @Param        body body models.UpdateAppRequest true "Fields to update (all optional)"
// @Success      200 {object} models.APIResponse
// @Failure      400,401 {object} models.APIResponse
// @Router       /admin/apps/{id} [put]
func (h *AdminHandler) UpdateApp(c *gin.Context) {
	id := c.Param("id")
	var body models.UpdateAppRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(err.Error()))
		return
	}

	if body.Name != nil {
		h.db.Exec(context.Background(), "UPDATE applications SET name = $1 WHERE id = $2", *body.Name, id)
	}
	if body.Slug != nil {
		h.db.Exec(context.Background(), "UPDATE applications SET slug = $1 WHERE id = $2", *body.Slug, id)
	}
	if body.MaxHWIDPerLicense != nil {
		h.db.Exec(context.Background(), "UPDATE applications SET max_hwid_per_license = $1 WHERE id = $2", *body.MaxHWIDPerLicense, id)
	}

	c.JSON(http.StatusOK, models.OK(nil))
}

// DeleteApp godoc
// @Summary      Delete an application
// @Tags         apps
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "App ID"
// @Success      200 {object} models.APIResponse
// @Failure      401 {object} models.APIResponse
// @Router       /admin/apps/{id} [delete]
func (h *AdminHandler) DeleteApp(c *gin.Context) {
	id := c.Param("id")
	h.db.Exec(context.Background(), "DELETE FROM applications WHERE id = $1", id)
	c.JSON(http.StatusOK, models.OK(nil))
}

// ListLicenses godoc
// @Summary      List licenses
// @Tags         licenses
// @Produce      json
// @Security     BearerAuth
// @Param        appId query string false "Filter by app ID"
// @Success      200 {object} models.APIResponse{data=[]models.License}
// @Failure      401 {object} models.APIResponse
// @Router       /admin/licenses [get]
func (h *AdminHandler) ListLicenses(c *gin.Context) {
	appID := c.Query("appId")
	args := []interface{}{}
	query := "SELECT id, key, app_id, user_id, plan, max_threads, max_hwid, expires_at, is_active, note, created_at, last_validated_at FROM licenses"
	if appID != "" {
		query += " WHERE app_id = $1"
		args = append(args, appID)
	}
	query += " ORDER BY created_at DESC LIMIT 200"

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("query error"))
		return
	}
	defer rows.Close()

	var licenses []models.License
	for rows.Next() {
		var l models.License
		rows.Scan(&l.ID, &l.Key, &l.AppID, &l.UserID, &l.Plan, &l.MaxThreads, &l.MaxHWID,
			&l.ExpiresAt, &l.IsActive, &l.Note, &l.CreatedAt, &l.LastValidatedAt)
		licenses = append(licenses, l)
	}
	c.JSON(http.StatusOK, models.OK(licenses))
}

// CreateLicense godoc
// @Summary      Issue a license key
// @Tags         licenses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body models.CreateLicenseRequest true "License details"
// @Success      201 {object} models.APIResponse{data=models.License}
// @Failure      400,401 {object} models.APIResponse
// @Router       /admin/licenses [post]
func (h *AdminHandler) CreateLicense(c *gin.Context) {
	var req models.CreateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(err.Error()))
		return
	}
	if req.MaxThreads == 0 {
		req.MaxThreads = 5
	}

	// Get app slug and default HWID limit
	var slug string
	var appMaxHWID int
	err := h.db.QueryRow(context.Background(),
		"SELECT slug, max_hwid_per_license FROM applications WHERE id = $1", req.AppID).
		Scan(&slug, &appMaxHWID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("app not found"))
		return
	}
	if req.MaxHWID == 0 {
		req.MaxHWID = appMaxHWID
	}

	prefix := slug
	if len(prefix) > 4 {
		prefix = prefix[:4]
	}

	key := services.GenerateLicenseKey(prefix)
	expiresAt := time.Now().AddDate(0, 0, req.Days)

	var license models.License
	err = h.db.QueryRow(context.Background(), `
		INSERT INTO licenses (key, app_id, plan, max_threads, max_hwid, expires_at, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, key, app_id, user_id, plan, max_threads, max_hwid, expires_at, is_active, note, created_at`,
		key, req.AppID, req.Plan, req.MaxThreads, req.MaxHWID, expiresAt, req.Note,
	).Scan(&license.ID, &license.Key, &license.AppID, &license.UserID,
		&license.Plan, &license.MaxThreads, &license.MaxHWID, &license.ExpiresAt, &license.IsActive,
		&license.Note, &license.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("failed to create license"))
		return
	}
	c.JSON(http.StatusCreated, models.OK(license))
}

// UpdateLicense godoc
// @Summary      Update a license
// @Tags         licenses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                      true "License ID"
// @Param        body body models.UpdateLicenseRequest true "Fields to update (all optional)"
// @Success      200 {object} models.APIResponse
// @Failure      400,401 {object} models.APIResponse
// @Router       /admin/licenses/{id} [put]
func (h *AdminHandler) UpdateLicense(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		IsActive   *bool   `json:"isActive"`
		MaxThreads *int    `json:"maxThreads"`
		MaxHWID    *int    `json:"maxHwid"`
		Days       *int    `json:"days"`
		Note       *string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(err.Error()))
		return
	}

	if body.IsActive != nil {
		h.db.Exec(context.Background(),
			"UPDATE licenses SET is_active = $1 WHERE id = $2", *body.IsActive, id)
	}
	if body.MaxThreads != nil {
		h.db.Exec(context.Background(),
			"UPDATE licenses SET max_threads = $1 WHERE id = $2", *body.MaxThreads, id)
	}
	if body.MaxHWID != nil {
		h.db.Exec(context.Background(),
			"UPDATE licenses SET max_hwid = $1 WHERE id = $2", *body.MaxHWID, id)
	}
	if body.Days != nil {
		h.db.Exec(context.Background(),
			"UPDATE licenses SET expires_at = NOW() + ($1 || ' days')::INTERVAL WHERE id = $2",
			fmt.Sprintf("%d", *body.Days), id)
	}
	if body.Note != nil {
		h.db.Exec(context.Background(),
			"UPDATE licenses SET note = $1 WHERE id = $2", *body.Note, id)
	}

	c.JSON(http.StatusOK, models.OK(nil))
}

// DeleteLicense godoc
// @Summary      Delete a license
// @Tags         licenses
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "License ID"
// @Success      200 {object} models.APIResponse
// @Failure      401 {object} models.APIResponse
// @Router       /admin/licenses/{id} [delete]
func (h *AdminHandler) DeleteLicense(c *gin.Context) {
	id := c.Param("id")
	h.db.Exec(context.Background(), "DELETE FROM licenses WHERE id = $1", id)
	c.JSON(http.StatusOK, models.OK(nil))
}

// ListHWIDs godoc
// @Summary      List devices bound to a license
// @Tags         hwids
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "License ID"
// @Success      200 {object} models.APIResponse{data=[]models.LicenseHWID}
// @Failure      401 {object} models.APIResponse
// @Router       /admin/licenses/{id}/hwids [get]
func (h *AdminHandler) ListHWIDs(c *gin.Context) {
	id := c.Param("id")
	rows, err := h.db.Query(context.Background(),
		"SELECT id, license_id, hwid, device_name, first_seen_at, last_seen_at FROM license_hwids WHERE license_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("query error"))
		return
	}
	defer rows.Close()

	var hwids []models.LicenseHWID
	for rows.Next() {
		var h models.LicenseHWID
		rows.Scan(&h.ID, &h.LicenseID, &h.HWID, &h.DeviceName, &h.FirstSeenAt, &h.LastSeenAt)
		hwids = append(hwids, h)
	}
	c.JSON(http.StatusOK, models.OK(hwids))
}

// AddHWID godoc
// @Summary      Manually add a device binding
// @Tags         hwids
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                true "License ID"
// @Param        body body models.AddHWIDRequest true "Device info"
// @Success      201 {object} models.APIResponse{data=models.LicenseHWID}
// @Failure      400,401,409 {object} models.APIResponse
// @Router       /admin/licenses/{id}/hwids [post]
func (h *AdminHandler) AddHWID(c *gin.Context) {
	licenseID := c.Param("id")
	var req models.AddHWIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(err.Error()))
		return
	}

	// Simple check if limit reached (for manual admin override, we might skip this,
	// but let's enforce it here unless they increase the maxHwid first)
	var maxHwid, currentCount int
	err := h.db.QueryRow(context.Background(),
		"SELECT max_hwid FROM licenses WHERE id = $1", licenseID).Scan(&maxHwid)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Fail("license not found"))
		return
	}

	h.db.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM license_hwids WHERE license_id = $1", licenseID).Scan(&currentCount)

	if currentCount >= maxHwid {
		c.JSON(http.StatusConflict, models.Fail("license has reached its device limit"))
		return
	}

	var hwid models.LicenseHWID
	err = h.db.QueryRow(context.Background(), `
		INSERT INTO license_hwids (license_id, hwid, device_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (license_id, hwid) DO UPDATE SET last_seen_at = NOW()
		RETURNING id, license_id, hwid, device_name, first_seen_at, last_seen_at`,
		licenseID, req.HWID, req.DeviceName,
	).Scan(&hwid.ID, &hwid.LicenseID, &hwid.HWID, &hwid.DeviceName, &hwid.FirstSeenAt, &hwid.LastSeenAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("failed to add hwid: "+err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.OK(hwid))
}

// DeleteHWID godoc
// @Summary      Remove a device binding
// @Tags         hwids
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "HWID record ID"
// @Success      200 {object} models.APIResponse
// @Failure      401 {object} models.APIResponse
// @Router       /admin/hwids/{id} [delete]
func (h *AdminHandler) DeleteHWID(c *gin.Context) {
	id := c.Param("id")
	h.db.Exec(context.Background(), "DELETE FROM license_hwids WHERE id = $1", id)
	c.JSON(http.StatusOK, models.OK(nil))
}
