import { useState, useEffect } from 'react'

interface Section {
  id: string
  type: string
  data: any
  order: number
}

interface SectionEditorProps {
  section: Section
  onSave: (section: Section) => void
  onCancel: () => void
}

export default function SectionEditor({ section, onSave, onCancel }: SectionEditorProps) {
  const [data, setData] = useState(section.data || {})

  useEffect(() => {
    setData(section.data || {})
  }, [section])

  const handleSave = () => {
    onSave({
      ...section,
      data,
    })
  }

  const renderEditor = () => {
    switch (section.type) {
      case 'hero':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Titre</label>
              <input
                type="text"
                value={data.title || ''}
                onChange={(e) => setData({ ...data, title: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                placeholder="Titre principal"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Sous-titre</label>
              <input
                type="text"
                value={data.subtitle || ''}
                onChange={(e) => setData({ ...data, subtitle: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                placeholder="Sous-titre"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">URL de l'image</label>
              <input
                type="text"
                value={data.imageUrl || ''}
                onChange={(e) => setData({ ...data, imageUrl: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                placeholder="https://example.com/image.jpg"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Texte du bouton CTA</label>
              <input
                type="text"
                value={data.ctaText || ''}
                onChange={(e) => setData({ ...data, ctaText: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                placeholder="Découvrir"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">URL du bouton CTA</label>
              <input
                type="text"
                value={data.ctaUrl || ''}
                onChange={(e) => setData({ ...data, ctaUrl: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                placeholder="/products"
              />
            </div>
          </div>
        )

      case 'banner':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Texte</label>
              <input
                type="text"
                value={data.text || ''}
                onChange={(e) => setData({ ...data, text: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                placeholder="Texte de la bannière"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">URL de l'image (optionnel)</label>
              <input
                type="text"
                value={data.imageUrl || ''}
                onChange={(e) => setData({ ...data, imageUrl: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                placeholder="https://example.com/image.jpg"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Couleur de fond</label>
              <input
                type="color"
                value={data.backgroundColor || '#3B82F6'}
                onChange={(e) => setData({ ...data, backgroundColor: e.target.value })}
                className="w-full h-10 border border-gray-300 rounded-md"
              />
            </div>
          </div>
        )

      case 'richText':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Contenu</label>
              <textarea
                value={data.content || ''}
                onChange={(e) => setData({ ...data, content: e.target.value })}
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                rows={10}
                placeholder="Contenu HTML ou texte"
              />
            </div>
          </div>
        )

      default:
        return (
          <div className="text-gray-500">
            Éditeur pour le type "{section.type}" - À implémenter
          </div>
        )
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Éditer la section</h2>
        <div className="flex space-x-2">
          <button
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50"
          >
            Annuler
          </button>
          <button
            onClick={handleSave}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Enregistrer
          </button>
        </div>
      </div>

      {renderEditor()}
    </div>
  )
}

