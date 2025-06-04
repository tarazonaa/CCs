import type React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate, useParams } from 'react-router-dom'
import { AuthProvider, useAuth } from '@/contexts/AuthContext'
import Dashboard from '@/pages/Dashboard'
import './i18n/config'
import { LanguageWrapper } from './LanguageWrapper'
import NotFound from './pages/NotFound'
import AuthPage from './pages/AuthPage'

// Protected Route Component
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth()
  const { lng = 'en' } = useParams()
  return isAuthenticated ? <>{children}</> : <Navigate to={`/${lng}/login`} replace />
}

const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth()
  const { lng = 'en' } = useParams()
  return !isAuthenticated ? <>{children}</> : <Navigate to={`/${lng}/dashboard`} replace />
}

const App: React.FC = () => {
  return (
    <AuthProvider>
      <Router>
        <div className="min-h-screen bg-gray-50">
          <Routes>
            <Route
              path="/:lng/login"
              element={
                <LanguageWrapper>
                  <PublicRoute>
                    <AuthPage />
                  </PublicRoute>
                </LanguageWrapper>
              }
            />
            <Route
              path="/:lng/dashboard"
              element={
                <LanguageWrapper>
                  <ProtectedRoute>
                    <Dashboard />
                  </ProtectedRoute>
                </LanguageWrapper>
              }
            />
            <Route
              path="/"
              element={
                <Navigate
                  to={`/${navigator.language.startsWith('es') ? 'es' : 'en'}/dashboard`}
                  replace
                />
              }
            />

            <Route path="*" element={<NotFound />} />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  )
}

export default App
