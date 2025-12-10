import React from 'react'

interface Cart {
  id: string
  items: Array<{
    product_id: string
    quantity: number
    price: number
    subtotal: number
  }>
  total: number
  currency: string
}

interface CartSummaryProps {
  cart: Cart
  shippingOption: string
}

const shippingCosts: Record<string, number> = {
  standard: 5.99,
  express: 12.99,
  free: 0,
}

export default function CartSummary({ cart, shippingOption }: CartSummaryProps) {
  const shippingCost = shippingCosts[shippingOption] || 0
  const total = cart.total + shippingCost

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: cart.currency || 'EUR',
    }).format(value)
  }

  return (
    <div className="bg-gray-50 rounded-lg p-6">
      <h2 className="text-xl font-semibold mb-4">Résumé de la commande</h2>

      <div className="space-y-2 mb-4">
        {cart.items.map((item, index) => (
          <div key={index} className="flex justify-between text-sm">
            <span>
              Produit × {item.quantity}
            </span>
            <span>{formatCurrency(item.subtotal)}</span>
          </div>
        ))}
      </div>

      <div className="border-t border-gray-300 pt-4 space-y-2">
        <div className="flex justify-between">
          <span>Sous-total</span>
          <span>{formatCurrency(cart.total)}</span>
        </div>
        <div className="flex justify-between">
          <span>Livraison ({shippingOption})</span>
          <span>{formatCurrency(shippingCost)}</span>
        </div>
        <div className="flex justify-between font-bold text-lg pt-2 border-t border-gray-300">
          <span>Total</span>
          <span>{formatCurrency(total)}</span>
        </div>
      </div>
    </div>
  )
}

