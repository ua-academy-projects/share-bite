import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { EditBrandProfileForm } from "./EditBrandProfileForm";
import type { BrandProfile } from "@/api/business";

type EditBrandProfileModalProps = {
  brand: BrandProfile;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: (updatedBrand: BrandProfile) => void;
  onRefresh?: () => void;
};

export function EditBrandProfileModal({
  brand,
  isOpen,
  onOpenChange,
  onSuccess,
  onRefresh,
}: EditBrandProfileModalProps) {
  const handleSuccess = (updatedBrand: BrandProfile) => {
    onSuccess(updatedBrand);
    onOpenChange(false);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-4xl border-white/10 shadow-2xl p-8 bg-[#0d241d]">
        <DialogHeader className="mb-6">
          <DialogTitle className="text-2xl font-black text-white">
            Brand Workspace
          </DialogTitle>
          <DialogDescription className="text-[#9fb2a7] text-base">
            Refine your digital presence and visual identity.
          </DialogDescription>
        </DialogHeader>
        <EditBrandProfileForm
          brand={brand}
          onSuccess={handleSuccess}
          onCancel={() => onOpenChange(false)}
          onRefresh={onRefresh}
        />
      </DialogContent>
    </Dialog>
  );
}
