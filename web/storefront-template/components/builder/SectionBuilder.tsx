import { useState, useEffect } from 'react'
import SectionLibrary from './SectionLibrary'
import SectionEditor from './SectionEditor'

interface Section {
  id: string
  type: string
  data: any
  order: number
}

interface SectionBuilderProps {
  sections?: Section[]
  onSectionsChange?: (sections: Section[]) => void
}

function SectionItem({ section, onEdit, onDelete, onMoveUp, onMoveDown, canMoveUp, canMoveDown }: { 
  section: Section
  onEdit: () => void
  onDelete: () => void
  onMoveUp: () => void
  onMoveDown: () => void
  canMoveUp: boolean
  canMoveDown: boolean
}) {
  const getSectionName = (type: string) => {
    const names: Record<string, string> = {
      hero: 'Hero',
      productGrid: 'Grille de produits',
      testimonials: 'Témoignages',
      features: 'Caractéristiques',
      banner: 'Bannière',
      richText: 'Texte enrichi',
      imageGallery: 'Galerie d\'images',
    }
    return names[type] || type
  }

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-4 mb-4 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="flex flex-col space-y-1">
            <button
              onClick={onMoveUp}
              disabled={!canMoveUp}
              className="text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed"
            >
              ↑
            </button>
            <button
              onClick={onMoveDown}
              disabled={!canMoveDown}
              className="text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed"
            >
              ↓
            </button>
          </div>
          <div>
            <h3 className="font-semibold">{getSectionName(section.type)}</h3>
            <p className="text-sm text-gray-500">Section {section.type}</p>
          </div>
        </div>
        <div className="flex space-x-2">
          <button
            onClick={onEdit}
            className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Éditer
          </button>
          <button
            onClick={onDelete}
            className="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700"
          >
            Supprimer
          </button>
        </div>
      </div>
    </div>
  )
}

export default function SectionBuilder({ sections: initialSections = [], onSectionsChange }: SectionBuilderProps) {
  const [sections, setSections] = useState<Section[]>(initialSections)
  const [showLibrary, setShowLibrary] = useState(false)
  const [editingSection, setEditingSection] = useState<Section | null>(null)

  useEffect(() => {
    if (onSectionsChange) {
      onSectionsChange(sections)
    }
  }, [sections, onSectionsChange])

  const handleMoveUp = (index: number) => {
    if (index > 0) {
      const newSections = [...sections]
      const temp = newSections[index]
      newSections[index] = newSections[index - 1]
      newSections[index - 1] = temp
      setSections(newSections.map((s, i) => ({ ...s, order: i })))
    }
  }

  const handleMoveDown = (index: number) => {
    if (index < sections.length - 1) {
      const newSections = [...sections]
      const temp = newSections[index]
      newSections[index] = newSections[index + 1]
      newSections[index + 1] = temp
      setSections(newSections.map((s, i) => ({ ...s, order: i })))
    }
  }

  const handleAddSection = (type: string) => {
    const newSection: Section = {
      id: `section-${Date.now()}`,
      type,
      data: {},
      order: sections.length,
    }
    setSections([...sections, newSection])
    setShowLibrary(false)
  }

  const handleEditSection = (section: Section) => {
    setEditingSection(section)
  }

  const handleSaveSection = (updatedSection: Section) => {
    setSections(sections.map((s) => (s.id === updatedSection.id ? updatedSection : s)))
    setEditingSection(null)
  }

  const handleDeleteSection = (sectionId: string) => {
    setSections(sections.filter((s) => s.id !== sectionId))
  }

  if (editingSection) {
    return (
      <SectionEditor
        section={editingSection}
        onSave={handleSaveSection}
        onCancel={() => setEditingSection(null)}
      />
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Builder de Sections</h2>
        <button
          onClick={() => setShowLibrary(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          + Ajouter une section
        </button>
      </div>

      {showLibrary && (
        <SectionLibrary
          onSelect={handleAddSection}
          onClose={() => setShowLibrary(false)}
        />
      )}

      {sections.length === 0 ? (
        <div className="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
          <p className="text-gray-500 mb-4">Aucune section configurée</p>
          <button
            onClick={() => setShowLibrary(true)}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
          >
            Ajouter votre première section
          </button>
        </div>
      ) : (
        <div>
          {sections.map((section, index) => (
            <SectionItem
              key={section.id}
              section={section}
              onEdit={() => handleEditSection(section)}
              onDelete={() => handleDeleteSection(section.id)}
              onMoveUp={() => handleMoveUp(index)}
              onMoveDown={() => handleMoveDown(index)}
              canMoveUp={index > 0}
              canMoveDown={index < sections.length - 1}
            />
          ))}
        </div>
      )}
    </div>
  )
}

