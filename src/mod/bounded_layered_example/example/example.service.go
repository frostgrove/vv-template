package example

import (
	"context"
	"strings"

	"github.com/frostgrove/vv/errs"
)

// service holds the use cases. It validates the command, applies the rules and
// talks to the repository; it knows nothing about HTTP.
type service struct {
	repo *repo
}

func newService(repo *repo) *service {
	return &service{repo: repo}
}

const openPageSize = 50

func (this *service) create(ctx context.Context, command createExampleCommand) (exampleDto, error) {
	title := strings.TrimSpace(command.Title)
	if title == "" {
		return exampleDto{}, errs.Validation().Code(errs.CodeRequired).Field("title").Message("title is required").Fault()
	}
	created, err := this.repo.create(ctx, title)
	if err != nil {
		return exampleDto{}, err
	}
	return toDto(created), nil
}

func (this *service) listOpen(ctx context.Context) ([]exampleDto, error) {
	open, err := this.repo.open(ctx, openPageSize)
	if err != nil {
		return nil, err
	}
	result := make([]exampleDto, 0, len(open))
	for _, model := range open {
		result = append(result, toDto(model))
	}
	return result, nil
}

func (this *service) complete(ctx context.Context, command completeExampleCommand) (exampleDto, error) {
	completed, err := this.repo.complete(ctx, command.Id)
	if err != nil {
		return exampleDto{}, err
	}
	return toDto(completed), nil
}
