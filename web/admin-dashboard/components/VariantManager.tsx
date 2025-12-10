import React, { useState } from 'react'

interface Variant {
  id?: string
  name: string
  sku: string
  price: number
  stock: number
}

interface VariantManagerProps {
  variants: Variant[]
  onChange: (variants: Variant[]) => void
}

export default function VariantManager({ variants, onChange }: VariantManagerProps) {
  const [localVariants, setLocalVariants] = useState<Variant[]>(variants)

  const handleAdd = () => {
    const newVariants = [...localVariants, { name: '', sku: '', price: 0, stock: 0 }]
    setLocalVariants(newVariants)
    onChange(newVariants)
  }

  const handleUpdate = (index: number, field: keyof Variant, value: any) => {
    const newVariants = [...localVariants]
    newVariants[index] = { ...newVariants[index], [field]: value }
    setLocalVariants(newVariants)
    onChange(newVariants)
  }

  const handleRemove = (index: number) => {
    const newVariants = localVariants.filter((_, i) => i !== index)
    setLocalVariants(newVariants)
    onChange(newVariants)
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold">Variantes de produit</h3>
        <button
          type="button"
          onClick={handleAdd}
          className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 text-sm"
        >
          + Ajouter une variante
        </button>
      </div>

      {localVariants.length === 0 ? (
        <p className="text-gray-500 text-sm">Aucune variante. Cliquez sur "Ajouter une variante" pour en créer.</p>
      ) : (
        <div className="space-y-3">
          {localVariants.map((variant, index) => (
            <div key={index} className="border border-gray-200 rounded-lg p-4">
              <div className="grid grid-cols-2 gap-4 mb-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nom *
                  </label>
                  <input
                    type="text"
                    required
                    value={variant.name}
                    onChange={(e) => handleUpdate(index, 'name', e.target.value)}
                    placeholder="Ex: Rouge, Taille M"
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    SKU *
                  </label>
                  <input
                    type="text"
                    required
                    value={variant.sku}
                    onChange={(e) => handleUpdate(index, 'sku', e.target.value)}
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4 mb-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Prix (optionnel)
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={variant.price}
                    onChange={(e) => handleUpdate(index, 'price', parseFloat(e.target.value) || 0)}
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Stock *
                  </label>
                  <input
                    type="number"
                    required
                    min="0"
                    value={variant.stock}
                    onChange={(e) => handleUpdate(index, 'stock', parseInt(e.target.value) || 0)}
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                  />
                </div>
              </div>
              <button
                type="button"
                onClick={() => handleRemove(index)}
                className="text-red-600 hover:text-red-800 text-sm"
              >
                Supprimer
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

