package errors

type(
	ServiceBusinessError struct{
		Msg string
	}

	ServiceNotFoundByIdError struct{
		ID string
	}

	ServiceAlreadyExistsError struct{
		Msg string
	}
)

func (e *ServiceBusinessError) Error() string {
	return e.Msg
}

func (e *ServiceNotFoundByIdError) Error() string {
	return e.ID
}

func (e *ServiceAlreadyExistsError) Error() string {
	return e.Msg
}