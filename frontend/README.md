# Drawing Canvas Frontend
A modern React-based web application that provides an interactive drawing canvas for digit recognition with AI inference capabilities. Users can draw digits (0-9) and receive real-time AI predictions with confidence scores.
Features
Core Functionality

# Interactive Drawing Canvas: canvas-based drawing interface with smooth touch and mouse support
AI Integration: Real-time digit recognition using machine learning inference API
Drawing History: View and manage previously created drawings with their AI predictions
User Authentication: Secure login/signup system with token management
Multi-language Support: Internationalization (i18n) with English and Spanish support
Responsive Design: Approach with adaptive layouts

# User Experience

Modern UI/UX: Modern design with glassmorphism effects and smooth animations
Real-time Feedback: Instant visual feedback during drawing and prediction processes
Grid Overlay: Optional grid system for better drawing alignment and number segmentation
Canvas Controls: Clear, save, and grid toggle functionality
Drawing Tips: Built-in guidance for optimal digit drawing

# Tech Stack
Frontend Framework

React 18 with TypeScript for type safety
Vite for fast development and optimized builds
React Router for client-side routing with language-specific URLs

# Styling & Animation

Tailwind CSS with custom Apple-inspired design system
Framer Motion for smooth animations and transitions
Custom CSS Variables for consistent theming

# State Management & API

React Context API for authentication state
Axios for HTTP requests with interceptors
React i18next for internationalization

# UI Components & Icons

Phosphor Icons for consistent iconography
Notistack for toast notifications
Custom Components with motion and accessibility features

# Project Structure
```
src/
├── components/
│   ├── DrawingCanvas.tsx          # Main drawing interface
│   ├── DrawingHistory.tsx         # History grid view
│   ├── DrawingPredictionCard.tsx  # Individual drawing cards
│   └── layout/
│   └── Navbar.tsx                 # Navigation header
├── contexts/
│   └── AuthContext.tsx    
|
├── i18n/
│   ├── config.ts                  # i18n configuration
│   └── locales/
│       ├── en.json                # English translations
│       └── es.json                # Spanish translations        # Authentication state management
├── pages/
│   ├── Dashboard.tsx              # Main dashboard with tabs
│   ├── AuthPage.tsx               # Login/signup container
│   ├── LoginForm.tsx              # Login form component
│   ├── SignupForm.tsx             # Registration form
│   └── NotFound.tsx               # 404 error page
|
├── types/
│   └── auth.ts                    # TypeScript type definitions
└── styles/
|    └── globals.css                # Global styles and CSS variables
├── App.tsx 
```

# API Integration
Authentication Endpoints

POST /oauth2/token - User login with email/password
POST /oauth2/introspect - Token validation and refresh
POST /auth/register - New user registration
POST /auth/logout - User logout

Drawing & Inference

POST /api/v1/inference - Submit drawing for AI prediction
POST /api/v1/images - Save drawing and prediction results
GET /api/v1/images - Retrieve user's drawing history
GET /api/v1/images/blob/{id} - Fetch individual drawing blobs


# Authentication System
Comprehensive auth implementation:

UUID Token Management: Automatic token refresh and validation
Form Validation: Client-side validation with helpful error messages
Password Security: Enforced strong password requirements
Session Persistence: Secure localStorage-based session management


# Local Development
Prerequisites

Node.js 16+ and npm/yarn
Access to backend API endpoints

# Installation & Running
bash# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Language Support

English: Default language with complete translations
Spanish: Full Spanish localization support
URL Structure: Language-specific routes (/en/dashboard, /es/dashboard)

Adding Translations

Add translation keys to src/i18n/locales/en.json
Provide Spanish translations in src/i18n/locales/es.json
Use useTranslation hook in components: const { t } = useTranslation()
Reference translations: {t('translation_key')}

Performance Optimizations
Canvas Performance

Event Throttling: Optimized drawing event handlers
Canvas Caching: Efficient redraw strategies
Image Processing: Client-side preprocessing reduces server load

# Network Optimization

Request Batching: Combined image upload operations
Token Management: Automatic refresh prevents unnecessary re-authentication
Blob Caching: URL.createObjectURL for efficient image display

# Deployment
Production Build
npm run build

# Output will be in dist/ directory
# Serve static files from dist/

# Environment Configuration

VITE_PROVISION_KEY
VITE_AUTH_ENDPOINT
