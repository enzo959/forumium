import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import axios from '../api/axios'

type Post = {
  id: number
  title: string
  content: string
  image: string
  created_at: string
  author: {
    id: number
    username: string
    avatar: string
  }
  categories: { id: number; name: string }[]
  likes: number
  dislikes: number
}

function PostDetail() {
  const { id } = useParams()
  const [post, setPost] = useState<Post | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const res = await axios.get(`/api/posts/${id}`)
        setPost(res.data)
      } catch {
        setError('Post introuvable ou erreur serveur.')
      } finally {
        setLoading(false)
      }
    }
    fetchPost()
  }, [id])

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
          <a href="/" className="text-2xl font-bold text-blue-600">Forum</a>
          <div className="flex gap-3">
            <a href="/login" className="px-4 py-2 text-sm text-blue-600 border border-blue-600 rounded hover:bg-blue-50">
              Connexion
            </a>
            <a href="/register" className="px-4 py-2 text-sm text-white bg-blue-600 rounded hover:bg-blue-700">
              Inscription
            </a>
          </div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        <a href="/" className="text-sm text-blue-600 hover:underline mb-6 inline-block">
          ← Retour aux posts
        </a>

        {loading && <p className="text-gray-500 text-sm">Chargement...</p>}
        {error && <p className="text-red-500 text-sm">{error}</p>}

        {post && (
          <div className="bg-white rounded shadow-sm p-6">
            <h1 className="text-2xl font-bold text-gray-800 mb-2">
              {post.title}
            </h1>
            <p className="text-sm text-gray-400 mb-4">
              Par <span className="font-medium text-gray-600">{post.author.username}</span>
              {' · '}
              {new Date(post.created_at).toLocaleDateString('fr-FR')}
            </p>
            {post.categories && post.categories.length > 0 && (
              <div className="flex flex-wrap gap-2 mb-6">
                {post.categories.map((cat) => (
                  <span
                    key={cat.id}
                    className="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                  >
                    {cat.name}
                  </span>
                ))}
              </div>
            )}
            {post.image && (
              <img
                src={`http://localhost:8080/${post.image}`}
                alt={post.title}
                className="w-full max-h-96 object-cover rounded mb-6"
              />
            )}
            <div
              className="prose max-w-none text-gray-700 border-t border-gray-100 pt-6"
              dangerouslySetInnerHTML={{ __html: post.content }}
            />
          </div>
        )}
      </main>
    </div>
  )
}

export default PostDetail