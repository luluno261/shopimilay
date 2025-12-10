import React, { useState } from 'react'

interface OrderStatusUpdaterProps {
  orderId: string
  currentStatus: string
  onUpdate: (orderId: string, status: string) => void
}

const STATUS_OPTIONS = [
  { value: 'pending', label: 'En attente' },
  { value: 'paid', label: 'Payé' },
  { value: 'shipped', label: 'Expédié' },
  { value: 'delivered', label: 'Livré' },
  { value: 'cancelled', label: 'Annulé' },
  { value: 'failed', label: 'Échoué' },
]

export default function OrderStatusUpdater({
  orderId,
  currentStatus,
  onUpdate,
}: OrderStatusUpdaterProps) {
  const [status, setStatus] = useState(currentStatus)
  const [updating, setUpdating] = useState(false)

  const handleUpdate = async () => {
    if (status === currentStatus) return

    setUpdating(true)
    try {
      onUpdate(orderId, status)
    } catch (error) {
      console.error('Erreur:', error)
      alert('Erreur lors de la mise à jour')
    } finally {
      setUpdating(false)
    }
  }

  return (
    <div className="flex items-center gap-4">
      <select
        value={status}
        onChange={(e) => setStatus(e.target.value)}
        className="border border-gray-300 rounded-md px-3 py-2"
        disabled={updating}
      >
        {STATUS_OPTIONS.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      {status !== currentStatus && (
        <button
          onClick={handleUpdate}
          disabled={updating}
          className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
        >
          {updating ? 'Mise à jour...' : 'Mettre à jour'}
        </button>
      )}
    </div>
  )
}

