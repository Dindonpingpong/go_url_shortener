package errors

import "fmt"

type(
	StorageEmptyResultError struct{
		ID string
	}

	StorageAlreadyExistsError struct{
		ShortURL string
	}
)

func (e *StorageEmptyResultError) Error() string {
	return fmt.Sprintf("%s not found in storage", e.ID)
}

func (e *StorageAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists", e.ShortURL)
}