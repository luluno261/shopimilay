import React, { useState } from 'react'
import api from '../lib/api'

interface Order {
  id: string
  total_amount: number
  currency: string
  payment_intent_id?: string
}

interface RefundModalProps {
  order: Order
  onClose: () => void
  onSuccess: () => void
}

export default function RefundModal({ order, onClose, onSuccess }: RefundModalProps) {
  const [refundType, setRefundType] = useState<'full' | 'partial'>('full')
  const [amount, setAmount] = useState<number>(order.total_amount)
  const [processing, setProcessing] = useState(false)

  const handleRefund = async () => {
    if (refundType === 'partial' && amount <= 0) {
      alert('Le montant doit être supérieur à 0')
      return
    }

    if (refundType === 'partial' && amount > order.total_amount) {
      alert('Le montant ne peut pas dépasser le montant total')
      return
    }

    setProcessing(true)
    try {
      const refundAmount = refundType === 'full' ? undefined : amount
      await api.post(`/orders/${order.id}/refund`, { amount: refundAmount })
      alert('Remboursement initié avec succès')
      onSuccess()
    } catch (error: any) {
      console.error('Erreur:', error)
      alert(error.response?.data?.error || 'Erreur lors du remboursement')
    } finally {
      setProcessing(false)
    }
  }

  const formatCurrency = (value: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: currency || 'EUR',
    }).format(value)
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-8 max-w-md w-full">
        <h2 className="text-2xl font-bold mb-6">Remboursement</h2>

        <div className="space-y-4">
          <div>
            <p className="text-sm text-gray-600 mb-2">Montant total de la commande</p>
            <p className="text-lg font-semibold">
              {formatCurrency(order.total_amount, order.currency)}
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Type de remboursement
            </label>
            <div className="space-y-2">
              <label className="flex items-center">
                <input
                  type="radio"
                  value="full"
                  checked={refundType === 'full'}
                  onChange={(e) => setRefundType(e.target.value as 'full' | 'partial')}
                  className="mr-2"
                />
                Remboursement complet
              </label>
              <label className="flex items-center">
                <input
                  type="radio"
                  value="partial"
                  checked={refundType === 'partial'}
                  onChange={(e) => setRefundType(e.target.value as 'full' | 'partial')}
                  className="mr-2"
                />
                Remboursement partiel
              </label>
            </div>
          </div>

          {refundType === 'partial' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Montant à rembourser
              </label>
              <input
                type="number"
                step="0.01"
                min="0"
                max={order.total_amount}
                value={amount}
                onChange={(e) => setAmount(parseFloat(e.target.value) || 0)}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              />
              <p className="text-sm text-gray-500 mt-1">
                Maximum: {formatCurrency(order.total_amount, order.currency)}
              </p>
            </div>
          )}

          <div className="flex gap-4 mt-6">
            <button
              onClick={handleRefund}
              disabled={processing}
              className="flex-1 bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50"
            >
              {processing ? 'Traitement...' : 'Confirmer le remboursement'}
            </button>
            <button
              onClick={onClose}
              className="flex-1 bg-gray-200 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-300"
            >
              Annuler
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

