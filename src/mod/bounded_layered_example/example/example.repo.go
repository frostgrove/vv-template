package example

import (
	"context"

	"github.com/frostgrove/vv/crud"
	"github.com/frostgrove/vv/crud/decorators/specs"
	"github.com/google/uuid"
)

// repo holds the queries this context asks, written against the generated
// repository (*ExampleRepo) and its typed metamodel (Example_). It never
// names a driver: the crud.Source behind it is chosen by the composition root.
type repo struct {
	examples *ExampleRepo
}

func newRepo(examples *ExampleRepo) *repo {
	return &repo{examples: examples}
}

func (this *repo) create(ctx context.Context, title string) (Example, error) {
	return this.examples.Save(ctx, &Example{Title: title})
}

func (this *repo) open(ctx context.Context, limit int) ([]Example, error) {
	return this.examples.GetAll(ctx,
		crud.Where(specs.Predicate(Example_.Done.Eq(false))),
		crud.OrderBy(crud.Desc("CreatedAt")),
		crud.Limit(limit),
	)
}

func (this *repo) countOpen(ctx context.Context) (int64, error) {
	return this.examples.Count(ctx, crud.Where(specs.Predicate(Example_.Done.Eq(false))))
}

func (this *repo) complete(ctx context.Context, id uuid.UUID) (Example, error) {
	done := true
	return this.examples.Update(ctx, id, ExampleUpdate{Done: &done})
}
