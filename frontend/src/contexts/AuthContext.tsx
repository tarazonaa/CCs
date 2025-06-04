// contexts/AuthContext.tsx
import axios from 'axios'
import { randomUUID } from 'node:crypto'
import type React from 'react'
import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

const authEndpoint = process.env.AUTH_ENDPOINT

interface User {
  id: string
  email: string
  name: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<boolean>
  logout: () => void
  loading: boolean
  signup: (email: string, password: string, username: string) => Promise<boolean>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Check for existing session on app load
    const checkAuth = async () => {
      try {
        const token = localStorage.getItem('auth_token')
        if (token) {
          // TODO: Validate token with your Go backend
          // For now, we'll simulate a logged-in user
          setUser({
            id: '1',
            email: 'user@example.com',
            name: 'Andres',
          })
        }
      } catch (error) {
        console.error('Auth check failed:', error)
        localStorage.removeItem('auth_token')
      } finally {
        setLoading(false)
      }
    }

    checkAuth()
  }, [])

  const login = async (email: string, password: string): Promise<boolean> => {
    try {
      if (email && password) {
        const consumer_id = await axios
          .post(`${authEndpoint}/admin/consumers`, {
            username: email,
            custom_id: randomUUID(),
          })
          .then((response) => response.id)
          .catch((err) => console.log(`There was an error fetching the consumer_id: ${err}`))

        const client = await axios
          .post(`${authEndpoint}/admin/clients`, {
            name: 'CCs',
            redirect_uri: '/dashboard',
            consumer_id,
          })
          .then((response) => response.id)
          .catch((err) => console.log(`There was an error fetching the client: ${err}`))

        const auth_token = await axios.post(`${authEndpoint}/oauth2/tokens`, {
          credential: {
            id: client,
          },
          token_type: 'bearer',
          expires_in: 7200,
          scope: 'read write',
          authenticated_userid: consumer_id,
        })
        setUser(user)
        localStorage.setItem('auth_token', 'mock_token_123')
        return true
      }
      return false
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  const logout = () => {
    setUser(null)
    localStorage.removeItem('auth_token')
  }

  const signup = async () => {
    return true
  }

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    login,
    logout,
    loading,
    signup,
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    )
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
