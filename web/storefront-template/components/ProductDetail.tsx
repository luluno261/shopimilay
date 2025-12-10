import React from 'react'

interface ProductDetailProps {
  product: {
    id: string
    name: string
    description: string
    price: number
    currency: string
    images: string[]
  }
}

export default function ProductDetail({ product }: ProductDetailProps) {
  return (
    <div className="product-detail">
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
    </div>
  )
}

