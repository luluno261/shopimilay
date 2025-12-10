interface HeroSectionProps {
  data: {
    title?: string
    subtitle?: string
    imageUrl?: string
    ctaText?: string
    ctaUrl?: string
  }
}

export default function HeroSection({ data }: HeroSectionProps) {
  return (
    <div
      className="relative h-96 flex items-center justify-center bg-cover bg-center"
      style={{
        backgroundImage: data.imageUrl ? `url(${data.imageUrl})` : 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}
    >
      <div className="absolute inset-0 bg-black bg-opacity-40"></div>
      <div className="relative z-10 text-center text-white px-4">
        <h1 className="text-4xl md:text-6xl font-bold mb-4">{data.title || 'Titre principal'}</h1>
        {data.subtitle && <p className="text-xl mb-6">{data.subtitle}</p>}
        {data.ctaText && (
          <a
            href={data.ctaUrl || '#'}
            className="inline-block bg-white text-gray-900 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors"
          >
            {data.ctaText}
          </a>
        )}
      </div>
    </div>
  )
}

