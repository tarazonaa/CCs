import DrawingCanvas from "@/components/DrawingCanvas";
import DrawingHistory from "@/components/DrawingHistory";
import { useAuth } from "@/contexts/AuthContext";
import {
  BracketsCurly,
  ClockCounterClockwise,
  PencilSimple,
  User,
} from "@phosphor-icons/react";
import axios from "axios";
import { AnimatePresence, motion } from "framer-motion";
import type React from "react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useLocation } from "react-router-dom";

export interface ImageMetadata {
  sent_image_id: string;
  received_image_id: string;
  created_at: Date;
}

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const [imagesData, setImagesData] = useState<ImageMetadata[]>([]);
  const [activeTab, setActiveTab] = useState<"draw" | "history">("draw");
  const location = useLocation();

  const { t } = useTranslation();

  useEffect(() => {
    // Check if URL hash is #history and update activeTab
    if (location.hash === "#history") {
      setActiveTab("history");
    } else {
      setActiveTab("draw");
    }
  }, [location.hash]);

  const handleDrawingComplete = (imageData: string) => {
    // const newDrawing: Drawing = {
    //   id: Date.now().toString(),
    //   imageData,
    //   timestamp: new Date(),
    //   prediction: '?',
    //   confidence: 0,
    // };
    // setDrawings((prev) => [newDrawing, ...prev]);
  };

  const clearHistory = () => {
    //   setDrawings([]);
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        duration: 0.3,
        when: "beforeChildren",
        staggerChildren: 0.1,
      },
    },
    exit: {
      opacity: 0,
      transition: { duration: 0.2 },
    },
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: { duration: 0.3 },
    },
    exit: {
      y: -20,
      opacity: 0,
      transition: { duration: 0.2 },
    },
  };

  useEffect(() => {
    axios
      .get(`${import.meta.env.VITE_AUTH_ENDPOINT}/api/v1/images`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("access_token")}`,
        },
      })
      .then((response) => {
        const fetchedImagesData: ImageMetadata[] = response.data.data.map(
          (item: ImageMetadata) => ({
            sent_image_id: item.sent_image_id,
            received_image_id: item.received_image_id,
            created_at: item.created_at,
          }),
        );
        setImagesData(fetchedImagesData);
      })
      .catch((error) => {
        console.error("Error fetching images data:", error);
      });
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <header className="pt-6 px-8">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
          className="flex justify-between items-start"
        >
          <div>
            <h1 className="text-3xl font-semibold text-text-primary">
              {t("dashboard_title")}
            </h1>
            <p className="text-text-secondary mt-1">
              {t("welcome_user", {
                name: user?.username || "artist",
              })}
            </p>
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={logout}
              className="flex items-center space-x-1 px-4 py-2 bg-primary text-black rounded-lg shadow-lg transition-all duration-300 hover:bg-primary/90 mt-3"
            >
              <User weight="bold" size={18} />
              <span>{t("logout")}</span>
            </motion.button>
          </div>

          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2 bg-surface/50 backdrop-blur-sm rounded-xl p-2 shadow-lg">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={() => setActiveTab("draw")}
                className={`relative flex items-center space-x-2 px-6 py-3 rounded-lg font-medium transition-all duration-300 ${
                  activeTab === "draw"
                    ? "bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-lg shadow-blue-500/30 scale-105"
                    : "text-text-secondary hover:text-text-primary hover:bg-background/80"
                }`}
              >
                <PencilSimple weight="bold" size={18} />
                <span>{t("tab_draw")}</span>
                {activeTab === "draw" && (
                  <motion.div
                    layoutId="activeTab"
                    className="absolute inset-0 bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg -z-10"
                    initial={false}
                    transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
                  />
                )}
              </motion.button>
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={() => setActiveTab("history")}
                className={`relative flex items-center space-x-2 px-6 py-3 rounded-lg font-medium transition-all duration-300 ${
                  activeTab === "history"
                    ? "bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-lg shadow-blue-500/30 scale-105"
                    : "text-text-secondary hover:text-text-primary hover:bg-background/80"
                }`}
              >
                <ClockCounterClockwise weight="bold" size={18} />
                <span>{t("tab_history", { count: imagesData.length })}</span>
                {activeTab === "history" && (
                  <motion.div
                    layoutId="activeTab"
                    className="absolute inset-0 bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg -z-10"
                    initial={false}
                    transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
                  />
                )}
              </motion.button>
            </div>
          </div>
        </motion.div>
      </header>

      <main className="max-w-7xl mx-auto px-8 py-8">
        <AnimatePresence mode="wait">
          {activeTab === "draw" ? (
            <motion.div
              key="draw"
              variants={containerVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
            >
              <motion.div variants={itemVariants}>
                <DrawingCanvas onDrawingComplete={handleDrawingComplete} />
              </motion.div>
            </motion.div>
          ) : (
            <motion.div
              key="history"
              variants={containerVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
            >
              <motion.div variants={itemVariants}>
                <DrawingHistory
                  imagesData={imagesData}
                  onClearHistory={clearHistory}
                />
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>
      </main>
    </div>
  );
};

export default Dashboard;
