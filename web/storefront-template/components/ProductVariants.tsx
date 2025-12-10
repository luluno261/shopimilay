import React from 'react'

interface Variant {
  id: string
  name: string
  price: number
  stock: number
}

interface ProductVariantsProps {
  variants: Variant[]
  selectedVariant: string | null
  onSelectVariant: (variantId: string) => void
}

export default function ProductVariants({
  variants,
  selectedVariant,
  onSelectVariant,
}: ProductVariantsProps) {
  return (
    <div className="mb-6">
      <label className="block text-sm font-medium text-gray-700 mb-2">
        Variante
      </label>
      <div className="flex flex-wrap gap-2">
        {variants.map((variant) => (
          <button
            key={variant.id}
            onClick={() => onSelectVariant(variant.id)}
            className={`px-4 py-2 border rounded-lg ${
              selectedVariant === variant.id
                ? 'border-blue-600 bg-blue-50 text-blue-600'
                : 'border-gray-300 hover:border-gray-400'
            } ${variant.stock === 0 ? 'opacity-50 cursor-not-allowed' : ''}`}
            disabled={variant.stock === 0}
          >
            {variant.name}
            {variant.stock === 0 && ' (Épuisé)'}
          </button>
        ))}
      </div>
    </div>
  )
}

