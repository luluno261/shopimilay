interface BannerSectionProps {
  data: {
    text?: string
    imageUrl?: string
    backgroundColor?: string
    textColor?: string
  }
}

export default function BannerSection({ data }: BannerSectionProps) {
  const backgroundColor = data.backgroundColor || '#3B82F6'
  const textColor = data.textColor || '#FFFFFF'

  return (
    <div
      className="py-8 px-4 text-center"
      style={{
        backgroundColor: data.imageUrl ? 'transparent' : backgroundColor,
        backgroundImage: data.imageUrl ? `url(${data.imageUrl})` : undefined,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        color: textColor,
      }}
    >
      {data.imageUrl && <div className="absolute inset-0 bg-black bg-opacity-40"></div>}
      <div className="relative z-10">
        <p className="text-xl md:text-2xl font-semibold">{data.text || 'Offre spéciale'}</p>
      </div>
    </div>
  )
}

