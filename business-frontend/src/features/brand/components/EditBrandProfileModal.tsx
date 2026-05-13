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
};

export function EditBrandProfileModal({
  brand,
  isOpen,
  onOpenChange,
  onSuccess,
}: EditBrandProfileModalProps) {
  const handleSuccess = (updatedBrand: BrandProfile) => {
    onSuccess(updatedBrand);
    onOpenChange(false);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-3xl border-white/10 shadow-2xl p-8">
        <DialogHeader className="mb-6">
          <DialogTitle className="text-2xl font-bold text-[#F9F7F2]">
            Edit Brand Profile
          </DialogTitle>
          <DialogDescription className="text-[#cbd5cf]">
            Update your public brand identity and visual presence.
          </DialogDescription>
        </DialogHeader>
        <EditBrandProfileForm
          brand={brand}
          onSuccess={handleSuccess}
          onCancel={() => onOpenChange(false)}
        />
      </DialogContent>
    </Dialog>
  );
}
