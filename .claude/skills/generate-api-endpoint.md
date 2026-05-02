# Skill: Generate API Endpoint

## Usage
`/generate-api-endpoint [METHOD] [/path] [description]`

## Example
`/generate-api-endpoint POST /api/v1/admin/users Create a new customer user`

## What It Does
1. Adds request/response structs to `models/models.go`
2. Adds handler method to `handlers/admin.go` (or `license.go` for public routes)
3. Registers route in `cmd/server/main.go` under the correct group
4. Tells you if middleware (JWTAuth or AppAPIKey) is needed

## Template

### Model (models.go)
```go
type CreateXRequest struct {
    Field string `json:"field" binding:"required"`
}
```

### Handler (handlers/admin.go)
```go
func (h *AdminHandler) CreateX(c *gin.Context) {
    var req models.CreateXRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, models.Fail(err.Error()))
        return
    }
    // db call
    c.JSON(201, models.OK(result))
}
```

### Route (main.go)
```go
// Inside protected group (JWT required)
protected.POST("/resource", adminH.CreateX)

// Inside public group (X-App-Key required)
public.POST("/resource", licenseH.Method)
```
