import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import { Bell } from 'lucide-react';

export const NotificationsPage: React.FC = () => {
  const { data: notifications, isLoading } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => apiClient.getNotifications(),
  });

  return (
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <div className="max-w-2xl w-full">
        <header className="mb-8 flex items-center gap-3">
          <div className="p-3 bg-primary/10 text-primary rounded-full">
            <Bell size={24} />
          </div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Notifications</h1>
        </header>

        {isLoading ? (
          <div className="flex items-center justify-center py-12 text-muted-foreground">Loading notifications...</div>
        ) : !notifications?.items?.length ? (
          <div className="flex flex-col items-center justify-center py-16 text-muted-foreground bg-card rounded-2xl border border-border shadow-sm">
            <Bell size={48} className="mb-4 opacity-20" />
            <p className="text-lg font-medium">No notifications yet</p>
            <p className="text-sm">When you get notifications, they'll show up here.</p>
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            {notifications.items.map((notification) => (
              <div 
                key={notification.id} 
                className={`p-4 rounded-xl border transition-colors ${notification.read ? 'bg-card border-border' : 'bg-primary/5 border-primary/20'}`}
              >
                <div className="flex justify-between items-start gap-4">
                  <p className="text-foreground text-sm leading-relaxed">{notification.message}</p>
                  <span className="text-xs text-muted-foreground whitespace-nowrap">
                    {new Date(notification.createdAt).toLocaleDateString()}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
