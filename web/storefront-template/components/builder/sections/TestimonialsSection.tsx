interface TestimonialsSectionProps {
  data: {
    title?: string
    testimonials?: Array<{
      name: string
      role: string
      content: string
      avatar?: string
    }>
  }
}

const defaultTestimonials = [
  {
    name: 'Jean Dupont',
    role: 'Client satisfait',
    content: 'Excellent service et produits de qualité. Je recommande vivement !',
  },
  {
    name: 'Marie Martin',
    role: 'Cliente fidèle',
    content: 'Une expérience d\'achat exceptionnelle. Livraison rapide et produits conformes.',
  },
  {
    name: 'Pierre Durand',
    role: 'Nouveau client',
    content: 'Très satisfait de mon premier achat. Je reviendrai certainement.',
  },
]

export default function TestimonialsSection({ data }: TestimonialsSectionProps) {
  const testimonials = data.testimonials || defaultTestimonials

  return (
    <div className="py-12 bg-white">
      <div className="container mx-auto px-4">
        {data.title && <h2 className="text-3xl font-bold mb-8 text-center">{data.title}</h2>}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {testimonials.map((testimonial, index) => (
            <div key={index} className="bg-gray-50 rounded-lg p-6 border border-gray-200">
              <p className="text-gray-700 mb-4 italic">&quot;{testimonial.content}&quot;</p>
              <div className="flex items-center">
                {testimonial.avatar && (
                  <img
                    src={testimonial.avatar}
                    alt={testimonial.name}
                    className="w-12 h-12 rounded-full mr-4"
                  />
                )}
                <div>
                  <p className="font-semibold">{testimonial.name}</p>
                  <p className="text-sm text-gray-500">{testimonial.role}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

