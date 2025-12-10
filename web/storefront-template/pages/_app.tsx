import type { AppProps } from 'next/app'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useEffect } from 'react'
import { useRouter } from 'next/router'
import '../styles/globals.css'
import '../styles/theme-variables.css'
import CookieBanner from '../components/CookieBanner'
import { loadTheme, applyTheme } from '../lib/theme'

const queryClient = new QueryClient()

export default function App({ Component, pageProps }: AppProps) {
  const router = useRouter()

  useEffect(() => {
    const loadAndApplyTheme = async () => {
      const merchantId = router.query.merchant_id as string || 'default'
      const theme = await loadTheme(merchantId)
      applyTheme(theme)
    }

    loadAndApplyTheme()
  }, [router.query.merchant_id])

  return (
    <QueryClientProvider client={queryClient}>
      <Component {...pageProps} />
      <CookieBanner />
    </QueryClientProvider>
  )
}

