# src

```
cmd/main.go          the service            -> boot.Run()
cmd/migrate          the vvgoose CLI        (make migrate …)
cmd/seed             the seed command       -> boot.Seed()
config/example.yml   copy to config/app.yml
resources/migrations Goose SQL files
src/
  boot/              the composition root: which pieces make up a process
  config/            the configuration tree, loaded strictly
  core/              infrastructure every module shares
    auth/            the API-key guard (authentication only)
    database/        the database/sql pool and the crud.Source over it
    observability/   slog logger, OpenTelemetry SDK
  mod/               bounded contexts (see mod/README.md)
  transports/
    http/            the Fiber app, root and health routes, serving
```

## boot

The only place that knows every piece. `core.go` lists the infrastructure every
process shares (config, logging and telemetry, the database, the guard, the
runner supervisor), `app.go` the catalog of bounded contexts, `run.go` the
processes: `Run` (the server, with the profile from `app.profile`) and `Seed`
(the seed command). The modules never import each other's internals, and
nothing outside `boot` decides what runs.

The same catalog builds every kind of process. What separates the API
replica, a worker replica and the seed command is a `module.Profile`
(`complete`, `serving`, `working`, `seeding`), not three hand-kept lists
(framework decision D-106).

## config

One `Config` struct with fixed top-level sections:

| Section | For |
|---|---|
| `app` | name, env (selects seeders), profile |
| `transports` | anything that handles communication: `fiber` (addr, prefix, CORS, dev API keys) |
| `infra` | what the service runs on: `db` is a `vvdb.Config` (engine, host, pool, migration) |
| `clients` | external APIs and SDKs |
| `internal` | other internal services or modules |
| `observability` | logger and OpenTelemetry |

`config.Load` uses `vvcfg.LoadFrom` with `Strict: true`: unknown keys refuse
the start. A block validates itself with a `ValidateSelf() error` method; the
loader walks the tree and reports every broken block at once.

## core

Shared infrastructure, provided to the fx graph:

- **database** — `crudsqlfx.Module`: opens one `*sql.DB` from `infra.db`
  (driver `pgx` registered by import), pings it at start, closes it at stop,
  and provides the `crud.Source` every generated repository binds to. It also
  contributes the `database.primary` readiness check. With OTel on, the source
  is wrapped by `vvotel.Source` through a declared layer.
- **auth** — an `auth.Guard` over static API keys, mounted before every API
  route as ordered middleware (`appfiber.Guarding`). Optional: a request
  without a key passes on anonymously and the route's own policy decides.
- **observability** — the `*slog.Logger` (also `slog.Default`) and the
  `Telemetry` (SDK providers, or no-ops when disabled).

## mod

Bounded contexts: simplified (layered, one package) or original DDD (several
packages). See [mod/README.md](mod/README.md).

## transports

`transports/http` owns the single `*fiber.App` (recover, OTel, request log,
CORS, the framework error envelope), the root routes (`/`, `/favicon.ico`), the
health routes and `appfiber.Serving`. Serving mounts every contributed route
under `transports.fiber.prefix` (`/api`), **compares Fiber's routing table with
every route's access declaration and refuses to start on any difference**,
waits for the port, and shuts the server down gracefully when fx stops.

| Path | Access |
|---|---|
| `GET /`, `GET /favicon.ico` | public |
| `GET /api/health/live` | public — runs nothing |
| `GET /api/health/ready` | public — 503 only when a required check is down |
| `GET /api/health/detail` | `health:detail` |
| module routes | whatever each module declares |

Add another transport (gRPC, a queue consumer) as a sibling package with its
own `Register`, and add it in `boot`.

## Rules the framework enforces here

- Every route declares its access, and enforcement is a projection of the same
  declaration (D-073, D-100). The template uses `appfiber.RefusingUnchecked`,
  so a hand-written `Requires` that nothing checks fails the start too.
- No bare goroutines: background work is a `runtime.Runner` contributed with
  `runtimefx.AsRunner` and supervised (D-092).
- The application opens the database and hands it over (D-057); repositories
  never name a driver.
- Framework code logs through the logger in the request context (D-062).
