export type Status = 'VALIDATING' | 'AWAITING_PAYMENT' | 'PROCESSING_PAYMENT' | 'CONFIRMED' | 'FAILED'

export interface ReservationInput {
  name: string
  email: string
  date: string
  guests: number
}

export interface ReservationStatus {
  status: Status
  confirmationCode?: string
  errorMessage?: string
}

export async function createReservation(input: ReservationInput): Promise<string> {
  const res = await fetch('/api/reservations', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  if (!res.ok) {
    throw new Error(await res.text())
  }
  const body: { id: string } = await res.json()
  return body.id
}

export async function confirmPayment(id: string): Promise<void> {
  const res = await fetch(`/api/reservations/${id}/pay`, { method: 'POST' })
  if (!res.ok) {
    throw new Error(await res.text())
  }
}

export async function getReservationStatus(id: string): Promise<ReservationStatus> {
  const res = await fetch(`/api/reservations/${id}`)
  if (!res.ok) {
    throw new Error(await res.text())
  }
  return res.json()
}
