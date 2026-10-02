package main

import (
	"errors"
	"fmt"
)

// Custom error types
type ValidationError struct {
	Field   string
	Message string
}
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       int
}
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id=%d not found", e.Resource, e.ID)
}

// Sentinel errors
var ErrUnauthorized = errors.New("unauthorized")

// Error wrapping with %w
func getUser(id int) (string, error) {
	if id <= 0 {
		return "", &ValidationError{Field: "id", Message: "must be positive"}
	}
	if id > 100 {
		return "", fmt.Errorf("getUser: %w", &NotFoundError{Resource: "User", ID: id})
	}
	return fmt.Sprintf("User-%d", id), nil
}

func processUser(id int) error {
	user, err := getUser(id)
	if err != nil {
		return fmt.Errorf("processUser: %w", err)
	}
	fmt.Println("Processing:", user)
	return nil
}

// panic + recover — recover must be in a deferred function
func safeDiv(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	result = a / b // panics if b == 0
	return
}

func main() {
	// errors.Is — match sentinel
	err := fmt.Errorf("access denied: %w", ErrUnauthorized)
	fmt.Println(errors.Is(err, ErrUnauthorized)) // true

	// errors.As — extract concrete type
	if err := processUser(999); err != nil {
		var nfe *NotFoundError
		if errors.As(err, &nfe) {
			fmt.Printf("Not found: %s #%d\n", nfe.Resource, nfe.ID)
		}
		fmt.Println("Full error:", err)
	}

	if err := processUser(-1); err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("Validation: field=%s msg=%s\n", ve.Field, ve.Message)
		}
	}

	// panic / recover
	res, err := safeDiv(10, 2)
	fmt.Printf("10/2 = %d, err=%v\n", res, err)

	res, err = safeDiv(10, 0)
	fmt.Printf("10/0 = %d, err=%v\n", res, err)
}