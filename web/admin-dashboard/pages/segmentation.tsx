import { useState } from 'react'
import Head from 'next/head'
import api from '../lib/api'

export default function SegmentationPage() {
  const [segments, setSegments] = useState<any[]>([])
  const [showForm, setShowForm] = useState(false)

  const handleCreateSegment = async (data: any) => {
    try {
      await api.post('/segmentation', data)
      // Recharger les segments
      setShowForm(false)
    } catch (error: any) {
      alert(error.response?.data?.error || 'Erreur lors de la création')
    }
  }

  return (
    <>
      <Head>
        <title>Segmentation - OmniSphere Admin</title>
      </Head>

      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold">Segmentation</h1>
          <button
            onClick={() => setShowForm(true)}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
          >
            + Nouvelle audience
          </button>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold mb-4">Créer une audience</h2>
          <p className="text-gray-600 mb-6">
            Créez des audiences basées sur les comportements et événements clients.
          </p>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Nom de l'audience
              </label>
              <input
                type="text"
                placeholder="Ex: Clients VIP"
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Critères de segmentation
              </label>
              <select className="w-full border border-gray-300 rounded-md px-3 py-2">
                <option>Clients qui ont acheté produit X</option>
                <option>N'a pas visité depuis 30 jours</option>
                <option>LTV &gt; €100</option>
                <option>Panier moyen &gt; €50</option>
              </select>
            </div>

            <button className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700">
              Créer l'audience
            </button>
          </div>
        </div>
      </div>
    </>
  )
}

