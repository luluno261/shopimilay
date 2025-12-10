import { useState, useEffect } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import Link from 'next/link'
import api from '../lib/api'
import DynamicSectionRenderer from '../components/DynamicSectionRenderer'

interface Section {
  id: string
  type: string
  data: any
  order: number
}

interface Product {
  id: string
  name: string
  description: string
  price: number
  currency: string
  images: string[]
}

export default function Home() {
  const router = useRouter()
  const [sections, setSections] = useState<Section[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [useSections, setUseSections] = useState(false)

  useEffect(() => {
    loadConfig()
  }, [router.query.merchant_id])

  const loadConfig = async () => {
    try {
      setLoading(true)
      const merchantId = router.query.merchant_id as string || 'default'
      
      // Charger la configuration du storefront
      // Note: Pour le storefront public, on passe merchant_id en paramètre
      const configResponse = await api.get('/store-builder/config', {
        params: { merchant_id: merchantId },
        headers: {
          'X-Merchant-ID': merchantId,
        },
      }).catch(() => ({ data: { sections: [] } }))

      const loadedSections = configResponse.data.sections || []
      
      if (loadedSections.length > 0) {
        setSections(loadedSections)
        setUseSections(true)
      } else {
        // Fallback: charger les produits si pas de sections configurées
        loadProducts(merchantId)
        setUseSections(false)
      }
    } catch (error) {
      console.error('Erreur lors du chargement de la configuration:', error)
      // Fallback: charger les produits
      const merchantId = router.query.merchant_id as string || 'default'
      loadProducts(merchantId)
      setUseSections(false)
    } finally {
      setLoading(false)
    }
  }

  const loadProducts = async (merchantId: string) => {
    try {
      const response = await api.get('/products', {
        params: { merchant_id: merchantId, limit: 12 },
      })
      setProducts(response.data.products || [])
    } catch (error) {
      console.error('Erreur lors du chargement des produits:', error)
    }
  }

  const formatCurrency = (value: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: currency || 'EUR',
    }).format(value)
  }

  return (
    <>
      <Head>
        <title>Boutique - OmniSphere</title>
        <meta name="description" content="Boutique en ligne OmniSphere" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        {loading ? (
          <div className="text-center py-12">
            <p className="text-gray-600">Chargement...</p>
          </div>
        ) : useSections ? (
          <DynamicSectionRenderer sections={sections} />
        ) : (
          <div className="container mx-auto px-4 py-8">
            <div className="text-center mb-12">
              <h1 className="text-4xl font-bold mb-4">Bienvenue sur notre boutique</h1>
              <p className="text-xl text-gray-600">
                Découvrez notre sélection de produits
              </p>
            </div>

            {products.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-600">Aucun produit disponible</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                {products.map((product) => (
                  <Link
                    key={product.id}
                    href={`/products/${product.id}`}
                    className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
                  >
                    {product.images && product.images.length > 0 && (
                      <img
                        src={product.images[0]}
                        alt={product.name}
                        className="w-full h-48 object-cover"
                      />
                    )}
                    <div className="p-4">
                      <h3 className="text-lg font-semibold mb-2 line-clamp-2">
                        {product.name}
                      </h3>
                      <p className="text-sm text-gray-600 mb-3 line-clamp-2">
                        {product.description}
                      </p>
                      <p className="text-xl font-bold text-blue-600">
                        {formatCurrency(product.price, product.currency)}
                      </p>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        )}
      </main>
    </>
  )
}

