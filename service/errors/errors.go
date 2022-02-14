package errors

type (
	ServiceBusinessError struct {
		Msg string
	}

	ServiceNotFoundByIDError struct {
		ID string
	}

	ServiceAlreadyExistsError struct {
		Msg string
	}
)

func (e *ServiceBusinessError) Error() string {
	return e.Msg
}

func (e *ServiceNotFoundByIDError) Error() string {
	return e.ID
}

func (e *ServiceAlreadyExistsError) Error() string {
	return e.Msg
}
