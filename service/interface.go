package service

import (
	"context"

	"github.com/Luftalian/shokai_golang_webapp/entity"
	"github.com/Luftalian/shokai_golang_webapp/store"
)

type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}

type UserRegisterer interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}