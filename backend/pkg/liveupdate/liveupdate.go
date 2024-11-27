package liveupdate

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Update struct {
	TotalNumberofRequestsDone int32 `json:"total_numberof_requests"`
	SucceededRequests         int32 `json:"succeeded_requests"`
	FailedRequests            int32 `json:"failed_requests"`
	TargetUsers               int32 `json:"target_users"`
}

type Updater interface {
	Set(id uuid.UUID, u *Update)
	Get(id uuid.UUID) (*Update, error)
	Delete(id uuid.UUID)
}

type updater struct {
	m sync.Map
}

func New() Updater {
	return &updater{
		m: sync.Map{},
	}
}

func (ur *updater) Set(id uuid.UUID, u *Update) {
	ur.m.Store(id, u)
}

func (ur *updater) Get(id uuid.UUID) (*Update, error) {
	u, ok := ur.m.Load(id)
	if !ok {
		logrus.Error("error in getting update")
		return nil, errors.New("error in getting update")
	}

	update, ok := u.(*Update)
	if !ok {
		logrus.Error("erorr in getting update")
		return nil, errors.New("error in getting update")

	}
	return update, nil
}

// Delete implements Updater.
func (ur *updater) Delete(id uuid.UUID) {
	ur.m.Delete(id)
}
