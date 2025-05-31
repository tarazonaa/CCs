import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { motion, AnimatePresence } from 'framer-motion';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  
  const { login } = useAuth();

  const handleSubmit = async (e: { preventDefault: () => void; }) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);
    
    try {
      await login(email, password);
    } catch (err) {
      setError('Invalid email or password');
    } finally {
      setIsLoading(false);
    }
  };

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
          transition={{ 
            duration: 0.8,
            type: "spring",
            stiffness: 100 
          }}
        >
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="bg-surface/90 backdrop-blur-xl shadow-2xl rounded-2xl p-8 space-y-6"
          >
            <div className="text-center">
              <motion.div
                initial={{ opacity: 0, y: -20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.3 }}
              >
                <h2 className="text-3xl font-bold tracking-tight text-text-primary">
                  Welcome back
                </h2>
                <p className="mt-2 text-sm text-text-secondary">
                  Sign in to continue to your account
                </p>
              </motion.div>
            </div>

            <AnimatePresence mode="wait">
              {error && (
                <motion.div
                  className="bg-error/10 border border-error/20 rounded-lg px-4 py-3"
                  initial={{ opacity: 0, y: -10, height: 0 }}
                  animate={{ opacity: 1, y: 0, height: 'auto' }}
                  exit={{ opacity: 0, y: -10, height: 0 }}
                  transition={{ duration: 0.2 }}
                >
                  <p className="text-sm text-error text-center">{error}</p>
                </motion.div>
              )}
            </AnimatePresence>

            <motion.form 
              className="space-y-6"
              onSubmit={handleSubmit}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.3, delay: 0.4 }}
            >
              <div className="space-y-2">
                <label htmlFor="email" className="block text-sm font-medium text-text-primary">
                  Email address
                </label>
                <motion.input
                  whileFocus={{ scale: 1.01 }}
                  type="email"
                  id="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="block w-full rounded-lg border border-border bg-surface/50 px-4 py-3 text-text-primary placeholder:text-text-secondary/60 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
                  placeholder="Enter your email"
                  disabled={isLoading}
                />
              </div>

              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <label htmlFor="password" className="block text-sm font-medium text-text-primary">
                    Password
                  </label>
                  <motion.a
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                    href="#"
                    className="text-sm font-medium text-primary hover:text-primary-dark transition-colors"
                  >
                    Forgot password?
                  </motion.a>
                </div>
                <motion.input
                  whileFocus={{ scale: 1.01 }}
                  type="password"
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="block w-full rounded-lg border border-border bg-surface/50 px-4 py-3 text-text-primary placeholder:text-text-secondary/60 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
                  placeholder="Enter your password"
                  disabled={isLoading}
                />
              </div>

              <motion.button
                whileHover={{ scale: 1.01 }}
                whileTap={{ scale: 0.98 }}
                type="submit"
                disabled={isLoading}
                className={`relative w-full rounded-lg bg-gradient-to-r from-primary to-secondary py-3 shadow-lg transition-all duration-200 ${
                  isLoading ? 'opacity-80' : 'hover:shadow-primary/20 hover:shadow-xl'
                }`}
              >
                {isLoading ? (
                  <div className="flex items-center justify-center">
                    <motion.div
                      className="h-5 w-5 border-2 border-white/30 border-t-white rounded-full"
                      animate={{ rotate: 360 }}
                      transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                    />
                    <span className="ml-2">Signing in...</span>
                  </div>
                ) : (
                  'Sign in'
                )}
              </motion.button>
            </motion.form>

            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.3, delay: 0.5 }}
              className="text-center"
            >
              <p className="text-sm text-text-secondary">
                Don't have an account?{' '}
                <motion.a
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  href="#"
                  className="font-medium text-primary hover:text-primary-dark transition-colors"
                >
                  Sign up now
                </motion.a>
              </p>
            </motion.div>
          </motion.div>
        </motion.div>
      </motion.div>
    </div>
  );
}

export default Login;