interface FeaturesSectionProps {
  data: {
    title?: string
    features?: Array<{
      icon?: string
      title: string
      description: string
    }>
  }
}

const defaultFeatures = [
  {
    icon: '🚚',
    title: 'Livraison rapide',
    description: 'Livraison en 24-48h partout en France',
  },
  {
    icon: '🔒',
    title: 'Paiement sécurisé',
    description: 'Transactions 100% sécurisées',
  },
  {
    icon: '↩️',
    title: 'Retours faciles',
    description: 'Retours gratuits sous 30 jours',
  },
  {
    icon: '💬',
    title: 'Support client',
    description: 'Assistance disponible 7j/7',
  },
]

export default function FeaturesSection({ data }: FeaturesSectionProps) {
  const features = data.features || defaultFeatures

  return (
    <div className="py-12 bg-gray-50">
      <div className="container mx-auto px-4">
        {data.title && <h2 className="text-3xl font-bold mb-8 text-center">{data.title}</h2>}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {features.map((feature, index) => (
            <div key={index} className="bg-white rounded-lg p-6 text-center border border-gray-200">
              {feature.icon && <div className="text-4xl mb-4">{feature.icon}</div>}
              <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
              <p className="text-gray-600">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

