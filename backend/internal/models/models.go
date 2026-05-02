package models

import "time"

type Application struct {
	ID                string    `db:"id" json:"id"`
	Name              string    `db:"name" json:"name"`
	Slug              string    `db:"slug" json:"slug"`
	APIKey            string    `db:"api_key" json:"apiKey"`
	MaxHWIDPerLicense int       `db:"max_hwid_per_license" json:"maxHwidPerLicense"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
}

type License struct {
	ID              string     `db:"id" json:"id"`
	Key             string     `db:"key" json:"key"`
	AppID           string     `db:"app_id" json:"appId"`
	UserID          *string    `db:"user_id" json:"userId"`
	Plan            string     `db:"plan" json:"plan"`
	MaxThreads      int        `db:"max_threads" json:"maxThreads"`
	MaxHWID         int        `db:"max_hwid" json:"maxHwid"`
	ExpiresAt       time.Time  `db:"expires_at" json:"expiresAt"`
	IsActive        bool       `db:"is_active" json:"isActive"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	LastValidatedAt *time.Time `db:"last_validated_at" json:"lastValidatedAt"`
	Note            string     `db:"note" json:"note"`
}

type LicenseHWID struct {
	ID          string    `db:"id" json:"id"`
	LicenseID   string    `db:"license_id" json:"licenseId"`
	HWID        string    `db:"hwid" json:"hwid"`
	DeviceName  string    `db:"device_name" json:"deviceName"`
	FirstSeenAt time.Time `db:"first_seen_at" json:"firstSeenAt"`
	LastSeenAt  time.Time `db:"last_seen_at" json:"lastSeenAt"`
}

type User struct {
	ID               string    `db:"id" json:"id"`
	Email            string    `db:"email" json:"email"`
	PasswordHash     string    `db:"password_hash" json:"-"`
	StripeCustomerID *string   `db:"stripe_customer_id" json:"stripeCustomerId"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

type Admin struct {
	ID           string `db:"id" json:"id"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	Role         string `db:"role" json:"role"`
}

// Request/Response DTOs

type ValidateRequest struct {
	AppSlug    string `json:"appSlug" binding:"required"`
	Key        string `json:"key" binding:"required,max=32"`
	HWID       string `json:"hwid" binding:"required,max=64"`
	DeviceName string `json:"deviceName" binding:"max=100"`
}

type ValidateResponse struct {
	Valid      bool          `json:"valid"`
	Error      string        `json:"error,omitempty"`
	ExpiresAt  *time.Time    `json:"expiresAt,omitempty"`
	MaxThreads int           `json:"maxThreads,omitempty"`
	Devices    []LicenseHWID `json:"devices,omitempty"`
}

type CreateLicenseRequest struct {
	AppID      string `json:"appId" binding:"required"`
	UserID     string `json:"userId"`
	Plan       string `json:"plan" binding:"required"`
	MaxThreads int    `json:"maxThreads"`
	MaxHWID    int    `json:"maxHwid"`
	Days       int    `json:"days" binding:"required"`
	Note       string `json:"note"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAppRequest struct {
	Name              string `json:"name" binding:"required"`
	Slug              string `json:"slug" binding:"required"`
	MaxHWIDPerLicense int    `json:"maxHwidPerLicense" example:"2"`
}

type UpdateAppRequest struct {
	Name              *string `json:"name"`
	Slug              *string `json:"slug"`
	MaxHWIDPerLicense *int    `json:"maxHwidPerLicense"`
}

type AddHWIDRequest struct {
	HWID       string `json:"hwid" binding:"required"`
	DeviceName string `json:"deviceName"`
}

type UpdateLicenseRequest struct {
	IsActive   *bool   `json:"isActive"`
	MaxThreads *int    `json:"maxThreads"`
	MaxHWID    *int    `json:"maxHwid"`
	Days       *int    `json:"days"`
	Note       *string `json:"note"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(data interface{}) APIResponse {
	return APIResponse{Success: true, Data: data}
}

func Fail(err string) APIResponse {
	return APIResponse{Success: false, Error: err}
}
