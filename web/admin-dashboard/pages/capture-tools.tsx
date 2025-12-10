import { useState, useEffect } from 'react'
import Head from 'next/head'
import api from '../lib/api'

interface Popup {
  id: string
  name: string
  type: 'exit-intent' | 'time-based' | 'scroll'
  is_active: boolean
  content: string
}

export default function CaptureToolsPage() {
  const [popups, setPopups] = useState<Popup[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    loadPopups()
  }, [])

  const loadPopups = async () => {
    try {
      setLoading(true)
      // TODO: Implémenter l'endpoint /capture/popups dans le backend
      const response = await api.get('/capture/popups').catch(() => ({
        data: { popups: [] },
      }))
      setPopups(response.data.popups || [])
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Head>
        <title>Outils de Capture - OmniSphere Admin</title>
      </Head>

      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold">Outils de Capture de Leads</h1>
          <button
            onClick={() => setShowForm(true)}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
          >
            + Nouveau pop-up
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-4">Pop-ups</h2>
            <p className="text-gray-600 mb-4">
              Créez des pop-ups pour capturer des leads (sortie, intention d'achat, etc.)
            </p>
            {loading ? (
              <p>Chargement...</p>
            ) : popups.length === 0 ? (
              <p className="text-gray-500">Aucun pop-up configuré</p>
            ) : (
              <div className="space-y-2">
                {popups.map((popup) => (
                  <div key={popup.id} className="border border-gray-200 rounded p-3">
                    <div className="flex justify-between items-center">
                      <div>
                        <p className="font-medium">{popup.name}</p>
                        <p className="text-sm text-gray-500">{popup.type}</p>
                      </div>
                      <span
                        className={`px-2 py-1 text-xs rounded-full ${
                          popup.is_active
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {popup.is_active ? 'Actif' : 'Inactif'}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-4">Barres d'annonce</h2>
            <p className="text-gray-600 mb-4">
              Créez des barres d'annonce pour promouvoir des offres spéciales
            </p>
            <button className="bg-gray-200 text-gray-800 px-4 py-2 rounded-lg hover:bg-gray-300">
              + Nouvelle barre
            </button>
          </div>
        </div>
      </div>
    </>
  )
}

