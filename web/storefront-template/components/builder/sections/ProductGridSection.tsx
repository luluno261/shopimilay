import React, { useState, useEffect } from 'react'
import api from '../../../lib/api'

interface Product {
  id: string
  name: string
  price: number
  currency: string
  images: string[]
}

interface ProductGridSectionProps {
  data: {
    title?: string
    limit?: number
    category?: string
  }
}

export default function ProductGridSection({ data }: ProductGridSectionProps) {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadProducts()
  }, [])

  const loadProducts = async () => {
    try {
      setLoading(true)
      const response = await api.get('/products', {
        params: { limit: data.limit || 8 },
      })
      setProducts(response.data.products || [])
    } catch (error) {
      console.error('Erreur:', error)
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (value: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: currency || 'EUR',
    }).format(value)
  }

  return (
    <div className="py-12 bg-gray-50">
      <div className="container mx-auto px-4">
        {data.title && <h2 className="text-3xl font-bold mb-8 text-center">{data.title}</h2>}
        {loading ? (
          <div className="text-center py-12">Chargement des produits...</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {products.map((product) => (
              <div key={product.id} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
                {product.images && product.images.length > 0 && (
                  <img
                    src={product.images[0]}
                    alt={product.name}
                    className="w-full h-48 object-cover"
                  />
                )}
                <div className="p-4">
                  <h3 className="text-lg font-semibold mb-2">{product.name}</h3>
                  <p className="text-xl font-bold text-blue-600">
                    {formatCurrency(product.price, product.currency)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

