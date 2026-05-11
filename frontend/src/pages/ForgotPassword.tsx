import { useState } from 'react'
import axios from '../api/axios'
 
function ForgotPassword() {
  const [email, setEmail] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
 
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setMessage('')
 
    try {
      await axios.post('/api/auth/forgot-password', { email })
      setMessage('Si cet email existe, un lien de reset a été envoyé.')
    } catch {
      setError('Une erreur est survenue, réessayez.')
    } finally {
      setLoading(false)
    }
  }
 
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="bg-white p-8 rounded shadow w-full max-w-md">
        <h2 className="text-2xl font-bold text-gray-800 mb-2">Mot de passe oublié</h2>
        <p className="text-gray-500 text-sm mb-6">
          Entrez votre email, vous recevrez un lien pour réinitialiser votre mot de passe.
        </p>
 
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              placeholder="votre@email.com"
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
 
          {message && (
            <p className="text-green-600 text-sm mb-4">{message}</p>
          )}
 
          {error && (
            <p className="text-red-500 text-sm mb-4">{error}</p>
          )}
 
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50 text-sm font-medium"
          >
            {loading ? 'Envoi en cours...' : 'Envoyer le lien'}
          </button>
        </form>
 
        <p className="text-center text-sm text-gray-500 mt-4">
          <a href="/login" className="text-blue-600 hover:underline">
            Retour à la connexion
          </a>
        </p>
      </div>
    </div>
  )
}
 
export default ForgotPassword