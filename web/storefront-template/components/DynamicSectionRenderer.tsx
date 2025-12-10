import HeroSection from './builder/sections/HeroSection'
import ProductGridSection from './builder/sections/ProductGridSection'
import TestimonialsSection from './builder/sections/TestimonialsSection'
import FeaturesSection from './builder/sections/FeaturesSection'
import BannerSection from './builder/sections/BannerSection'
import RichTextSection from './builder/sections/RichTextSection'
import ImageGallerySection from './builder/sections/ImageGallerySection'

interface Section {
  id: string
  type: string
  data: any
  order: number
}

interface DynamicSectionRendererProps {
  sections: Section[]
}

export default function DynamicSectionRenderer({ sections }: DynamicSectionRendererProps) {
  const renderSection = (section: Section) => {
    switch (section.type) {
      case 'hero':
        return <HeroSection key={section.id} data={section.data} />
      case 'productGrid':
        return <ProductGridSection key={section.id} data={section.data} />
      case 'testimonials':
        return <TestimonialsSection key={section.id} data={section.data} />
      case 'features':
        return <FeaturesSection key={section.id} data={section.data} />
      case 'banner':
        return <BannerSection key={section.id} data={section.data} />
      case 'richText':
        return <RichTextSection key={section.id} data={section.data} />
      case 'imageGallery':
        return <ImageGallerySection key={section.id} data={section.data} />
      default:
        return (
          <div key={section.id} className="py-8 text-center text-gray-500">
            Section de type "{section.type}" non reconnue
          </div>
        )
    }
  }

  // Trier les sections par ordre
  const sortedSections = [...sections].sort((a, b) => a.order - b.order)

  return (
    <div>
      {sortedSections.map((section) => renderSection(section))}
    </div>
  )
}

