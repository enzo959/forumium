import { useState, useEffect } from 'react'
import { useParams, useLocation } from 'react-router-dom'
import axios from '../api/axios'
import Navbar from '../components/Navbar'
import { useAuth } from '../context/AuthContext'

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
  const { isAuthenticated, user } = useAuth()
  const [post, setPost] = useState<Post | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [commentError, setCommentError] = useState('')
  const [commentLoading, setCommentLoading] = useState(false)
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

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault()
    setCommentError('')

    if (!newComment.trim()) {
      setCommentError('Le commentaire ne peut pas être vide.')
      return
    }

    setCommentLoading(true)

    const optimisticComment: Comment = {
      id: Date.now(),
      content: newComment,
      created_at: new Date().toISOString(),
      author: {
        id: user?.id || 0,
        username: user?.username || '',
        avatar: '',
      },
    }
    setComments((prev) => [...prev, optimisticComment])
    setNewComment('')

    try {
      const res = await axios.post(`/api/posts/${id}/comments`, { content: optimisticComment.content })
      setComments((prev) =>
        prev.map((c) => (c.id === optimisticComment.id ? { ...optimisticComment, id: res.data.id } : c))
      )
    } catch {
      setComments((prev) => prev.filter((c) => c.id !== optimisticComment.id))
      setCommentError("Erreur lors de l'envoi du commentaire.")
    } finally {
      setCommentLoading(false)
    }
  }

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

        <div className="bg-white rounded shadow-sm p-6">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">
            Commentaires ({comments.length})
          </h2>

          {isAuthenticated ? (
            <form onSubmit={handleSubmitComment} className="mb-6">
              <textarea
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                placeholder="Écrire un commentaire..."
                rows={3}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 mb-2"
              />
              {commentError && <p className="text-red-500 text-sm mb-2">{commentError}</p>}
              <button
                type="submit"
                disabled={commentLoading}
                className="px-4 py-2 text-sm text-white bg-blue-600 rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {commentLoading ? 'Envoi...' : 'Commenter'}
              </button>
            </form>
          ) : (
            <p className="text-sm text-gray-400 mb-6">
              <a href="/login" className="text-blue-600 hover:underline">Connectez-vous</a> pour commenter.
            </p>
          )}

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