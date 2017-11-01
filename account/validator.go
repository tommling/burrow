package account

import (
	"encoding/json"

	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

type Validator interface {
	Addressable
	// The validator's voting power
	Power() uint64
	// Alter the validator's voting power by amount that can be negative or positive.
	// A power of 0 effectively unbonds the validator
	WithNewPower(uint64) Validator
}

// Neither abci_types or tm_types has quite the representation we want
type ConcreteValidator struct {
	Address Address
	PubKey  crypto.PubKey
	Power   uint64
}

type concreteValidatorWrapper struct {
	*ConcreteValidator `json:"unwrap"`
}

var _ Validator = concreteValidatorWrapper{}

var _ = wire.RegisterInterface(struct{ Validator }{}, wire.ConcreteType{concreteValidatorWrapper{}, 0x01})

func AsValidator(account Account) Validator {
	return ConcreteValidator{
		Address: account.Address(),
		PubKey:  account.PubKey(),
		Power:   account.Balance(),
	}.Validator()
}

func AsConcreteValidator(validator Validator) *ConcreteValidator {
	if validator == nil {
		return nil
	}
	if ca, ok := validator.(concreteValidatorWrapper); ok {
		return ca.ConcreteValidator
	}
	return &ConcreteValidator{
		Address: validator.Address(),
		PubKey:  validator.PubKey(),
		Power:   validator.Power(),
	}
}

func (cvw concreteValidatorWrapper) Address() Address {
	return cvw.ConcreteValidator.Address
}

func (cvw concreteValidatorWrapper) PubKey() crypto.PubKey {
	return cvw.ConcreteValidator.PubKey
}

func (cvw concreteValidatorWrapper) Power() uint64 {
	return cvw.ConcreteValidator.Power
}

func (cvw concreteValidatorWrapper) WithNewPower(power uint64) Validator {
	cv := cvw.Copy()
	cv.Power = power
	return concreteValidatorWrapper{
		ConcreteValidator: cv,
	}
}

func (cv ConcreteValidator) Validator() Validator {
	return concreteValidatorWrapper{
		ConcreteValidator: &cv,
	}
}

func (cv *ConcreteValidator) Copy() *ConcreteValidator {
	cvCopy := *cv
	return &cvCopy
}

func (cv *ConcreteValidator) String() string {
	if cv == nil {
		return "Nil Validator"
	}

	bs, err := json.Marshal(cv)
	if err != nil {
		return "error serialising Validator"
	}

	return string(bs)
}