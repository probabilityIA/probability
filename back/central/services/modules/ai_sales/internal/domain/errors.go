package domain

import "fmt"

type ErrSessionNotFound struct {
	PhoneNumber string
}

func (e *ErrSessionNotFound) Error() string {
	return fmt.Sprintf("session not found for phone: %s", e.PhoneNumber)
}

type ErrBedrockUnavailable struct {
	Cause error
}

func (e *ErrBedrockUnavailable) Error() string {
	return fmt.Sprintf("bedrock unavailable: %v", e.Cause)
}

type ErrToolExecutionFailed struct {
	ToolName string
	Cause    error
}

func (e *ErrToolExecutionFailed) Error() string {
	return fmt.Sprintf("tool %s execution failed: %v", e.ToolName, e.Cause)
}

type ErrMaxToolIterations struct {
	Count int
}

func (e *ErrMaxToolIterations) Error() string {
	return fmt.Sprintf("max tool iterations reached: %d", e.Count)
}

type ErrProductNotFound struct {
	SKU string
}

func (e *ErrProductNotFound) Error() string {
	return fmt.Sprintf("product not found: %s", e.SKU)
}

type ErrProductOutOfStock struct {
	SKU string
}

func (e *ErrProductOutOfStock) Error() string {
	return fmt.Sprintf("product out of stock: %s", e.SKU)
}
