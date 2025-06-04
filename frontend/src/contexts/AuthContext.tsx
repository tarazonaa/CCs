import axios from 'axios'
import type React from 'react'
import { createContext, useContext, useEffect, useState } from 'react'
import { useParams } from 'react-router'

const authEndpoint = import.meta.env.VITE_AUTH_ENDPOINT
const provisionKey = import.meta.env.VITE_PROVISION_KEY

interface User {
  id: string
  email: string
  name: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  refreshToken: () => Promise<void>
  logout: () => void
  loading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within an AuthProvider')
  return context
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const { lang } = useParams<{ lang: string }>()

  useEffect(() => {
    const checkSession = async () => {
      const token = localStorage.getItem('access_token')
      if (!token) return setLoading(false)

      try {
        const { data } = await axios.post(`${authEndpoint}/oauth2/introspect`, { token })
        if (data.active) {
          setUser({
            id: data.authenticated_userid,
            email: data.email ?? 'unknown@email.com',
            name: data.username ?? 'unknown',
          })
        } else {
          localStorage.removeItem('access_token')
        }
      } catch (err) {
        console.error('Token validation failed:', err)
        localStorage.removeItem('access_token')
      } finally {
        setLoading(false)
      }
    }

    checkSession()
  }, [])

  /*
   *  '{
    "grant_type": "password",
    "client_id": "CCs-client-id",
    "client_secret": "holajorge",
        "email": "andres.tara.so@gmail.com",
        "password": "holaJorge@123",
        "scope": "read write"
  }'
   * */

  const login = async (email: string, password: string) => {
    await axios
      .post(`${authEndpoint}/oauth2/token`, {
        client_id: 'CCs-client-id',
        client_secret: 'holajorge',
        grant_type: 'password',
        provision_key: provisionKey,
        scope: 'read write',
        email,
        password,
      })
      .then((res) => {
        localStorage.setItem('access_token', res.data.access_token)
        localStorage.setItem('refresh_token', res.data.refresh_token)
      })
      .finally(() => {
        const currLang = lang || 'en'
        window.location.href = `/${currLang}/dashboard`
      })
  }
  const refreshToken = async () => {
    const refreshToken = localStorage.getItem('refresh_token')
    if (!refreshToken) {
      logout() // No refresh token, force logout
      return null
    }

    try {
      const response = await axios.post(`${authEndpoint}/oauth2/token`, {
        grant_type: 'refresh_token',
        refresh_token: refreshToken,
        client_id: 'CCs-client-id',
        client_secret: 'holajorge',
      })

      // Store new tokens
      localStorage.setItem('access_token', response.data.access_token)
      if (response.data.refresh_token) {
        localStorage.setItem('refresh_token', response.data.refresh_token)
      }

      return response.data.access_token
    } catch (error) {
      console.error('Token refresh failed:', error)
      logout() // Refresh failed, force logout
      return null
    }
  }

  const logout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setUser(null)
  }

  const value = {
    user,
    isAuthenticated: !!user,
    login,
    logout,
    loading,
    refreshToken,
  }

  return loading ? (
    <div className="flex justify-center items-center min-h-screen">
      <div className="animate-spin h-10 w-10 border-b-2 border-blue-500 rounded-full" />
    </div>
  ) : (
    <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
  )
}
