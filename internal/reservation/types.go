package reservation

const TaskQueue = "reservation-task-queue"

type Status string

const (
	StatusValidating        Status = "VALIDATING"
	StatusAwaitingPayment   Status = "AWAITING_PAYMENT"
	StatusProcessingPayment Status = "PROCESSING_PAYMENT"
	StatusConfirmed         Status = "CONFIRMED"
	StatusFailed            Status = "FAILED"
)

// Input is the workflow's starting parameter, submitted by the API.
type Input struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Date   string `json:"date"`
	Guests int    `json:"guests"`
}

// QueryResult is returned by the "status" workflow query so the API can
// poll progress without waiting on workflow completion.
type QueryResult struct {
	Status           Status `json:"status"`
	ConfirmationCode string `json:"confirmationCode,omitempty"`
	ErrorMessage     string `json:"errorMessage,omitempty"`
}
