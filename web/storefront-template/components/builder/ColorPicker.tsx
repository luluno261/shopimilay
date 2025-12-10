interface Colors {
  primary: string
  secondary: string
  background: string
  text: string
  link: string
  button: string
}

interface ColorPickerProps {
  colors: Colors
  onChange: (colors: Colors) => void
}

const colorLabels: Record<keyof Colors, string> = {
  primary: 'Couleur primaire',
  secondary: 'Couleur secondaire',
  background: 'Couleur de fond',
  text: 'Couleur de texte',
  link: 'Couleur des liens',
  button: 'Couleur des boutons',
}

export default function ColorPicker({ colors, onChange }: ColorPickerProps) {
  const handleColorChange = (key: keyof Colors, value: string) => {
    onChange({
      ...colors,
      [key]: value,
    })
  }

  return (
    <div className="space-y-6">
      {Object.entries(colorLabels).map(([key, label]) => (
        <div key={key} className="flex items-center justify-between">
          <label className="text-sm font-medium text-gray-700 w-48">{label}</label>
          <div className="flex items-center space-x-4 flex-1">
            <input
              type="color"
              value={colors[key as keyof Colors]}
              onChange={(e) => handleColorChange(key as keyof Colors, e.target.value)}
              className="w-16 h-10 border border-gray-300 rounded cursor-pointer"
            />
            <input
              type="text"
              value={colors[key as keyof Colors]}
              onChange={(e) => handleColorChange(key as keyof Colors, e.target.value)}
              className="flex-1 border border-gray-300 rounded-md px-3 py-2"
              placeholder="#000000"
            />
          </div>
        </div>
      ))}
    </div>
  )
}

