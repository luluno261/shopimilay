import { useState, useEffect } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import CheckoutForm from '../components/CheckoutForm'
import CartSummary from '../components/CartSummary'
import DiscountCode from '../components/DiscountCode'
import ShippingOptions from '../components/ShippingOptions'
import api from '../lib/api'

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

export default function CheckoutPage() {
  const router = useRouter()
  const [cart, setCart] = useState<Cart | null>(null)
  const [loading, setLoading] = useState(true)
  const [discountCode, setDiscountCode] = useState('')
  const [discountApplied, setDiscountApplied] = useState(false)
  const [shippingOption, setShippingOption] = useState('standard')

  useEffect(() => {
    loadCart()
  }, [])

  const loadCart = async () => {
    try {
      setLoading(true)
      const response = await api.get('/cart')
      setCart(response.data)
    } catch (error) {
      console.error('Erreur lors du chargement du panier:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleApplyDiscount = async () => {
    if (!discountCode.trim()) return

    try {
      // TODO: Implémenter l'application du code de réduction
      setDiscountApplied(true)
      alert('Code de réduction appliqué')
    } catch (error: any) {
      alert(error.response?.data?.error || 'Code de réduction invalide')
    }
  }

  const handleCheckout = async (formData: any) => {
    try {
      const response = await api.post('/checkout', {
        cart_id: cart?.id,
        payment_method: 'card',
        shipping_address: formData.shippingAddress,
        billing_address: formData.billingAddress,
        discount_code: discountCode || undefined,
        shipping_option: shippingOption,
      })

      if (response.data.client_secret) {
        // Rediriger vers la page de paiement Stripe
        router.push(`/payment?client_secret=${response.data.client_secret}&order_id=${response.data.order_id}`)
      } else {
        router.push(`/order-confirmation?order_id=${response.data.order_id}`)
      }
    } catch (error: any) {
      console.error('Erreur:', error)
      alert(error.response?.data?.error || 'Erreur lors du checkout')
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p>Chargement...</p>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p>Votre panier est vide</p>
        <button
          onClick={() => router.push('/')}
          className="mt-4 bg-blue-600 text-white px-4 py-2 rounded-lg"
        >
          Retour à l'accueil
        </button>
      </div>
    )
  }

  return (
    <>
      <Head>
        <title>Checkout - OmniSphere</title>
      </Head>

      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">Finaliser la commande</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2">
            <CheckoutForm onSubmit={handleCheckout} />
            
            <div className="mt-8">
              <ShippingOptions
                selected={shippingOption}
                onSelect={setShippingOption}
              />
            </div>

            <div className="mt-8">
              <DiscountCode
                value={discountCode}
                onChange={setDiscountCode}
                onApply={handleApplyDiscount}
                applied={discountApplied}
              />
            </div>
          </div>

          <div>
            <CartSummary cart={cart} shippingOption={shippingOption} />
          </div>
        </div>
      </div>
    </>
  )
}

