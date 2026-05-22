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
      </main>
    </div>
  )
}

export default PostDetail