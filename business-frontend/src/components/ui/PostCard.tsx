import * as React from "react";
import { 
  MoreHorizontal, 
  Heart, 
  MessageCircle, 
  Share2, 
  Bookmark 
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

// Експортуємо тип, щоб інші сторінки могли його використовувати
export type PostData = {
  id: number;
  content: string;
  created_at: string;
  org: {
    id: number;
    name: string;
    profileType: string;
  };
  images: string[];
};

interface PostCardProps {
  post: PostData;
}

export function PostCard({ post }: PostCardProps) {
  const getInitials = (name: string) => {
    return name ? name.substring(0, 2).toUpperCase() : "SB";
  };

  const formatDate = (dateStr: string) => {
    if (!dateStr || dateStr.startsWith("0001")) return "Just now";
    
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 1) return "Just now";
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  return (
    <Card className="w-full max-w-[460px] bg-white dark:bg-[#112f26] border border-gray-200 dark:border-[#2f5e50]/60 shadow-2xl rounded-[2rem] overflow-hidden flex flex-col mx-auto">
      <CardContent className="p-0 flex flex-col">
        
        {/* Шапка поста */}
        <div className="px-4 py-3 sm:px-5 sm:py-4 flex items-center justify-between bg-transparent">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-tr from-emerald-500 to-[#1A3C34] rounded-full flex items-center justify-center text-white font-bold text-sm shadow-inner shrink-0 cursor-pointer hover:opacity-90 transition-opacity">
              {getInitials(post.org?.name)}
            </div>
            <div className="flex flex-col">
              <h3 className="text-[#1A3C34] dark:text-white font-bold text-[0.95rem] leading-none hover:underline cursor-pointer">
                {post.org?.name || "Your Venue"}
              </h3>
              <span className="text-gray-500 dark:text-gray-400 text-[0.8rem] mt-1 font-medium">
                {formatDate(post.created_at)}
              </span>
            </div>
          </div>
          <Button variant="ghost" size="icon" className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded-full h-8 w-8">
            <MoreHorizontal size={20} />
          </Button>
        </div>
        
        {/* Фотографії */}
        {post.images && post.images.length > 0 && (
          <div className={`w-full bg-gray-100 dark:bg-black/40 grid gap-0.5 ${post.images.length === 1 ? 'grid-cols-1' : 'grid-cols-2'}`}>
            {post.images.map((imgUrl, idx) => (
              <div key={idx} className="relative aspect-square sm:aspect-[4/5] w-full overflow-hidden">
                <img 
                  src={imgUrl} 
                  alt={`Post image ${idx + 1}`} 
                  className="w-full h-full object-cover"
                  onError={(e) => {
                    e.currentTarget.onerror = null;
                    e.currentTarget.src = "https://placehold.co/600x600/112f26/FFF?text=Image+Not+Found";
                  }}
                />
              </div>
            ))}
          </div>
        )}

        {/* Панель дій */}
        <div className="px-4 py-3 sm:px-5 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Heart size={26} className="text-gray-700 dark:text-gray-200 hover:text-rose-500 dark:hover:text-rose-500 cursor-pointer transition-colors stroke-[1.5]" />
            <MessageCircle size={26} className="text-gray-700 dark:text-gray-200 hover:text-emerald-500 cursor-pointer transition-colors stroke-[1.5]" />
            <Share2 size={24} className="text-gray-700 dark:text-gray-200 hover:text-emerald-500 cursor-pointer transition-colors stroke-[1.5]" />
          </div>
          <Bookmark size={24} className="text-gray-700 dark:text-gray-200 hover:text-emerald-500 cursor-pointer transition-colors stroke-[1.5]" />
        </div>

        {/* Опис поста */}
        <div className="px-4 pb-5 sm:px-5">
          <p className="text-[#1A3C34] dark:text-gray-100 text-[0.95rem] leading-relaxed">
            <span className="font-extrabold mr-2 cursor-pointer hover:underline">
              {post.org?.name || "Your Venue"}:
            </span>
            {post.content}
          </p>
        </div>
        
      </CardContent>
    </Card>
  );
}