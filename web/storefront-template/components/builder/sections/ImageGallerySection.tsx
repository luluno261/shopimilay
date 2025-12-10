interface ImageGallerySectionProps {
  data: {
    title?: string
    images?: string[]
  }
}

const defaultImages = [
  'https://via.placeholder.com/400x300',
  'https://via.placeholder.com/400x300',
  'https://via.placeholder.com/400x300',
  'https://via.placeholder.com/400x300',
]

export default function ImageGallerySection({ data }: ImageGallerySectionProps) {
  const images = data.images || defaultImages

  return (
    <div className="py-12 bg-white">
      <div className="container mx-auto px-4">
        {data.title && <h2 className="text-3xl font-bold mb-8 text-center">{data.title}</h2>}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {images.map((image, index) => (
            <div key={index} className="overflow-hidden rounded-lg">
              <img
                src={image}
                alt={`Gallery image ${index + 1}`}
                className="w-full h-64 object-cover hover:scale-105 transition-transform"
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

