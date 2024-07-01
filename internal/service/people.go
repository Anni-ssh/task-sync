package service

import (
	"TaskSync/internal/entities"
	"TaskSync/internal/storage"
	"context"
)

// PeopleService представляет сервис для работы с данными пользователей.
type PeopleService struct {
	storage storage.PeopleManage
}

// NewPeopleService создает новый экземпляр PeopleService.
func NewPeopleService(s storage.PeopleManage) *PeopleService {
	return &PeopleService{storage: s}
}

// Create создает новую запись пользователя.
func (p *PeopleService) Create(ctx context.Context, people entities.People) (int, error) {
	return p.storage.Create(ctx, people)
}

// GetByID возвращает данные пользователя по его ID.
func (p *PeopleService) GetByID(ctx context.Context, peopleID int) (entities.People, error) {
	return p.storage.GetByID(ctx, peopleID)
}

// GetByFilter возвращает список пользователей, отфильтрованных по указанным параметрам.
func (p *PeopleService) GetByFilter(ctx context.Context, filterPeople entities.People, limit, offset int) ([]entities.People, error) {
	return p.storage.GetByFilter(ctx, filterPeople, limit, offset)
}

// List возвращает список всех пользователей.
func (p *PeopleService) List(ctx context.Context) ([]entities.People, error) {
	return p.storage.List(ctx)
}

// Update обновляет данные пользователя.
func (p *PeopleService) Update(ctx context.Context, people entities.People) error {
	return p.storage.Update(ctx, people)
}

// Delete удаляет пользователя по его ID.
func (p *PeopleService) Delete(ctx context.Context, peopleID int) error {
	return p.storage.Delete(ctx, peopleID)
}
