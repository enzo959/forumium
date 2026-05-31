import { useState, useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import axios from '../api/axios'
import { useAuth } from '../context/AuthContext'
import Navbar from '../components/Navbar'

type Post = {
  ID: number
  title: string
  content: string
  User: { username: string }
  Categories: { name: string }[]
  CreatedAt: string
}

type Category = {
  ID: number
  name: string
}

const LIMIT = 10

function Home() {
  const location = useLocation()
  const [posts, setPosts] = useState<Post[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [selectedCategory, setSelectedCategory] = useState<string>('')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
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

  useEffect(() => {
    const fetchPosts = async () => {
      setLoading(true)
      setError('')
      try {
        const res = await axios.get(`/api/posts?page=${page}&limit=${LIMIT}`)
        setPosts(res.data.posts)
        setTotal(res.data.pagination.total)
      } catch {
        setError('Erreur lors du chargement des posts.')
      } finally {
        setLoading(false)
      }
    }
    fetchPosts()
  }, [page, location.key])

  const filteredPosts = selectedCategory
    ? posts.filter((post) =>
        post.Categories?.some((cat) => cat.name === selectedCategory)
      )
    : posts

  const totalPages = Math.ceil(total / LIMIT)

  const getExcerpt = (html: string) => {
    const text = html?.replace(/<[^>]*>/g, '') || ''
    return text.length > 200 ? text.slice(0, 200) + '...' : text
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-4xl mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">
          Derniers posts
        </h2>

        <div className="flex flex-wrap gap-2 mb-6">
          <button
            onClick={() => { setSelectedCategory(''); setPage(1) }}
            className={`px-3 py-1 text-sm rounded-full border ${
              selectedCategory === ''
                ? 'bg-blue-600 text-white border-blue-600'
                : 'bg-white text-gray-600 border-gray-300 hover:bg-gray-100'
            }`}
          >
            Tous
          </button>
          {categories.map((cat) => (
            <button
              key={cat.ID}
              onClick={() => { setSelectedCategory(cat.name); setPage(1) }}
              className={`px-3 py-1 text-sm rounded-full border ${
                selectedCategory === cat.name
                  ? 'bg-blue-600 text-white border-blue-600'
                  : 'bg-white text-gray-600 border-gray-300 hover:bg-gray-100'
              }`}
            >
              {cat.name}
            </button>
          ))}
        </div>

        {loading && <p className="text-gray-500 text-sm">Chargement des posts...</p>}
        {error && <p className="text-red-500 text-sm">{error}</p>}
        {!loading && !error && filteredPosts.length === 0 && (
          <p className="text-gray-400 text-sm">Aucun post pour le moment.</p>
        )}

        <div className="flex flex-col gap-4">
          {filteredPosts.map((post) => (
            <div key={post.ID} className="bg-white rounded shadow-sm p-5 hover:shadow-md transition">
              <h3 className="text-lg font-semibold text-gray-800 mb-1">
                {post.title}
              </h3>
              <p className="text-sm text-gray-600 mb-3">
                {getExcerpt(post.content)}
              </p>
              <p className="text-sm text-gray-400 mb-3">
                Par <span className="font-medium text-gray-600">{post.User?.username}</span>
                {' · '}
                {new Date(post.CreatedAt).toLocaleDateString('fr-FR')}
              </p>
              {post.Categories && post.Categories.length > 0 && (
                <div className="flex flex-wrap gap-2 mb-3">
                  {post.Categories.map((cat) => (
                    <span
                      key={cat.name}
                      className="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                    >
                      {cat.name}
                    </span>
                  ))}
                </div>
              )}
              <div className="flex justify-end">
                <a
                  href={`/posts/${post.ID}`}
                  className="px-3 py-1 text-sm text-blue-600 border border-blue-600 rounded hover:bg-blue-50"
                >
                  Voir le post
                </a>
              </div>
            </div>
          ))}
        </div>

        {totalPages > 1 && (
          <div className="flex justify-center gap-2 mt-8">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-4 py-2 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-40"
            >
              Précédent
            </button>
            <span className="px-4 py-2 text-sm text-gray-600">
              Page {page} / {totalPages}
            </span>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              className="px-4 py-2 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-40"
            >
              Suivant
            </button>
          </div>
        )}
      </main>
    </div>
  )
}

export default Home