import { useState } from 'react'
import Head from 'next/head'
import api from '../lib/api'

export default function MigrationPage() {
  const [source, setSource] = useState('csv')
  const [file, setFile] = useState<File | null>(null)
  const [loading, setLoading] = useState(false)
  const [migrationId, setMigrationId] = useState<string | null>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0])
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!file) {
      alert('Veuillez sélectionner un fichier')
      return
    }

    setLoading(true)
    try {
      // Convertir le fichier en base64
      const reader = new FileReader()
      reader.onload = async (event) => {
        const fileData = (event.target?.result as string).split(',')[1] // Enlever le préfixe data:...

        const response = await api.post('/migration/import', {
          source,
          file_data: fileData,
        })

        setMigrationId(response.data.migration_id)
        alert('Import démarré ! ID: ' + response.data.migration_id)
      }
      reader.readAsDataURL(file)
    } catch (error: any) {
      alert(error.response?.data?.error || 'Erreur lors de l\'import')
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Head>
        <title>Migration de données - OmniSphere Admin</title>
      </Head>

      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Migration de données</h1>

        <div className="bg-white rounded-lg shadow-md p-8 max-w-2xl">
          <h2 className="text-2xl font-semibold mb-6">Importer des données</h2>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Source *
              </label>
              <select
                value={source}
                onChange={(e) => setSource(e.target.value)}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              >
                <option value="csv">CSV générique</option>
                <option value="json">JSON générique</option>
                <option value="shopify">Shopify</option>
                <option value="woocommerce">WooCommerce</option>
              </select>
            </div>

            {source === 'shopify' && (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Store URL
                  </label>
                  <input
                    type="url"
                    placeholder="https://votre-boutique.myshopify.com"
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    API Key
                  </label>
                  <input
                    type="text"
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
              </div>
            )}

            {source === 'woocommerce' && (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Store URL
                  </label>
                  <input
                    type="url"
                    placeholder="https://votre-boutique.com"
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Consumer Key
                  </label>
                  <input
                    type="text"
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Consumer Secret
                  </label>
                  <input
                    type="text"
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
              </div>
            )}

            {(source === 'csv' || source === 'json') && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Fichier *
                </label>
                <input
                  type="file"
                  accept={source === 'csv' ? '.csv' : '.json'}
                  onChange={handleFileChange}
                  className="w-full border border-gray-300 rounded-md px-3 py-2"
                />
                <p className="text-sm text-gray-500 mt-1">
                  Format attendu: {source === 'csv' ? 'CSV avec colonnes: name, description, price, sku' : 'JSON avec tableau de produits'}
                </p>
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 text-white py-3 px-6 rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? 'Import en cours...' : 'Démarrer l\'import'}
            </button>
          </form>

          {migrationId && (
            <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
              <p className="text-green-800">
                Import démarré ! ID: <strong>{migrationId}</strong>
              </p>
              <p className="text-sm text-green-600 mt-2">
                Vous pouvez suivre la progression avec cet ID.
              </p>
            </div>
          )}
        </div>
      </div>
    </>
  )
}

