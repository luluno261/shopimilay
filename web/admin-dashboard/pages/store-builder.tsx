import { useState, useEffect } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
// Note: Les composants SectionBuilder et ThemeCustomizer doivent être importés depuis storefront-template
// Pour l'instant, on les importe directement depuis le chemin relatif
// TODO: Créer des composants partagés ou copier les composants dans admin-dashboard
import dynamic from 'next/dynamic'

const SectionBuilder = dynamic(() => import('../../storefront-template/components/builder/SectionBuilder'), { ssr: false })
const ThemeCustomizer = dynamic(() => import('../../storefront-template/components/builder/ThemeCustomizer'), { ssr: false })
import api from '../lib/api'

interface Section {
  id: string
  type: string
  data: any
  order: number
}

interface Theme {
  colors: {
    primary: string
    secondary: string
    background: string
    text: string
    link: string
    button: string
  }
  typography: {
    fontFamily: string
    baseSize: number
    headingSizes: {
      h1: number
      h2: number
      h3: number
      h4: number
      h5: number
      h6: number
    }
    lineHeight: number
    fontWeight: number
  }
  layout: {
    containerWidth: number
    padding: number
    margin: number
    borderRadius: number
    shadow: string
  }
}

export default function StoreBuilderPage() {
  const router = useRouter()
  const [sections, setSections] = useState<Section[]>([])
  const [theme, setTheme] = useState<Theme | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'sections' | 'theme'>('sections')

  useEffect(() => {
    loadConfig()
  }, [])

  const loadConfig = async () => {
    try {
      setLoading(true)
      // TODO: Implémenter les endpoints backend
      const [sectionsResponse, themeResponse] = await Promise.all([
        api.get('/store-builder/config').catch(() => ({ data: { sections: [] } })),
        api.get('/store-builder/theme').catch(() => ({ data: null })),
      ])
      setSections(sectionsResponse.data.sections || [])
      setTheme(themeResponse.data || null)
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSectionsChange = async (newSections: Section[]) => {
    setSections(newSections)
    try {
      await api.post('/store-builder/config', {
        sections: newSections,
      })
    } catch (error) {
      console.error('Erreur lors de la sauvegarde:', error)
    }
  }

  const handleThemeChange = async (newTheme: Theme) => {
    setTheme(newTheme)
    try {
      await api.post('/store-builder/theme', {
        theme: newTheme,
      })
    } catch (error) {
      console.error('Erreur lors de la sauvegarde du thème:', error)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p>Chargement...</p>
      </div>
    )
  }

  return (
    <>
      <Head>
        <title>Store Builder - OmniSphere Admin</title>
      </Head>

      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Store Builder</h1>
        <p className="text-gray-600 mb-8">
          Personnalisez votre boutique avec des sections modulaires et un thème personnalisé
        </p>

        {/* Tabs */}
        <div className="border-b border-gray-200 mb-6">
          <nav className="flex space-x-1">
            <button
              onClick={() => setActiveTab('sections')}
              className={`
                px-4 py-2 text-sm font-medium border-b-2 transition-colors
                ${
                  activeTab === 'sections'
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }
              `}
            >
              Sections
            </button>
            <button
              onClick={() => setActiveTab('theme')}
              className={`
                px-4 py-2 text-sm font-medium border-b-2 transition-colors
                ${
                  activeTab === 'theme'
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }
              `}
            >
              Thème
            </button>
          </nav>
        </div>

        {/* Tab Content */}
        <div>
          {activeTab === 'sections' && (
            <SectionBuilder sections={sections} onSectionsChange={handleSectionsChange} />
          )}
          {activeTab === 'theme' && (
            <ThemeCustomizer theme={theme || undefined} onThemeChange={handleThemeChange} />
          )}
        </div>
      </div>
    </>
  )
}

