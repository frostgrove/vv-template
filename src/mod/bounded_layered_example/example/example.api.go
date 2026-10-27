package example

import "context"

// API is what other bounded contexts may use. They depend on this interface
// (provided by the module), never on the service, the repository or the model.
type API interface {
	CountOpen(ctx context.Context) (int64, error)
}

type api struct {
	repo *repo
}

func NewAPI(repo *repo) API {
	return &api{repo: repo}
}

func (this *api) CountOpen(ctx context.Context) (int64, error) {
	return this.repo.countOpen(ctx)
}
