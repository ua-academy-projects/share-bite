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
      <DialogContent className="sm:max-w-[1000px] w-[95vw] max-h-[90vh] overflow-y-auto border-white/10 shadow-2xl p-10 md:p-12">
        <DialogHeader className="space-y-4 mb-8">
          <DialogTitle className="text-4xl font-bold tracking-tight text-[#F9F7F2]">
            Public Identity
          </DialogTitle>
          <DialogDescription className="text-[#cbd5cf] text-xl max-w-2xl leading-relaxed">
            Manage how your brand appears to customers. A complete profile increases engagement and trust within the ShareBite community.
          </DialogDescription>
        </DialogHeader>
        <div className="mt-4">
          <EditBrandProfileForm
            brand={brand}
            onSuccess={handleSuccess}
            onCancel={() => onOpenChange(false)}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}
