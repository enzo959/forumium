import { useAuth } from '../context/AuthContext'

function Navbar() {
  const { user, isAuthenticated, logout } = useAuth()

  return (
    <header className="bg-white shadow-sm">
      <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
        <a href="/" className="text-2xl font-bold text-blue-600">Forum</a>
        <div className="flex gap-3 items-center">
          {isAuthenticated ? (
            <>
              <span className="text-sm text-gray-600 font-medium">
                {user?.username}
              </span>
              <button
                onClick={logout}
                className="px-4 py-2 text-sm text-red-600 border border-red-600 rounded hover:bg-red-50"
              >
                Déconnexion
              </button>
            </>
          ) : (
            <>
              <a href="/login" className="px-4 py-2 text-sm text-blue-600 border border-blue-600 rounded hover:bg-blue-50">
                Connexion
              </a>
              <a href="/register" className="px-4 py-2 text-sm text-white bg-blue-600 rounded hover:bg-blue-700">
                Inscription
              </a>
            </>
          )}
        </div>
      </div>
    </header>
  )
}

export default Navbar