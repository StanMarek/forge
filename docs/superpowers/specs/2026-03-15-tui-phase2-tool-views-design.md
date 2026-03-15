# TUI Phase 2 — Remaining Tool Views Design Spec

**Date:** 2026-03-15
**Status:** Approved
**Scope:** 5 new tool views, update createView factory, delete placeholder

---

## Goal

Implement the remaining 5 TUI tool views (JWT, JSON, Hash, URL, UUID) following the Base64View pattern established in Phase 1. After this, every tool in the sidebar is fully functional.

## Pattern

Every view follows the same structure:
- Implements `views.ToolView` interface: `Init()`, `Update()`, `View()`, `SetSize()`
- `textarea.Model` for user input (except UUID generate mode)
- `viewport.Model` for read-only output
- `ctrl+` keybindings for mode/option toggles (avoids textarea conflicts)
- Live processing: calls core tool function on every input change
- Error display: `err string` field, styled with `styles.ErrorStyle`

## Views

### JWT View (`ui/tui/views/jwt.go`)

```go
type jwtMode int
const (
    jwtModeFull    jwtMode = iota  // show header + payload + signature
    jwtModeHeader                   // header only
    jwtModePayload                  // payload only
)

type JWTView struct {
    input  textarea.Model
    output viewport.Model
    mode   jwtMode
    width, height int
    err    string
}
```

- **Keybindings:** `ctrl+h` header only, `ctrl+p` payload only, `ctrl+f` full output
- **Processing:** `tools.JWTDecode(input)` → display based on mode:
  - Full: `result.Output` (formatted header + payload + signature)
  - Header: `result.Header`
  - Payload: `result.Payload`
- **Status bar:** `ctrl+f: full  ctrl+h: header  ctrl+p: payload  tab: switch panel`

### JSON View (`ui/tui/views/json.go`)

```go
type jsonMode int
const (
    jsonModeFormat   jsonMode = iota
    jsonModeMinify
    jsonModeValidate
)

type JSONView struct {
    input    textarea.Model
    output   viewport.Model
    mode     jsonMode
    sortKeys bool
    width, height int
    err      string
}
```

- **Keybindings:** `ctrl+f` format, `ctrl+m` minify, `ctrl+v` validate, `ctrl+s` toggle sort-keys
- **Processing:**
  - Format: `tools.JSONFormat(input, 2, sortKeys, false)`
  - Minify: `tools.JSONMinify(input)`
  - Validate: `tools.JSONValidate(input)` → shows "valid" or error
- **Status bar:** `ctrl+f: format  ctrl+m: minify  ctrl+v: validate  ctrl+s: sort keys  tab: switch panel`

### Hash View (`ui/tui/views/hash.go`)

```go
type HashView struct {
    input     textarea.Model
    output    viewport.Model
    uppercase bool
    width, height int
}
```

- **Keybindings:** `ctrl+u` toggle uppercase
- **Processing:** Calls `tools.Hash(input, algo, uppercase)` for all 4 algorithms (md5, sha1, sha256, sha512). Displays all results:
  ```
  MD5:    <hash>
  SHA1:   <hash>
  SHA256: <hash>
  SHA512: <hash>
  ```
- **No error field** — Hash never fails on valid input (empty string is valid).
- **Status bar:** `ctrl+u: uppercase  tab: switch panel`

### URL View (`ui/tui/views/url.go`)

```go
type urlMode int
const (
    urlModeParse  urlMode = iota
    urlModeEncode
    urlModeDecode
)

type URLView struct {
    input     textarea.Model
    output    viewport.Model
    mode      urlMode
    component bool
    width, height int
    err       string
}
```

- **Keybindings:** `ctrl+p` parse (default), `ctrl+e` encode, `ctrl+d` decode, `ctrl+o` toggle component mode
- **Processing:**
  - Parse: `tools.URLParse(input)` → `result.Output` (formatted components)
  - Encode: `tools.URLEncode(input, component)` → `result.Output`
  - Decode: `tools.URLDecode(input)` → `result.Output`
- **Status bar:** `ctrl+p: parse  ctrl+e: encode  ctrl+d: decode  ctrl+o: component  tab: switch panel`

### UUID View (`ui/tui/views/uuid.go`)

```go
type uuidMode int
const (
    uuidModeGenerate uuidMode = iota
    uuidModeValidate
    uuidModeParse
)

type UUIDView struct {
    input     textarea.Model
    output    viewport.Model
    mode      uuidMode
    version   int  // 4 or 7
    uppercase bool
    noHyphens bool
    width, height int
    err       string
}
```

- **Keybindings:**
  - `ctrl+g` generate a new UUID (re-generates on each press)
  - `ctrl+4` switch to v4, `ctrl+7` switch to v7 (only in generate mode, ignored otherwise to avoid conflict with typing in textarea)
  - `ctrl+u` toggle uppercase
  - `ctrl+n` toggle no-hyphens
  - `ctrl+v` switch to validate mode (textarea active for input)
  - `ctrl+p` switch to parse mode (textarea active for input)
- **Processing:**
  - Generate: `tools.UUIDGenerate(version, uppercase, noHyphens)` → show in output. Auto-generates one on init.
  - Validate: `tools.UUIDValidate(input)` → "valid (version N)" or error
  - Parse: `tools.UUIDParse(input)` → `result.Output`
- **Generate mode:** textarea is hidden/unfocused. Output shows the generated UUID. `ctrl+g` generates a new one.
- **Validate/Parse mode:** textarea is visible and focused for input.
- **Status bar:** `ctrl+g: generate  ctrl+v: validate  ctrl+p: parse  ctrl+u: uppercase  tab: switch panel`

## Changes to Existing Files

### `ui/tui/app.go` — `createView` factory

```go
func createView(toolID, toolName string, width, height int) views.ToolView {
    var view views.ToolView
    switch toolID {
    case "base64":
        view = views.NewBase64View()
    case "jwt":
        view = views.NewJWTView()
    case "json":
        view = views.NewJSONView()
    case "hash":
        view = views.NewHashView()
    case "url":
        view = views.NewURLView()
    case "uuid":
        view = views.NewUUIDView()
    default:
        view = views.NewPlaceholder(toolName)
    }
    view.SetSize(width, height)
    return view
}
```

### `ui/tui/views/placeholder.go` — Keep for safety

Keep placeholder.go for the `default` case in createView. It costs nothing and handles unexpected tool IDs gracefully.

## File Manifest

```
Create: ui/tui/views/jwt.go
Create: ui/tui/views/json.go
Create: ui/tui/views/hash.go
Create: ui/tui/views/url.go
Create: ui/tui/views/uuid.go
Modify: ui/tui/app.go (createView factory)
```

## Testing Strategy

- Manual testing: `go run . tui`, select each tool, verify input/output works
- Core logic is already tested (200 tests). Views are thin presentation.
