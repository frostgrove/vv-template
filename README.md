# vv service template

A Go service on the [vv framework](https://github.com/frostgrove/vv): uber/fx
composition, Fiber v3 HTTP with a boot-time access gate, driver-neutral
repositories generated from plain Go models, Goose migrations, health probes,
supervised background workers, slog logging and optional OpenTelemetry.

## Quick start

```bash
cp config/example.yml config/app.yml   # git-ignored; edit as needed
make up                                # postgres 18 on localhost:5432 (compose.dev.yml)
make migrate                           # apply resources/migrations
make seed                              # optional: demo products (env: dev only, safe to rerun)
make dev                               # http://localhost:8009
```

```bash
curl localhost:8009/api/health/live
curl localhost:8009/api/health/ready
curl -H 'X-Api-Key: dev-reader-key' 'localhost:8009/api/products?sort=-createdAt'
curl -H 'X-Api-Key: dev-admin-key' -H 'Content-Type: application/json' \
     -d '{"name":"Hex bolt"}' localhost:8009/api/products
curl localhost:8009/api/products                       # 401: no key
curl -H 'X-Api-Key: dev-reader-key' -X DELETE localhost:8009/api/products/<id>   # 403
```

`Ctrl+C` stops the server gracefully: the listener drains, workers stop, the
pool closes, telemetry flushes. `make down` stops postgres.

| Target | Does |
|---|---|
| `make dev` | run the service (`cmd/main.go`) with `config/app.yml` |
| `make migrate [args]` | the vvgoose CLI (`cmd/migrate`); bare = apply pending. `make migrate status`, `make migrate rollback`, `make migrate migration add_x`, `make migrate table users`. `make goose …` is an alias |
| `make seed` | run every seeder the modules contribute for `app.env` (`cmd/seed`) |
| `make generate` / `make generate-check` | regenerate / verify `vv_gen.go` beside every `*.model.go` under `src/mod` |
| `make build` `make vet` `make test` `make lint` | the usual; `lint` needs golangci-lint |
| `make up` / `make down` | the dev postgres |
| `make workspace` | build against a local framework checkout (see below) |

The layout and the conventions are in [src/README.md](src/README.md); how a
bounded context is written and added is in [src/mod/README.md](src/mod/README.md).

## Configuration

`config/app.yml` (or `--config-path`, or `CONFIG_PATH`) is loaded **strictly**
by `vvcfg`: a key no field declares refuses the start (framework decision
D-086), and every block validates itself. Every scalar has an environment
variable (see `src/config/config.go`; the database block uses `DB_*`).

The dev config uses `sslmode: disable`. Omitting `sslmode` means verified TLS
(`verify-full`, D-081), which is what a deployed environment should get.

`app.profile` decides which roles of the modules this process wires:
`complete` (routes + workers, the default), `serving` (routes only) or
`working` (workers only). Health routes are always served.

## Auth

The dev config issues two static API keys (`X-Api-Key`): `dev-admin-key` holds
every permission, `dev-reader-key` only the `:read` ones. Every route declares
what reaching it requires — `Requires(perm)`, `Authenticated(why)` or
`Public(why)` — and the service refuses to start if Fiber's routing table and
those declarations disagree in either direction (D-073, D-100). Replace the
static keys (`src/core/auth`) with a hashed key store, `auth/authjwt` or
`auth/access` before deploying; routes do not change.

## Observability

Logs are `log/slog` (text or JSON, `observability.logger`), wrapped by
`vvotel.TraceHandler` so a line logged with a request context carries
`trace_id`/`span_id`. The logger is also put into every request context for
framework code (`port.Logger`, D-062).

`observability.otel.enabled: true` (or `OTEL_ENABLED=true`) starts the
OpenTelemetry SDK with OTLP/gRPC export of traces and metrics to
`observability.otel.endpoint` (or the standard `OTEL_EXPORTER_OTLP_*`
variables), instruments Fiber requests and wraps the repository source with
`vvotel.Source`. It needs a collector listening there; without one, shutdown
waits for the flush timeout (about 10 s) and exits non-zero. Off (the default),
nothing is started or exported.

## Framework version

The framework (`github.com/frostgrove/vv` and its satellite modules) has no
release tags yet. `go.mod` therefore requires every vv module at the
pseudo-version of one framework commit, and adds two `replace` lines that pin
the placeholder version `v0.0.0-00010101000000-000000000000` — which the
satellites use for the root module and for `authfiber` — to that same commit.
Everything resolves through the public Go proxy; no framework checkout is
needed.

To move to another framework commit `<sha>`:

```bash
OLD=$(go list -m -f '{{.Version}}' github.com/frostgrove/vv)
NEW=$(go list -m -f '{{.Version}}' github.com/frostgrove/vv@<sha>)
sed -i "s/$OLD/$NEW/g" go.mod
go mod tidy && make build
```

Check afterwards that no frostgrove requirement is left at another version and
that a satellite has not started requiring a new placeholder module (`go mod
tidy` fails with `unknown revision 000000000000`; add a third replace line in
that case). When the framework is tagged, drop both replace lines and require
the tag.

To work against a local framework checkout (sibling directory `../framework`
by default), run `make workspace` (`FRAMEWORK=path` to override). It writes a
`go.work` that is git-ignored; delete it to go back to the pinned commit.
