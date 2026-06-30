# ChefScript Dictionary

## Identity

- **Name:** ChefScript
- **Version:** 1.0.0
- **What:** A JSON-to-UI engine for desktop apps. You describe intent in JSON atoms — the Go engine renders HTML/CSS/JS into a WebView2 shell. No browser, no client framework, no build step.
- **Replaces:** React, Vue, Angular, Electron. One Go binary, one JSON page per screen, one runtime. SPA behavior via partial navigation without the SPA overhead.
- **You write:** JSON atoms: `[tag, props?, ...children]`. That's it.
- **You don't write:** HTML, CSS, JS, or client-side state management. The engine and runtime handle all of that.

---

## Principles

Non-negotiable design rules. Every decision flows from these.

- **json_is_intent:** JSON declares WHAT. The engine decides HOW. JSON is dumb, the engine is smart. Never embed logic, loops, or conditionals in JSON.
- **no_raw_html:** No raw HTML or inline styles in page JSON. Use dictionary atoms or register a new component in Go.
- **no_inline_scripts:** All JS lives in engine/runtime.go. Components declare intent via data attributes, the runtime provides behavior. Required for partial navigation to work.
- **binary_data_pipe:** Schema is the pipe, binary is the transport, database is the state. All MongoDB collections use schemas/binary/ for encoding. No raw Col() calls for schematized data.
- **no_client_framework:** No React, no Vue, no client-side state. Partial navigation (csNavigate) gives SPA behavior — body swap, history update, no reload — without a framework eating memory.
- **trace_dont_guess:** Built-in diagnostics panel (F8) and flight recorder give full visibility into render and runtime. Use them before guessing at bugs.
- **components_are_self_contained:** A component receives props + children, returns complete HTML/CSS/JS. No external dependencies unless explicitly composed.
- **logic_outside_json:** JSON uses template variables (`{{data.key}}`) as insertion points. Logic lives in Go — page loaders, action handlers, component render functions.

---

## Architecture

How the layers connect. Every request flows through this.

### Render Flow

1. `GET /page/{name}` hits the Go server
2. Session middleware loads user session from MongoDB cookie
3. Page loader (if registered) injects live data or redirects
4. Template engine resolves `{{variables}}` against page context
5. ChefScript engine walks the JSON body array, resolves each atom against the component registry
6. Components render to HTML/CSS/JS. Diagnostics collector traces every step.
7. Runtime JS (csRuntime) is injected into every page — powers all interactive behavior
8. HTML response served to WebView2 window

### Action Flow

1. User interaction triggers `csPost(action, data)` or `csAction(action)`
2. `POST /api/{action}` hits the Go server
3. Registered handler runs, returns ActionResult
4. JS runtime processes the result:
   - Data → DOM patches: target (innerHTML by data-id) or setAttr (attributes by selector)
   - Toast → floating notification
   - Error → error toast + field highlighting

### Navigation

- **Full page:** `window.location.href` — complete reload, new page context, re-runs everything
- **Partial:** `csNavigate(url)` — fetches `/partial/{name}`, swaps `document.body.innerHTML`, updates title + history. No reload. Runtime persists. Flight recorder client buffer persists.

### Data Pipe

- **Binary schema:** Component declares schema prop → engine loads binary from MongoDB → schema decodes at fixed offsets → native Go values passed as props. Zero JSON serialization.
- **Template vars:** `{{user.*}}`, `{{data.*}}`, `{{session.*}}` resolved by template engine before render. Page loaders inject `{{data.*}}` values.
- **Static JSON:** JSON files in `data/` loaded at startup. Use for reference data, lookup tables, descriptions — editable without recompiling.

### Session

Cookie-backed MongoDB sessions with 24h TTL. `CreateSession`, `UpdateSession`, `DestroySession`, `GetSessionFromCtx`. Auto-loaded into every request.

---

## Runtime

The JS runtime (`engine/runtime.go`) is embedded in every page. These are the framework-level functions available to all components.

### csPost

- **Signature:** `window.csPost(action, data, opts)`
- **Does:** POST `/api/{action}` with JSON body. Handles button loading state, form error clearing.
- **Response handling:**
  - `result.error` → error toast, field errors if `result.fields`, `onError` callback
  - `result.toast` → floating notification
  - `result.data` → DOM patching loop: iterates patches, applies target or setAttr

### csAction

- **Signature:** `window.csAction(action, triggerEl)`
- **Dispatches:**
  - `post:{action}:{formId}` — Collects form data, calls csPost
  - `navigate:{url}` — csNavigate or window.location
  - `modal:{id}` — Opens modal
  - `drawer:{id}` — Opens drawer
  - `snackbar:{msg}:{variant}` — Shows toast
  - `close-modal:{id}` — Closes modal
  - `close-drawer:{id}` — Closes drawer

### csNavigate

- **Signature:** `window.csNavigate(url)`
- **Does:** Converts `/page/` to `/partial/`, fetches JSON, swaps body, updates title + history. Falls back to full page load on error.

### DOM Patching

#### Target Patch

- **Fields:** `{ target: data-id, html: string, append: bool, scroll: bool }`
- **Does:** Finds element by data-id, sets or appends innerHTML. Scroll auto-scrolls to bottom.

#### setAttr Patch

- **Fields:** `{ setAttr: css-selector, attrs: { key: value, ... } }`
- **Does:** Finds element by CSS selector, sets each attribute. Use for data attributes, classes, visibility toggles.

### on_enter

Inputs with `data-on-enter` auto-wire keydown listener. Enter collects value + all `data-*` attributes, POSTs via csPost. ESC dispatches `cs:escape` CustomEvent.

### Autosave

Forms with `data-autosave` auto-checkpoint to localStorage on input (600ms debounce). Restore on load, clear on successful submit.

### Flight Recorder (Client)

`window.__flight` ring buffer (500 events). `flightRecord(cat, level, msg, detail)` appends. Survives partial navigation. Reset on full page load only if buffer is empty.

### Client State

Client-side state system for multi-step interactions. A plain JS object on window holds state. Elements declare which keys they read via `data-state`. Actions mutate state and the runtime updates all bound elements. Server is only touched at boundaries — page load (get data) and completion (save result).

#### State Object

One flat object per feature on window. This is the single source of truth for all client-side state during the interaction.

- **Naming:** `window.__wire`, `window.__quiz`, `window.__heist` — prefixed with `__` to separate from framework globals.
- **Structure:** Flat keys. No nesting. Every key is directly addressable as `feature.key`.
- **Example:** `window.__wire = { batches: [], batchIdx: 0, connections: {}, selectedAnswer: null, locked: false, switchPulled: false, results: [], module: '' };`

#### Get (data-state)

Elements declare what state they display. The runtime reads the value and sets the element content or attribute.

- **Syntax:** `data-state="feature.key"` on any element.
- **Behaviour:** Runtime reads `window.__feature.key` and sets element textContent to the value. Element updates automatically when the key changes.
- **Attribute binding:** `data-state-attr="feature.key:attrName"` binds a state key to an element attribute instead of textContent. Use for class toggles, visibility, data attributes.
- **Examples:**
  - `<div data-state="wire.batchIdx"></div>` — displays current batch index
  - `<div data-state-attr="wire.locked:data-locked"></div>` — sets data-locked attribute from state
  - `<div data-state="wire.module"></div>` — displays module name

#### Set (csState.set)

Actions mutate state. After mutation, the runtime finds every element bound to the changed keys and updates them. O(changed keys) — no tree walk, no diffing.

- **JS API:** `csState.set(feature, {key: value, key2: value2})`
- **Behaviour:**
  1. Merges key-value pairs into `window.__feature`
  2. For each changed key, finds all elements with matching `data-state` or `data-state-attr`
  3. Updates textContent or attribute value
  4. No re-render, no virtual DOM, no diffing — direct element updates
- **From actions:** Action functions call `csState.set()` after doing their work. This is the only way state changes propagate to the DOM.
- **Example:** `csState.set('wire', { batchIdx: 2, locked: false, switchPulled: false });`

#### Init (csState.init)

Go component embeds initial state on page load. The runtime registers the state object and builds the reverse index of element bindings.

- **Syntax:** `csState.init(feature, {key: value, ...})`
- **Behaviour:**
  1. Creates `window.__feature` with the given values
  2. Scans DOM for all `data-state` and `data-state-attr` attributes matching the feature
  3. Builds reverse index: `key → [element, element, ...]` for O(1) lookup on set
  4. Applies initial values to all bound elements
- **Example:** `csState.init('wire', { batches: [...], batchIdx: 0, module: 'Network Analysis', results: [], locked: false });`

#### Actions

Functions that contain game/interaction logic. They read state, apply rules, mutate state via `csState.set()`. Validation is here — client-side, against data already in state.

- **Pattern:** Read from `window.__feature` → apply logic → `csState.set(feature, changes)`
- **Example flow:**

```javascript
function wireClickAnswer(aIdx) {
  var s = window.__wire;
  if (s.locked) return;
  if (s.connections[aIdx] !== undefined) { delete s.connections[aIdx]; }
  else { s.selectedAnswer = aIdx; }
  csState.set('wire', { connections: s.connections, selectedAnswer: s.selectedAnswer });
}
```

#### Save

On completion, one POST sends the final result to the server. Server records it via binary schema. This is the only server write.

- **Pattern:** `csPost(action, {module, cleared, total})` — server encodes and stores. No mid-interaction server calls.
- **Example:** `csPost('wire/save', { module: __wire.module, cleared: 3, total: 5 });`

#### Client State Rules

- State is global — no component tree, no prop passing. Any function can read or write it.
- Validation happens client-side. The data is already in the state object. Don't round-trip to validate.
- The server delivers data and records results. It does not validate game logic or manage mid-interaction state.
- Binary pipe and client state are separate channels. Binary moves data at page boundaries. Client state lives in JS between boundaries.
- One state object per feature. Keep it flat.
- `data-state` binds display. `csState.set()` triggers updates. Actions contain logic. That's the full model.

#### Performance

- **set:** O(changed keys) — reverse index lookup, direct DOM update per bound element.
- **comparison:** React diffs entire virtual DOM tree O(tree size) per state change. This touches only elements bound to changed keys. No tree walk, no reconciliation, no component re-execution.

---

## Diagnostics

Built-in observability for every ChefScript app. No setup required.

### Panel

- **Toggle:** F8 — 3-state: hidden → open → closed
- **Tabs:** All, Errors, Warnings, Template, Render, Actions, Rules, Flight
- **Per render:** JSON validation, template variable tracing (resolved/nil), atom rendering (tag, depth, parent, props count), structural rule enforcement
- **Export:** Download filtered entries or flight events as JSON

### Flight Recorder

- **What:** Passive, always-on event recorder at framework boundaries. Two ring buffers (Go + JS), 500 events each. RAM only — gone on app close.
- **Zero latency:** Array push in JS, mutex append in Go. Nanosecond operations. Read only when panel opens.
- **Server events:** page:load, partial:load, action:enter/exit, session:create/update/destroy, template resolution, render diagnostics
- **Client events:** csPost send/response/error, DOM patches (target + setAttr), csNavigate start/done/fallback, JS errors, unhandled rejections
- **Viewing:** Flight tab in diagnostics panel. Fetches server events via `/api/flight/snapshot`, reads client events from `window.__flight` directly. Merges chronologically. Polls every 2.5s while tab is active.
- **Source badges:** SRV = server event, CLI = client event

---

## Data Patterns

How data moves between JSON pages, Go handlers, and the DOM.

### Action Result Fields

- **Error:** string — show error toast, re-enable submit button.
- **Fields:** `map[string]string` — per-field error messages, highlights specific inputs.
- **Toast:** string — floating notification message.
- **ToastVariant:** `info | success | warning | error`
- **Data:** object or array — DOM patches. If contains `target`, patches element by data-id. If contains `setAttr`, sets attributes by CSS selector. Array = multi-target.

### Page Changes

Page-to-page navigation happens via links. The runtime intercepts `<a href='/page/...'>` clicks and uses csNavigate to swap the body. Handlers never return redirects.

### Action Responses

Handlers return patches (Data), toasts, or errors. Use target patches to update UI regions, setAttr patches to toggle state, toasts for notifications.

### Template Variables

- `{{user.id}}` — Logged-in user ID
- `{{user.hackerName}}` — User display name
- `{{user.role}}` — student | instructor
- `{{data.key}}` — Arbitrary data from page loader
- `{{session.key}}` — Raw session values

### Binary Schema

- **Declare:** Component uses schema prop → engine intercepts before render
- **Flow:** `GetBinarySchema → BinaryFindAll → Decode at fixed offsets → native values as props`
- **Speed:** O(fields) per record. No parsing, no scanning. Values reach component as native types.

---

## Anti-Patterns

What NOT to do. These break the architecture.

- Don't write raw HTML in page JSON — use atoms or register a component
- Don't add `<script>` tags or inline JS — all JS goes in `engine/runtime.go`
- Don't build client-side state management — state lives server-side in session or MongoDB
- Don't use raw `Col()` calls for schematized data — go through binary schema pipe
- Don't guess at bugs — F8, check Flight tab, trace the actual flow
- Don't embed logic in JSON — JSON is the map, Go is the logic
- Don't handcraft CSS/JS for things Bootstrap already provides — pipe it through atoms
- Don't hardcode px for spacing, radius, shadows, or transitions in component CSS — use design tokens from `engine/tokens.go`
- Don't use `toStringMap()` for data that round-trips through MongoDB — use `toBsonSafeMap()` which handles `bson.D` conversion
- Don't add app-specific instrumentation to the flight recorder — it instruments framework boundaries only

---

## Conventions

### Files

- **pages:** `pages/*.json` — one file per page
- **components:** `engine/components.go` (registry + helpers), `engine/components_atom_*.go` (render functions by category: layout, navigation, inputs, actions, display, feedback, overlay, data, media, chat) — registered in `RegisterDefaults()`
- **actions:** `engine/` or `handlers/` — registered via `engine.RegisterAction(name, handler)`
- **data:** `data/*.json` for static reference data, loaded at startup
- **schemas:** `schemas/binary/*.json` — field-level schemas for binary-encoded MongoDB collections
- **icons:** `public/icons/*.svg` — referenced by filename without `.svg` extension
- **runtime:** `engine/js/*.js` — client-side JS split by concern. `engine/runtime.go` assembles them via `go:embed`.

### Naming

- **data_id:** All components auto-generate a data-id. Override with the `data-id` prop for DOM targeting.
- **flags:** Engine auto-generates flag names as `cmd_target`. Account for system context to avoid duplicates across systems.
- **actions:** Namespace with slash: `auth/login`, `dirtbyte/command`, `admin/export`

### Design Tokens

All component CSS uses MUI-derived design tokens defined in `engine/tokens.go` `tokenCSS()`. No hardcoded px for spacing, radius, shadows, or transitions in component CSS.

- **Source:** `engine/tokens.go` — `tokenCSS()` emits `:root` with all tokens, `[data-theme="dark"]` with dark overrides
- **Spacing:** 8px base. Use `var(--spacing-N)`: 1=4px, 2=8px, 3=12px, 4=16px, 5=20px, 6=24px, 8=32px, 10=40px, 12=48px
- **Typography:** MUI rem scale. `var(--font-body2)` for 14px, `var(--font-body1)` for 16px, `var(--font-caption)` for 12px. Shorthand includes weight/size/line-height/family.
- **Shape:** `var(--radius-sm)` 4px, `var(--radius-md)` 8px, `var(--radius-lg)` 12px, `var(--radius-chip)` 16px, `var(--radius-full)` pill
- **Elevation:** `var(--elevation-N)`: 0=none, 1-4 graduated, 6/8/12/24 for overlays. Replaces box-shadow hardcoding.
- **Transitions:** `var(--duration-shortest)` 150ms, `var(--duration-shorter)` 200ms, `var(--duration-short)` 250ms, `var(--duration-standard)` 300ms, `var(--duration-complex)` 375ms. Easing: `var(--ease-standard)`, `var(--ease-decelerate)`, `var(--ease-accelerate)`, `var(--ease-sharp)`.
- **Z-index:** `var(--z-drawer)` 1200, `var(--z-modal)` 1300, `var(--z-snackbar)` 1400, `var(--z-tooltip)` 1500
- **Dark mode:** `[data-theme="dark"]` flips palette, text, surface, and elevation tokens. Component CSS unchanged — only variables swap. Any ChefScript app gets dark mode by setting `data-theme="dark"` on the root element.

#### Theme Cascade

1. `:root` — MUI tokens (engine light default)
2. `[data-theme="dark"]` — engine dark mode (palette + surface flip)
3. `[data-bs-theme="dark"]` — Bootstrap utilities layer
4. App theme overrides (themeDark/themeStretch) — app-specific palette
5. Scoped overrides (`.lxlab`, `.pylab`, `.sqlab`, `.dirtbyte`) — per-app identity colors

#### Token Rules

- Never hardcode px for spacing in component CSS — use `var(--spacing-N)`
- Never hardcode border-radius — use `var(--radius-sm/md/lg/chip/full)`
- Never hardcode box-shadow — use `var(--elevation-N)`
- Never hardcode transition duration or easing — use `var(--duration-*)` and `var(--ease-*)`
- Font sizes use rem from the MUI scale, not px
- Colors that vary by theme must use CSS variables, not hex codes. Scoped app themes override only the palette variables they need.
- Layout constants (`--container-max-width`) are named variables but not part of the spacing scale

### Air Dev

Air watches the project for file changes and instantly rebuilds the binary on save. No manual rebuild step during development.

- **Config:** `.air.toml` in project root
- **Run:** `air`
- **Build cmd:** `go build -o ./tmp/app.exe .`
- **Watched extensions:** `.go`, `.json`, `.js`
- **Excluded dirs:** `tmp`, `python-runtime`, `node_modules`, `.git`
- **Behaviour:**
  - On any `.go`, `.json`, or `.js` file save → kills running process → rebuilds → restarts
  - 1 second delay before rebuild (debounce rapid saves)
  - Console clears on each rebuild
  - Test files (`*_test.go`) are excluded from watch

### Adding a Component

Write a Go render function in the appropriate `engine/components_atom_*.go` file (layout, navigation, inputs, actions, display, feedback, overlay, data, media, chat). Register in `RegisterDefaults()` in `engine/components.go`. Use from JSON immediately.

### Adding an Action

Write a Go function matching `APIHandler` signature. Call `engine.RegisterAction(name, handler)` in `app.go`. Call from JSON via `on:click post:{name}`.

### Self-Contained Binary

The exe must be self-contained. All assets are embedded at compile time via `go:embed`. The only external dependency is `python-runtime/` which ships alongside the exe.

**How it works:** `embed.go` in the project root embeds `pages/`, `schemas/binary/`, `data/`, `public/`, `python/`, and `dictionary.json` into a single `embed.FS`. `engine/embed.go` provides `ReadEmbedFile` and `ReadEmbedDir` that read from the embedded FS first, falling back to OS for CLI tools.

**Rules:**

- Never use `os.ReadFile` for application assets — use `engine.ReadEmbedFile` instead.
- Never use `os.ReadDir` for application directories — use `engine.ReadEmbedDir` instead.
- All file paths must use forward slashes (`embed.FS` requirement). Convert with `strings.ReplaceAll(path, "\\", "/")`.
- New data files must be covered by the `go:embed` directive in `embed.go`. If you add a new directory, add it to the directive.
- New JS files go in `engine/js/` and get a `go:embed` directive in `engine/runtime.go` — they are automatically part of the binary.
- Python scripts are embedded but extracted to a temp directory at startup because Python needs real files on disk.
- Static assets (`public/`) are served via `http.FS` from the embedded FS.
- CLI tools (`cmd/*`) run without the embedded FS — they use the OS fallback path. This is intentional.

**Build:** `go build -ldflags "-H windowsgui -X main.mongoURIEnc=<hex>" -o cybersecurity.exe`

The `-X main.mongoURIEnc=<hex>` ldflag is **mandatory**. The MongoDB server is locked behind credentials; the connection string (creds + host) is XOR-obfuscated with `mongoKey` and hex-encoded, then injected **only at build time** (`main.go`, `mongoURI()`) — never stored in source. Without the flag the binary falls back to `mongodb://localhost:27017` and **cannot reach the protected server — no scores, coins, sessions, or updates**. This is also the app's write protection: because only the built app holds the creds, external POST/PUT to the API can't write to the DB — game writes (`coins/earn`, `*/save`, completions) are safe by transport, not per-handler validation. If the server endpoint/creds change, regenerate the hex and update this command.

**Current hex (do NOT lose again):**
`ca9cafbeda86e63c16deabcbb683d279c6919eb8c592be606dc2a3e0b9dffa64958596a182acf4307582f0cb9cd2c35796929086d1a5c4370ec3e683e2c8913b91c2fbeb82d2b53116cea9c7a08ef37ad281a2bc8883e06b509f`

- This is `hex(XOR(connection-URI, mongoKey))`; `mongoKey` is the byte array in `main.go`.
- Decodes to `mongodb://cyberlab_app:…@172.16.1.61:27017/…` (verified via XOR with `mongoKey`).
- Stored nowhere in source/binary by design. If lost: XOR-decode this hex with `mongoKey`, or recover from prior session transcripts under `convo/`. To regenerate: XOR the full URI with `mongoKey`, then hex-encode.

**Deploy:** `cybersecurity.exe` + `python-runtime/` folder. Nothing else needed.

### Window Close

Use `navigator.sendBeacon` for any data that must reach the server when a window closes. Normal fetch/csPost gets killed mid-flight — sendBeacon survives window close.

**Pattern:** `window.addEventListener('beforeunload', function() { navigator.sendBeacon('/api/endpoint', JSON.stringify(payload)); });`

**Rules:**

- sendBeacon is fire-and-forget — no response, no callback. Use only for writes the server can process without confirming back.
- Pair with a normal save on game-end/completion so sendBeacon only covers the early-close case.
- Idempotent payloads — if the normal save already fired, sendBeacon sends a zero/no-op. Server must handle both.
- Content-Type for sendBeacon is automatically `text/plain`. Server handler must accept JSON from sendBeacon (parse body directly, don't rely on Content-Type).

### Screenshot

F9 captures a native screenshot of the webview window using Win32 `PrintWindow` with `PW_RENDERFULLCONTENT`. Captures actual GPU-rendered pixels — works with all CSS including backdrop-filter, custom properties, and design tokens.

- **Shortcut:** F9
- **Mechanism:** JS sends `csPost('screenshot/save', {})` → Go handler calls `engine.CaptureWindow()` → Win32 PrintWindow → BGRA→RGBA → PNG encode → saves timestamped file to `tmp/` folder.
- **Rules:**
  - Do not use html2canvas or any DOM-to-canvas library — they cannot handle CSS custom properties or backdrop-filter.
  - `engine.WindowHandle` must be set from `main.go` after webview creation.
  - Screenshots save to `tmp/` folder with timestamp filenames (`screenshot_YYYYMMDD_HHMMSS.png`).

### Data Queries

- **Rule:** If a view displays a list of records from binary, use `BinaryFindPage` for paginated queries. If it loads a fixed data set (game questions, challenge definitions, startup data), use `BinaryFindAll`.
- **Filtering:** Per-user lookups use `BinaryFind(filter)` — never `BinaryFindAll` + loop. Filter keys must be indexed fields. Exact match: `{"userId": "abc"}`. Range: `{"score": {"$gte": 80}}`. Compound: multiple keys in one filter. Operators: `$gt`, `$gte`, `$lt`, `$lte`, `$ne`, `$in`.
- **Pagination:** `BinaryFindPage(filter, QueryOpts{Page, PageSize, SortBy, SortDir})` for any user-facing record list. Returns results + total count for prev/next UI.
- **Sorting:** Sort fields must be indexed (`index: true` in schema). Default sort is `_id` (insertion order).
- **Indexing:** Mark fields with `index: true` if they are used for filtering, sorting, or pagination. One rule — three capabilities.
- **Type safety:** Index values are type-coerced by the engine based on schema field type. Numeric fields stored as numbers, strings as strings, timestamps as int64. Callers pass whatever they have — the engine enforces correctness. Range queries compare correctly because types are guaranteed.

---

## Platform

ChefScript is a batteries-included desktop app engine. Beyond UI rendering it provides a full server, database, session management, security, and a Python bridge for data-heavy tasks.

### HTTP Server

`net/http` server on `:7070`. WebView2 navigates to localhost — no browser needed.

**Routing:**

- `/page/:name` — Loads `pages/:name.json`, runs page loader, applies template context, renders HTML.
- `/api/:action` — Calls registered Go action handler, returns ActionResult JSON.
- `/public/` — Static file server.

### MongoDB

`go.mongodb.org/mongo-driver/v2`. Call `engine.Col(name)` to get a collection anywhere.

**Helpers:** `engine.ConnectMongo(uri, db)`, `engine.Col(name)`

### Sessions

Cookie-backed sessions stored in MongoDB with TTL index (24h). Automatically loaded into every request context.

**Helpers:** `engine.CreateSession(w, data)`, `engine.DestroySession(w, r)`, `engine.GetSessionFromCtx(r)`

### Security

- **Sanitize:** `SanitizeMiddleware` runs on every POST/PUT/PATCH — strips HTML, null bytes, non-printable unicode, trims whitespace, enforces 1024 char max per field.
- **Hashing:** `engine.HashPassword(plain)` and `engine.CheckPassword(hash, plain)` — bcrypt cost 12. `models.CreateUser` auto-hashes passwords. `models.UpdatePassword` is the only write path.

### Page System

Register a loader for any page to inject live data or redirect before render.

**Helpers:** `engine.RegisterPage(name, func(r) *PageContext)`, `engine.NewPageContext()`

**Context variables:**

- `{{user.id}}` — Logged-in user ID from session.
- `{{user.hackerName}}` — Logged-in user hacker name.
- `{{user.firstName}}` — First name.
- `{{user.lastName}}` — Last name.
- `{{user.role}}` — Role — student | instructor.
- `{{data.key}}` — Arbitrary data injected by the page loader.
- `{{session.key}}` — Raw session values.

### Action System

Register Go functions as named API actions. Called from JSON pages via `on:click`.

**Helpers:** `engine.RegisterAction(name, handler)`, `engine.APIHandler(fn)`

**Result fields:**

- **Error:** Show error toast + re-enable submit button.
- **Fields:** `map[string]string` — highlight specific form inputs with per-field error messages.
- **Toast:** Show a floating notification.
- **ToastVariant:** `info | success | warning | error`
- **Data:** Return data to the client. If Data contains a `target` field, the JS runtime patches the DOM element matching that data-id. See data_targeting.

### Form Handling

Full lifecycle managed by the engine JS runtime. No custom JS needed.

**on_click syntax:** `post:actionName:formId`

**Lifecycle:**

1. Button clicked → disabled + label becomes Loading…
2. Form fields collected, previous field errors cleared
3. POST `/api/:action` with JSON body
4. On error → error toast shown, field errors highlighted, button re-enabled
5. On `data.target` → DOM element patched in place
6. On toast → notification shown, button re-enabled

#### Data Targeting

In-place DOM updates without page reload. When an action returns Data, the JS runtime finds elements by data-id and patches their content. Supports single or multi-target patching.

**When to use:** All action responses. Terminal commands, live feedback, chat responses, challenge validation, form submissions — every handler returns patches.

**Patch types:**

##### Target Patch

- `target` — Required. The data-id of the DOM element to update.
- `html` — Required. HTML string to inject into the element.
- `append` — Optional boolean. true = append to existing content. false/omitted = replace content.
- `scroll` — Optional boolean. true = scroll element to bottom after update. Useful for terminal/chat output.

##### setAttr Patch

- `setAttr` — Required. CSS selector for the target element.
- `attrs` — Required. Object of attribute key-value pairs to set on the element.
- Sets attributes on an element by CSS selector. Use for toggling data attributes, visibility, classes.
- Example: `{"setAttr": ".db", "attrs": {"data-file-open": "true"}}`

##### Single Target

Return Data as a single object to patch one element.

**Go example:** `return ActionResult{Data: map[string]interface{}{"target": "terminal-output", "html": "<div>New line</div>", "append": true, "scroll": true}}`

**Effect:** Appends to the element with `data-id='terminal-output'` and scrolls to bottom.

##### Multi Target

Return Data as an array of patch objects to update multiple elements from a single action. The runtime iterates the array and applies each patch to its target.

**Go example:** `return ActionResult{Data: []map[string]interface{}{{"target": "terminal-output", "html": "<div>response</div>", "append": true, "scroll": true}, {"target": "status-bar", "html": "Connected: VAULT"}, {"target": "prompt", "html": "vault>"}}}`

**Effect:** One command updates terminal output, status bar, and prompt simultaneously — no page reload.

**When to use:** When a single user action affects multiple parts of the UI. Terminal + status bar, chat + notification count, form + preview pane.

##### JSON Example

```json
["div", { "data-id": "terminal-output", "class": "terminal" }]
```

##### Data Targeting Rules

- The target element must exist on the page with a matching data-id.
- The action returns Data — DOM elements are patched in place, no page reload.
- Data can be a single object (one target) or an array of objects (multi-target).
- Toast and Data can be combined — show a notification AND update element(s).
- The Go action is responsible for rendering HTML in `Data.html`. The JSON page just declares the container.
- State should be persisted server-side (session or DB) so a page refresh reconstructs the same view.

#### Field Errors

Return `Fields: map[string]string` from ActionResult to highlight specific inputs with per-field messages.

#### Autosave

Automatic localStorage checkpointing for any form. Off by default — zero overhead unless enabled.

**Enable:** Add `data-autosave="unique-key"` to any form element.

**Behaviour:**

- On input/change → debounced 600ms save to localStorage
- On page load → silently restores all field values
- On successful submit → checkpoint cleared
- On error → checkpoint preserved, no data lost

**Example:** `["form", { "id": "form-profile", "data-autosave": "profile-checkpoint" }]`

**JS API:**

- `csAutosave.save(key, formEl)` — Manually trigger a save.
- `csAutosave.restore(key, formEl)` — Manually restore a checkpoint.
- `csAutosave.clear(key)` — Manually clear a checkpoint.

### Excel Export

Stream an xlsx file download from any Go handler. Python openpyxl builds the file — no Go Excel dependency needed.

- **Go API:** `engine.DownloadExcel(w http.ResponseWriter, filename string, rows []map[string]interface{}) error`
- **Rows format:** Slice of maps — keys become column headers (derived from first row), values become cells. Missing keys default to empty string.
- **Styling:** Headers styled with dark fill and cyan font. Columns auto-fit to content width (max 60 chars).
- **Go example:** `engine.DownloadExcel(w, "students.xlsx", []map[string]interface{}{{"name": "Alice", "score": 92}})`
- **Python module:** `export.to_excel` — callable directly via `engine.CallPython("export", "to_excel", args)` if needed

### Chat Widget

Floating chat bubble (bottom-right) backed by an n8n webhook. Glass morphism styled. Session-scoped chatId for memory tracking. Multiple instances supported.

- **Component:** `chat-widget`
- **Props:**
  - `webhook` — Required. Full n8n webhook URL.
  - `route` — Optional. Passed as route field in the POST body. Default: general.
  - `title` — Optional. Header label. Default: Chat.
  - `data-id` — Required for testing.
- **Behaviour:**
  - Bubble click → opens glass panel
  - Message sent → POST `{chatId, message, route}` to webhook
  - chatId scoped to sessionStorage per widget instance
  - Typing indicator shown while awaiting response
  - Response rendered from `data.output` field
  - Network error shown inline
- **Example:** `["chat-widget", { "webhook": "https://your-n8n/webhook/.../chat", "route": "general", "title": "Chat", "data-id": "chat--support" }]`

### Python Bridge

Python subprocess runs alongside the Go server. Go sends JSON over stdin, reads JSON from stdout. Calls are serialized (mutex). The bridge is an open pipe — any Python module can be added.

- **Start:** `engine.StartPython(scriptPath)` — called once in `main.go`
- **Stop:** `engine.StopPython()` — deferred in `main.go`
- **Call:** `engine.CallPython(module, method, args) → (interface{}, error)`
- **Args format:** Pass `map[string]interface{}` for keyword args, `[]interface{}` for positional.
- **Adding modules:** Create a Python file in `python/`. Define functions that accept JSON-serializable args and return JSON-serializable results. Call via `engine.CallPython(module, method, args)` from Go.

**Built-in modules:**

- `export.to_excel` — Build an xlsx file from a list of dicts. Args: rows (list of dicts), filename (string), sheet (string, optional). Returns base64-encoded file — use `engine.DownloadExcel` in Go instead of calling this directly.

**Example:**

- **Go:** `data, err := engine.CallPython("my_module", "my_method", map[string]interface{}{"key": "value"})`
- **Response:** `data` is `interface{}` — type-assert to `map[string]interface{}` or `[]interface{}` as needed

### Partial Navigation

SPA-like page transitions without a client framework. When navigating to `/page/*`, the runtime fetches `/partial/*` instead, swaps `document.body.innerHTML`, updates title and history. No full reload.

**How it works:**

1. `csNavigate(url)` converts `/page/name` to `/partial/name`
2. Fetches JSON: `{ html, title, url }`
3. Swaps `document.body.innerHTML` with response html
4. Updates `document.title` and `pushState`
5. Scrolls to top. MutationObserver re-wires `on:enter` listeners.
6. On error, falls back to full page load

**What persists:** `window.*` globals, JS closures, flight recorder client buffer. The runtime IIFE does NOT re-execute.

**What resets:** All DOM state, diagnostics panel visibility, scroll position.

**Link interception:** All `<a href="/page/...">` clicks are intercepted and routed through csNavigate automatically. No special markup needed.

### Data References

Two systems for injecting data into pages.

#### Template Variables

- **Syntax:** `{{path}}` in any string value in page JSON
- **Sources:**
  - `{{user.*}}` — Session user fields — id, hackerName, firstName, lastName, role
  - `{{data.*}}` — Arbitrary data from page loader
  - `{{session.*}}` — Raw session values
- **Resolution:** Template engine resolves before render. Single-var strings preserve type (number, array, object). Mixed strings interpolate to string.

#### @data References

- **Syntax:** `@data/filename.json` in any prop value
- **Does:** Engine reads the JSON file from the data directory and passes the parsed result as the prop value.
- **Rule:** Invalid references produce an engine error, not a silent failure.

---

## Data Flow

Schema-based data flow. The schema IS the connection between database and component. No page loaders, no template variable string-replacement, no JSON round-trips. Binary decodes to native values, schema routes them directly to the component.

**Principle:** Parse once at the boundary. Then it's just variable assignment — O(1). Data stays as native values from decode to render.

### Get

Component requests data. Schema connects it directly to the binary collection.

- **Flow:** component declares schema → engine loads binary → schema decodes → native values passed as props
- **Syntax:** Any component can declare a `schema` prop. The engine resolves it before render.
- **Props:**
  - `schema` — Binary schema name (matches collection field in `schemas/binary/*.json`).
  - `filter` — Optional. Object of field:value pairs to filter records.
  - `sort` — Optional. Field name to sort by.
  - `limit` — Optional. Max records to return.
- **Example:** `["challenge-list", { "schema": "challenges", "filter": { "track": "linux" }, "data-id": "cl--linux" }]`

### Post

Form sends data. Schema defines what fields to collect and where to store.

- **Flow:** form declares schema → engine collects fields matching schema → schema encodes → binary insert
- **Syntax:** Form can declare a `schema` prop. On submit, the engine maps form fields to schema fields automatically.
- **Example:** `["form", { "id": "join-form", "schema": "users", "action": "post:auth/join" }]`

### How It Works

1. Component has schema prop → engine intercepts before render
2. Engine calls `GetBinarySchema(name)` → gets the schema with field definitions
3. Schema calls `BinaryFindAll()` or filtered variant → raw bytes from MongoDB
4. `Schema.Decode()` reads bytes at fixed offsets → native Go values in one pass
5. Values passed directly to component as props — no JSON marshal/unmarshal, no string replacement
6. Component renders with native values. Zero intermediate serialization.
7. For POST: reverse the flow. Form values → `Schema.Encode()` → binary bytes → MongoDB.

### Speed

Binary decode is fixed-offset reads + in-memory lookup table. O(fields) per record. No parsing, no scanning. Values reach the component as native types — int, string, []string — never re-serialized.

---

## Binary Schema Format

Schema files define the binary encoding for a MongoDB collection. Located in `schemas/binary/*.json`. Loaded at engine startup via `LoadBinarySchemas()`.

### File Structure

- `collection` — Logical collection name used in component schema props
- `binaryCollection` — Physical MongoDB collection storing the binary records
- `fields` — Array of field definitions

### Field Definition

- `name` — Field name — used in data access and component props
- `type` — One of: `uint8`, `uint16`, `uint32`, `uint64`, `lookup`, `lookup16`, `timestamp`
- `values` — Optional array of strings — lookup table values for lookup/lookup16 types
- `dynamic` — Optional boolean — if true, lookup values are discovered at runtime and persisted to `_binary_lookups` collection
- `index` — Optional boolean — if true, field value is duplicated as a top-level MongoDB field alongside the binary blob. Enables `BinaryFind(filter)` for targeted queries instead of `BinaryFindAll()` + loop. MongoDB index is created automatically on schema load. Values are type-coerced by the engine: numeric fields (uint8/16/32/64) stored as numbers, lookup/lookup16 as strings, timestamps as int64. This ensures range queries (`$gt`, `$gte`, `$lt`, `$lte`) compare correctly.

### Types

- `uint8` — 1 byte, 0–255
- `uint16` — 2 bytes, 0–65535
- `uint32` — 4 bytes, 0–4.2B
- `uint64` — 8 bytes, large integer
- `lookup` — 1 byte, maps integer ID → string via lookup table (max 255 values)
- `lookup16` — 2 bytes, maps integer ID → string via lookup table (max 65535 values)
- `timestamp` — 8 bytes, Unix epoch, decoded to RFC3339 string at render

### Lookup Types

- **Fixed:** Values known at schema time. Specified in values array. Unknown values map to ID 0.
- **Dynamic:** Values discovered at runtime. Each unique string gets a new ID. Table persisted to `_binary_lookups` MongoDB collection across restarts.

### Encode/Decode

- **Encode:** `map[string]interface{}` → fixed-size byte array → stored as MongoDB `Binary(0x00, bytes)` in field `d`. Indexed fields also written as top-level MongoDB fields.
- **Decode:** Binary bytes from MongoDB → fixed-offset reads → native Go values in one pass → passed as component props

### Operations

- `BinaryInsert(doc)` — Encode + insert. Indexed fields written alongside blob.
- `BinaryInsertMany(docs)` — Batch encode + insert.
- `BinaryFindAll()` — Load all records, decode all blobs. Returns `[]map` with `_id`.
- `BinaryFind(filter)` — Query MongoDB using indexed top-level fields, decode only matching blobs. O(matched) not O(all). Filter keys must be indexed fields. Exact match: `{"userId": "abc"}`. Range: `{"score": {"$gte": 80}}`. Operators: `$gt`, `$gte`, `$lt`, `$lte`, `$ne`, `$in`. Values are coerced to correct types by the engine — numeric fields stored as numbers, strings as strings, timestamps as int64.
- `BinaryFindPage(filter, opts)` — Paginated + sorted query. Returns `(results, totalCount, error)`. QueryOpts: Page (0-based), PageSize (default 20), SortBy (indexed field, default `_id`), SortDir (1=asc, -1=desc). Filter can be nil for all records. Sort field must be indexed.
- `BinaryUpdate(id, doc)` — Re-encode + replace by `_id`. Index fields updated.
- `BinaryDelete(id)` — Remove by `_id`.

### Example Schema

```json
{
  "collection": "challenges",
  "binaryCollection": "challenges_bin",
  "fields": [
    { "name": "track", "type": "lookup", "values": ["linux", "python", "sql", "network"], "dynamic": false },
    { "name": "difficulty", "type": "uint8" },
    { "name": "title", "type": "lookup16", "values": [], "dynamic": true },
    { "name": "created", "type": "timestamp" }
  ]
}
```

---

## Data Resolution

Generalized pattern: component declares a data dependency, engine resolves it. The source can be binary schema (MongoDB), static JSON file, or Go map — the component doesn't know or care. Same contract as schema pipe, extended to all data.

### Sources

- **binary_schema:** Dynamic data from MongoDB via binary encode/decode. Use for: user data, progress, challenges — anything that changes at runtime.
- **static_json:** JSON files loaded at engine startup. Use for: reference data, lookup tables, descriptions — anything that doesn't change at runtime but should be editable without recompiling.
- **go_map:** Constants defined in Go code. Use for: truly fixed data that never changes. Prefer static_json unless there's a reason to hardcode.

**Principle:** Component declares need → engine fills it. If it needs to be editable, put it in `data/`. If it needs to be queryable, put it in binary schema. If it's a constant, put it in Go. The component code is the same either way.

### Example

Command descriptions for DirtByte — loaded from `data/command_descriptions.json` at startup, resolved by component via `GetCommandDescription()`.

- **Data file:** `data/command_descriptions.json`
- **Load:** `engine.LoadCommandDescriptions(path)` — called once in `app.go`
- **Resolve:** `engine.GetCommandDescription(cmd) → (syntax, description)`
- **Usage:** Component calls `GetCommandDescription` in Go render function. No props needed — engine resolves directly.

---

## Atom Format

**Syntax:** `[tag, props?, ...children]`

**Rules:**

- Position 0: string — the component name.
- Position 1: if an object, it is props. If not an object, it is a child.
- Position 2+: children — strings or nested atoms.
- Props are always optional. Omit the object entirely if none needed.
- Children can be strings (raw text) or arrays (nested atoms).
- Nesting is unlimited.

**Examples:**

```json
["text", "Hello world"]
["button", { "label": "Go" }]
["button", { "variant": "solid" }, "Click me"]
["stack", { "direction": "row", "spacing": 3 }, ["text", "Left"], ["text", "Right"]]
["card", ["heading", { "level": 3 }, "Title"], ["text", "Body"]]
```

---

## Page Format

Top-level JSON object. Only `body` is required.

- `title` — string, optional, default "ChefScript". Window title.
- `theme` — string, optional, default "dark". dark = black bg, cyan accents. light = warm white, navy accents.
- `body` — array, required. Array of atoms. The entire UI tree.

**Example:**

```json
{
  "title": "My App",
  "theme": "dark",
  "body": [
    ["nav", { "brand": "MyApp" }, ["nav-link", { "href": "#" }, "Home"]],
    ["box", { "p": 4 }, ["header", { "title": "Dashboard" }]],
    ["footer", "Built with ChefScript"]
  ]
}
```

---

## Agent Rules

Critical rules. Follow these exactly.

- DIALOGS and DRAWERS must be placed at the top level of body, NOT inside box or any other wrapper. They are hidden overlays that position themselves fixed on screen.
- TABS: set `active:true` on exactly one tab per tabs group. `panel` values must be unique within the group.
- ACCORDION: `accordion-item` must be a direct child of `accordion`.
- BREADCRUMB: `breadcrumb-item` must be a direct child of `breadcrumb`. Last item has no href.
- AUTOCOMPLETE and SELECT dropdowns open over other content — do not wrap them in overflow:hidden containers.
- GRID: use `container`/`item` props with breakpoint column spans (xs/sm/md/lg/xl). STACK: use for simple flex layouts with spacing.
- All interactive elements automatically get ripple, floating labels, and keyboard navigation — no extra setup needed.
- SNACKBAR: triggered by `on:click` action strings, not placed in the JSON as a visible component.
- `data-id`: all components auto-generate a data-id. Override with the `data-id` prop for specific targeting.
- JSON is the map. Never embed logic in JSON. If a component needs looping or conditions, that is the engine's job.

---

## Events

Wire behaviors via `on:click` and `on:enter` action strings. No JS required.

### Actions

- `modal:{id}` — Open modal with matching id prop.
- `close-modal:{id}` — Close a specific modal.
- `drawer:{id}` — Open drawer with matching id prop.
- `close-drawer:{id}` — Close a specific drawer.
- `snackbar:{message}:{variant}` — Show a snackbar toast. Variant: `info | success | warning | error`.
- `post:{action}` — POST to `/api/{action}`. Collects form data if formId provided. Used by `on:click` and `on:enter`.
- `post:{action}:{formId}` — POST to `/api/{action}` with form data from formId.
- `navigate:{url}` — Navigate to url.

### on:click

Fires when a button is clicked. Value is an action string.

**Example:** `["button", { "on:click": "post:save-record:form-edit" }, "Save"]`

### on:enter

Fires when Enter is pressed on an input. Runtime wires the keydown listener automatically. Input value is sent as `value` in the POST body. All `data-*` attributes on the input are included as extra fields.

**Props:**

- `on:enter` — Action string (e.g. `post:dirtbyte/command`). Required.
- `clear` — Boolean. If true, input is cleared after successful action.
- `data-*` — Any `data-*` attributes are sent as extra fields in the POST body.

**Flow:**

1. User presses Enter on input with `on:enter` prop.
2. Runtime collects input value + all `data-*` attributes.
3. If action starts with `post:` → POST `/api/{action}` with merged data.
4. Response is handled by csPost — data patches, toasts, and errors.
5. If clear prop is set, input value is cleared on success.

**Examples:**

```json
["input", { "data-id": "terminal-input", "on:enter": "post:dirtbyte/command", "data-mission": "m01", "clear": true }]
["input", { "data-id": "search-input", "on:enter": "post:search", "clear": false }]
["input", { "data-id": "chat-input", "on:enter": "post:chat/send", "data-channel": "support", "clear": true }]
```

**Escape:** When ESC is pressed on an `on:enter` input, the runtime dispatches a `cs:escape` CustomEvent that bubbles. Components can listen for this to close panels.

### Event Examples

```json
["button", { "on:click": "modal:confirm-delete" }, "Delete"]
["button", { "on:click": "close-modal:confirm-delete" }, "Cancel"]
["button", { "on:click": "drawer:settings" }, "Settings"]
["button", { "on:click": "snackbar:Record saved:success" }, "Save"]
["button", { "on:click": "snackbar:Connection failed:error" }, "Retry"]
["input", { "on:enter": "post:dirtbyte/command", "data-mission": "m01", "clear": true }]
```

---

## Quick Reference

All 100+ atom components grouped by category. Each category maps to `engine/components_atom_<category>.go`.

- **Layout:** header, footer, card, sidebar, section, divider, split-view, split-pane, paper, app-bar, box, grid, stack, image-list, image-list-item, masonry
- **Navigation:** nav, nav-link, tabs, tab, accordion, accordion-item, breadcrumb, breadcrumb-item, pagination, stepper, stepper-step, toolbar, bottom-nav, bottom-nav-action, menubar, menubar-item, speed-dial, speed-dial-action
- **Inputs:** input, textarea, select, native-select, autocomplete, form-field, multi-select, checkbox, radio, switch, form, slider, number-input, file-upload, tag-input, date-input, search, color-input, fab, toggle-group, toggle-button, transfer-list
- **Actions:** button, icon, icon-button, button-group, copy-button
- **Display:** heading, text, avatar, avatar-group, empty-state, kbd, code, code-block, timeline, timeline-item, rating, callout, image, link, tag
- **Feedback:** alert, badge, chip, spinner, skeleton, progress, tooltip, banner
- **Overlay:** menu, menu-item, popover, drawer, snackbar, confirm, notification, notification-item, command, command-item, context-menu, hover-card, backdrop, dialog, dialog-title, dialog-content, dialog-actions
- **Data:** stat-card, table, list, list-item, kv-list, kv-item, data-grid, tree, tree-item, virtual-list, chart, calendar
- **Media:** video, audio, iframe, aspect-ratio, carousel, rich-text
- **Chat:** chat-widget, data-chat

---

## Icons

Icons live in `/public/icons/`. Use the filename without `.svg` as the name prop.

- **Naming:** Files are named descriptively: `rocket.svg`, `star.svg`, `gear.svg`, `trash.svg`, `plus.svg`, `check.svg`, `bell.svg`, `search.svg`, `home.svg`, `heart.svg`, `lightning.svg`, `shield.svg`, `arrow-back.svg`, `bell-fill.svg`, etc.
- **Variants:** Many icons have `-fill` suffix for filled variant: `bell.svg` vs `bell-fill.svg`.
- **Usage:** `["icon", { "name": "rocket", "size": 24 }]`

---

## Build Recipes

Step-by-step patterns for solving common build problems. Each recipe is a checklist. Follow the steps, use the pieces, get the result. Reference implementation: DirtByte (`engine/dirtbyte.go`, `engine/components_dirtbyte.go`).

### Recipe 1: Same-Page Interaction

**Problem:** User does something, server responds, UI updates without navigation.

1. Add `data-id` to every element that will be updated
2. Add `onclick='csAction("post:action/name:formId",this)'` to the trigger (button inside component) OR `on:click='post:action/name:formId'` on an atom
3. Wrap inputs in a `<form id='formId'>` with named fields
4. Register handler: `engine.RegisterAction('action/name', handler)`
5. Handler reads body, does work, returns `ActionResult{Data: map[string]interface{}{"target": "data-id", "html": "..."}}`
6. Runtime applies patch — element found by data-id, innerHTML replaced
7. For multiple updates, return Data as `[]interface{}` with multiple patch objects

**Reference:** `engine/dirtbyte.go` — `dirtbyteCommandHandler` returns multi-target patches

### Recipe 2: Stateful Multi-Step Flow

**Problem:** User progresses through steps, state persists between interactions.

1. Define state structure: `map[string]interface{}` with flags, history, phase, etc.
2. Store in session: `sess.Data['feature_state'][itemId] = state`
3. Write `getState(sess, id)` helper — returns state with defaults if missing
4. Write `saveState(sess, id, state)` helper — writes back + calls `UpdateSession`
5. Handler loads state, checks progression (gate logic), mutates state, saves, returns patches
6. Page loader reads state to render current view on page load
7. Gate logic: check required flags/conditions before allowing action. Wrong time = hint message. Right time = grant flag + advance

**Reference:** `engine/dirtbyte.go` — `getDirtbyteState/saveDirtbyteState`, `checkProgression` gate logic

### Recipe 3: Custom Full-Page Layout

**Problem:** Page needs its own look beyond standard atoms.

1. Create `engine/components_feature.go`
2. Write render function: `func(props, children, engine) (string, error)`
3. Build HTML with `strings.Builder` — use data-id on patchable elements
4. Write CSS function returning scoped styles with unique prefix (e.g. `.pylab__`)
5. Append `<style>` tag: `b.WriteString("<style>" + featureCSS() + "</style>")`
6. Register: `e.Register('feature-name', ComponentFunc(renderFeature))`
7. Page JSON: single atom `['feature-name', { props }]`
8. Register component in `app.go`: `engine.RegisterFeatureComponents(e)`

**Reference:** `engine/components_dirtbyte.go` — `renderDirtbyteActive`, `dirtbyteActiveCSS()`

### Recipe 4: Data-Driven List with Selection

**Problem:** User picks from a list, page updates to show the selected item.

1. Page loader reads query param: `r.URL.Query().Get('c')`
2. Load full list from cache/binary
3. Load selected item by ID, inject both into `ctx.Data`
4. Component renders list — each item is a link `<a href='/page/feature?c=itemId'>`
5. Runtime intercepts link click → partial nav → body swap with new selection
6. For dropdowns: `<select data-navigate='/page/feature?c={value}'>` — runtime wires change event to csNavigate

**Reference:** `engine/challenge_validator.go` — `challengePageLoader`, `engine/runtime.go` — `csWireNavigateSelect`

### Recipe 5: Interactive Input with Live Feedback

**Problem:** User types, presses enter or clicks, gets immediate response in place.

1. Render input with `data-on-enter='post:action/name'` — runtime auto-wires Enter key
2. Add `data-*` attributes for extra payload (`data-mission`, `data-challenge`, etc.)
3. Add `data-clear` to auto-clear input after submit
4. Handler receives `{value: 'user input', ...extraData}`
5. Handler returns target patch to output area with `append:true`, `scroll:true`
6. For buttons instead of Enter: use `onclick='csAction("post:action:formId",this)'` with `type='button'`

**Reference:** `engine/components_dirtbyte.go` line 1435-1440 — terminal input with `data-on-enter`

### Recipe 6: Layout Switching Without Re-render

**Problem:** Toggle visibility or position of UI regions without a server call.

1. Add data attributes to root element: `data-file-open='false'`, `data-active-panel=''`
2. Write CSS selectors that match attribute values: `.root[data-file-open='true'] .panel { display: flex; }`
3. Toggle via JS: `onclick='element.dataset.fileOpen = "true"'` or via handler setAttr patch
4. For handler-driven: return `{setAttr: '.root', attrs: {data-active-panel: 'cmds'}}`
5. CSS handles all layout changes — no re-render, no server call for pure UI toggles

**Reference:** `engine/components_dirtbyte.go` — `dirtbytePanelCSS`, `data-file-open`/`data-active-panel` selectors

### Recipe 7: Multi-Target Response

**Problem:** One action needs to update multiple parts of the UI simultaneously.

1. Give each updatable region a unique `data-id`
2. Handler returns Data as `[]interface{}` array of patch objects
3. Each patch: `{target: 'data-id', html: '...', append: bool, scroll: bool}`
4. Mix target patches and setAttr patches in the same array
5. Runtime iterates and applies all patches in order

**Example return:** `ActionResult{Data: []interface{}{{"target": "output", "html": "<div>result</div>", "append": true, "scroll": true}, {"target": "status", "html": "Connected"}, {"setAttr": ".root", "attrs": {"data-state": "active"}}}}`

**Reference:** `engine/dirtbyte.go` — `dirtbyteCommandHandler` returns terminal + hint + prompt + file patches

### Recipe 8: Binary Data into Component

**Problem:** Component needs records from the database.

1. Define schema in `schemas/binary/collection.json` — field names, types, dynamic flags
2. Seed data: `cmd/seedchallenges` pattern — read JSON, serialize nested fields, `BinaryInsertMany`
3. Load at startup: `LoadFromBinary()` reads schema, decodes all records, caches in memory
4. Page loader calls `GetData()` from cache, injects into `ctx.Data`
5. Template var `{{data.items}}` resolves to the array
6. Component receives as prop, iterates to render
7. For schema-driven resolution: component declares schema prop, engine intercepts before render

**Reference:** `engine/challenges.go` — `LoadChallengesFromBinary`, `ChallengeToClientMap`

### Recipe 9: Context-Preserving Navigation

**Problem:** Moving between related pages without losing the user's place.

1. Links carry context as query params: `<a href='/page/feature-mode?c=itemId'>`
2. Page loader reads query param and loads the relevant data
3. Both pages share the same page loader (or use `ChallengePageLoaderFor` pattern)
4. Runtime intercepts link → partial nav → body swap. URL updates, context preserved
5. Component renders the active state based on props (mode, selectedId, etc.)
6. Register shared loader: `engine.RegisterPage('feature-mode', loaderFunc)`

**Reference:** `pages/python.json` + `python-sandbox.json` — same loader, different mode prop, linked via mode toggle

### Recipe 10: Overlay Toggling

**Problem:** Show/hide drawers, modals, or side panels.

1. Use built-in atoms: `['drawer', {id: 'name', title: '...', side: 'right'}]`
2. Trigger with `on:click='drawer:name'` on a button atom
3. For dialogs: `['dialog', {}, ['dialog-title', ...]]` + `on:click='modal:name'`
4. For custom panels inside components: use `data-active-panel` attribute + CSS (recipe 6)
5. Runtime handles open/close via `csDrawer.open(id)` / `csModal.open(id)`
6. Backdrop click closes automatically

**Reference:** `pages/linux.json` — drawer for blackbook, `engine/runtime.go` — `csDrawer`/`csModal`

### Recipe 11: Client-State Game

**Problem:** Multi-step interaction (game, quiz, builder) needs state that persists across many user actions without server round trips.

1. Go component embeds all data and calls `csState.init('feature', {...})` with initial state
2. Add `data-state='feature.key'` to elements that display state values
3. Add `data-state-attr='feature.key:attrName'` to elements that bind state to attributes
4. Write action functions that read `window.__feature`, apply logic, call `csState.set('feature', {changes})`
5. Validation is client-side — compare against data already in the state object, not a server call
6. On completion, one csPost sends final result. Server writes via binary schema.
7. Page loader only runs once (get data). Action handler only runs once (save result). Everything between is JS.

**Reference:** `engine/wire.go` — `wireGameLoader` embeds batches, `engine/js/wire.js` — wire game interaction, `dictionary.json` `runtime.client_state` — full API

### Recipe 12: JS Runtime Split

**Problem:** The JS runtime was a single 5000-line Go raw string. Needed splitting for maintainability without losing the single-output behavior.

**Solution:** `go:embed` — each JS section lives in its own file under `engine/js/`. At compile time, Go embeds them as strings. `runtime.go` concatenates them inside one `<script>` IIFE.

**Structure:**

- `engine/js/core_helpers.js` — Shared helpers — escH, shuffle, flight recorder (~13 lines). Loads first. All other files depend on these.
- `engine/js/core_components.js` — UI component behaviors — ripple, tabs, accordion, modal, drawer, snackbar, autocomplete, select, rating, slider, number input, tag input, file upload, code copy, menu, popover, stepper, autosave (~485 lines)
- `engine/js/core_forms.js` — Form utilities — csPost, csAction, form serialization, button loading, on:enter (~202 lines)
- `engine/js/core_nav.js` — Navigation — csNavigate, partial page swap, DirtByte panel wiring, ESC handler, notepad autosave, terminal scroll (~211 lines)
- `engine/js/core_data.js` — Data components — error capture, chat widget, DataGrid, tree, virtualList, notification, command palette, context menu, split view, calendar, multi-select (~405 lines)
- `engine/js/core_desktop.js` — Desktop OS — window management, taskbar, clock, data chat, binary fetch, page data renderers (~438 lines)
- `engine/js/state.js` — csState client-side state system + splitEmoji helper (~77 lines)
- `engine/js/[game].js` — One file per game. Each is a self-contained IIFE or set of window functions. Games: wire, cybermoji, cyberstrike, interpreter, decode_alert, cyber_defense, social_engineer, escalation_ladder, defense_tower, network_builder, digital_heist, cyber_reigns
- `engine/js/admin.js` — Admin dashboard JS — navigation, student detail renderer, accordions, CRUD
- `engine/js/quiz.js` — Quiz JS — answer selection, submit, scorecard renderer
- `engine/runtime.go` — Assembler — go:embed directives + `buildRuntime()` concatenates all files inside `<script>(function(){'use strict'; ... })();</script>`

**Rules:**

- All JS lives in `engine/js/` — never inline in components or page templates
- Each game gets its own file. New game = new `.js` file + new embed directive in `runtime.go`
- Game IIFEs are self-contained: `(function(){...})();` — local vars don't leak
- Shared functions use `window.xxx = function(){}` — accessible across files
- Shared helpers (`escH`, `shuffle`) are defined ONCE in `core.js` — do NOT redefine locally in game files. Use them directly. They are in scope because all files share the outer IIFE.
- Exception: if a game needs different semantics (e.g. `wire.js` shuffle mutates in-place instead of returning a copy), define locally with a comment explaining why.
- Alias pattern: if a game used a different name historically (e.g. `escHtml`), alias it: `var escHtml=escH;` — don't duplicate the implementation.
- The outer IIFE wrapper is in `runtime.go` `buildRuntime()` — not in any `.js` file
- `go build` handles everything — no external bundler, no npm, no watch process. All assets embedded into the binary via `go:embed`. See conventions.self_contained_binary.

---

## Patterns

Common page compositions. Copy and adapt.

### Dashboard Page

Nav + stat cards + table.

```json
{
  "title": "Dashboard",
  "theme": "dark",
  "body": [
    ["nav", { "brand": "App" }, ["nav-link", { "href": "#" }, "Dashboard"], ["nav-link", { "href": "#" }, "Jobs"]],
    ["box", { "p": 4 },
      ["header", { "title": "Dashboard", "subtitle": "Today" }],
      ["grid", { "container": true, "spacing": 3 },
        ["grid", { "item": true, "xs": 12, "md": 4 }, ["stat-card", { "label": "Jobs Today", "value": "14", "trend": "up" }]],
        ["grid", { "item": true, "xs": 12, "md": 4 }, ["stat-card", { "label": "Revenue", "value": "$4,200", "trend": "up" }]],
        ["grid", { "item": true, "xs": 12, "md": 4 }, ["stat-card", { "label": "Open", "value": "3" }]]
      ],
      ["table", {
        "columns": ["Client", "Service", "Status"],
        "rows": [["HVAC Co", "AC Install", "Active"], ["Pipe Masters", "Drain Clean", "Scheduled"]]
      }]
    ],
    ["footer", "ChefScript"]
  ]
}
```

### Form Page

A data entry form with inputs, selects, and submit.

```json
{
  "title": "New Job",
  "theme": "dark",
  "body": [
    ["box", { "p": 4 },
      ["header", { "title": "New Job" }],
      ["grid", { "container": true, "spacing": 3 },
        ["grid", { "item": true, "xs": 12, "md": 6 }, ["input", { "label": "Client Name", "name": "client", "required": true }]],
        ["grid", { "item": true, "xs": 12, "md": 6 }, ["input", { "label": "Phone", "name": "phone", "type": "tel" }]]
      ],
      ["grid", { "container": true, "spacing": 3 },
        ["grid", { "item": true, "xs": 12, "md": 6 }, ["select", { "label": "Service Type", "name": "service", "options": ["AC Repair", "Installation", "Maintenance"] }]],
        ["grid", { "item": true, "xs": 12, "md": 6 }, ["select", { "label": "Technician", "name": "tech", "options": ["Alice", "Bob", "Carol"] }]]
      ],
      ["textarea", { "label": "Notes", "name": "notes", "rows": 4 }],
      ["stack", { "direction": "row", "spacing": 2 },
        ["button", { "variant": "solid", "color": "success", "on:click": "snackbar:Job created:success" }, "Create Job"],
        ["button", { "variant": "ghost" }, "Cancel"]
      ]
    ]
  ]
}
```

### Detail with Drawer

List of records. Click opens a drawer with detail.

**Note:** Drawer must be at body level, not inside box or any wrapper.

**Trigger:** `["button", { "variant": "ghost", "size": "sm", "on:click": "drawer:job-detail" }, "View"]`

**Drawer:** `["drawer", { "id": "job-detail", "title": "Job Detail", "side": "right" }, ["text", "Content here"]]`

### Confirm Dialog

Destructive action with confirmation dialog.

**Trigger:** `["button", { "variant": "outline", "color": "danger", "on:click": "modal:confirm-delete" }, "Delete"]`

**Modal:**

```json
["dialog", {},
  ["dialog-title", {}, "Delete Record?"],
  ["dialog-content", {}, ["text", "This cannot be undone."]],
  ["dialog-actions", {},
    ["button", { "variant": "ghost" }, "Cancel"],
    ["button", { "variant": "solid", "color": "danger" }, "Delete"]
  ]
]
```

### Tabbed Detail

Detail view with tabbed sections.

```json
["tabs",
  ["tab", { "label": "Overview", "panel": "ov", "active": true }, ["text", "Overview content"]],
  ["tab", { "label": "History", "panel": "hist" }, ["text", "History content"]],
  ["tab", { "label": "Notes", "panel": "notes" }, ["textarea", { "label": "Add note", "name": "note" }]]
]
```

### Search List

Autocomplete search above a results table.

```json
["stack", { "spacing": 3 },
  ["autocomplete", { "label": "Search clients", "name": "q", "options": ["HVAC Co", "Pipe Masters", "Cool Air"] }],
  ["table", {
    "columns": ["Name", "Phone", "Jobs"],
    "rows": [["HVAC Co", "555-0100", "12"], ["Pipe Masters", "555-0101", "7"]]
  }]
]
```

---

## Components

### card

Elevated block with border, background, padding. Children stack vertically with gap. Primary content container.

**Children:** Any atoms.

```json
["card", ["heading", { "level": 3 }, "Title"], ["text", "Body text"], ["button", { "variant": "ghost", "size": "sm" }, "Action"]]
```

### divider

Horizontal separator line. Optional centered label for section breaks.

**Props:**

- `label` (string) — Text centered in the line. e.g. OR, Section Name.

```json
["divider"]
["divider", { "label": "OR" }]
["divider", { "label": "Personal Information" }]
```

### nav

Top navigation bar. Brand name left, links right. Fixed height 48px with bottom border.

**Props:**

- `brand` (string) — App name displayed in accent color on the left.

**Children:** `nav-link` atoms.

```json
["nav", { "brand": "MyApp" }, ["nav-link", { "href": "#" }, "Home"], ["nav-link", { "href": "#" }, "Jobs"], ["nav-link", { "href": "#" }, "Settings"]]
```

### nav-link

Link inside nav. Muted text that brightens on hover.

**Props:**

- `href` (string, default "#") — Link destination.

**Children:** Link text.

```json
["nav-link", { "href": "/clients" }, "Clients"]
```

### header

Full-width page title bar with optional subtitle. Use at top of page body.

**Props:**

- `title` (string) — Large title text.
- `subtitle` (string) — Small muted text beside title — date, version, context.

```json
["header", { "title": "Job List", "subtitle": "47 records" }]
```

### footer

Page footer. Centered, small, muted text.

**Children:** Footer text.

```json
["footer", "ChefScript v0.2.0"]
```

### heading

Section title. Bold, clean. Use level to control size hierarchy.

**Props:**

- `level` (number, default 2) — 1=largest (page title), 2=section, 3=subsection, 4-6=smaller.

**Children:** Heading text.

```json
["heading", { "level": 2 }, "Section Title"]
```

### text

Body paragraph. Secondary muted color, 1.5 line height. The default readable text block.

**Props:**

- `class` (string) — Extra CSS classes.

**Children:** Text content.

```json
["text", "This is a paragraph of body text."]
```

### button

Interactive button. Ripple effect on click. Supports 4 variants, 4 sizes, 4 semantic colors.

**Props:**

- `label` (string) — Button text. Falls back to children if omitted.
- `variant` (string, default "outline") — `solid | outline | ghost | link`
- `size` (string, default "md") — `xs | sm | md | lg`
- `color` (string) — `primary | success | danger | warning`
- `on:click` (string) — Action string.
- `data-id` (string, default "button")

**Children:** Used as label if label prop is omitted.

```json
["button", { "variant": "solid" }, "Save"]
["button", { "variant": "solid", "color": "danger", "on:click": "modal:confirm-delete" }, "Delete"]
["button", { "variant": "ghost", "size": "sm" }, "Cancel"]
["button", { "variant": "outline", "color": "success" }, "Approve"]
```

### icon

Inline SVG icon. Reads from `/public/icons/{name}.svg`. Scales cleanly, inherits color.

**Props:**

- `name` (string, required) — Filename without .svg.
- `size` (number, default 20) — Pixel size (width and height).
- `color` (string, default "currentColor") — Any CSS color.
- `class` (string)
- `data-id` (string)

```json
["icon", { "name": "rocket", "size": 24 }]
["icon", { "name": "check", "size": 16, "color": "var(--color-success)" }]
```

### stat-card

Metric display card. Small uppercase label, large value, optional trend arrow.

**Props:**

- `label` (string) — Metric name in small uppercase.
- `value` (string, default "0") — The number or value. Pass as string.
- `trend` (string) — `up | down`. up = green arrow, down = accent arrow.

```json
["stat-card", { "label": "Revenue", "value": "$18,240", "trend": "up" }]
```

### input

Text input with floating label. Label sits inside the field and rises on focus (MUI pattern). Supports error and hint states.

**Props:**

- `label` (string) — Floating label.
- `type` (string, default "text") — `text | email | password | number | tel | url`
- `name` (string) — Form field name.
- `value` (string) — Pre-filled value.
- `hint` (string) — Helper text shown below the field.
- `error` (string) — Error message. Turns border red.
- `required` (boolean, default false)
- `disabled` (boolean, default false)
- `data-id` (string)

```json
["input", { "label": "Email", "type": "email", "name": "email", "required": true }]
["input", { "label": "Phone", "type": "tel", "name": "phone", "hint": "Include country code" }]
["input", { "label": "Username", "name": "user", "error": "Username is already taken" }]
```

### textarea

Multi-line text input with floating label. Vertically resizable.

**Props:**

- `label` (string)
- `name` (string)
- `rows` (number, default 4) — Initial visible rows.
- `hint` (string)
- `disabled` (boolean, default false)

```json
["textarea", { "label": "Job Notes", "name": "notes", "rows": 5 }]
```

### select

Custom styled dropdown. Click to open, click option to select. Keyboard navigable.

**Props:**

- `label` (string) — Label shown above the trigger.
- `name` (string) — Form field name for the hidden input.
- `placeholder` (string, default "Select...") — Default display text.
- `options` (array) — Array of option strings.

```json
["select", { "label": "Status", "name": "status", "options": ["Active", "Inactive", "Pending", "Cancelled"] }]
```

### autocomplete

Search input with live-filtered dropdown. Type to filter options. Arrow keys navigate, Enter selects, Escape closes.

**Props:**

- `label` (string)
- `name` (string)
- `placeholder` (string) — Placeholder text in the input.
- `options` (array) — Array of option strings to filter through.

```json
["autocomplete", { "label": "Search Clients", "name": "client", "options": ["HVAC Solutions", "Plumbing Pro", "Electric Plus", "Cool Air Co"] }]
```

### checkbox

Styled checkbox with label. Uses custom box — no browser default styling.

**Props:**

- `label` (string) — Text beside the checkbox.
- `name` (string)
- `checked` (boolean, default false)
- `disabled` (boolean, default false)

```json
["checkbox", { "label": "Send confirmation email", "name": "confirm-email", "checked": true }]
```

### radio

Radio button with label. Group by using the same name across multiple radio atoms.

**Props:**

- `label` (string)
- `name` (string) — Group name — must match across all radios in a group.
- `value` (string) — Value submitted when selected.
- `checked` (boolean, default false)

```json
["radio", { "label": "Residential", "name": "job-type", "value": "residential", "checked": true }]
["radio", { "label": "Commercial", "name": "job-type", "value": "commercial" }]
["radio", { "label": "Industrial", "name": "job-type", "value": "industrial" }]
```

### switch

Toggle switch with label. Animated thumb slides on check.

**Props:**

- `label` (string)
- `name` (string)
- `checked` (boolean, default false)

```json
["switch", { "label": "Enable notifications", "name": "notifications", "checked": true }]
```

### alert

Inline feedback message with icon. Four semantic variants. Optional dismiss button.

**Props:**

- `variant` (string, default "info") — `info | success | warning | error`
- `title` (string) — Bold heading line inside the alert.
- `dismissible` (boolean, default false)

**Children:** Alert message body.

```json
["alert", { "variant": "success", "title": "Job Created" }, "The job has been added to the schedule."]
["alert", { "variant": "error", "title": "Save Failed", "dismissible": true }, "Check your connection and try again."]
["alert", { "variant": "warning" }, "This client has an overdue invoice."]
```

### badge

Small inline status pill. Six color variants.

**Props:**

- `variant` (string, default "default") — `default | primary | success | warning | error | info`

**Children:** Badge label text.

```json
["badge", { "variant": "success" }, "Active"]
["badge", { "variant": "warning" }, "Pending"]
["badge", { "variant": "error" }, "Overdue"]
```

### chip

Tag-style label. Often used for filters or selected items. Optional dismiss button removes the chip.

**Props:**

- `label` (string) — Chip text.
- `dismissible` (boolean, default false)
- `color` (string)

```json
["chip", { "label": "Plumbing" }]
["chip", { "label": "HVAC", "dismissible": true }]
```

### spinner

Animated circular loading indicator. Four sizes.

**Props:**

- `size` (string, default "md") — `xs | sm | md | lg`

```json
["spinner", { "size": "md" }]
```

### skeleton

Shimmer placeholder for loading states. Mimics text lines and optional avatar circle.

**Props:**

- `lines` (number, default 1) — Number of text line placeholders. Last line is shorter.
- `avatar` (boolean, default false) — Show a circular avatar placeholder.

```json
["skeleton", { "lines": 1 }]
["skeleton", { "lines": 3, "avatar": true }]
```

### progress

Horizontal progress bar with optional label and percentage readout.

**Props:**

- `value` (number) — Current progress value.
- `max` (number, default 100) — Maximum value.
- `label` (string) — Text on the left. Percentage auto-calculated on the right.
- `color` (string) — CSS color for the fill bar. Uses accent color by default.

```json
["progress", { "value": 65, "label": "Monthly Target" }]
["progress", { "value": 3, "max": 10, "label": "Jobs Remaining", "color": "var(--color-warning)" }]
```

### tooltip

Hover tooltip. Wraps a child element. Appears on hover in the specified position.

**Props:**

- `text` (string, required) — Tooltip text content.
- `position` (string, default "top") — `top | bottom | left | right`

**Children:** The element the tooltip is attached to.

```json
["tooltip", { "text": "Delete this record", "position": "top" }, ["button", { "variant": "ghost" }, "Delete"]]
["tooltip", { "text": "View details", "position": "right" }, ["icon", { "name": "search", "size": 18 }]]
```

### tabs

Tabbed interface. Children are tab atoms. Clicking a tab shows its panel, hides the rest.

**Props:**

- `data-id` (string)

**Children:** `tab` atoms only.

**Rule:** Set `active:true` on exactly one tab.

```json
["tabs",
  ["tab", { "label": "Overview", "panel": "ov", "active": true }, ["text", "Overview content"]],
  ["tab", { "label": "Details", "panel": "dt" }, ["text", "Details content"]],
  ["tab", { "label": "Notes", "panel": "nt" }, ["textarea", { "label": "Add note", "name": "note" }]]
]
```

### tab

Single tab inside tabs. Renders the trigger button and its panel content together.

**Props:**

- `label` (string, required) — Text shown on the tab button.
- `panel` (string, required) — Unique id for this panel within the tabs group.
- `active` (boolean, default false)

**Children:** Panel content — any atoms.

```json
["tab", { "label": "History", "panel": "hist" }, ["table", { "columns": ["Date", "Action"], "rows": [] }]]
```

### accordion

Vertically stacked expandable sections. Click header to expand/collapse with smooth animation.

**Props:**

- `data-id` (string)

**Children:** `accordion-item` atoms only.

```json
["accordion",
  ["accordion-item", { "title": "What is included?", "open": true }, "All parts and labour for the agreed scope."],
  ["accordion-item", { "title": "How long does it take?" }, "Typically 2–4 hours depending on complexity."]
]
```

### accordion-item

A single collapsible section inside accordion.

**Props:**

- `title` (string, required) — Header text shown at all times.
- `open` (boolean, default false) — Start expanded.

**Children:** Content revealed on expand. Any atoms.

```json
["accordion-item", { "title": "Warranty Policy", "open": true }, ["text", "All parts carry a 12-month warranty."]]
```

### dialog

MUI-style dialog overlay. Composed from dialog-title, dialog-content, and dialog-actions sub-components.

**Props:**

- `maxWidth` (string, default "sm") — `xs | sm | md | lg | xl | false`
- `fullWidth` (boolean, default false)
- `data-id` (string)

**Children:** `dialog-title`, `dialog-content`, `dialog-actions` atoms.

**Placement:** MUST be a direct child of body. Not inside box, grid, or card.

```json
["dialog", { "maxWidth": "sm", "data-id": "dlg--confirm" },
  ["dialog-title", {}, "Confirm Delete"],
  ["dialog-content", {}, ["text", "This action cannot be undone. Are you sure?"]],
  ["dialog-actions", {},
    ["button", { "variant": "ghost" }, "Cancel"],
    ["button", { "variant": "solid", "color": "danger" }, "Delete"]
  ]
]
```

### dialog-title

Title bar for a dialog.

**Children:** Title text.

### dialog-content

Body content area for a dialog.

**Children:** Any atoms.

### dialog-actions

Action button row for a dialog. Buttons align right.

**Children:** `button` atoms.

### drawer

Slide-in side panel. Hidden until triggered. Backdrop darkens. MUST be placed at top level of body.

**Props:**

- `id` (string, required) — Unique id. Referenced in trigger: `on:click drawer:{id}`
- `title` (string) — Panel header with auto-generated close button.
- `side` (string, default "right") — `left | right`

**Children:** Drawer body content.

**Placement:** MUST be a direct child of body. Not inside box, grid, or card.

```json
["drawer", { "id": "job-detail", "title": "Job Detail", "side": "right" },
  ["heading", { "level": 3 }, "JOB-042"],
  ["text", "AC Installation — 42 Maple Street"],
  ["divider"],
  ["button", { "variant": "solid", "color": "success" }, "Mark Complete"]
]
```

### breadcrumb

Navigation trail. Items separated by /. Last item is current page — no href, bold text.

**Children:** `breadcrumb-item` atoms.

```json
["breadcrumb",
  ["breadcrumb-item", { "href": "#" }, "Home"],
  ["breadcrumb-item", { "href": "#" }, "Clients"],
  ["breadcrumb-item", "John Smith"]
]
```

### breadcrumb-item

One step in a breadcrumb. Omit href for the current (last) item.

**Props:**

- `href` (string) — Link URL. Omit for the current page item.

**Children:** Item label text.

```json
["breadcrumb-item", { "href": "/jobs" }, "Jobs"]
```

### pagination

Page number controls with prev/next. Shows ellipsis for large page counts.

**Props:**

- `total` (number) — Total number of items.
- `page` (number, default 1) — Currently active page.
- `per-page` (number, default 10) — Items per page.
- `on:change` (string) — Action name fired with page number appended.

```json
["pagination", { "total": 150, "page": 5, "per-page": 10 }]
```

### table

Data table. Columns from array, rows from array of arrays. Striped and hoverable by default.

**Props:**

- `columns` (array) — Array of column header strings.
- `rows` (array) — Array of row arrays.
- `striped` (boolean, default true)
- `hoverable` (boolean, default true)

```json
["table", {
  "columns": ["Job ID", "Client", "Status", "Value"],
  "rows": [
    ["JOB-001", "HVAC Co", "Completed", "$850"],
    ["JOB-002", "Pipe Masters", "In Progress", "$220"]
  ]
}]
```

### list

Styled list. Unordered by default. Use divided for border-separated items.

**Props:**

- `ordered` (boolean, default false)
- `divided` (boolean, default false) — Add border between items.

**Children:** `list-item` atoms.

```json
["list",
  ["list-item", "First item"],
  ["list-item", "Second item"]
]
["list", { "divided": true },
  ["list-item", { "href": "#" }, "Clickable link item"],
  ["list-item", "Plain item"]
]
```

### list-item

A single item in a list.

**Props:**

- `href` (string) — Makes the item a navigation link.
- `on:click` (string) — Action string.
- `data-id` (string)

**Children:** Item content — text or atoms.

```json
["list-item", "Plain text item"]
["list-item", { "href": "/jobs/42" }, "JOB-042 — HVAC Installation"]
["list-item", { "on:click": "modal:job-detail" }, "Click to view"]
```

### avatar

Image or initials circle. Shows initials, an image, or a default person icon.

**Props:**

- `src` (string) — Image URL.
- `alt` (string) — Alt text.
- `initials` (string) — Up to 2 characters.
- `size` (string, default "md") — `xs | sm | md | lg`
- `color` (string, default "default") — `default | success | danger | info | warning`
- `data-id` (string)

```json
["avatar", { "initials": "JD", "size": "md", "color": "info", "data-id": "avatar--jd" }]
["avatar", { "src": "/photos/mike.jpg", "alt": "Mike", "size": "md", "data-id": "avatar--mike" }]
```

### avatar-group

Stacked overlapping avatars. Place avatar atoms as children.

**Props:**

- `size` (string, default "md")
- `data-id` (string)

**Children:** `avatar` atoms.

```json
["avatar-group", { "data-id": "avatar-group--team" },
  ["avatar", { "initials": "JD", "size": "md", "color": "info" }],
  ["avatar", { "initials": "AR", "size": "md", "color": "success" }],
  ["avatar", { "initials": "+3", "size": "md" }]
]
```

### empty-state

Centered placeholder when content is absent. Icon + title + description + optional CTA button.

**Props:**

- `icon` (string, default "inbox") — Icon name.
- `title` (string, default "Nothing here yet")
- `description` (string) — Supporting text.
- `action` (string) — Label for optional CTA button.
- `on:click` (string) — Action string for CTA.
- `data-id` (string)

```json
["empty-state", { "icon": "inbox", "title": "No jobs found", "description": "Adjust your filters or create a new job.", "action": "New Job", "on:click": "modal:new-job", "data-id": "empty-state--jobs" }]
```

### kbd

Keyboard shortcut display. Use + to separate keys.

**Props:**

- `keys` (string) — Key combination string.
- `data-id` (string)

```json
["kbd", { "keys": "Ctrl+S", "data-id": "kbd--save" }]
["kbd", { "keys": "Ctrl+Shift+P", "data-id": "kbd--palette" }]
```

### code

Inline monospace code. Renders inline within text.

**Children:** Code string.

```json
["text", "Call ", ["code", "csAction('modal:x')"], " to open a modal."]
```

### code-block

Multiline code panel with language label and copy button.

**Props:**

- `lang` (string) — Language label.
- `content` (string) — Code content. Use `\n` for newlines.
- `data-id` (string)

```json
["code-block", { "lang": "json", "content": "{\"key\": \"value\"}", "data-id": "code-block--example" }]
```

### timeline

Vertical event list with colored dots and connectors.

**Props:**

- `data-id` (string)

**Children:** `timeline-item` atoms.

```json
["timeline", { "data-id": "timeline--activity" },
  ["timeline-item", { "time": "2 min ago", "title": "Job assigned", "color": "info", "data-id": "tli--1" }],
  ["timeline-item", { "time": "1 hr ago", "title": "Invoice sent", "color": "success", "data-id": "tli--2" }]
]
```

### timeline-item

A single event in a timeline. Must be a direct child of timeline.

**Props:**

- `time` (string) — Relative or absolute time label.
- `title` (string) — Event title.
- `description` (string) — Supporting detail text.
- `color` (string, default "default") — `default | success | danger | info | warning`
- `data-id` (string)

```json
["timeline-item", { "time": "Now", "title": "Job #482 created", "color": "info", "data-id": "tli--1" }]
```

### rating

Star rating. Interactive by default, read-only optional.

**Props:**

- `value` (number, default 0) — Current rating.
- `max` (number, default 5) — Number of stars.
- `readonly` (boolean, default false)
- `name` (string, default "rating")
- `data-id` (string)

```json
["rating", { "value": 4, "max": 5, "name": "score", "data-id": "rating--score" }]
["rating", { "value": 3, "max": 5, "readonly": true, "data-id": "rating--display" }]
```

### slider

Range slider with optional label and live value display.

**Props:**

- `label` (string)
- `name` (string, default "slider")
- `min` (number, default 0)
- `max` (number, default 100)
- `value` (number, default 50)
- `step` (number, default 1)
- `data-id` (string)

```json
["slider", { "label": "Volume", "name": "volume", "min": 0, "max": 100, "value": 65, "data-id": "slider--volume" }]
```

### number-input

Numeric stepper with +/- buttons. Respects min, max, and step.

**Props:**

- `label` (string)
- `name` (string, default "number")
- `value` (number, default 0)
- `min` (string)
- `max` (string)
- `step` (number, default 1)
- `data-id` (string)

```json
["number-input", { "label": "Quantity", "name": "qty", "min": "1", "max": "99", "value": 3, "data-id": "number-input--qty" }]
```

### file-upload

Drag-and-drop file upload zone.

**Props:**

- `label` (string) — Text in the drop zone.
- `hint` (string) — Small hint text.
- `name` (string, default "file")
- `accept` (string, default "*")
- `multiple` (boolean, default false)
- `data-id` (string)

```json
["file-upload", { "label": "Drop reports here or click to upload", "hint": "PDF, JPG up to 20MB", "accept": ".pdf,.jpg", "multiple": true, "name": "reports", "data-id": "file-upload--reports" }]
```

### tag-input

Multi-value chip input. Press Enter or comma to add a tag. Click × to remove.

**Props:**

- `label` (string)
- `placeholder` (string, default "Add tag...")
- `name` (string, default "tags")
- `tags` (string) — Comma-separated initial tags.
- `data-id` (string)

```json
["tag-input", { "label": "Tags", "name": "tags", "tags": "HVAC,Urgent", "placeholder": "Add tag...", "data-id": "tag-input--tags" }]
```

### date-input

Styled native date picker.

**Props:**

- `label` (string)
- `name` (string, default "date")
- `value` (string) — YYYY-MM-DD format.
- `min` (string)
- `max` (string)
- `data-id` (string)

```json
["date-input", { "label": "Scheduled Date", "name": "scheduled", "data-id": "date-input--scheduled" }]
```

### menu

Dropdown context menu. Button trigger opens a floating list of actions.

**Props:**

- `id` (string, required)
- `label` (string, default "Actions")
- `variant` (string, default "outline")
- `size` (string, default "md")
- `data-id` (string)

**Children:** `menu-item` atoms.

```json
["menu", { "id": "job-actions", "label": "Actions", "variant": "outline", "data-id": "menu--job" },
  ["menu-item", { "label": "Edit", "icon": "pencil", "on:click": "modal:edit-job", "data-id": "menu-item--edit" }],
  ["menu-item", { "label": "Delete", "icon": "trash", "color": "danger", "on:click": "snackbar:Deleted:error", "data-id": "menu-item--delete" }]
]
```

### menu-item

A single item inside a menu dropdown.

**Props:**

- `label` (string)
- `icon` (string) — Icon name.
- `color` (string) — `danger | success`
- `on:click` (string)
- `disabled` (boolean, default false)
- `data-id` (string)

```json
["menu-item", { "label": "Edit", "icon": "pencil", "on:click": "modal:edit", "data-id": "menu-item--edit" }]
```

### popover

Floating panel anchored to a trigger button.

**Props:**

- `id` (string, required)
- `label` (string, default "More info")
- `variant` (string, default "outline")
- `size` (string, default "md")
- `placement` (string, default "bottom") — `bottom | top | left | right`
- `data-id` (string)

**Children:** Any atoms rendered inside the panel.

```json
["popover", { "id": "job-info", "label": "Details", "placement": "bottom", "data-id": "popover--job" },
  ["heading", { "level": 3 }, "Job #482"],
  ["text", "Assigned: Mike K."],
  ["text", "Status: In Progress"]
]
```

### stepper

Multi-step wizard indicator.

**Props:**

- `current` (number, default 1) — The active step number.
- `orientation` (string, default "horizontal") — `horizontal | vertical`
- `data-id` (string)

**Children:** `stepper-step` atoms.

```json
["stepper", { "current": 2, "orientation": "horizontal", "data-id": "stepper--job" },
  ["stepper-step", { "step": 1, "title": "Create Job", "data-id": "step--1" }],
  ["stepper-step", { "step": 2, "title": "Assign Tech", "data-id": "step--2" }],
  ["stepper-step", { "step": 3, "title": "Schedule", "data-id": "step--3" }]
]
```

### stepper-step

A single step in a stepper. Must be a direct child of stepper.

**Props:**

- `step` (number, required) — Step number (1-based).
- `title` (string) — Step label.
- `description` (string) — Short supporting text.
- `data-id` (string)

```json
["stepper-step", { "step": 1, "title": "Details", "description": "Enter job info", "data-id": "step--1" }]
```

### toolbar

Horizontal action bar. Title on the left, children on the right.

**Props:**

- `title` (string) — Left-side title label.
- `bordered` (boolean, default false)
- `data-id` (string)

**Children:** Any atoms rendered on the right side.

```json
["toolbar", { "title": "Jobs", "bordered": true, "data-id": "toolbar--jobs" },
  ["button", { "variant": "ghost", "size": "sm", "label": "Filter", "data-id": "toolbar--filter" }],
  ["button", { "variant": "solid", "size": "sm", "label": "New Job", "data-id": "toolbar--new" }]
]
```

### chat-widget

Floating chat bubble + glass-theme chat panel. Connects to an n8n webhook.

**Props:**

- `webhook` (string, required) — Full n8n webhook URL.
- `route` (string, default "general") — Logical route name.
- `title` (string, default "Chat") — Header label.
- `data-id` (string, required)

**Children:** none

```json
["chat-widget", { "webhook": "https://your-n8n/webhook/abc/chat", "route": "support", "title": "Support", "data-id": "chat--support" }]
```

### data-chat

Natural language query widget for any MongoDB collection. Embeds a chat panel where the user types plain English questions. A local LLM (Ollama/Mistral) translates the question into a structured JSON query intent. The intent executes deterministically against MongoDB — the LLM never touches the data.

**Props:**

- `schema` (string, required) — Path to the query schema JSON file.
- `placeholder` (string, default "Ask a question about the data...")
- `id` (string, default "data-chat")

**Children:** none

**How it works:**

1. User types a question in plain English
2. POST `/api/data/ask` sends `{ question, schema }` to the Go server
3. Go calls Python data_query module via the Python bridge
4. Python loads the query schema and builds an LLM prompt from it
5. Local Ollama (Mistral, temperature: 0) translates the question into a JSON query intent
6. Python executes the intent deterministically against MongoDB
7. Computed result returned to the chat panel — no hallucination possible

**Key principle:** The LLM is a translator, not an oracle. It converts natural language to structured intent. The database provides the truth.

```json
["data-chat", { "schema": "schemas/tickets.json", "placeholder": "Ask about tickets...", "id": "ticket-chat" }]
```

### form

Form wrapper. Renders a `<form>` with cs-form class.

**Props:**

- `id` (string) — Form ID — referenced by button `on:click` post actions.
- `data-autosave` (string) — localStorage key.
- `data-id` (string)

**Children:** Any input, select, textarea, button atoms.

```json
["form", { "id": "login-form", "data-autosave": "login-draft", "data-id": "form--login" },
  ["input", { "name": "email", "label": "Email", "type": "email", "data-id": "form--login--email" }],
  ["input", { "name": "password", "label": "Password", "type": "password", "data-id": "form--login--password" }],
  ["button", { "label": "Sign in", "variant": "solid", "on:click": "post:auth/login:login-form", "data-id": "form--login--submit" }]
]
```

### sidebar

App-shell left navigation sidebar. Fixed width (240px), full height.

**Props:**

- `brand` (string) — App name.
- `data-id` (string)

**Children:** `nav-link` atoms.

```json
["sidebar", { "brand": "ChefScript", "data-id": "sidebar--main" },
  ["nav-link", { "href": "/page/dashboard", "data-id": "nav--dashboard" }, "Dashboard"],
  ["nav-link", { "href": "/page/quizzes", "data-id": "nav--quizzes" }, "Quizzes"],
  ["nav-link", { "href": "/page/settings", "data-id": "nav--settings" }, "Settings"]
]
```

### section

Titled content block. Optional heading + description above a body region.

**Props:**

- `title` (string) — Section heading (h2).
- `description` (string) — Muted subtitle.
- `data-id` (string)

**Children:** Any atoms.

```json
["section", { "title": "Recent Activity", "description": "Your last 7 days", "data-id": "section--activity" },
  ["list", { "data-id": "list--activity" }]
]
```

### callout

Inline editorial block — tips, notes, warnings. Left-bordered with variant colour.

**Props:**

- `variant` (string, default "info") — `info | warning | tip | danger`
- `title` (string) — Optional bold title.
- `icon` (string) — Override the default variant icon.
- `data-id` (string)

**Children:** Text or atoms.

```json
["callout", { "variant": "warning", "title": "Heads up", "data-id": "callout--warn" }, "This action will affect all users."]
["callout", { "variant": "tip", "data-id": "callout--tip" }, "Use keyboard shortcut Ctrl+S to save."]
```

### image

Image element with optional size and border-radius.

**Props:**

- `src` (string, required) — Image URL.
- `alt` (string, default "")
- `width` (string) — Width in px or percent.
- `height` (string)
- `rounded` (boolean, default false)
- `data-id` (string)

```json
["image", { "src": "/public/logo.png", "alt": "Logo", "width": "120", "data-id": "image--logo" }]
```

### link

Inline anchor link.

**Props:**

- `href` (string, default "#")
- `label` (string)
- `target` (string) — `_blank` to open in new tab.
- `variant` (string, default "default") — `default | muted | danger`
- `data-id` (string)

**Children:** Text if label is not provided.

```json
["link", { "href": "/page/dashboard", "label": "Go to dashboard", "data-id": "link--dashboard" }]
```

### search

Search input with icon and clear button.

**Props:**

- `placeholder` (string, default "Search...")
- `name` (string, default "q")
- `value` (string)
- `on:search` (string) — Action fired on Enter key.
- `data-id` (string)

```json
["search", { "placeholder": "Search users...", "on:search": "users/search", "data-id": "search--users" }]
```

### color-input

Native colour picker with hex value label.

**Props:**

- `label` (string)
- `name` (string, default "color")
- `value` (string, default "#000000")
- `data-id` (string)

```json
["color-input", { "label": "Accent Colour", "name": "accent", "value": "#00b4d8", "data-id": "color--accent" }]
```

### snackbar

Toast notification container. Place once per page. Driven by `csSnackbar(message, variant)` or automatically by `ActionResult.Toast`.

**Props:**

- `position` (string, default "bottom-right") — `bottom-right | bottom-left | top-right | top-left | bottom-center`
- `data-id` (string)

**Children:** none

```json
["snackbar", { "position": "bottom-right", "data-id": "snackbar--main" }]
```

### confirm

Confirmation dialog. Pre-wired yes/no dialog.

**Props:**

- `id` (string, required)
- `title` (string, default "Are you sure?")
- `message` (string)
- `confirm-label` (string, default "Confirm")
- `cancel-label` (string, default "Cancel")
- `on:confirm` (string) — Action fired on confirm.
- `variant` (string, default "danger") — `danger | warning | default`
- `data-id` (string)

```json
["confirm", { "id": "del-user", "title": "Delete user?", "message": "This cannot be undone.",
  "on:confirm": "users/delete", "variant": "danger", "data-id": "confirm--del-user" }]
["button", { "label": "Delete", "variant": "outline", "on:click": "modal:del-user", "data-id": "btn--delete" }]
```

### carousel

Horizontal scroll carousel. CSS scroll-snap based.

**Props:**

- `id` (string, required)
- `data-id` (string)

**Children:** Slide elements with class `cs-carousel__slide`.

```json
["carousel", { "id": "hero", "data-id": "carousel--hero" },
  ["div", { "class": "cs-carousel__slide" }, ["card", {}, ["heading", { "level": 2 }, "Slide 1"]]],
  ["div", { "class": "cs-carousel__slide" }, ["card", {}, ["heading", { "level": 2 }, "Slide 2"]]],
  ["div", { "class": "cs-carousel__slide" }, ["card", {}, ["heading", { "level": 2 }, "Slide 3"]]]
]
```

### rich-text

Renders pre-built HTML content (CMS output, rendered markdown).

**Props:**

- `content` (string) — Raw HTML string to inject.
- `data-id` (string)

**Children:** Fallback if content prop is not provided.

```json
["rich-text", { "content": "{{data.articleHtml}}", "data-id": "rich-text--article" }]
```

### chart

SVG data chart. Native rendering — no external libraries. Supports bar, line, and pie types.

**Props:**

- `type` (string, default "bar") — `bar | line | pie`
- `data` (array, required) — Array of `{label, value}` objects.
- `height` (number, default 180)
- `title` (string) — Optional label above the chart.
- `data-id` (string)

```json
["chart", { "type": "bar", "title": "Monthly Revenue", "height": 200,
  "data": [{"label":"Jan","value":4200},{"label":"Feb","value":6700},{"label":"Mar","value":5100},{"label":"Apr","value":8300}],
  "data-id": "chart--revenue" }]
["chart", { "type": "pie",
  "data": [{"label":"Go","value":60},{"label":"Python","value":25},{"label":"JS","value":15}],
  "data-id": "chart--langs" }]
```

### banner

Full-width site-level message strip.

**Props:**

- `variant` (string, default "info") — `info | success | warning | danger`
- `title` (string) — Bold text prefix.
- `dismissible` (boolean, default true)
- `data-id` (string)

**Children:** Banner message text.

```json
["banner", { "variant": "warning", "title": "Scheduled maintenance:", "data-id": "banner--maint" }, "System will be offline Saturday 2am–4am UTC."]
```

### form-field

Consistent label + hint + error wrapper for ANY child atom.

**Props:**

- `label` (string)
- `hint` (string)
- `error` (string) — Overrides hint and colours label red.
- `required` (boolean, default false)
- `data-id` (string)

**Children:** Any single input atom.

```json
["form-field", { "label": "Skill tags", "hint": "Select all that apply", "required": true, "data-id": "field--tags" },
  ["multi-select", { "name": "tags", "options": ["Go","Python","Rust"], "data-id": "tags--skills" }]
]
```

### kv-list

Key-value detail list.

**Props:**

- `divided` (boolean, default true)
- `data-id` (string)

**Children:** `kv-item` atoms.

```json
["kv-list", { "data-id": "kv--user-detail" },
  ["kv-item", { "key": "Username", "value": "kibbyd", "data-id": "kv--username" }],
  ["kv-item", { "key": "Role", "value": "Admin", "value-variant": "success", "data-id": "kv--role" }],
  ["kv-item", { "key": "Joined", "value": "Feb 2026", "data-id": "kv--joined" }]
]
```

### kv-item

Single key-value row inside a kv-list.

**Props:**

- `key` (string) — Label on the left.
- `value` (string) — Value on the right.
- `value-variant` (string) — `success | warning | danger | muted`
- `href` (string) — Makes the value a link.
- `data-id` (string)

```json
["kv-item", { "key": "Status", "value": "Active", "value-variant": "success", "data-id": "kv--status" }]
```

### button-group

Connects multiple buttons into a single visual unit.

**Props:**

- `data-id` (string)

**Children:** `button` atoms.

```json
["button-group", { "data-id": "btn-group--view" },
  ["button", { "variant": "outline", "size": "sm", "label": "Grid", "data-id": "btn--grid" }],
  ["button", { "variant": "outline", "size": "sm", "label": "List", "data-id": "btn--list" }],
  ["button", { "variant": "outline", "size": "sm", "label": "Kanban", "data-id": "btn--kanban" }]
]
```

### copy-button

Button that copies a value to clipboard. Shows ✓ Copied for 2 seconds.

**Props:**

- `value` (string, required) — Text to copy.
- `label` (string, default "Copy")
- `variant` (string, default "outline")
- `size` (string, default "md")
- `data-id` (string)

```json
["copy-button", { "value": "{{data.apiKey}}", "label": "Copy API key", "size": "sm", "data-id": "copy--api-key" }]
```

### multi-select

Multi-value select with chip output.

**Props:**

- `label` (string)
- `name` (string, required)
- `options` (array, required)
- `placeholder` (string, default "Select...")
- `data-id` (string)

```json
["multi-select", { "label": "Technologies", "name": "tech", "placeholder": "Pick stack...",
  "options": ["Go","Python","JavaScript","Rust","TypeScript"], "data-id": "ms--tech" }]
```

### icon-button

Square icon-only button.

**Props:**

- `icon` (string, required)
- `size` (string, default "md") — `sm | md | lg`
- `variant` (string, default "ghost") — `ghost | outline | solid`
- `on:click` (string)
- `data-id` (string)

```json
["icon-button", { "icon": "trash", "size": "sm", "variant": "ghost", "on:click": "items/delete", "data-id": "btn--delete" }]
```

### tag

Small inline label for categorisation.

**Props:**

- `label` (string, required)
- `color` (string, default "default") — `default | primary | success | warning | danger`
- `removable` (boolean, default false)
- `on:remove` (string) — Action string fired when × is clicked.
- `data-id` (string)

```json
["tag", { "label": "Published", "color": "success", "data-id": "tag--status" }]
["tag", { "label": "React", "removable": true, "on:remove": "tags/remove:react", "data-id": "tag--react" }]
```

### data-grid

Sortable, filterable data table.

**Props:**

- `columns` (array, required) — Array of `{ key, label, sortable?, filterable? }`.
- `rows` (array, required) — Array of row objects.
- `data-id` (string)

```json
["data-grid", {
  "data-id": "grid--users",
  "columns": [
    { "key": "name", "label": "Name", "sortable": true, "filterable": true },
    { "key": "role", "label": "Role", "sortable": true },
    { "key": "email", "label": "Email", "filterable": true }
  ],
  "rows": [
    { "name": "Alice", "role": "Admin", "email": "alice@example.com" },
    { "name": "Bob", "role": "Editor", "email": "bob@example.com" }
  ]
}]
```

### tree

Hierarchical tree view with expand/collapse.

**Props:**

- `data-id` (string)

**Children:** `tree-item` atoms.

```json
["tree", { "data-id": "tree--files" },
  ["tree-item", { "label": "src", "icon": "folder", "data-id": "tree--src" },
    ["tree-item", { "label": "main.go", "icon": "file", "data-id": "tree--main" }],
    ["tree-item", { "label": "engine", "icon": "folder", "data-id": "tree--engine" }]
  ],
  ["tree-item", { "label": "go.mod", "icon": "file", "data-id": "tree--gomod" }]
]
```

### virtual-list

Windowed list renderer for large datasets.

**Props:**

- `columns` (array, required) — Column definition objects.
- `rows` (array, required) — Row data objects.
- `row-height` (number, default 40)
- `data-id` (string)

```json
["virtual-list", {
  "data-id": "vl--logs",
  "row-height": 36,
  "columns": [{ "key": "ts", "label": "Time" }, { "key": "msg", "label": "Message" }],
  "rows": "{{data.logRows}}"
}]
```

### notification

Bell icon button that opens a dropdown notification panel.

**Props:**

- `id` (string, required)
- `count` (number, default 0) — Unread badge count.
- `data-id` (string)

**Children:** `notification-item` atoms.

```json
["notification", { "id": "main", "count": 3, "data-id": "notif--main" },
  ["notification-item", { "title": "Deploy complete", "body": "v1.4.2 is live.", "time": "2m ago", "unread": true, "data-id": "ni--deploy" }],
  ["notification-item", { "title": "PR merged", "body": "feat: dark mode", "time": "1h ago", "data-id": "ni--pr" }]
]
```

### command

Command palette modal. Opens with Cmd+K / Ctrl+K.

**Props:**

- `id` (string, required)
- `placeholder` (string, default "Type a command or search...")
- `items` (array) — Array of `{ label, action, description?, icon? }`.
- `data-id` (string)

**Children:** `command-item` atoms (optional).

```json
["command", { "id": "main", "placeholder": "Search commands...", "data-id": "cmd--main" },
  ["command-item", { "label": "Dashboard", "on:click": "nav:dashboard", "description": "Go home", "icon": "home", "data-id": "ci--dashboard" }],
  ["command-item", { "label": "Settings", "on:click": "nav:settings", "icon": "settings", "data-id": "ci--settings" }]
]
```

### context-menu

Right-click context menu.

**Props:**

- `id` (string, required) — Bind targets use `data-context-menu="this-id"`.
- `data-id` (string)

**Children:** `menu-item` atoms.

```json
["context-menu", { "id": "file-ctx", "data-id": "ctx--file" },
  ["menu-item", { "label": "Open", "icon": "folder", "on:click": "file/open", "data-id": "mi--open" }],
  ["menu-item", { "label": "Delete", "icon": "trash", "on:click": "file/delete", "color": "danger", "data-id": "mi--delete" }]
]
```

### hover-card

Pure-CSS hover tooltip panel.

**Props:**

- `placement` (string, default "top") — `top | bottom | left | right`
- `data-id` (string)

**Children:** First child = trigger. Remaining children = panel content.

```json
["hover-card", { "placement": "top", "data-id": "hc--user" },
  ["badge", { "label": "Alice" }],
  ["card", {},
    ["text", {}, "Admin — joined 2024"]
  ]
]
```

### split-view

Resizable split panel.

**Props:**

- `id` (string, required)
- `direction` (string, default "horizontal") — `horizontal | vertical`
- `default-size` (number, default 50) — Initial size % for first pane.
- `data-id` (string)

**Children:** Exactly two `split-pane` atoms.

```json
["split-view", { "id": "editor", "direction": "horizontal", "default-size": 30, "data-id": "sv--editor" },
  ["split-pane", { "data-id": "sp--sidebar" }, ["text", {}, "File tree"]],
  ["split-pane", { "data-id": "sp--main" }, ["text", {}, "Editor area"]]
]
```

### calendar

Month-view date picker.

**Props:**

- `id` (string, required)
- `name` (string) — Hidden input name.
- `value` (string) — Initial selected date (YYYY-MM-DD).
- `on:change` (string) — Action fired with selected date.
- `data-id` (string)

```json
["calendar", { "id": "booking", "name": "date", "value": "2026-03-15", "on:change": "booking/setDate", "data-id": "cal--booking" }]
```

### video

HTML5 video player.

**Props:**

- `src` (string, required) — Video file URL.
- `poster` (string) — Preview image URL.
- `controls` (boolean, default true)
- `autoplay` (boolean, default false)
- `loop` (boolean, default false)
- `muted` (boolean, default false)
- `width` (string) — CSS width.
- `data-id` (string)

```json
["video", { "src": "/media/demo.mp4", "poster": "/media/thumb.jpg", "controls": true, "width": "100%", "data-id": "vid--demo" }]
```

### audio

HTML5 audio player with native controls.

**Props:**

- `src` (string, required)
- `controls` (boolean, default true)
- `autoplay` (boolean, default false)
- `loop` (boolean, default false)
- `data-id` (string)

```json
["audio", { "src": "/media/track.mp3", "data-id": "aud--track" }]
```

### iframe

Embedded iframe.

**Props:**

- `src` (string, required)
- `title` (string, required)
- `width` (string, default "100%")
- `height` (string, default "400px")
- `sandbox` (string) — sandbox attribute value.
- `data-id` (string)

```json
["iframe", { "src": "https://www.openstreetmap.org/export/embed.html?bbox=-0.1,51.4,-0.0,51.5", "title": "Map", "height": "300px", "data-id": "frame--map" }]
```

### aspect-ratio

Wrapper that enforces a fixed aspect ratio.

**Props:**

- `ratio` (string, default "16:9") — Ratio in W:H format.
- `data-id` (string)

**Children:** Any content — typically image, video, or iframe.

```json
["aspect-ratio", { "ratio": "16:9", "data-id": "ar--hero" },
  ["iframe", { "src": "https://www.youtube.com/embed/dQw4w9WgXcQ", "title": "Video", "data-id": "frame--yt" }]
]
```

---

## Section Index

Where to find everything else. Search by key name.

- **platform:** Server, MongoDB, sessions, security, page system, action system, form handling, autosave, Excel export, chat widget, Python bridge
- **data_flow:** Binary schema pipe deep-dive — encode/decode flow, how_it_works steps, speed characteristics
- **data_resolution:** Component data dependency resolution — binary schema vs static JSON vs Go map
- **atom_format:** Atom syntax rules and examples
- **page_format:** Page JSON structure — title, theme, body
- **agent_rules:** Structural rules for dialogs, drawers, tabs, accordion, breadcrumbs, autocomplete
- **events:** on:click and on:enter action strings, escape handling, full examples
- **quick_reference:** All 100+ atom components grouped by category
- **icons:** Icon naming, variants, usage
- **build_recipes:** 11 step-by-step patterns for solving common build problems
- **patterns:** Copy-paste page compositions — dashboard, form, drawer detail, confirm dialog, tabbed detail, search list
- **components:** Full reference for every component — props, children, examples

### Also See

- `engine/runtime.go` — Full JS runtime source (~870 lines)
- `engine/flight_recorder.go` — Flight recorder — ring buffer, API endpoints
- `engine/diagnostics.go` — Diagnostic collector — validation, template tracing, feeds flight recorder
- `engine/diagnostics_panel.go` — Panel UI — filters, flight tab, export, CSS
