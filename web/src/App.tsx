import { useEffect, useRef, useState } from 'react'
import {
  confirmPayment,
  createReservation,
  getReservationStatus,
  type ReservationInput,
  type ReservationStatus,
} from './api'

const STEPS: { key: ReservationStatus['status']; label: string }[] = [
  { key: 'VALIDATING', label: 'Validating' },
  { key: 'AWAITING_PAYMENT', label: 'Awaiting payment' },
  { key: 'PROCESSING_PAYMENT', label: 'Processing payment' },
  { key: 'CONFIRMED', label: 'Confirmed' },
]

function stepIndex(status: ReservationStatus['status']): number {
  if (status === 'FAILED') return -1
  return STEPS.findIndex((s) => s.key === status)
}

const initialForm: ReservationInput = { name: '', email: '', date: '', guests: 2 }

export default function App() {
  const [form, setForm] = useState<ReservationInput>(initialForm)
  const [reservationId, setReservationId] = useState<string | null>(null)
  const [status, setStatus] = useState<ReservationStatus | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [paying, setPaying] = useState(false)
  const [payError, setPayError] = useState<string | null>(null)
  const pollRef = useRef<number | null>(null)

  useEffect(() => {
    return () => {
      if (pollRef.current) window.clearInterval(pollRef.current)
    }
  }, [])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSubmitError(null)
    setSubmitting(true)
    try {
      const id = await createReservation(form)
      setReservationId(id)
      setStatus({ status: 'VALIDATING' })
      pollRef.current = window.setInterval(async () => {
        try {
          const s = await getReservationStatus(id)
          setStatus(s)
          if (s.status === 'CONFIRMED' || s.status === 'FAILED') {
            if (pollRef.current) window.clearInterval(pollRef.current)
          }
        } catch {
          // transient — keep polling
        }
      }, 1000)
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'failed to submit reservation')
    } finally {
      setSubmitting(false)
    }
  }

  async function handlePay() {
    if (!reservationId) return
    setPayError(null)
    setPaying(true)
    try {
      await confirmPayment(reservationId)
    } catch (err) {
      setPayError(err instanceof Error ? err.message : 'failed to submit payment')
    } finally {
      setPaying(false)
    }
  }

  function reset() {
    setReservationId(null)
    setStatus(null)
    setPayError(null)
    setForm(initialForm)
  }

  return (
    <div className="page">
      <div className="card">
        <h1>Make a reservation</h1>

        {!reservationId && (
          <form onSubmit={handleSubmit}>
            <label>
              Name
              <input
                required
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
              />
            </label>
            <label>
              Email
              <input
                required
                type="email"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
              />
            </label>
            <label>
              Date
              <input
                required
                type="date"
                value={form.date}
                onChange={(e) => setForm({ ...form, date: e.target.value })}
              />
            </label>
            <label>
              Guests
              <input
                required
                type="number"
                min={1}
                max={20}
                value={form.guests}
                onChange={(e) => setForm({ ...form, guests: Number(e.target.value) })}
              />
            </label>

            {submitError && <p className="error">{submitError}</p>}

            <button type="submit" disabled={submitting}>
              {submitting ? 'Submitting…' : 'Reserve'}
            </button>
          </form>
        )}

        {reservationId && status && (
          <div className="progress">
            <p className="reservation-id">Reservation {reservationId}</p>

            {status.status !== 'FAILED' ? (
              <ol className="steps">
                {STEPS.map((step, i) => {
                  const current = stepIndex(status.status)
                  const done = i < current || status.status === 'CONFIRMED' && i <= current
                  const active = i === current && status.status !== 'CONFIRMED'
                  return (
                    <li key={step.key} className={done ? 'done' : active ? 'active' : ''}>
                      {step.label}
                    </li>
                  )
                })}
              </ol>
            ) : (
              <p className="error">Reservation failed: {status.errorMessage}</p>
            )}

            {status.status === 'AWAITING_PAYMENT' && (
              <div className="pay-box">
                <p>Your reservation is validated. Pay now to confirm it.</p>
                {payError && <p className="error">{payError}</p>}
                <button onClick={handlePay} disabled={paying}>
                  {paying ? 'Submitting…' : 'Pay now'}
                </button>
              </div>
            )}

            {status.status === 'CONFIRMED' && (
              <p className="confirmation">Confirmation code: {status.confirmationCode}</p>
            )}

            {(status.status === 'CONFIRMED' || status.status === 'FAILED') && (
              <button onClick={reset}>Make another reservation</button>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
