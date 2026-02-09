**Purpose**
- **Goal**: Help a coding AI become immediately productive in this Go REST API repository.
- **Scope**: Concrete, repository-specific patterns, conventions, and integration points — not generic advice.

**Big Picture**
- **Architecture**: Simple monolith using standard library `net/http` for routing, manual dependency injection in `main.go` and layered repo/service/handler separation. See [main.go](main.go) and the directories `handlers/`, `services/`, `repositories/`, `models/`.
- **Data flow**: HTTP -> handler (parses request) -> service (business logic thin wrapper) -> repository (SQL, returns models) -> handler -> JSON response. Example: `handlers/product_handler.go` -> `services/product_service.go` -> `repositories/product_repository.go`.

**Key files to inspect**
- **Entry & wiring**: [main.go](main.go) — config load (`viper`), DB init (`database.InitDB`), DI, route registration.
- **Handlers**: [handlers/product_handler.go](handlers/product_handler.go), [handlers/category_handler.go](handlers/category_handler.go) — use `strings.TrimPrefix` to extract IDs from URL paths (no external router).
- **Services**: [services/product_service.go](services/product_service.go) — thin service layer delegating to repository.
- **Repositories**: [repositories/product_repository.go](repositories/product_repository.go) — hand-written SQL, joins categories for `category_name`.
- **Models**: [models/product.go](models/product.go), [models/category.go](models/category.go) — JSON tags define canonical request/response field names.
- **Database helper**: [database/database.go](database/database.go) — uses `github.com/lib/pq` and returns `*sql.DB`.
- **Migrations**: `migrations/` — SQL migration files exist, but note `main.go` also auto-creates tables and sample data at startup via `createTablesAndData()`.

**Project-specific conventions & gotchas**
- **Routing style**: The project does not use a router library — handlers inspect `r.Method` and `r.URL.Path` directly. When changing endpoints, update `main.go` registration and any `TrimPrefix` logic in handlers.
- **Field naming mismatch**: The `README.md` uses Indonesian keys like `nama`, `harga`, `stok` in examples, but model JSON tags use `name`, `price`, `stock`. Prefer model tags (`models/*.go`) as source of truth for JSON I/O.
- **Language/strings**: Error messages and some responses are in Indonesian (e.g., repository errors like "produk tidak ditemukan"). Preserve locale consistency when editing messages.
- **DB lifecycle**: `main.go` will run without a DB if `DB_CONN` is empty; when DB is present it runs `createTablesAndData()` to create tables and sample data — changes to schema should be mirrored in both `migrations/` and this function.
- **Dependency injection**: DI is manual and happens inline in `main.go`. Tests or new components should follow the same constructor patterns: `NewXRepository(db)`, `NewXService(repo)`, `NewXHandler(service)`.

**Build / run / debug**
- **Run locally**: `go run .` (also documented in [README.md](README.md)).
- **Build binary**: `go build -o kasir-api` then `./kasir-api`.
- **Config**: Environment-driven via `viper`. Use `.env` or env vars `PORT` and `DB_CONN`. Example: `DB_CONN=postgres://user:pass@host:5432/dbname?sslmode=disable PORT=8080 go run .`.
- **Health**: `GET /health` to verify server is alive.

**Integration points & deps**
- **Database**: PostgreSQL via `github.com/lib/pq` (see `go.mod`). SQL is executed with `database/sql` — repositories use parameterized queries ($1 style).
- **Deployments**: README references Railway and Zeabur deployments — CI/CD is not present in repo; inspect `zeabur.yaml` if modifying deployment behavior.

**Editing guidance for AI agents**
- **Preserve protocols**: Keep the simple handler/service/repo layering. When adding a new resource, follow the same file layout and constructor patterns.
- **SQL changes**: Update `migrations/` SQL, repository queries, and `createTablesAndData()` in `main.go` together.
- **Endpoint examples**: Use existing endpoints as templates — `product_handler.go` shows parsing, validation, JSON decode/encode patterns and status codes.
- **Tests**: No tests present. If adding tests, inject `*sql.DB` with a test database or use interfaces for repos to allow mocking.

**Examples (copy-paste guidance)**
- Extract ID in handlers: use the same pattern as products: `idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")` then `strconv.Atoi(idStr)` ([handlers/product_handler.go](handlers/product_handler.go)).
- Return JSON and set header: `w.Header().Set("Content-Type", "application/json")` then `json.NewEncoder(w).Encode(payload)` (used across handlers).

**When uncertain**
- Prefer model JSON tags over README examples when deciding request/response shapes.
- For schema changes, search for SQL strings in `main.go` and `repositories/` to avoid divergence.

If anything here is unclear or you want different emphasis (tests, CI, or API docs), tell me what to expand or correct.
