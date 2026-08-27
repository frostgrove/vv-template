# App Guide

## Bounded Contexts

each folder is bounded context that can be simple (from only 1 package) or complex (from 2 or more packages). Bounded Context can have only 1 root aggregate.

package default structure:
.api - external api. Public interfaces & types. Only from root package
.repo - repository
.service - usecases
.dto - data transfer objects. Usually for output types
.commands - value objects that passes to usecases
.module - register package in fx. declares func Register() *fx.Option.
.model - db model
.http-handler - http requests handler

## Examples

See the `example` bounded context for a simple exmaple
See the `product` bounded context for a mixed example
See the `wirehouse` bounded context for a complex example