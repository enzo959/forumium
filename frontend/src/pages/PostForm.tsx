import { useParams } from 'react-router-dom'

function PostForm() {
  const { id } = useParams()
  const isEditing = !!id

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
          <a href="/" className="text-2xl font-bold text-blue-600">Forum</a>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-6">
          {isEditing ? 'Modifier le post' : 'Créer un post'}
        </h2>
      </main>
    </div>
  )
}

export default PostForm