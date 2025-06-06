import React from 'react';
import Sidebar from './SideBar';

interface AppLayoutProps {
  children: React.ReactNode;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  return (
    <div className="flex min-h-screen bg-background">
      <Sidebar />
      <div className="flex-1 ml-[80px]">
        {children}
      </div>
    </div>
  );
};

export default AppLayout;