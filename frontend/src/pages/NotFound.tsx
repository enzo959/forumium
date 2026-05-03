import { Link } from 'react-router-dom'
 
function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">404</h1>
        <Link to="/" className="text-blue-600 underline">Retour à l'accueil</Link>
      </div>
    </div>
  )
}
 
export default NotFound
 