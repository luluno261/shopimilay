import { useState } from 'react'
import ColorPicker from './ColorPicker'
import TypographyEditor from './TypographyEditor'
import LayoutSettings from './LayoutSettings'

interface Theme {
  colors: {
    primary: string
    secondary: string
    background: string
    text: string
    link: string
    button: string
  }
  typography: {
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
  layout: {
    containerWidth: number
    padding: number
    margin: number
    borderRadius: number
    shadow: string
  }
}

interface ThemeCustomizerProps {
  theme?: Theme
  onThemeChange?: (theme: Theme) => void
}

const defaultTheme: Theme = {
  colors: {
    primary: '#3B82F6',
    secondary: '#8B5CF6',
    background: '#FFFFFF',
    text: '#1F2937',
    link: '#3B82F6',
    button: '#3B82F6',
  },
  typography: {
    fontFamily: 'Inter, sans-serif',
    baseSize: 16,
    headingSizes: {
      h1: 48,
      h2: 36,
      h3: 30,
      h4: 24,
      h5: 20,
      h6: 18,
    },
    lineHeight: 1.5,
    fontWeight: 400,
  },
  layout: {
    containerWidth: 1200,
    padding: 16,
    margin: 16,
    borderRadius: 8,
    shadow: '0 1px 3px rgba(0,0,0,0.1)',
  },
}

export default function ThemeCustomizer({ theme: initialTheme, onThemeChange }: ThemeCustomizerProps) {
  const [theme, setTheme] = useState<Theme>(initialTheme || defaultTheme)
  const [activeTab, setActiveTab] = useState<'colors' | 'typography' | 'layout'>('colors')

  const handleThemeUpdate = (updates: Partial<Theme>) => {
    const newTheme = { ...theme, ...updates }
    setTheme(newTheme)
    if (onThemeChange) {
      onThemeChange(newTheme)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold mb-6">Personnalisation du thème</h2>

      {/* Tabs */}
      <div className="border-b border-gray-200 mb-6">
        <nav className="flex space-x-1">
          {[
            { id: 'colors' as const, label: 'Couleurs', icon: '🎨' },
            { id: 'typography' as const, label: 'Typographie', icon: '✍️' },
            { id: 'layout' as const, label: 'Mise en page', icon: '📐' },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                px-4 py-2 text-sm font-medium border-b-2 transition-colors
                ${
                  activeTab === tab.id
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }
              `}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div>
        {activeTab === 'colors' && (
          <ColorPicker
            colors={theme.colors}
            onChange={(colors) => handleThemeUpdate({ colors })}
          />
        )}
        {activeTab === 'typography' && (
          <TypographyEditor
            typography={theme.typography}
            onChange={(typography) => handleThemeUpdate({ typography })}
          />
        )}
        {activeTab === 'layout' && (
          <LayoutSettings
            layout={theme.layout}
            onChange={(layout) => handleThemeUpdate({ layout })}
          />
        )}
      </div>
    </div>
  )
}

