# Skill: Add a New App to LicenseManager

## Usage
`/add-app [app-name] [slug]`

## Example
`/add-app "Amazon Tool" amazon-tool`

## What It Does
Guides you through connecting a new desktop app to LicenseManager.

## Steps

### 1. Create App in Admin Panel
```
Open http://localhost:5173 -> Apps -> New App
Fill in: Name, Slug, Max HWIDs per license
Copy the generated API Key
```

### 2. Add SDK to Your Go/Wails App
```bash
# In your app's go.mod directory
go get github.com/cheo/licensemanager/sdk/go/license
```

### 3. Use the SDK
```go
import "github.com/cheo/licensemanager/sdk/go/license"

client := license.New(license.Config{
    ServerURL:  "https://your-license-server.com",
    AppSlug:    "amazon-tool",       // must match slug in admin panel
    AppAPIKey:  "<api-key>",         // from admin panel
    DeviceName: "Home PC",
})

info, err := client.Validate(licenseKey)
if err != nil {
    // block the app: show error, exit
}
// info.MaxThreads, info.ExpiresAt are available
```

### 4. Issue a Test License
```
Admin u2192 Licenses u2192 Issue License
Select your new app, set days=30, maxThreads=5
Copy the key
```

### 5. Test
```go
info, err := client.Validate("AMAZ-XXXX-XXXX-XXXX")
// err == nil, info.Valid == true
```
