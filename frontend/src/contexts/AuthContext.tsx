import axios from 'axios'
import { useSnackbar } from 'notistack'
import { useTranslation } from 'react-i18next'
import type React from 'react'
import { createContext, useCallback, useContext, useEffect, useState } from 'react'

const authEndpoint = import.meta.env.VITE_API_URL
const provisionKey = import.meta.env.VITE_PROVISION_KEY

interface User {
  id: string
  email: string
  name: string
  username: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  refreshToken: () => Promise<void>
  register: (username: string, name: string, email: string, password: string) => Promise<void>
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

  const { enqueueSnackbar } = useSnackbar()
  const { t } = useTranslation()

  const checkSession = useCallback(async () => {
    const token = localStorage.getItem('access_token')
    if (!token) return setLoading(false)

    try {
      const { data } = await axios.post(`${authEndpoint}/oauth2/introspect`, { token })
      if (data.should_refresh) {
        const newToken = await refreshToken()
        if (newToken) {
          setUser({
            id: data.authenticated_userid,
            email: data.email,
            name: data.name,
            username: data.username,
          })
        }
      } else if (data.active) {
        setUser({
          id: data.authenticated_userid,
          email: data.email,
          name: data.name,
          username: data.username,
        })
      } else {
        localStorage.removeItem('access_token')
      }
      return data
    } catch (err) {
      console.error('Token validation failed:', err)
      localStorage.removeItem('access_token')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    checkSession()
  }, [checkSession])

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
      .then(async (res) => {
        localStorage.setItem('access_token', res.data.access_token)
        localStorage.setItem('refresh_token', res.data.refresh_token)
        const user = await checkSession()
        enqueueSnackbar(`${t('welcome')}, ${user?.username}`, {
          variant: 'success',
        })
      })
  }
  const refreshToken = async () => {
    const refreshToken = localStorage.getItem('refresh_token')
    if (!refreshToken) {
      logout()
      return null
    }

    try {
      const response = await axios.post(`${authEndpoint}/oauth2/token`, {
        grant_type: 'refresh_token',
        refresh_token: refreshToken,
        client_id: 'CCs-client-id',
        client_secret: 'holajorge',
      })

      localStorage.setItem('access_token', response.data.access_token)
      if (response.data.refresh_token) {
        localStorage.setItem('refresh_token', response.data.refresh_token)
      }

      return response.data.access_token
    } catch (error) {
      console.error('Token refresh failed:', error)
      logout()
      return null
    }
  }

  const register = async (username: string, name: string, email: string, password: string) => {
    const response = await axios.post(`${authEndpoint}/auth/register`, {
      email,
      username,
      name,
      password,
    })

    if (response.status === 201) {
      await login(email, password)
    }
  }

  const logout = async () => {
    const access_token = localStorage.getItem('access_token')
    const response = await axios.post(
      `${authEndpoint}/auth/logout`,
      {},
      {
        headers: {
          Authorization: `Bearer ${access_token}`,
        },
      }
    )
    if (response.status === 200) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      setUser(null)
    }
  }

  const value = {
    user,
    isAuthenticated: !!user,
    login,
    logout,
    loading,
    refreshToken,
    register,
  }

  return loading ? (
    <div className="flex justify-center items-center min-h-screen">
      <div className="animate-spin h-10 w-10 border-b-2 border-blue-500 rounded-full" />
    </div>
  ) : (
    <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
  )
}
