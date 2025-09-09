package observers

import (
	"github.com/goravel/framework/contracts/database/orm"
)

type UserObserver struct{}

func (u *UserObserver) Created(event orm.Event) error {
	return nil
}

func (u *UserObserver) Updated(event orm.Event) error {
	return nil
}

func (u *UserObserver) Deleted(event orm.Event) error {
	return nil
}

func (u *UserObserver) ForceDeleted(event orm.Event) error {
	return nil
}
