import type React from 'react'
import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { motion, AnimatePresence } from 'framer-motion'
import { useTranslation } from 'react-i18next'
import { useNavigate, useParams } from 'react-router'

interface SignupFormProps {
  switchToLogin: () => void
}

const SignupForm: React.FC<SignupFormProps> = ({ switchToLogin }) => {
  const { t } = useTranslation()
  const { register } = useAuth()

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [username, setUsername] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)

  const { lang } = useParams()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    try {
      if (password !== confirmPassword) {
        throw new Error(t('passwords_do_not_match'))
      }
      await register(username, name, email, password)
      setIsLoading(false)
      const currLang = lang || 'en'
      navigate(`/${currLang}/dashboard`)
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message || t('signup_failed'))
      }
    } finally {
      setIsLoading(false)
    }
  }

  const EyeIcon = () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
      />
    </svg>
  )

  const EyeOffIcon = () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"
      />
    </svg>
  )

  return (
    <>
      <div className="text-center">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3 }}
        >
          <h2 className="text-3xl font-bold tracking-tight text-text-primary">
            {t('create_account')}
          </h2>
          <p className="mt-2 text-sm text-text-secondary">{t('signup_to_continue')}</p>
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
          <label htmlFor="name" className="block text-sm font-medium text-text-primary">
            {t('username')}
          </label>
          <motion.input
            whileFocus={{ scale: 1.01 }}
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            className="block w-full rounded-lg border border-border bg-surface/50 px-4 py-3 text-text-primary placeholder:text-text-secondary/60 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
            placeholder={t('username_placeholder')}
            disabled={isLoading}
          />
        </div>
        <div className="space-y-2">
          <label htmlFor="name" className="block text-sm font-medium text-text-primary">
            {t('full_name')}
          </label>
          <motion.input
            whileFocus={{ scale: 1.01 }}
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            className="block w-full rounded-lg border border-border bg-surface/50 px-4 py-3 text-text-primary placeholder:text-text-secondary/60 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
            placeholder={t('full_name_placeholder')}
            disabled={isLoading}
          />
        </div>

        <div className="space-y-2">
          <label htmlFor="email" className="block text-sm font-medium text-text-primary">
            {t('email_address')}
          </label>
          <motion.input
            whileFocus={{ scale: 1.01 }}
            type="email"
            id="email"
            pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="peer invalid:[&:not(:placeholder-shown):not(:focus)]:border-red-500 block w-full rounded-lg border border-border bg-surface/50 px-4 py-3 text-text-primary placeholder:text-text-secondary/60 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
            placeholder="Enter your email"
            disabled={isLoading}
          />
          <span className="mt-2 hidden text-sm text-red-500 peer-[&:not(:placeholder-shown):not(:focus):invalid]:block">
            {t('invalid_email')}
          </span>
        </div>

        <div className="space-y-2">
          <label htmlFor="password" className="block text-sm font-medium text-text-primary">
            {t('password')}
          </label>
          <div className="relative">
            <motion.input
              whileFocus={{ scale: 1.01 }}
              type={showPassword ? 'text' : 'password'}
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@!#$%^&*])[A-Za-z\d@!#$%^&*]{8,}$"
              className="peer block w-full invalid:[&:not(:placeholder-shown):not(:focus)]:border-red-500 rounded-lg border border-border bg-surface/50 px-4 py-3 pr-12 text-text-primary placeholder:text-text-secondary/60 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="*********************"
              disabled={isLoading}
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute inset-y-0 right-0 flex items-center pr-3 text-text-secondary hover:text-text-primary transition-colors"
              disabled={isLoading}
            >
              {showPassword ? <EyeOffIcon /> : <EyeIcon />}
            </button>
          </div>
          <span className="mt-2 hidden text-sm text-red-500 peer-[&:not(:placeholder-shown):not(:focus):invalid]:block">
            {t('password_requirements_intro')}
            <ul className="list-disc list-inside mt-1 space-y-1">
              <li>{t('requirement_uppercase')}</li>
              <li>{t('requirement_8_characters')}</li>
              <li>{t('requirement_lowercase')}</li>
              <li>{t('requirement_number')}</li>
              <li>{t('requirement_special')}</li>
            </ul>
          </span>
        </div>

        <div className="space-y-2">
          <label htmlFor="confirmPassword" className="block text-sm font-medium text-text-primary">
            {t('confirm_password')}
          </label>
          <div className="relative">
            <motion.input
              whileFocus={{ scale: 1.01 }}
              type={showConfirmPassword ? 'text' : 'password'}
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              className="block w-full rounded-lg border border-border bg-surface/50 px-4 py-3 pr-12 text-text-primary placeholder:text-text-secondary/60 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="*********************"
              disabled={isLoading}
            />
            <button
              type="button"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              className="absolute inset-y-0 right-0 flex items-center pr-3 text-text-secondary hover:text-text-primary transition-colors"
              disabled={isLoading}
            >
              {showConfirmPassword ? <EyeOffIcon /> : <EyeIcon />}
            </button>
          </div>
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
                transition={{ duration: 1, repeat: Number.POSITIVE_INFINITY, ease: 'linear' }}
              />
              <span className="ml-2">{t('signing_up')}</span>
            </div>
          ) : (
            t('sign_up_button')
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
          {t('already_have_account')}{' '}
          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            onClick={switchToLogin}
            className="font-medium text-primary hover:text-primary-dark transition-colors"
          >
            {t('sign_in_button')}
          </motion.button>
        </p>
      </motion.div>
    </>
  )
}

export default SignupForm
