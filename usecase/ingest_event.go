package usecase

import (
	"context"
	"github.com/wahyusumartha/analytics-open-search/repository"
)

type IngestEvent interface {
	Execute(ctx context.Context, input EventInput) error
}

func NewIngestEventUseCase(repo repository.EventRepository) IngestEventUseCase {
	return IngestEventUseCase{repo: repo}
}

type IngestEventUseCase struct {
	repo repository.EventRepository
}

func (uc IngestEventUseCase) Execute(ctx context.Context, input EventInput) error {
	err := input.Validate()
	if err != nil {
		return err
	}

	event := repository.Event{
		ID:         input.ID,
		Name:       input.Name,
		CreatedAt:  input.Timestamp,
		Properties: input.Properties,
	}
	err = uc.repo.Save(ctx, event)
	return err
}
