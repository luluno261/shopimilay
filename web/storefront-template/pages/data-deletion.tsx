import { useState } from 'react'
import Head from 'next/head'
import api from '../lib/api'

export default function DataDeletionPage() {
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [submitted, setSubmitted] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      await api.post('/auth/delete-account', { email })
      setSubmitted(true)
    } catch (error: any) {
      alert(error.response?.data?.error || 'Erreur lors de la demande de suppression')
    } finally {
      setLoading(false)
    }
  }

  if (submitted) {
    return (
      <>
        <Head>
          <title>Demande de suppression - OmniSphere</title>
        </Head>
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          <div className="bg-green-50 border border-green-200 rounded-lg p-6">
            <h1 className="text-2xl font-bold mb-4 text-green-800">
              Demande reçue
            </h1>
            <p className="text-green-700">
              Votre demande de suppression de données a été reçue. Nous traiterons votre
              demande dans les 30 jours conformément au RGPD.
            </p>
          </div>
        </div>
      </>
    )
  }

  return (
    <>
      <Head>
        <title>Suppression de données - OmniSphere</title>
      </Head>

      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <h1 className="text-3xl font-bold mb-6">Demande de suppression de données</h1>

        <div className="bg-white rounded-lg shadow-md p-6">
          <p className="mb-6 text-gray-700">
            Conformément au RGPD, vous pouvez demander la suppression de vos données
            personnelles. Cette action est irréversible.
          </p>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Adresse email *
              </label>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="votre@email.com"
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-red-600 text-white py-3 px-6 rounded-lg hover:bg-red-700 disabled:opacity-50"
            >
              {loading ? 'Traitement...' : 'Demander la suppression'}
            </button>
          </form>
        </div>
      </div>
    </>
  )
}

