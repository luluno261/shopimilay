import { useState, useEffect } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import styles from '../styles/Dashboard.module.css'
import KPICard from '../components/KPICard'
import Chart from '../components/Charts'
import { dashboardApi, DashboardStats } from '../lib/api'

export default function Dashboard() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('overview')
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (activeTab === 'overview') {
      loadStats()
    }
  }, [activeTab])

  const loadStats = async () => {
    try {
      setLoading(true)
      const data = await dashboardApi.getStats()
      setStats(data)
      setError(null)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Erreur lors du chargement des statistiques')
      console.error('Erreur:', err)
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: 'EUR',
    }).format(value)
  }

  return (
    <>
      <Head>
        <title>Dashboard - OmniSphere Admin</title>
        <meta name="description" content="Interface d'administration OmniSphere" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <div className={styles.container}>
        <header className={styles.header}>
          <h1>OmniSphere Admin Dashboard</h1>
          <nav className={styles.nav}>
            <button
              className={activeTab === 'overview' ? styles.active : ''}
              onClick={() => setActiveTab('overview')}
            >
              Vue d'ensemble
            </button>
            <button
              className={activeTab === 'orders' ? styles.active : ''}
              onClick={() => {
                setActiveTab('orders')
                router.push('/orders')
              }}
            >
              Commandes
            </button>
            <button
              className={activeTab === 'products' ? styles.active : ''}
              onClick={() => {
                setActiveTab('products')
                router.push('/products')
              }}
            >
              Produits
            </button>
            <button
              className={activeTab === 'customers' ? styles.active : ''}
              onClick={() => {
                setActiveTab('customers')
                router.push('/customers')
              }}
            >
              Clients
            </button>
            <button
              className={activeTab === 'growth' ? styles.active : ''}
              onClick={() => {
                setActiveTab('growth')
                router.push('/growth-command-center')
              }}
            >
              Growth Command Center
            </button>
            <button
              className={activeTab === 'store-builder' ? styles.active : ''}
              onClick={() => {
                setActiveTab('store-builder')
                router.push('/store-builder')
              }}
            >
              Store Builder
            </button>
          </nav>
        </header>

        <main className={styles.main}>
          {activeTab === 'overview' && (
            <div className={styles.content}>
              <h2 className="text-2xl font-bold mb-6">Vue d'ensemble</h2>
              
              {loading && <p className="text-gray-600">Chargement...</p>}
              {error && <p className="text-red-600">Erreur: {error}</p>}
              
              {stats && !loading && (
                <>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                    <KPICard
                      title="MRR (Monthly Recurring Revenue)"
                      value={formatCurrency(stats.mrr)}
                      subtitle="Ce mois"
                    />
                    <KPICard
                      title="GMV (Gross Merchandise Value)"
                      value={formatCurrency(stats.gmv)}
                      subtitle="Total"
                    />
                    <KPICard
                      title="Nombre de commandes"
                      value={stats.order_count}
                      subtitle="Total"
                    />
                    <KPICard
                      title="Panier moyen"
                      value={formatCurrency(stats.avg_cart)}
                      subtitle="Par commande"
                    />
                  </div>
                  
                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <Chart
                      title="Évolution des commandes"
                      data={[
                        { label: 'Jan', value: 10 },
                        { label: 'Fév', value: 15 },
                        { label: 'Mar', value: 12 },
                        { label: 'Avr', value: 20 },
                        { label: 'Mai', value: 18 },
                        { label: 'Juin', value: 25 },
                      ]}
                      type="bar"
                    />
                    <Chart
                      title="Évolution du chiffre d'affaires"
                      data={[
                        { label: 'Jan', value: 1000 },
                        { label: 'Fév', value: 1500 },
                        { label: 'Mar', value: 1200 },
                        { label: 'Avr', value: 2000 },
                        { label: 'Mai', value: 1800 },
                        { label: 'Juin', value: 2500 },
                      ]}
                      type="line"
                    />
                  </div>
                </>
              )}
            </div>
          )}

          {activeTab === 'orders' && (
            <div className={styles.content}>
              <h2>Commandes</h2>
              <p>Gestion des commandes - À implémenter</p>
            </div>
          )}

          {activeTab === 'products' && (
            <div className={styles.content}>
              <h2>Produits</h2>
              <p>Gestion des produits - À implémenter</p>
            </div>
          )}

          {activeTab === 'customers' && (
            <div className={styles.content}>
              <h2>Clients</h2>
              <p>Gestion des clients (CRM) - À implémenter</p>
            </div>
          )}
        </main>
      </div>
    </>
  )
}

