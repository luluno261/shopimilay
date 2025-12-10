import { useState } from 'react'
import Head from 'next/head'
import GrowthCommandCenter from '../components/GrowthCommandCenter'

export default function GrowthCommandCenterPage() {
  return (
    <>
      <Head>
        <title>Growth Command Center - OmniSphere Admin</title>
        <meta name="description" content="Centre de commande marketing OmniSphere" />
      </Head>

      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Growth Command Center</h1>
        <p className="text-gray-600 mb-8">
          Centre de commande centralisé pour gérer toutes vos activités marketing
        </p>
        <GrowthCommandCenter />
      </div>
    </>
  )
}

