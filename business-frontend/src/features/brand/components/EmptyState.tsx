type EmptyStateProps = {
  title: string;
  description?: string;
};

export function EmptyState({ title, description }: EmptyStateProps) {
  return (
    <div className="rounded-3xl border border-white/10 bg-[#0f1b17]/70 px-6 py-8 text-center">
      <p className="text-sm font-semibold text-[#F9F7F2]">{title}</p>
      {description ? (
        <p className="mt-2 text-xs text-[#9fb2a7]">{description}</p>
      ) : null}
    </div>
  );
}
