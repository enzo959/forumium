import { useState, useEffect } from 'react'
import axios from '../api/axios'

type Post = {
  ID: number
  Title: string
  Content: string
  User: { username: string }
  Categories: { Name: string }[]
  CreatedAt: string
}

function Home() {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const res = await axios.get('/api/posts')
        setPosts(res.data.posts)
      } catch {
        setError('Erreur lors du chargement des posts.')
      } finally {
        setLoading(false)
      }
    }

    fetchPosts()
  }, [])

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-blue-600">Forum</h1>
          <div className="flex gap-3">
            <a
              href="/login"
              className="px-4 py-2 text-sm text-blue-600 border border-blue-600 rounded hover:bg-blue-50"
            >
              Connexion
            </a>
            <a
              href="/register"
              className="px-4 py-2 text-sm text-white bg-blue-600 rounded hover:bg-blue-700"
            >
              Inscription
            </a>
          </div>
        </div>
      </header>
      <main className="max-w-4xl mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-6">
          Derniers posts
        </h2>

        {loading && (
          <p className="text-gray-500 text-sm">Chargement des posts...</p>
        )}

        {error && (
          <p className="text-red-500 text-sm">{error}</p>
        )}

        {!loading && !error && posts.length === 0 && (
          <p className="text-gray-400 text-sm">Aucun post pour le moment.</p>
        )}
      </main>
    </div>
  )
}

export default Home