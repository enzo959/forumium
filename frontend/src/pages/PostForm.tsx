import { useState } from 'react'
import { useParams } from 'react-router-dom'

function PostForm() {
  const { id } = useParams()
  const isEditing = !!id

  const [title, setTitle] = useState('')
  const [error, setError] = useState('')

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
          <a href="/" className="text-2xl font-bold text-blue-600">Forum</a>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-6">
          {isEditing ? 'Modifier le post' : 'Créer un post'}
        </h2>

        <div className="bg-white rounded shadow-sm p-6">
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Titre
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              placeholder="Titre de votre post"
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {error && (
            <p className="text-red-500 text-sm mb-4">{error}</p>
          )}
        </div>
      </main>
    </div>
  )
}

export default PostForm