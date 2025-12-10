import Head from 'next/head'

export default function PrivacyPage() {
  return (
    <>
      <Head>
        <title>Politique de Confidentialité - OmniSphere</title>
      </Head>

      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <h1 className="text-3xl font-bold mb-6">Politique de Confidentialité</h1>

        <div className="prose prose-lg">
          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">1. Collecte des données</h2>
            <p>
              Nous collectons les données personnelles que vous nous fournissez lors de
              l'inscription, de la commande ou de l'utilisation de nos services.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">2. Utilisation des données</h2>
            <p>
              Vos données sont utilisées pour traiter vos commandes, améliorer nos services
              et vous envoyer des communications marketing (avec votre consentement).
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">3. Protection des données</h2>
            <p>
              Nous mettons en œuvre des mesures de sécurité appropriées pour protéger vos
              données personnelles contre tout accès non autorisé.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">4. Vos droits</h2>
            <p>
              Conformément au RGPD, vous avez le droit d'accéder, de rectifier, de supprimer
              vos données personnelles. Vous pouvez exercer ces droits en nous contactant ou
              en utilisant notre{' '}
              <a href="/data-deletion" className="text-blue-600 underline">
                page de suppression de données
              </a>
              .
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">5. Cookies</h2>
            <p>
              Nous utilisons des cookies pour améliorer votre expérience. Vous pouvez gérer
              vos préférences de cookies dans les paramètres de votre navigateur.
            </p>
          </section>
        </div>
      </div>
    </>
  )
}

