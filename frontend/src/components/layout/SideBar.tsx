import React from 'react';
import { useNavigate, useLocation, useParams } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useAuth } from '@/contexts/AuthContext';
import {
  House,
  ClockCounterClockwise,
  User,
  SignOut,
  Palette,
  Moon,
  Sun,
  Images
} from '@phosphor-icons/react';

const Sidebar = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { lng = 'en' } = useParams();
  const { user, logout } = useAuth();
  const [isDarkMode, setIsDarkMode] = React.useState(false);

  const profileImage = '/css.png'; // Path to the profile image

  const navItems = [
    { 
      name: 'Dashboard', 
      path: `/${lng}/dashboard`, 
      icon: <House weight="fill" size={24} />
    },
    { 
      name: 'History', 
      path: '#history', // This will be handled in the dashboard
      icon: <ClockCounterClockwise weight="fill" size={24} /> 
    },
    { 
      name: 'Gallery', 
      path: `/${lng}/gallery`, 
      icon: <Images weight="fill" size={24} /> 
    },
    { 
      name: 'Theme', 
      path: `/${lng}/theme`, 
      icon: <Palette weight="fill" size={24} /> 
    },
  ];

  const isActive = (path: string) => {
    if (path === '#history') {
      return location.pathname.includes('dashboard') && location.hash === '#history';
    }
    return location.pathname === path;
  };

  const handleNav = (path: string) => {
    if (path === '#history') {
      // Handle the history tab in the dashboard
      if (location.pathname.includes('dashboard')) {
        window.location.hash = 'history';
      } else {
        navigate(`/${lng}/dashboard`);
        setTimeout(() => {
          window.location.hash = 'history';
        }, 100);
      }
    } else {
      navigate(path);
    }
  };

  const toggleDarkMode = () => {
    if (isDarkMode) {
      document.documentElement.classList.remove('dark');
    } else {
      document.documentElement.classList.add('dark');
    }
    setIsDarkMode(!isDarkMode);
  };

  const handleLogout = () => {
    logout();
    navigate(`/${lng}/login`);
  };

  return (
    <aside className="fixed left-0 top-0 h-full w-[80px] bg-surface/80 backdrop-blur-lg border-r border-border/50 flex flex-col items-center py-8 z-10">
      <div className="mb-10">
        <motion.div 
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          className="relative w-12 h-12 rounded-full overflow-hidden border-2 border-primary/30"
        >
          <img 
            src={profileImage} 
            alt="Profile" 
            className="w-full h-full object-cover"
            onError={(e) => {
              const target = e.target as HTMLImageElement;
              target.src = `https://ui-avatars.com/api/?name=${user?.name || 'User'}&background=0D8ABC&color=fff`;
            }}
          />
        </motion.div>
      </div>

      <nav className="flex flex-col items-center space-y-6 flex-1">
        {navItems.map((item) => (
          <motion.button
            key={item.name}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => handleNav(item.path)}
            className={`relative group w-12 h-12 rounded-xl flex items-center justify-center transition-all duration-300 ${
              isActive(item.path) 
                ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-xl shadow-blue-500/40 scale-110 border-2 border-blue-300' 
                : 'text-text-secondary hover:bg-surface hover:text-text-primary hover:scale-105'
            }`}
          >
            {item.icon}
            {/* Active indicator elements */}
            {isActive(item.path) && (
              <>
                {/* Left border indicator */}
                <motion.div
                  initial={{ scaleY: 0 }}
                  animate={{ scaleY: 1 }}
                  className="absolute -left-6 top-1/2 -translate-y-1/2 w-1 h-8 bg-blue-500 rounded-r-full"
                />
                {/* Pulsing effect */}
                <motion.div
                  animate={{ scale: [1, 1.2, 1] }}
                  transition={{ duration: 2, repeat: Infinity }}
                  className="absolute inset-0 rounded-xl bg-blue-400/20"
                />
              </>
            )}
            <span className="absolute left-14 px-2 py-1 rounded-md bg-surface shadow-apple text-text-primary text-xs whitespace-nowrap opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200">
              {item.name}
            </span>
          </motion.button>
        ))}
      </nav>

      <div className="flex flex-col items-center space-y-4 mb-6">
        <motion.button
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
          onClick={toggleDarkMode}
          className="w-10 h-10 rounded-full flex items-center justify-center text-text-secondary hover:bg-background transition-colors duration-200"
        >
          {isDarkMode ? <Sun size={22} /> : <Moon size={22} />}
        </motion.button>
        
        <motion.button
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
          onClick={handleLogout}
          className="w-10 h-10 rounded-full flex items-center justify-center text-error hover:bg-error/10 transition-colors duration-200"
        >
          <SignOut size={22} />
        </motion.button>
      </div>
    </aside>
  );
};

export default Sidebar;