import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

export const LanguageWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { lng } = useParams()
  const { i18n } = useTranslation()
  const [ready, setReady] = useState(false)

  useEffect(() => {
    if (!lng || !['en', 'es'].includes(lng)) {
      window.location.pathname = `/en${window.location.pathname.replace(/^\/[^/]+/, '')}`
      return
    }

    if (i18n.language !== lng) {
      i18n.changeLanguage(lng).then(() => setReady(true))
    } else {
      setReady(true)
    }
  }, [lng, i18n])

  if (!ready) return null

  return <>{children}</>
}
