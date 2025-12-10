import { useState } from 'react'
import api from '../lib/api'

export default function SegmentationTab() {
  const [segments, setSegments] = useState<any[]>([])
  const [showForm, setShowForm] = useState(false)

  const handleCreateSegment = async (data: any) => {
    try {
      await api.post('/segmentation', data)
      setShowForm(false)
      // Recharger les segments
    } catch (error: any) {
      alert(error.response?.data?.error || 'Erreur lors de la création')
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Segmentation</h2>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          + Nouvelle audience
        </button>
      </div>

      <div className="bg-gray-50 rounded-lg shadow-md p-6">
        <h3 className="text-xl font-semibold mb-4">Créer une audience</h3>
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
  )
}

