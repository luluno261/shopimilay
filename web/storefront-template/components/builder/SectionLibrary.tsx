interface SectionLibraryProps {
  onSelect: (type: string) => void
  onClose: () => void
}

const sectionTypes = [
  {
    id: 'hero',
    name: 'Hero',
    description: 'Section hero avec image de fond, titre et bouton CTA',
    icon: '🎯',
  },
  {
    id: 'productGrid',
    name: 'Grille de produits',
    description: 'Affichage de produits en grille avec filtres',
    icon: '🛍️',
  },
  {
    id: 'testimonials',
    name: 'Témoignages',
    description: 'Carrousel de témoignages clients',
    icon: '💬',
  },
  {
    id: 'features',
    name: 'Caractéristiques',
    description: 'Liste de caractéristiques avec icônes',
    icon: '✨',
  },
  {
    id: 'banner',
    name: 'Bannière',
    description: 'Bannière promotionnelle avec texte/image',
    icon: '📢',
  },
  {
    id: 'richText',
    name: 'Texte enrichi',
    description: 'Section de texte enrichi',
    icon: '📝',
  },
  {
    id: 'imageGallery',
    name: 'Galerie d\'images',
    description: 'Galerie d\'images',
    icon: '🖼️',
  },
]

export default function SectionLibrary({ onSelect, onClose }: SectionLibraryProps) {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-4xl w-full max-h-[80vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">Bibliothèque de sections</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {sectionTypes.map((type) => (
            <button
              key={type.id}
              onClick={() => onSelect(type.id)}
              className="bg-gray-50 hover:bg-gray-100 border border-gray-200 rounded-lg p-4 text-left transition-colors"
            >
              <div className="text-3xl mb-2">{type.icon}</div>
              <h3 className="font-semibold mb-1">{type.name}</h3>
              <p className="text-sm text-gray-600">{type.description}</p>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}

