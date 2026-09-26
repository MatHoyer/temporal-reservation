package reservation

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
)

type Activities struct{}

// ValidateReservation does basic sanity checks. No external calls — this is
// where a real system would check availability, blackout dates, etc.
func (a *Activities) ValidateReservation(_ context.Context, in Input) error {
	if in.Guests <= 0 || in.Guests > 20 {
		return fmt.Errorf("guest count %d out of range (1-20)", in.Guests)
	}
	if _, err := time.Parse("2006-01-02", in.Date); err != nil {
		return fmt.Errorf("invalid date %q: %w", in.Date, err)
	}
	return nil
}

// ChargePayment mocks a payment provider call: it takes a random amount of
// time and fails roughly a third of the time, so the activity's retry
// policy (configured on the workflow side) actually gets exercised.
func (a *Activities) ChargePayment(ctx context.Context, in Input) error {
	activity.RecordHeartbeat(ctx, "charging payment")
	time.Sleep(time.Duration(800+rand.Intn(2200)) * time.Millisecond)

	if rand.Intn(3) == 0 {
		return fmt.Errorf("payment provider declined the charge for %s", in.Email)
	}
	return nil
}

// SendConfirmation mocks a confirmation email/SMS and returns a fake code.
func (a *Activities) SendConfirmation(_ context.Context, in Input) (string, error) {
	time.Sleep(300 * time.Millisecond)
	code := fmt.Sprintf("RSV-%04d", rand.Intn(10000))
	return code, nil
}
