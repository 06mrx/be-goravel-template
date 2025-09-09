package providers

import (
	"goravel_api/app/rules"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type ValidationServiceProvider struct {
}

func (receiver *ValidationServiceProvider) Register(app foundation.Application) {

}
func (receiver *ValidationServiceProvider) Boot(app foundation.Application) {
	if err := facades.Validation().AddRules(receiver.rules()); err != nil {
		facades.Log().Errorf("add rules error: %+v", err)
	}
	if err := facades.Validation().AddFilters(receiver.filters()); err != nil {
		facades.Log().Errorf("add filters error: %+v", err)
	}
}
func (receiver *ValidationServiceProvider) rules() []validation.Rule {
	return []validation.Rule{
		&rules.Same{},
		&rules.Unique{},
		&rules.Digits{},
		&rules.DigitsBetween{},
		// &rules.MinLen{},
		// &rules.MaxLen{},
		&rules.Filetype{},
	}
}
func (receiver *ValidationServiceProvider) filters() []validation.Filter {
	return []validation.Filter{}
}
