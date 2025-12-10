import React, { useState, useEffect } from 'react'
import api from '../lib/api'
import CustomerTimeline from './CustomerTimeline'

interface Customer {
  id: string
  email: string
  name?: string
  total_orders: number
  total_spent: number
  ltv?: number
  last_order_date?: string
}

interface CustomerDetailProps {
  customer: Customer
  onClose: () => void
}

export default function CustomerDetail({ customer, onClose }: CustomerDetailProps) {
  const [orders, setOrders] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadCustomerData()
  }, [customer.id])

  const loadCustomerData = async () => {
    try {
      setLoading(true)
      // TODO: Implémenter l'endpoint /customers/:id dans le backend
      const response = await api.get(`/customers/${customer.id}`).catch(() => ({
        data: { orders: [], events: [] },
      }))
      setOrders(response.data.orders || [])
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: 'EUR',
    }).format(value)
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-8 max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">Détails du client</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="space-y-6">
          <div>
            <h3 className="text-lg font-semibold mb-4">Informations générales</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-600">Email</p>
                <p className="font-medium">{customer.email}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Total commandes</p>
                <p className="font-medium">{customer.total_orders}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Total dépensé</p>
                <p className="font-medium">{formatCurrency(customer.total_spent)}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">LTV (Lifetime Value)</p>
                <p className="font-medium">
                  {customer.ltv ? formatCurrency(customer.ltv) : 'N/A'}
                </p>
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-semibold mb-4">Historique des commandes</h3>
            {loading ? (
              <p>Chargement...</p>
            ) : orders.length === 0 ? (
              <p className="text-gray-500">Aucune commande</p>
            ) : (
              <div className="space-y-2">
                {orders.map((order) => (
                  <div
                    key={order.id}
                    className="border border-gray-200 rounded p-4 flex justify-between items-center"
                  >
                    <div>
                      <p className="font-medium">Commande #{order.id.substring(0, 8)}</p>
                      <p className="text-sm text-gray-500">
                        {new Date(order.created_at).toLocaleDateString('fr-FR')}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">
                        {formatCurrency(order.total_amount, order.currency)}
                      </p>
                      <p className="text-sm text-gray-500">{order.status}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          <CustomerTimeline customerId={customer.id} />
        </div>
      </div>
    </div>
  )
}

