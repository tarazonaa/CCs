import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import DrawingCanvas from '@/components/DrawingCanvas';
import DrawingHistory from '@/components/DrawingHistory';

export interface Drawing {
  id: string;
  imageData: string;
  prediction?: string;
  confidence?: number;
  timestamp: Date;
}

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const [drawings, setDrawings] = useState<Drawing[]>([]);
  const [activeTab, setActiveTab] = useState<'draw' | 'history'>('draw');

  const handleDrawingComplete = (imageData: string) => {
    const newDrawing: Drawing = {
      id: Date.now().toString(),
      imageData,
      timestamp: new Date(),
      prediction: '?',
      confidence: 0
    };
    
    setDrawings(prev => [newDrawing, ...prev]);
  };

  const clearHistory = () => {
    setDrawings([]);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
            <div className="flex items-center space-x-4">
              <span className="text-gray-600">Welcome, {user?.name}</span>
              <button
                onClick={logout}
                className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors duration-200"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      <nav className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex space-x-8">
            <button
              onClick={() => setActiveTab('draw')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors duration-200 ${
                activeTab === 'draw'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Draw
            </button>
            <button
              onClick={() => setActiveTab('history')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors duration-200 ${
                activeTab === 'history'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              History ({drawings.length})
            </button>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'draw' && (
          <DrawingCanvas onDrawingComplete={handleDrawingComplete} />
        )}
        {activeTab === 'history' && (
          <DrawingHistory drawings={drawings} onClearHistory={clearHistory} />
        )}
      </main>
    </div>
  );
};

export default Dashboard;