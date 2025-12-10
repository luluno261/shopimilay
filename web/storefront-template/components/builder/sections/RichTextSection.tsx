interface RichTextSectionProps {
  data: {
    content?: string
  }
}

export default function RichTextSection({ data }: RichTextSectionProps) {
  return (
    <div className="py-12 bg-white">
      <div className="container mx-auto px-4">
        <div
          className="prose max-w-none"
          dangerouslySetInnerHTML={{ __html: data.content || '<p>Contenu à ajouter</p>' }}
        />
      </div>
    </div>
  )
}

