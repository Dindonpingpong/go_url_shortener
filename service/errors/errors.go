package errors

type(
	ServiceBusinessError struct{
		Msg string
	}

)

func (e *ServiceBusinessError) Error() string {
	return e.Msg
}