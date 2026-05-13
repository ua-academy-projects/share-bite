import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import { Search, MapPin } from 'lucide-react';
import { Link } from 'react-router-dom';

export const ExplorePage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  
  const { data: venues, isLoading } = useQuery({
    queryKey: ['explore', 'nearby'],
    queryFn: () => apiClient.getExploreNearby(50.4501, 30.5234, 20), // Defaulting to Kyiv coordinates
  });

  const filteredVenues = venues?.filter(v => 
    v.name.toLowerCase().includes(searchQuery.toLowerCase())
  ) || [];

  return (
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <div className="max-w-4xl w-full">
        <header className="mb-10 text-center">
          <h1 className="text-5xl font-serif font-bold tracking-tight text-foreground mb-3">Explore</h1>
          <p className="text-muted-foreground text-lg font-medium">Discover trending places near you</p>
        </header>

        <div className="relative max-w-xl mx-auto mb-12">
          <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
            <Search className="h-5 w-5 text-muted-foreground" />
          </div>
          <input
            type="text"
            className="block w-full pl-12 pr-4 py-4 border border-border/50 rounded-full leading-5 bg-card/80 backdrop-blur-sm placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary sm:text-base shadow-sm transition-all"
            placeholder="Search for amazing venues..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center py-16 text-muted-foreground font-serif text-lg">Finding nearby venues...</div>
        ) : filteredVenues.length === 0 ? (
          <div className="text-center py-20 bg-card rounded-3xl border border-border shadow-sm">
            <p className="text-2xl font-serif text-muted-foreground">No venues found matching "{searchQuery}"</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-8">
            {filteredVenues.map((venue) => (
              <Link 
                key={venue.id} 
                to={`/restaurant/${venue.id}`}
                className="group flex flex-col bg-card rounded-3xl border border-border/50 shadow-sm overflow-hidden hover:shadow-xl dark:hover:shadow-primary/5 hover:-translate-y-1 transition-all duration-300"
              >
                <div className="h-40 bg-muted relative overflow-hidden">
                  {venue.avatar ? (
                    <img src={venue.avatar} alt={venue.name} className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center bg-primary/10 text-primary">
                      <MapPin size={32} />
                    </div>
                  )}
                  {/* Subtle Gradient Overlay */}
                  <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-80"></div>
                </div>
                <div className="p-5 flex flex-col flex-1 relative bg-card z-10">
                  <h3 className="font-bold text-xl text-foreground group-hover:text-primary transition-colors line-clamp-1">{venue.name}</h3>
                  <div className="flex items-center gap-1.5 mt-auto pt-4 text-sm font-semibold text-muted-foreground uppercase tracking-wider">
                    <MapPin size={14} />
                    <span>Nearby</span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
