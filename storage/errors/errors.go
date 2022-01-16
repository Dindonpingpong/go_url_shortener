package errors

import "fmt"

type(
	StorageEmptyResultError struct{
		ID string
	}
)

func (e *StorageEmptyResultError) Error() string {
	return fmt.Sprintf("%s not found in storage", e.ID)
}