import React, { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import api from '../lib/api'

interface Order {
  id: string
  status: string
  total_amount: number
  currency: string
  created_at: string
}

export default function OrderHistory() {
  const router = useRouter()
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadOrders()
  }, [])

  const loadOrders = async () => {
    try {
      setLoading(true)
      const response = await api.get('/account/orders')
      setOrders(response.data.orders || [])
    } catch (error) {
      console.error('Erreur lors du chargement des commandes:', error)
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (value: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: currency || 'EUR',
    }).format(value)
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'paid':
        return 'bg-green-100 text-green-800'
      case 'pending':
        return 'bg-yellow-100 text-yellow-800'
      case 'shipped':
        return 'bg-blue-100 text-blue-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  if (loading) {
    return <p>Chargement...</p>
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-xl font-semibold mb-4">Historique des commandes</h2>
      {orders.length === 0 ? (
        <p className="text-gray-500">Aucune commande</p>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div
              key={order.id}
              className="border border-gray-200 rounded p-4 flex justify-between items-center hover:bg-gray-50 cursor-pointer"
              onClick={() => router.push(`/account/orders/${order.id}`)}
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
                <span
                  className={`px-2 py-1 text-xs rounded-full ${getStatusColor(order.status)}`}
                >
                  {order.status}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

