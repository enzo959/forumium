import { createContext, useContext, useState, useEffect, ReactNode } from "react"

type User = {
  id: number
  username: string
  email: string
}

type LoginResponse = {
  user: User
  access_token: string
  refresh_token: string
}

type AuthContextType = {
  user: User | null
  accessToken: string | null
  login: (data: LoginResponse) => void
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | null>(null)

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null)
  const [accessToken, setAccessToken] = useState<string | null>(null)

  const login = (data: LoginResponse) => {
    setUser(data.user)
    setAccessToken(data.access_token)

    localStorage.setItem("access_token", data.access_token)
    localStorage.setItem("refresh_token", data.refresh_token)
    localStorage.setItem("user_id", data.user.id.toString())
  }

  const logout = async () => {
    const refreshToken = localStorage.getItem("refresh_token")
    const userId = localStorage.getItem("user_id")

    try {
      if (refreshToken && userId) {
        await fetch("/api/auth/logout", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            refresh_token: refreshToken,
            user_id: Number(userId),
          }),
        })
      }
    } catch (error) {
      console.error("Erreur logout :", error)
    }

    setUser(null)
    setAccessToken(null)
    localStorage.removeItem("access_token")
    localStorage.removeItem("refresh_token")
    localStorage.removeItem("user_id")
  }

  useEffect(() => {
    const refreshToken = localStorage.getItem("refresh_token")
    const userId = localStorage.getItem("user_id")

    if (!refreshToken || !userId) return

    const refresh = async () => {
      try {
        const res = await fetch("/api/auth/refresh", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            refresh_token: refreshToken,
            user_id: Number(userId),
          }),
        })

        if (!res.ok) throw new Error("Refresh failed")

        const data: LoginResponse = await res.json()
        login(data)
      } catch (error) {
        logout()
      }
    }

    refresh()
  }, [])

  const isAuthenticated = user !== null

  return (
    <AuthContext.Provider
      value={{
        user,
        accessToken,
        login,
        logout,
        isAuthenticated,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error("useAuth doit être utilisé dans un AuthProvider")
  }

  return context
}