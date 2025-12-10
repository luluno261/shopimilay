interface Layout {
  containerWidth: number
  padding: number
  margin: number
  borderRadius: number
  shadow: string
}

interface LayoutSettingsProps {
  layout: Layout
  onChange: (layout: Layout) => void
}

export default function LayoutSettings({ layout, onChange }: LayoutSettingsProps) {
  const handleChange = (key: keyof Layout, value: number | string) => {
    onChange({
      ...layout,
      [key]: value,
    })
  }

  return (
    <div className="space-y-6">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Largeur du conteneur: {layout.containerWidth}px
        </label>
        <input
          type="range"
          min="800"
          max="1920"
          step="20"
          value={layout.containerWidth}
          onChange={(e) => handleChange('containerWidth', parseInt(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Padding: {layout.padding}px
        </label>
        <input
          type="range"
          min="0"
          max="64"
          step="4"
          value={layout.padding}
          onChange={(e) => handleChange('padding', parseInt(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Margin: {layout.margin}px
        </label>
        <input
          type="range"
          min="0"
          max="64"
          step="4"
          value={layout.margin}
          onChange={(e) => handleChange('margin', parseInt(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Rayon des bordures: {layout.borderRadius}px
        </label>
        <input
          type="range"
          min="0"
          max="32"
          step="2"
          value={layout.borderRadius}
          onChange={(e) => handleChange('borderRadius', parseInt(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Ombre</label>
        <select
          value={layout.shadow}
          onChange={(e) => handleChange('shadow', e.target.value)}
          className="w-full border border-gray-300 rounded-md px-3 py-2"
        >
          <option value="none">Aucune</option>
          <option value="0 1px 3px rgba(0,0,0,0.1)">Légère</option>
          <option value="0 4px 6px rgba(0,0,0,0.1)">Moyenne</option>
          <option value="0 10px 20px rgba(0,0,0,0.15)">Forte</option>
        </select>
      </div>
    </div>
  )
}

