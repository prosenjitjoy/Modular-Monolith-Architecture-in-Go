package domain

import "errors"

var (
	ErrNameCannotBeBlank       = errors.New("customer name cannot be blank")
	ErrCustomerIDCannotBeBlank = errors.New("customer id cannot be blank")
	ErrSmsNumberCannotBeBlank  = errors.New("SMS number cannot be blank")
	ErrCustomerAlreadyEnabled  = errors.New("customer is already enabled")
	ErrCustomerAlreadyDisabled = errors.New("customer is already disabled")
)

type Customer struct {
	ID        string
	Name      string
	SmsNumber string
	Enabled   bool
}

func RegisterCustomer(id string, name string, smsNumber string) (*Customer, error) {
	if id == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if name == "" {
		return nil, ErrNameCannotBeBlank
	}

	if smsNumber == "" {
		return nil, ErrSmsNumberCannotBeBlank
	}

	return &Customer{
		ID:        id,
		Name:      name,
		SmsNumber: smsNumber,
		Enabled:   true,
	}, nil
}

func (c *Customer) Enable() error {
	if c.Enabled {
		return ErrCustomerAlreadyEnabled
	}

	c.Enabled = true

	return nil
}

func (c *Customer) Disable() error {
	if !c.Enabled {
		return ErrCustomerAlreadyDisabled
	}

	c.Enabled = false

	return nil
}
