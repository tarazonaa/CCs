import React, { useState } from 'react'
import SignupForm from './SignupForm'
import { motion, AnimatePresence } from 'framer-motion'
import LoginForm from './LoginForm'

const AuthPage = () => {
  const [mode, setMode] = useState<'login' | 'signup'>('login')

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-light via-secondary to-accent">
      <motion.div
        className="min-h-screen flex items-center justify-center px-4 py-12 sm:px-6 lg:px-8 backdrop-blur-sm"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.6 }}
      >
        <motion.div
          className="w-full max-w-md space-y-8"
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, type: 'spring', stiffness: 100 }}
        >
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="bg-surface/90 backdrop-blur-xl shadow-2xl rounded-2xl p-8 space-y-6"
          >
            <AnimatePresence mode="wait">
              {mode === 'login' ? (
                <LoginForm key="login" switchToSignup={() => setMode('signup')} />
              ) : (
                <SignupForm key="signup" switchToLogin={() => setMode('login')} />
              )}
            </AnimatePresence>
          </motion.div>
        </motion.div>
      </motion.div>
    </div>
  )
}

export default AuthPage
