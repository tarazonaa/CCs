import type React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate, useParams } from 'react-router-dom'
import { AuthProvider, useAuth } from '@/contexts/AuthContext'
import Dashboard from '@/pages/Dashboard'
import './i18n/config'
import { LanguageWrapper } from './LanguageWrapper'
import NotFound from './pages/NotFound'
import AuthPage from './pages/AuthPage'
import { SnackbarProvider } from 'notistack'
import AppLayout from '@/components/layout/AppLayout'
import { motion } from 'framer-motion'

// Protected Route Component
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const { lng = 'en' } = useParams();
  return isAuthenticated ? (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.3 }}
    >
      {children}
    </motion.div>
  ) : (
    <Navigate to={`/${lng}/login`} replace />
  );
};

const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const { lng = 'en' } = useParams();
  return !isAuthenticated ? (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.3 }}
    >
      {children}
    </motion.div>
  ) : (
    <Navigate to={`/${lng}/dashboard`} replace />
  );
};

const App: React.FC = () => {
  return (
    <SnackbarProvider
      maxSnack={3}
      anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      autoHideDuration={3000}
    >
      <AuthProvider>
        <Router>
          <div className="min-h-screen bg-background">
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
    </SnackbarProvider>
  )
}

export default App;