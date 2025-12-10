interface Typography {
  fontFamily: string
  baseSize: number
  headingSizes: {
    h1: number
    h2: number
    h3: number
    h4: number
    h5: number
    h6: number
  }
  lineHeight: number
  fontWeight: number
}

interface TypographyEditorProps {
  typography: Typography
  onChange: (typography: Typography) => void
}

const fontFamilies = [
  'Inter, sans-serif',
  'Roboto, sans-serif',
  'Open Sans, sans-serif',
  'Lato, sans-serif',
  'Montserrat, sans-serif',
  'Playfair Display, serif',
  'Merriweather, serif',
]

export default function TypographyEditor({ typography, onChange }: TypographyEditorProps) {
  const handleChange = (key: keyof Typography, value: any) => {
    onChange({
      ...typography,
      [key]: value,
    })
  }

  const handleHeadingSizeChange = (heading: keyof Typography['headingSizes'], value: number) => {
    onChange({
      ...typography,
      headingSizes: {
        ...typography.headingSizes,
        [heading]: value,
      },
    })
  }

  return (
    <div className="space-y-6">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Famille de polices</label>
        <select
          value={typography.fontFamily}
          onChange={(e) => handleChange('fontFamily', e.target.value)}
          className="w-full border border-gray-300 rounded-md px-3 py-2"
        >
          {fontFamilies.map((font) => (
            <option key={font} value={font}>
              {font}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Taille de base: {typography.baseSize}px
        </label>
        <input
          type="range"
          min="12"
          max="24"
          value={typography.baseSize}
          onChange={(e) => handleChange('baseSize', parseInt(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Tailles des titres</label>
        <div className="space-y-3">
          {Object.entries(typography.headingSizes).map(([heading, size]) => (
            <div key={heading} className="flex items-center justify-between">
              <label className="text-sm text-gray-600 w-12 uppercase">{heading}</label>
              <div className="flex items-center space-x-4 flex-1">
                <input
                  type="range"
                  min="16"
                  max="72"
                  value={size}
                  onChange={(e) => handleHeadingSizeChange(heading as keyof Typography['headingSizes'], parseInt(e.target.value))}
                  className="flex-1"
                />
                <span className="text-sm text-gray-600 w-12 text-right">{size}px</span>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Hauteur de ligne: {typography.lineHeight}
        </label>
        <input
          type="range"
          min="1"
          max="2"
          step="0.1"
          value={typography.lineHeight}
          onChange={(e) => handleChange('lineHeight', parseFloat(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Poids de la police: {typography.fontWeight}
        </label>
        <input
          type="range"
          min="300"
          max="700"
          step="100"
          value={typography.fontWeight}
          onChange={(e) => handleChange('fontWeight', parseInt(e.target.value))}
          className="w-full"
        />
      </div>
    </div>
  )
}

