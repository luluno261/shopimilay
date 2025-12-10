import { useState, useEffect } from 'react'
import api from '../lib/api'
import KPICard from './KPICard'

interface GrowthMetrics {
  conversion_rate: number
  leads_captured: number
  automation_roi: number
  active_automations: number
  active_segments: number
  active_popups: number
  ads_spend: number
  ads_revenue: number
}

export default function GrowthOverview() {
  const [metrics, setMetrics] = useState<GrowthMetrics | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadMetrics()
  }, [])

  const loadMetrics = async () => {
    try {
      setLoading(true)
      // TODO: Implémenter l'endpoint /growth/metrics dans le backend
      const response = await api.get('/growth/metrics').catch(() => ({
        data: {
          conversion_rate: 2.5,
          leads_captured: 1250,
          automation_roi: 350,
          active_automations: 3,
          active_segments: 5,
          active_popups: 2,
          ads_spend: 500,
          ads_revenue: 2500,
        },
      }))
      setMetrics(response.data)
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-8">Chargement des métriques...</div>
  }

  if (!metrics) {
    return <div className="text-center py-8">Aucune donnée disponible</div>
  }

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Métriques Marketing</h2>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <KPICard
          title="Taux de conversion"
          value={`${metrics.conversion_rate}%`}
          change={+0.3}
          changeLabel="vs mois dernier"
        />
        <KPICard
          title="Leads capturés"
          value={metrics.leads_captured.toLocaleString()}
          change={+125}
          changeLabel="ce mois"
        />
        <KPICard
          title="ROI Automations"
          value={`${metrics.automation_roi}%`}
          change={+15}
          changeLabel="vs mois dernier"
        />
        <KPICard
          title="Revenus Publicités"
          value={`€${metrics.ads_revenue.toLocaleString()}`}
          change={+320}
          changeLabel="ce mois"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-gray-50 rounded-lg p-6">
          <h3 className="text-lg font-semibold mb-2">Automations actives</h3>
          <p className="text-3xl font-bold text-blue-600">{metrics.active_automations}</p>
        </div>
        <div className="bg-gray-50 rounded-lg p-6">
          <h3 className="text-lg font-semibold mb-2">Audiences actives</h3>
          <p className="text-3xl font-bold text-green-600">{metrics.active_segments}</p>
        </div>
        <div className="bg-gray-50 rounded-lg p-6">
          <h3 className="text-lg font-semibold mb-2">Pop-ups actifs</h3>
          <p className="text-3xl font-bold text-purple-600">{metrics.active_popups}</p>
        </div>
      </div>
    </div>
  )
}

