import React, { useState } from 'react'

interface Address {
  street: string
  city: string
  state?: string
  zip_code: string
  country: string
}

interface CheckoutFormProps {
  onSubmit: (data: { shippingAddress: Address; billingAddress: Address }) => void
}

export default function CheckoutForm({ onSubmit }: CheckoutFormProps) {
  const [shippingAddress, setShippingAddress] = useState<Address>({
    street: '',
    city: '',
    zip_code: '',
    country: 'FR',
  })
  const [billingAddress, setBillingAddress] = useState<Address>({
    street: '',
    city: '',
    zip_code: '',
    country: 'FR',
  })
  const [useSameAddress, setUseSameAddress] = useState(true)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    const finalBillingAddress = useSameAddress ? shippingAddress : billingAddress
    onSubmit({ shippingAddress, billingAddress: finalBillingAddress })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold mb-4">Adresse de livraison</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Rue *
            </label>
            <input
              type="text"
              required
              value={shippingAddress.street}
              onChange={(e) =>
                setShippingAddress({ ...shippingAddress, street: e.target.value })
              }
              className="w-full border border-gray-300 rounded-md px-3 py-2"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Code postal *
              </label>
              <input
                type="text"
                required
                value={shippingAddress.zip_code}
                onChange={(e) =>
                  setShippingAddress({ ...shippingAddress, zip_code: e.target.value })
                }
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Ville *
              </label>
              <input
                type="text"
                required
                value={shippingAddress.city}
                onChange={(e) =>
                  setShippingAddress({ ...shippingAddress, city: e.target.value })
                }
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Pays *
            </label>
            <select
              required
              value={shippingAddress.country}
              onChange={(e) =>
                setShippingAddress({ ...shippingAddress, country: e.target.value })
              }
              className="w-full border border-gray-300 rounded-md px-3 py-2"
            >
              <option value="FR">France</option>
              <option value="BE">Belgique</option>
              <option value="CH">Suisse</option>
            </select>
          </div>
        </div>
      </div>

      <div>
        <label className="flex items-center">
          <input
            type="checkbox"
            checked={useSameAddress}
            onChange={(e) => setUseSameAddress(e.target.checked)}
            className="mr-2"
          />
          Utiliser la même adresse pour la facturation
        </label>
      </div>

      {!useSameAddress && (
        <div>
          <h2 className="text-xl font-semibold mb-4">Adresse de facturation</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Rue *
              </label>
              <input
                type="text"
                required
                value={billingAddress.street}
                onChange={(e) =>
                  setBillingAddress({ ...billingAddress, street: e.target.value })
                }
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Code postal *
                </label>
                <input
                  type="text"
                  required
                  value={billingAddress.zip_code}
                  onChange={(e) =>
                    setBillingAddress({ ...billingAddress, zip_code: e.target.value })
                  }
                  className="w-full border border-gray-300 rounded-md px-3 py-2"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Ville *
                </label>
                <input
                  type="text"
                  required
                  value={billingAddress.city}
                  onChange={(e) =>
                    setBillingAddress({ ...billingAddress, city: e.target.value })
                  }
                  className="w-full border border-gray-300 rounded-md px-3 py-2"
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Pays *
              </label>
              <select
                required
                value={billingAddress.country}
                onChange={(e) =>
                  setBillingAddress({ ...billingAddress, country: e.target.value })
                }
                className="w-full border border-gray-300 rounded-md px-3 py-2"
              >
                <option value="FR">France</option>
                <option value="BE">Belgique</option>
                <option value="CH">Suisse</option>
              </select>
            </div>
          </div>
        </div>
      )}

      <button
        type="submit"
        className="w-full bg-blue-600 text-white py-3 px-6 rounded-lg hover:bg-blue-700 text-lg font-semibold"
      >
        Procéder au paiement
      </button>
    </form>
  )
}

