export interface Theme {
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

export const defaultTheme: Theme = {
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

export function applyTheme(theme: Theme) {
  if (typeof document === 'undefined') return

  const root = document.documentElement

  // Apply colors
  root.style.setProperty('--color-primary', theme.colors.primary)
  root.style.setProperty('--color-secondary', theme.colors.secondary)
  root.style.setProperty('--color-background', theme.colors.background)
  root.style.setProperty('--color-text', theme.colors.text)
  root.style.setProperty('--color-link', theme.colors.link)
  root.style.setProperty('--color-button', theme.colors.button)

  // Apply typography
  root.style.setProperty('--font-family', theme.typography.fontFamily)
  root.style.setProperty('--font-size-base', `${theme.typography.baseSize}px`)
  root.style.setProperty('--font-size-h1', `${theme.typography.headingSizes.h1}px`)
  root.style.setProperty('--font-size-h2', `${theme.typography.headingSizes.h2}px`)
  root.style.setProperty('--font-size-h3', `${theme.typography.headingSizes.h3}px`)
  root.style.setProperty('--font-size-h4', `${theme.typography.headingSizes.h4}px`)
  root.style.setProperty('--font-size-h5', `${theme.typography.headingSizes.h5}px`)
  root.style.setProperty('--font-size-h6', `${theme.typography.headingSizes.h6}px`)
  root.style.setProperty('--line-height', theme.typography.lineHeight.toString())
  root.style.setProperty('--font-weight', theme.typography.fontWeight.toString())

  // Apply layout
  root.style.setProperty('--container-width', `${theme.layout.containerWidth}px`)
  root.style.setProperty('--padding', `${theme.layout.padding}px`)
  root.style.setProperty('--margin', `${theme.layout.margin}px`)
  root.style.setProperty('--border-radius', `${theme.layout.borderRadius}px`)
  root.style.setProperty('--shadow', theme.layout.shadow)
}

export async function loadTheme(merchantId: string): Promise<Theme> {
  try {
    const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
    const response = await fetch(`${API_BASE_URL}/store-builder/theme?merchant_id=${merchantId}`)
    if (response.ok) {
      const data = await response.json()
      return data.theme || defaultTheme
    }
  } catch (error) {
    console.error('Erreur lors du chargement du thème:', error)
  }
  return defaultTheme
}

export async function saveTheme(merchantId: string, theme: Theme): Promise<void> {
  try {
    const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
    await fetch(`${API_BASE_URL}/store-builder/theme`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        merchant_id: merchantId,
        theme,
      }),
    })
  } catch (error) {
    console.error('Erreur lors de la sauvegarde du thème:', error)
    throw error
  }
}

