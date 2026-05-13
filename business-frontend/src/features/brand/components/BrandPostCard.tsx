import type { BrandPost } from "@/api/business";
import { Link } from "react-router-dom";
import { ExternalLink, Calendar } from "lucide-react";

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "2-digit",
    year: "numeric"
  }).format(date);
}

type BrandPostCardProps = {
  post: BrandPost;
};

export function BrandPostCard({ post }: BrandPostCardProps) {
  const image = post.images?.[0];
  const formattedDate = formatDate(post.created_at);
  
  // If the post belongs to a venue, link to that venue's profile
  const linkTo = post.org?.profileType === "venue" ? `/venue/${post.org.id}` : null;

  const content = (
    <article className="group relative flex h-full flex-col overflow-hidden rounded-[32px] border border-white/10 bg-[#0f1b17]/60 backdrop-blur-sm transition-all duration-300 hover:border-[#98FF98]/30 hover:bg-[#0f1b17]/80 hover:shadow-[0_20px_40px_-15px_rgba(0,0,0,0.4)]">
      <div className="relative aspect-[4/3] overflow-hidden">
        {image ? (
          <img 
            src={image} 
            alt={post.org?.name || "Brand post"} 
            className="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110" 
          />
        ) : (
          <div className="h-full w-full bg-[radial-gradient(circle_at_top,_#2a5b4c_0%,_#0b0f0e_70%)]" />
        )}
        
        {/* Overlays */}
        <div className="absolute inset-0 bg-gradient-to-t from-[#0b0f0e] via-[#0b0f0e]/20 to-transparent opacity-80" />
        
        {/* Organization Tag */}
        <div className="absolute left-4 top-4 rounded-full bg-black/40 backdrop-blur-md border border-white/10 px-3 py-1.5">
          <span className="text-[10px] font-bold uppercase tracking-widest text-[#98FF98]">
            {post.org?.name}
          </span>
        </div>

        {linkTo && (
          <div className="absolute right-4 top-4 flex h-8 w-8 items-center justify-center rounded-full bg-[#98FF98] text-[#0b0f0e] opacity-0 transition-all transform scale-75 group-hover:opacity-100 group-hover:scale-100">
            <ExternalLink className="h-4 w-4" />
          </div>
        )}

        {/* Date on image */}
        <div className="absolute bottom-4 left-4 flex items-center gap-1.5 text-[11px] font-medium text-[#9fb2a7]">
          <Calendar className="h-3.5 w-3.5 text-[#98FF98]" />
          {formattedDate}
        </div>
      </div>

      <div className="flex flex-1 flex-col p-6">
        <p className="text-sm leading-relaxed text-[#e2e8e2] line-clamp-4 group-hover:text-white transition-colors">
          {post.content}
        </p>
        
        {linkTo && (
          <div className="mt-auto pt-4 flex items-center gap-2 text-[10px] font-bold uppercase tracking-[0.2em] text-[#98FF98] opacity-60 group-hover:opacity-100 transition-opacity">
            View Venue Profile
          </div>
        )}
      </div>
    </article>
  );

  if (linkTo) {
    return (
      <Link to={linkTo} className="block h-full transition-transform hover:-translate-y-1.5 active:scale-[0.98]">
        {content}
      </Link>
    );
  }

  return content;
}
