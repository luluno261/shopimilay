import { useState, useEffect } from 'react'
import Head from 'next/head'
import api from '../lib/api'

interface Automation {
  id: string
  name: string
  type: string
  is_active: boolean
  trigger: string
}

export default function AutomationPage() {
  const [automations, setAutomations] = useState<Automation[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadAutomations()
  }, [])

  const loadAutomations = async () => {
    try {
      setLoading(true)
      const response = await api.get('/automation').catch(() => ({
        data: { automations: [] },
      }))
      setAutomations(response.data.automations || [])
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  const automationTypes = [
    { id: 'cart-abandoned', name: 'Panier abandonné', description: 'Envoie des emails de rappel pour les paniers abandonnés' },
    { id: 'welcome', name: 'Email de bienvenue', description: 'Séquence de bienvenue pour nouveaux clients' },
    { id: 'win-back', name: 'Réactivation', description: 'Réactive les clients inactifs' },
  ]

  return (
    <>
      <Head>
        <title>Marketing Automation - OmniSphere Admin</title>
      </Head>

      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Marketing Automation</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {automationTypes.map((type) => (
            <div key={type.id} className="bg-white rounded-lg shadow-md p-6 border border-gray-200">
              <h3 className="text-xl font-semibold mb-2">{type.name}</h3>
              <p className="text-gray-600 mb-4">{type.description}</p>
              <button className="w-full bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
                Configurer
              </button>
            </div>
          ))}
        </div>
      </div>
    </>
  )
}

