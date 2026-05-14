import React, { useState, useRef, useEffect } from 'react';
import { Bell } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

export const NotificationBell: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const { data: notifications } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => apiClient.getNotifications(),
    refetchInterval: 30000,
    refetchIntervalInBackground: false,
    enabled: !!localStorage.getItem('token'),
  });

  const unreadCount = notifications?.items?.filter((n: any) => !n.read).length || 0;

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className="relative flex items-center gap-2 bg-accent text-accent-foreground hover:bg-[#e6c200] dark:hover:bg-accent/90 font-bold px-4 py-2 rounded-full shadow-md dark:shadow-lg dark:shadow-accent/20 transition-all hover:scale-105"
      >
        <Bell size={18} />
        <span className="text-sm">Notifications</span>
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 bg-destructive text-destructive-foreground text-[10px] font-black w-5 h-5 flex items-center justify-center rounded-full border-2 border-background">
            {unreadCount > 99 ? '99+' : unreadCount}
          </span>
        )}
      </button>

      {/* Dropdown Content */}
      {isOpen && (
        <div className="absolute right-0 mt-3 w-80 bg-popover border border-border rounded-3xl shadow-2xl overflow-hidden z-50 transform origin-top-right transition-all animate-in fade-in slide-in-from-top-4">
          <div className="p-5 border-b border-border bg-muted/30 backdrop-blur-md flex justify-between items-center">
            <h3 className="text-2xl font-serif font-bold text-accent">Latest Stories</h3>
            <span className="text-xs font-bold text-accent bg-accent/10 px-2 py-1 rounded-full">{unreadCount} New</span>
          </div>
          <div className="max-h-80 overflow-y-auto">
            {(!notifications?.items?.length) ? (
              <div className="p-8 text-center text-muted-foreground">
                <Bell className="mx-auto mb-3 opacity-20" size={32} />
                <p className="font-semibold text-sm">No new notifications</p>
                <p className="text-xs mt-1 opacity-70">You're all caught up!</p>
              </div>
            ) : (
              notifications.items.map((notification: any) => (
                <Link key={notification.id} to="/notifications" className="block p-4 border-b border-border hover:bg-muted/40 transition-colors group">
                  <div className="flex gap-3 items-start">
                    <div className="w-8 h-8 rounded-full bg-accent/10 flex items-center justify-center flex-shrink-0 mt-0.5">
                      <span className="text-accent font-bold text-xs">
                        {notification.type === 'like' ? '♥' : notification.type === 'comment' ? '💬' : '🔔'}
                      </span>
                    </div>
                    <div>
                      <p className="text-sm text-foreground font-medium leading-tight group-hover:text-primary transition-colors">
                        {notification.content || "You have a new notification"}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1 font-semibold">
                        {new Date(notification.createdAt).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                </Link>
              ))
            )}
          </div>
          <Link to="/notifications" className="block w-full text-center p-3 text-sm font-bold text-accent hover:bg-accent hover:text-accent-foreground transition-colors bg-muted/80">
            View All Notifications
          </Link>
        </div>
      )}
    </div>
  );
};
