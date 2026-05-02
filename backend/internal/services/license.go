package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/cheo/licensemanager/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LicenseService struct {
	db *pgxpool.Pool
}

func NewLicenseService(db *pgxpool.Pool) *LicenseService {
	return &LicenseService{db: db}
}

func (s *LicenseService) Validate(ctx context.Context, req models.ValidateRequest) models.ValidateResponse {
	// Load license + app in one query
	var (
		licenseID  string
		isActive   bool
		expiresAt  time.Time
		maxThreads int
		maxHWID    int
	)

	err := s.db.QueryRow(ctx, `
		SELECT l.id, l.is_active, l.expires_at, l.max_threads, l.max_hwid
		FROM licenses l
		JOIN applications a ON l.app_id = a.id
		WHERE l.key = $1 AND a.slug = $2
	`, req.Key, req.AppSlug).Scan(&licenseID, &isActive, &expiresAt, &maxThreads, &maxHWID)

	if err == pgx.ErrNoRows {
		return models.ValidateResponse{Valid: false, Error: "invalid license key"}
	}
	if err != nil {
		return models.ValidateResponse{Valid: false, Error: "server error"}
	}
	if !isActive {
		return models.ValidateResponse{Valid: false, Error: "license revoked"}
	}
	if time.Now().After(expiresAt) {
		return models.ValidateResponse{Valid: false, Error: "license expired"}
	}

	// Check if HWID is already registered
	var hwidID string
	hwidErr := s.db.QueryRow(ctx,
		"SELECT id FROM license_hwids WHERE license_id = $1 AND hwid = $2",
		licenseID, req.HWID,
	).Scan(&hwidID)

	if hwidErr == nil {
		// Known device — update last_seen
		s.db.Exec(ctx,
			"UPDATE license_hwids SET last_seen_at = NOW() WHERE id = $1",
			hwidID)
		s.touchValidated(ctx, licenseID)
		return models.ValidateResponse{Valid: true, ExpiresAt: &expiresAt, MaxThreads: maxThreads}
	}

	// New device — check slot availability
	var deviceCount int
	s.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM license_hwids WHERE license_id = $1",
		licenseID).Scan(&deviceCount)

	if deviceCount >= maxHWID {
		devices := s.getDevices(ctx, licenseID)
		return models.ValidateResponse{
			Valid:   false,
			Error:   fmt.Sprintf("max devices reached (%d/%d)", deviceCount, maxHWID),
			Devices: devices,
		}
	}

	// Register new device
	deviceName := req.DeviceName
	if deviceName == "" {
		deviceName = fmt.Sprintf("Device %d", deviceCount+1)
	}
	s.db.Exec(ctx,
		"INSERT INTO license_hwids (license_id, hwid, device_name) VALUES ($1, $2, $3)",
		licenseID, req.HWID, deviceName)
	s.touchValidated(ctx, licenseID)

	return models.ValidateResponse{Valid: true, ExpiresAt: &expiresAt, MaxThreads: maxThreads}
}

func (s *LicenseService) touchValidated(ctx context.Context, licenseID string) {
	s.db.Exec(ctx,
		"UPDATE licenses SET last_validated_at = NOW() WHERE id = $1",
		licenseID)
}

func (s *LicenseService) getDevices(ctx context.Context, licenseID string) []models.LicenseHWID {
	rows, err := s.db.Query(ctx,
		"SELECT id, license_id, hwid, device_name, first_seen_at, last_seen_at FROM license_hwids WHERE license_id = $1",
		licenseID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var devices []models.LicenseHWID
	for rows.Next() {
		var d models.LicenseHWID
		rows.Scan(&d.ID, &d.LicenseID, &d.HWID, &d.DeviceName, &d.FirstSeenAt, &d.LastSeenAt)
		devices = append(devices, d)
	}
	return devices
}

func (s *LicenseService) RemoveHWID(ctx context.Context, licenseID, hwidID string) error {
	_, err := s.db.Exec(ctx,
		"DELETE FROM license_hwids WHERE id = $1 AND license_id = $2",
		hwidID, licenseID)
	return err
}

func GenerateLicenseKey(prefix string) string {
	b := make([]byte, 12)
	rand.Read(b)
	hex := fmt.Sprintf("%X", b)
	// Format: PREFIX-XXXX-XXXX-XXXX
	parts := []string{
		strings.ToUpper(prefix),
		hex[0:4],
		hex[4:8],
		hex[8:12],
	}
	return strings.Join(parts, "-")
}
