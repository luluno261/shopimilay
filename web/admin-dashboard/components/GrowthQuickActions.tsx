import { useRouter } from 'next/router'

export default function GrowthQuickActions() {
  const router = useRouter()

  const quickActions = [
    {
      label: 'Créer une automation',
      description: 'Configurer une nouvelle séquence marketing',
      icon: '🤖',
      onClick: () => router.push('/automation'),
      color: 'bg-blue-600 hover:bg-blue-700',
    },
    {
      label: 'Créer un pop-up',
      description: 'Créer un pop-up de capture de leads',
      icon: '📥',
      onClick: () => router.push('/capture-tools'),
      color: 'bg-purple-600 hover:bg-purple-700',
    },
    {
      label: 'Créer une audience',
      description: 'Segmenter vos clients',
      icon: '🎯',
      onClick: () => router.push('/segmentation'),
      color: 'bg-green-600 hover:bg-green-700',
    },
    {
      label: 'Connecter publicités',
      description: 'Connecter Facebook ou Google Ads',
      icon: '📢',
      onClick: () => router.push('/ads'),
      color: 'bg-orange-600 hover:bg-orange-700',
    },
  ]

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">Actions rapides</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {quickActions.map((action, index) => (
          <button
            key={index}
            onClick={action.onClick}
            className={`
              ${action.color} text-white rounded-lg p-6 text-left transition-colors
              transform hover:scale-105 transition-transform
            `}
          >
            <div className="text-3xl mb-2">{action.icon}</div>
            <h3 className="font-semibold mb-1">{action.label}</h3>
            <p className="text-sm opacity-90">{action.description}</p>
          </button>
        ))}
      </div>
    </div>
  )
}

