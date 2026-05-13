import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { X } from "lucide-react";
import * as QRCodeLib from "qrcode";
import { useQRCodeModal } from "@/contexts/QRCodeModalContext";

export function QRCodeModalContainer() {
  const { isOpen, boxCode, closeModal } = useQRCodeModal();
  const [qrCodeDataUrl, setQrCodeDataUrl] = useState<string | null>(null);
  const [isGenerating, setIsGenerating] = useState(false);

  // Генерувати QR-код коли модаль відкривається
  useEffect(() => {
    if (isOpen && boxCode) {
      setIsGenerating(true);
      QRCodeLib.toDataURL(boxCode, {
        width: 256,
        color: {
          dark: "#1A3C34",
          light: "#FFFFFF"
        }
      })
        .then((url) => {
          setQrCodeDataUrl(url);
          setIsGenerating(false);
        })
        .catch((err) => {
          console.error("QR Code generation error:", err);
          setIsGenerating(false);
        });
    }
  }, [isOpen, boxCode]);

  // Очистити QR-код при закритті
  useEffect(() => {
    if (!isOpen) {
      setQrCodeDataUrl(null);
    }
  }, [isOpen]);

  // Запобігти закриттю при кліку на контент модалі
  const handleModalClick = (e: React.MouseEvent) => {
    e.stopPropagation();
  };

  if (!isOpen || !boxCode) return null;

  return createPortal(
    <>
      {/* Backdrop з блюром */}
      <div
        className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 transition-opacity duration-300"
        onClick={closeModal}
      />

      {/* Модальне вікно */}
      <div className="fixed inset-0 flex items-center justify-center z-50 p-4 pointer-events-none">
        <div
          className="bg-white dark:bg-[#163d32] rounded-3xl shadow-2xl p-8 max-w-md w-full border border-gray-100 dark:border-[#2f5e50] pointer-events-auto animate-in fade-in zoom-in-95 duration-300"
          onClick={handleModalClick}
        >
          {/* Header */}
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-[#1A3C34] dark:text-white">
              Бокс зарезервовано! 🎉
            </h2>
            <button
              onClick={closeModal}
              className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors p-1 hover:bg-gray-100 dark:hover:bg-[#0d241d] rounded-lg"
              aria-label="Закрити"
            >
              <X size={24} />
            </button>
          </div>

          {/* QR Code Container */}
          <QRCodeContainer qrCodeDataUrl={qrCodeDataUrl} isGenerating={isGenerating} />

          {/* Box Code */}
          <BoxCodeDisplay boxCode={boxCode} />

          {/* Instructions */}
          <InstructionsSection />

          {/* Copy Button */}
          <CopyButton boxCode={boxCode} />

          {/* Close Button */}
          <button
            onClick={closeModal}
            className="w-full bg-gray-200 hover:bg-gray-300 dark:bg-[#2f5e50] dark:hover:bg-[#3a7a66] text-[#1A3C34] dark:text-white font-bold py-3 rounded-xl transition-colors"
          >
            Готово
          </button>
        </div>
      </div>
    </>,
    document.body
  );
}

// Subcomponents для декомпозиції

interface QRCodeContainerProps {
  qrCodeDataUrl: string | null;
  isGenerating: boolean;
}

function QRCodeContainer({ qrCodeDataUrl, isGenerating }: QRCodeContainerProps) {
  return (
    <div className="flex flex-col items-center bg-gray-50 dark:bg-[#0d241d] rounded-2xl p-6 mb-6 border border-gray-200 dark:border-[#2f5e50]">
      <div className="bg-white p-4 rounded-xl shadow-sm">
        {qrCodeDataUrl && !isGenerating ? (
          <img src={qrCodeDataUrl} alt="QR Code" width={256} height={256} className="rounded" />
        ) : (
          <div className="w-64 h-64 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-700 dark:to-gray-600 rounded flex items-center justify-center">
            <div className="text-center">
              <div className="w-12 h-12 border-4 border-gray-300 dark:border-gray-500 border-t-emerald-500 rounded-full animate-spin mx-auto mb-2" />
              <span className="text-gray-500 dark:text-gray-400 text-sm">Генерування...</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

interface BoxCodeDisplayProps {
  boxCode: string | null;
}

function BoxCodeDisplay({ boxCode }: BoxCodeDisplayProps) {
  return (
    <div className="mb-6">
      <p className="text-gray-500 dark:text-gray-400 text-sm mb-2 font-medium">Код боксу:</p>
      <div className="bg-gray-50 dark:bg-[#0d241d] p-4 rounded-xl border border-gray-200 dark:border-[#2f5e50] font-mono font-bold text-lg text-[#1A3C34] dark:text-[#98FF98] text-center break-all select-all cursor-pointer hover:bg-gray-100 dark:hover:bg-[#163d32] transition-colors">
        {boxCode}
      </div>
    </div>
  );
}

function InstructionsSection() {
  return (
    <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4 mb-6">
      <p className="text-sm text-blue-700 dark:text-blue-300 leading-relaxed">
        📱 Збережіть або скопіюйте цей код. Він знадобиться вам при отриманні боксу.
      </p>
    </div>
  );
}

interface CopyButtonProps {
  boxCode: string | null;
}

function CopyButton({ boxCode }: CopyButtonProps) {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = async () => {
    if (boxCode) {
      try {
        await navigator.clipboard.writeText(boxCode);
        setIsCopied(true);
        setTimeout(() => setIsCopied(false), 2000);
      } catch (err) {
        console.error("Failed to copy:", err);
      }
    }
  };

  return (
    <button
      onClick={handleCopy}
      className={`w-full font-bold py-3 rounded-xl transition-all mb-3 ${
        isCopied
          ? "bg-emerald-500 dark:bg-emerald-600 text-white"
          : "bg-emerald-500 hover:bg-emerald-600 dark:bg-emerald-600 dark:hover:bg-emerald-700 text-white"
      }`}
    >
      {isCopied ? "✓ Скопійовано!" : "Скопіювати код"}
    </button>
  );
}
