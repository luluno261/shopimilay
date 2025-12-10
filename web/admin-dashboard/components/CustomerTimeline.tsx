import React, { useState, useEffect } from 'react'
import api from '../lib/api'

interface CustomerTimelineProps {
  customerId: string
}

interface TimelineEvent {
  id: string
  type: string
  description: string
  created_at: string
}

export default function CustomerTimeline({ customerId }: CustomerTimelineProps) {
  const [events, setEvents] = useState<TimelineEvent[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadEvents()
  }, [customerId])

  const loadEvents = async () => {
    try {
      setLoading(true)
      // TODO: Implémenter l'endpoint /customers/:id/events dans le backend
      const response = await api.get(`/customers/${customerId}/events`).catch(() => ({
        data: { events: [] },
      }))
      setEvents(response.data.events || [])
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'order_created':
        return '🛒'
      case 'order_paid':
        return '💳'
      case 'order_shipped':
        return '📦'
      case 'email_sent':
        return '📧'
      case 'cart_abandoned':
        return '🛑'
      default:
        return '📝'
    }
  }

  return (
    <div>
      <h3 className="text-lg font-semibold mb-4">Timeline</h3>
      {loading ? (
        <p>Chargement...</p>
      ) : events.length === 0 ? (
        <p className="text-gray-500">Aucun événement</p>
      ) : (
        <div className="space-y-4">
          {events.map((event) => (
            <div key={event.id} className="flex items-start gap-4">
              <div className="text-2xl">{getEventIcon(event.type)}</div>
              <div className="flex-1">
                <p className="font-medium">{event.description}</p>
                <p className="text-sm text-gray-500">
                  {new Date(event.created_at).toLocaleString('fr-FR')}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

