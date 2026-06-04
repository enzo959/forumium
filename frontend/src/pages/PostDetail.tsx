import { useState, useEffect } from 'react'
import { useParams, useLocation } from 'react-router-dom'
import axios from '../api/axios'
import Navbar from '../components/Navbar'

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

type Comment = {
  id: number
  content: string
  created_at: string
  author: {
    id: number
    username: string
    avatar: string
  }
}

function PostDetail() {
  const { id } = useParams()
  const location = useLocation()
  const [post, setPost] = useState<Post | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchPost = async () => {
      setLoading(true)
      setError('')
      try {
        const [postRes, commentsRes] = await Promise.all([
          axios.get(`/api/posts/${id}`),
          axios.get(`/api/posts/${id}/comments`),
        ])
        setPost(postRes.data)
        setComments(commentsRes.data.comments)
      } catch {
        setError('Post introuvable ou erreur serveur.')
      } finally {
        setLoading(false)
      }
    }
    fetchPost()
  }, [id, location.key])

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-4xl mx-auto px-4 py-8">
        <a href="/" className="text-sm text-blue-600 hover:underline mb-6 inline-block">
          ← Retour aux posts
        </a>

        {loading && <p className="text-gray-500 text-sm">Chargement...</p>}
        {error && <p className="text-red-500 text-sm">{error}</p>}

        {post && (
          <div className="bg-white rounded shadow-sm p-6 mb-6">
            <h1 className="text-2xl font-bold text-gray-800 mb-2">{post.title}</h1>
            <p className="text-sm text-gray-400 mb-4">
              Par <span className="font-medium text-gray-600">{post.author.username}</span>
              {' · '}
              {new Date(post.created_at).toLocaleDateString('fr-FR')}
            </p>
            {post.categories && post.categories.length > 0 && (
              <div className="flex flex-wrap gap-2 mb-6">
                {post.categories.map((cat) => (
                  <span key={cat.id} className="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full">
                    {cat.name}
                  </span>
                ))}
              </div>
            )}
            {post.image && (
              <img src={post.image} alt={post.title} className="w-full max-h-96 object-cover rounded mb-6" />
            )}
            <div
              className="prose max-w-none text-gray-700 border-t border-gray-100 pt-6"
              dangerouslySetInnerHTML={{ __html: post.content }}
            />
          </div>
        )}

        {/* Section commentaires */}
        <div className="bg-white rounded shadow-sm p-6">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">
            Commentaires ({comments.length})
          </h2>

          {comments.length === 0 && (
            <p className="text-gray-400 text-sm">Aucun commentaire pour le moment.</p>
          )}

          <div className="flex flex-col gap-4">
            {comments.map((comment) => (
              <div key={comment.id} className="border-b border-gray-100 pb-4">
                <p className="text-sm font-medium text-gray-700 mb-1">
                  {comment.author.username}
                  <span className="text-gray-400 font-normal ml-2">
                    {new Date(comment.created_at).toLocaleDateString('fr-FR')}
                  </span>
                </p>
                <p className="text-sm text-gray-600">{comment.content}</p>
              </div>
            ))}
          </div>
        </div>
      </main>
    </div>
  )
}

export default PostDetail