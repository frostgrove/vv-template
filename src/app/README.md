# App Guide

## Bounded Contexts

each folder is bounded context that can be simple (from only 1 package) or complex (from 2 or more packages). Bounded Context can have only 1 root aggregate.

package default structure:
.api - external api. Public interfaces & types. Only from root package
.model - ordinary Go database model. Put exported models in `*.model.go`;
vv derives snake_case columns, the ID/Id primary key and the default table name
without ORM or `db` tags.
.repo - generated repository blueprint and binding factory in `vv_gen.go`
.service - usecases
.dto - data transfer objects. Usually for output types
.commands - value objects that passes to usecases
.module - register package in fx. declares func Register() *fx.Option.
.http-handler - http requests handler

## Examples

See the `example` bounded context for a simple exmaple
See the `product` bounded context for a mixed example

Run `make generate` after changing models. It recursively finds model files
and writes DTOs, typed metamodels and driver-neutral repository factories next
to their packages. Bind the factory to whichever `crud.Source` the application
uses: `database/sql`, native pgx, a test source, or another vv adapter.
See the `wirehouse` bounded context for a complex example
