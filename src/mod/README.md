# Bounded contexts

Each folder under `src/mod` is a bounded context. It can be simple (one
package, layered by file) or complex (several packages, DDD). A bounded context
has one root aggregate and talks to others only through what its `.api` exposes.

The examples:

| Folder | Shows |
|---|---|
| `bounded_layered_example/product` | a CRUD resource: generated repository, the ten crudfiber routes behind declared permissions, a seeder |
| `bounded_layered_example/example` | a hand-written use case through every layer: commands, service, repo queries over the metamodel, DTOs, a public API, a background worker |
| `ddd_example/inventory` | a richer domain model kept in a `domain` package (models only so far) |

## File roles in a package

Files are named `<context>.<role>.go`.

| File | Holds |
|---|---|
| `.model.go` | persistence models: ordinary Go structs. vv derives the table (plural snake_case), the columns (snake_case) and the `Id`/`ID` primary key. Tags only for exceptions: `db:"id,pk,auto"` (the database generates the key), `db:",generated"` (never written, read back: defaults, triggers), `db:",immutable"`, `rel:"…"`; `json` for the wire shape |
| `vv_gen.go` | **generated, do not edit**: the update DTO (`<M>Update`), the typed metamodel (`<M>_`), the repository blueprint and `New<M>Repository(crud.Source) *<M>Repo`, and an init check that fails the start if the model and the file drifted |
| `.repo.go` | the queries this context asks, written against `*<M>Repo` and `<M>_` (`crud.Where(specs.Predicate(Example_.Done.Eq(false)))`). Never a driver or an ORM |
| `.command.go` | values a use case takes: what the caller asked for, decoded |
| `.service.go` | use cases: validate, apply the rules, call the repo. Validation errors are `errs` faults (`errs.Validation().Code(…).Field(…)`), which every transport renders as 422 with the field |
| `.dto.go` | what leaves the context; the model stays inside |
| `.api.go` | the interface other contexts may depend on |
| `.access.go` | the permissions (`auth.Permission`) this context's routes require |
| `.http-handler.go` | a constructor returning `appfiber.Route`: decode, call the service, encode |
| `.worker.go` | background work: a constructor returning a `runtime.Runner` (`runtime.NewPeriodic`, or your own `Run(ctx)` that blocks until ctx is done). Never `go func()` |
| `.seeder.go` | an `app.Seeder` that is safe to run twice, limited to the envs it names |
| `.module.go` | `Module() module.Definition`: everything the context contributes, filed by role |

### The module file

```go
func Module() module.Definition {
	return module.New("example").
		Order(200).                                              // position in the catalog
		Provide(NewExampleRepository, newRepo, newService, NewAPI). // every process
		Routes(appfiber.AsRoute(newHttpHandler)).                 // profiles with the api role
		Workers(runtimefx.AsRunner(newOpenReporter)).             // profiles with the worker role
		Seeders(appfx.AsSeeder(newSeeder)).                      // the seed command
		Checks(healthfx.AsCheck(newCheck)).                       // readiness, every process
		MustBuild()
}
```

Each constructor is wrapped with the annotation its subsystem owns
(`appfiber.AsRoute`, `runtimefx.AsRunner`, `appfx.AsSeeder`,
`healthfx.AsCheck`); `module` only files it under a role. The profile chosen in
`boot` decides which roles a process wires, so an API replica never starts a
worker and the seed command mounts no route.

### Routes declare their access

Every route says what reaching it takes, and the service refuses to start if a
mounted route has no declaration or a declaration has no route:

```go
return appfiber.Routes("/examples").
	GET("", appfiber.Requires(PermRead), handler.listOpen).
	POST("", appfiber.Requires(PermWrite), handler.create).
	POST("/:id/complete", appfiber.Requires(PermWrite), handler.complete).
	Route()
```

`Requires(perms…)` mounts the check in front of the handler (401 without a
caller, 403 without the permission), `Authenticated(why)` asks only for a
caller, `Public(why)` mounts none and needs the reason. Paths are relative to
the API prefix (`/api`).

A CRUD resource goes through the same registrar, using crudfiber's handlers one
by one (see `product.http-handler.go`). The alternative the framework also
supports — `crudfiber.Register` with the declaration derived from the
repository's `security.Gate` policy (`crudhttp.Table.GuardedBy`, D-107) — needs
an `appfiber.ExcusingUnchecked` rule in `transports/http`, because the template
otherwise refuses declarations nothing in front of the handler checks.

## Add a bounded context

1. Create `src/mod/<area>/<name>/` with `<name>.model.go`.
2. `make generate` writes `vv_gen.go` beside it.
3. Add a migration. `make migrate table <names>` renders one from the model,
   but it refuses a UUID key tagged `auto`; then use
   `make migrate migration create_<names>` and write the SQL by hand
   (`DEFAULT gen_random_uuid()`, `DEFAULT now()` for generated columns).
   `make migrate` applies it.
4. Write the layers you need (`.repo`, `.service`, `.http-handler`, …) and the
   permissions in `.access.go`.
5. Write `Module()` in `<name>.module.go`.
6. Add it to `boot/app.go`'s `module.MustCatalog(...)`.
7. Grant the new permissions to an API key in `config/app.yml` (dev only).
8. `make test` checks every profile over the catalog; `make dev`.

## Generation

`make generate` runs `go run github.com/frostgrove/vv/cmd/vv generate -dir
./src/mod`. It walks the tree, reads exported structs in `*.model.go`, and
writes `vv_gen.go` beside each package. Run it after every model change;
`make generate-check` fails if any `vv_gen.go` is behind (use it in CI). Even
without it, a model changed and not regenerated refuses to start.

The framework can also generate the public wire bodies
(`vv generate resource`) and the module definition itself
(`vv generate module`, from a reviewed `module.manifest.yml`). The template
writes `Module()` by hand: a generated definition lists the constructors
without the subsystem annotations, so its routes and checks do not reach
`appfiber` and `healthfx`.
