import { useState, useEffect } from 'react'
import Head from 'next/head'
import api from '../lib/api'

interface AdAccount {
  id: string
  platform: 'facebook' | 'google'
  account_id: string
  is_connected: boolean
}

export default function AdsPage() {
  const [accounts, setAccounts] = useState<AdAccount[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadAccounts()
  }, [])

  const loadAccounts = async () => {
    try {
      setLoading(true)
      // TODO: Implémenter l'endpoint /ads/accounts dans le backend
      const response = await api.get('/ads/accounts').catch(() => ({
        data: { accounts: [] },
      }))
      setAccounts(response.data.accounts || [])
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleConnect = async (platform: 'facebook' | 'google') => {
    try {
      // TODO: Implémenter l'OAuth flow pour Facebook/Google Ads
      alert(`Connexion ${platform} - À implémenter`)
    } catch (error) {
      console.error('Erreur:', error)
    }
  }

  const handleSyncAudience = async (accountId: string) => {
    try {
      await api.post(`/ads/accounts/${accountId}/sync`)
      alert('Synchronisation démarrée')
    } catch (error: any) {
      alert(error.response?.data?.error || 'Erreur lors de la synchronisation')
    }
  }

  return (
    <>
      <Head>
        <title>Gestion des Publicités - OmniSphere Admin</title>
      </Head>

      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Gestion des Publicités</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-4">Facebook Ads</h2>
            <p className="text-gray-600 mb-4">
              Connectez votre compte Facebook Ads pour synchroniser vos audiences
            </p>
            {accounts.find((a) => a.platform === 'facebook' && a.is_connected) ? (
              <div className="space-y-4">
                <p className="text-green-600">✓ Compte connecté</p>
                <button
                  onClick={() => handleSyncAudience(accounts.find((a) => a.platform === 'facebook')!.id)}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
                >
                  Synchroniser l'audience
                </button>
              </div>
            ) : (
              <button
                onClick={() => handleConnect('facebook')}
                className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
              >
                Connecter Facebook Ads
              </button>
            )}
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-4">Google Ads</h2>
            <p className="text-gray-600 mb-4">
              Connectez votre compte Google Ads pour synchroniser vos audiences
            </p>
            {accounts.find((a) => a.platform === 'google' && a.is_connected) ? (
              <div className="space-y-4">
                <p className="text-green-600">✓ Compte connecté</p>
                <button
                  onClick={() => handleSyncAudience(accounts.find((a) => a.platform === 'google')!.id)}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
                >
                  Synchroniser l'audience
                </button>
              </div>
            ) : (
              <button
                onClick={() => handleConnect('google')}
                className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
              >
                Connecter Google Ads
              </button>
            )}
          </div>
        </div>
      </div>
    </>
  )
}

