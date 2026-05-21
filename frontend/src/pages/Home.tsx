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
        <div className="flex flex-col gap-4">
          {posts.map((post) => (
            <div key={post.ID} className="bg-white rounded shadow-sm p-5 hover:shadow-md transition">
              <h3 className="text-lg font-semibold text-gray-800 mb-1">
                {post.Title}
              </h3>
              <p className="text-sm text-gray-400 mb-3">
                Par <span className="font-medium text-gray-600">{post.User?.username}</span>
                {' · '}
                {new Date(post.CreatedAt).toLocaleDateString('fr-FR')}
              </p>
              {post.Categories && post.Categories.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {post.Categories.map((cat) => (
                    <span
                      key={cat.Name}
                      className="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                    >
                      {cat.Name}
                    </span>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      </main>
    </div>
  )
}

export default Home