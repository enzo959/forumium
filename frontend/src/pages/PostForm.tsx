import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import RichTextEditor from '../components/RichTextEditor'
import axios from '../api/axios'

type Category = {
  ID: number
  name: string
}

function PostForm() {
  const { id } = useParams()
  const isEditing = !!id

  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [categories, setCategories] = useState<Category[]>([])
  const [selectedCategoryIDs, setSelectedCategoryIDs] = useState<number[]>([])
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await axios.get('/api/categories')
        setCategories(res.data.categories)
      } catch {
        console.error('Erreur chargement catégories')
      }
    }
    fetchCategories()
  }, [])

  const toggleCategory = (id: number) => {
    setSelectedCategoryIDs((prev) =>
      prev.includes(id) ? prev.filter((c) => c !== id) : [...prev, id]
    )
  }

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

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Contenu
            </label>
            <RichTextEditor content={content} onChange={setContent} />
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Catégories
            </label>
            <div className="flex flex-wrap gap-2">
              {categories.map((cat) => (
                <button
                  key={cat.ID}
                  type="button"
                  onClick={() => toggleCategory(cat.ID)}
                  className={`px-3 py-1 text-sm rounded-full border ${
                    selectedCategoryIDs.includes(cat.ID)
                      ? 'bg-blue-600 text-white border-blue-600'
                      : 'bg-white text-gray-600 border-gray-300 hover:bg-gray-100'
                  }`}
                >
                  {cat.name}
                </button>
              ))}
            </div>
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