import React from 'react'

interface ShippingOptionsProps {
  selected: string
  onSelect: (option: string) => void
}

const shippingOptions = [
  { id: 'standard', label: 'Livraison standard', price: 5.99, days: '3-5 jours' },
  { id: 'express', label: 'Livraison express', price: 12.99, days: '1-2 jours' },
  { id: 'free', label: 'Livraison gratuite', price: 0, days: '5-7 jours' },
]

export default function ShippingOptions({ selected, onSelect }: ShippingOptionsProps) {
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: 'EUR',
    }).format(value)
  }

  return (
    <div className="bg-gray-50 rounded-lg p-4">
      <h3 className="text-lg font-semibold mb-4">Options de livraison</h3>
      <div className="space-y-2">
        {shippingOptions.map((option) => (
          <label
            key={option.id}
            className={`flex items-center p-3 border rounded-lg cursor-pointer ${
              selected === option.id
                ? 'border-blue-600 bg-blue-50'
                : 'border-gray-300 hover:border-gray-400'
            }`}
          >
            <input
              type="radio"
              name="shipping"
              value={option.id}
              checked={selected === option.id}
              onChange={(e) => onSelect(e.target.value)}
              className="mr-3"
            />
            <div className="flex-1">
              <div className="font-medium">{option.label}</div>
              <div className="text-sm text-gray-600">{option.days}</div>
            </div>
            <div className="font-semibold">
              {option.price === 0 ? 'Gratuit' : formatCurrency(option.price)}
            </div>
          </label>
        ))}
      </div>
    </div>
  )
}

