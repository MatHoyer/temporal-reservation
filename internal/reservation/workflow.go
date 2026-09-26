package reservation

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// activities is only used here to reference activity methods by name when
// building ExecuteActivity calls — the actual receiver runs on the worker.
var activities = &Activities{}

const StatusQuery = "status"
const ConfirmPaymentSignal = "confirm_payment"

func ReservationWorkflow(ctx workflow.Context, in Input) (QueryResult, error) {
	result := QueryResult{Status: StatusValidating}

	if err := workflow.SetQueryHandler(ctx, StatusQuery, func() (QueryResult, error) {
		return result, nil
	}); err != nil {
		return result, err
	}

	validateCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 1},
	})
	if err := workflow.ExecuteActivity(validateCtx, activities.ValidateReservation, in).Get(validateCtx, nil); err != nil {
		result.Status = StatusFailed
		result.ErrorMessage = err.Error()
		return result, nil
	}

	// Wait for the user to actually hit "pay" instead of charging
	// automatically — the workflow blocks here for however long that takes,
	// which is the point: it's durable, so the worker can restart mid-wait
	// with no state lost.
	result.Status = StatusAwaitingPayment
	workflow.GetSignalChannel(ctx, ConfirmPaymentSignal).Receive(ctx, nil)

	result.Status = StatusProcessingPayment
	paymentCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    4,
		},
	})
	if err := workflow.ExecuteActivity(paymentCtx, activities.ChargePayment, in).Get(paymentCtx, nil); err != nil {
		result.Status = StatusFailed
		result.ErrorMessage = err.Error()
		return result, nil
	}

	confirmCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 2},
	})
	var code string
	if err := workflow.ExecuteActivity(confirmCtx, activities.SendConfirmation, in).Get(confirmCtx, &code); err != nil {
		result.Status = StatusFailed
		result.ErrorMessage = err.Error()
		return result, nil
	}

	result.Status = StatusConfirmed
	result.ConfirmationCode = code
	return result, nil
}
