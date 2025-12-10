import React, { useState } from 'react'
import api from '../lib/api'
import OrderStatusUpdater from './OrderStatusUpdater'
import RefundModal from './RefundModal'

interface Order {
  id: string
  user_id: string
  merchant_id: string
  status: string
  total_amount: number
  currency: string
  payment_intent_id?: string
  shipping_address: any
  billing_address: any
  created_at: string
  updated_at: string
}

interface OrderDetailProps {
  order: Order
  onClose: () => void
  onStatusUpdate: (orderId: string, status: string) => void
}

export default function OrderDetail({ order, onClose, onStatusUpdate }: OrderDetailProps) {
  const [showRefundModal, setShowRefundModal] = useState(false)

  const formatCurrency = (value: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: currency || 'EUR',
    }).format(value)
  }

  const formatAddress = (address: any) => {
    if (!address) return 'N/A'
    if (typeof address === 'string') {
      try {
        address = JSON.parse(address)
      } catch {
        return address
      }
    }
    return `${address.street || ''}, ${address.city || ''} ${address.zip_code || ''}, ${address.country || ''}`
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-8 max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">Détails de la commande</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="space-y-6">
          <div>
            <h3 className="text-lg font-semibold mb-2">Informations générales</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-600">ID Commande</p>
                <p className="font-medium">{order.id}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Date</p>
                <p className="font-medium">
                  {new Date(order.created_at).toLocaleString('fr-FR')}
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Statut</p>
                <p className="font-medium">{order.status}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Montant total</p>
                <p className="font-medium">
                  {formatCurrency(order.total_amount, order.currency)}
                </p>
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-semibold mb-2">Adresse de livraison</h3>
            <p className="text-gray-700">{formatAddress(order.shipping_address)}</p>
          </div>

          <div>
            <h3 className="text-lg font-semibold mb-2">Adresse de facturation</h3>
            <p className="text-gray-700">{formatAddress(order.billing_address)}</p>
          </div>

          <div className="flex gap-4">
            <OrderStatusUpdater
              orderId={order.id}
              currentStatus={order.status}
              onUpdate={onStatusUpdate}
            />
            {order.status === 'paid' && order.payment_intent_id && (
              <button
                onClick={() => setShowRefundModal(true)}
                className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
              >
                Rembourser
              </button>
            )}
          </div>
        </div>

        {showRefundModal && (
          <RefundModal
            order={order}
            onClose={() => setShowRefundModal(false)}
            onSuccess={onClose}
          />
        )}
      </div>
    </div>
  )
}

