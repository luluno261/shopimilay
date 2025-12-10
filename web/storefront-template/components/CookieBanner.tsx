import React, { useState, useEffect } from 'react'

export default function CookieBanner() {
  const [show, setShow] = useState(false)

  useEffect(() => {
    const consent = localStorage.getItem('cookie_consent')
    if (!consent) {
      setShow(true)
    }
  }, [])

  const handleAccept = () => {
    localStorage.setItem('cookie_consent', 'accepted')
    setShow(false)
  }

  const handleReject = () => {
    localStorage.setItem('cookie_consent', 'rejected')
    setShow(false)
  }

  if (!show) return null

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-gray-900 text-white p-4 z-50 shadow-lg">
      <div className="container mx-auto flex items-center justify-between">
        <div className="flex-1 mr-4">
          <p className="text-sm">
            Nous utilisons des cookies pour améliorer votre expérience sur notre site.
            En continuant à naviguer, vous acceptez notre{' '}
            <a href="/privacy" className="underline">
              politique de confidentialité
            </a>
            .
          </p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleReject}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded"
          >
            Refuser
          </button>
          <button
            onClick={handleAccept}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded"
          >
            Accepter
          </button>
        </div>
      </div>
    </div>
  )
}

