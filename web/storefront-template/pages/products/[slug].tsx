import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import Head from 'next/head'
import ProductDetail from '../../components/ProductDetail'
import ProductVariants from '../../components/ProductVariants'
import Breadcrumbs from '../../components/Breadcrumbs'
import api from '../../lib/api'

interface Product {
  id: string
  name: string
  description: string
  price: number
  currency: string
  images: string[]
  variants?: Array<{
    id: string
    name: string
    price: number
    stock: number
  }>
}

export default function ProductPage() {
  const router = useRouter()
  const { slug } = router.query
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [selectedVariant, setSelectedVariant] = useState<string | null>(null)

  useEffect(() => {
    if (slug) {
      loadProduct(slug as string)
    }
  }, [slug])

  const loadProduct = async (productSlug: string) => {
    try {
      setLoading(true)
      const response = await api.get(`/products/${productSlug}`)
      setProduct(response.data)
    } catch (error) {
      console.error('Erreur lors du chargement du produit:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAddToCart = async () => {
    if (!product) return

    try {
      await api.post('/cart/items', {
        product_id: product.id,
        variant_id: selectedVariant,
        quantity: 1,
      })
      alert('Produit ajouté au panier')
    } catch (error: any) {
      console.error('Erreur:', error)
      alert(error.response?.data?.error || 'Erreur lors de l\'ajout au panier')
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p>Chargement...</p>
      </div>
    )
  }

  if (!product) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p>Produit introuvable</p>
      </div>
    )
  }

  return (
    <>
      <Head>
        <title>{product.name} - OmniSphere</title>
        <meta name="description" content={product.description} />
      </Head>

      <div className="container mx-auto px-4 py-8">
        <Breadcrumbs
          items={[
            { label: 'Accueil', href: '/' },
            { label: 'Produits', href: '/products' },
            { label: product.name, href: null },
          ]}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-8">
          <div>
            {product.images && product.images.length > 0 && (
              <img
                src={product.images[0]}
                alt={product.name}
                className="w-full rounded-lg"
              />
            )}
          </div>

          <div>
            <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
            <p className="text-2xl font-semibold text-blue-600 mb-6">
              {new Intl.NumberFormat('fr-FR', {
                style: 'currency',
                currency: product.currency || 'EUR',
              }).format(product.price)}
            </p>

            <div className="mb-6">
              <p className="text-gray-700 whitespace-pre-line">{product.description}</p>
            </div>

            {product.variants && product.variants.length > 0 && (
              <ProductVariants
                variants={product.variants}
                selectedVariant={selectedVariant}
                onSelectVariant={setSelectedVariant}
              />
            )}

            <button
              onClick={handleAddToCart}
              className="w-full bg-blue-600 text-white py-3 px-6 rounded-lg hover:bg-blue-700 text-lg font-semibold"
            >
              Ajouter au panier
            </button>
          </div>
        </div>
      </div>
    </>
  )
}

