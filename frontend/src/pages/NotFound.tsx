import React from 'react'
import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-light via-secondary to-accent">
      <motion.div
        className="min-h-screen flex items-center justify-center px-4 py-12 sm:px-6 lg:px-8 backdrop-blur-sm"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.6 }}
      >
        <motion.div
          className="w-full max-w-md space-y-6 bg-surface/90 backdrop-blur-xl shadow-2xl rounded-2xl p-8 text-center"
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{
            duration: 0.8,
            type: 'spring',
            stiffness: 100,
          }}
        >
          <motion.h1
            className="text-4xl font-bold text-text-primary"
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ delay: 0.3, duration: 0.5 }}
          >
            404
          </motion.h1>
          <motion.p
            className="text-text-secondary text-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.4 }}
          >
            Oops! The page you're looking for doesn't exist.
          </motion.p>
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.5 }}
          >
            <Link
              to="/"
              className="inline-block mt-4 rounded-lg bg-gradient-to-r from-primary to-secondary px-5 py-3 text-black font-medium shadow-lg transition-transform hover:scale-105"
            >
              Go Home
            </Link>
          </motion.div>
        </motion.div>
      </motion.div>
    </div>
  )
}
