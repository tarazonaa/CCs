import i18n from "i18next";
import ICU from "i18next-icu";
import { initReactI18next } from "react-i18next";
import en from "./en/translation.json";
import es from "./es/translation.json";

export const defaultNS = "translation";

export const resources = {
  en: { translation: en },
  es: { translation: es },
} as const;

i18n
  .use(initReactI18next)
  .use(ICU)
  .init({
    resources,
    fallbackLng: "en",
    supportedLngs: ["en", "es"],
    ns: ["translation"],
    defaultNS,
    lng: undefined,
    interpolation: {
      escapeValue: false,
    },
    detection: {
      order: [],
      caches: [],
    },
  });

export default i18n;
