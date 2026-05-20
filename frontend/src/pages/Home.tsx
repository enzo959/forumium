function Home() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-blue-600">Forumim</h1>
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
      </main>
    </div>
  )
}

export default Home