import { createContext, useContext, useState, ReactNode, useRef, useEffect } from "react";

interface QRCodeModalContextType {
  isOpen: boolean;
  boxCode: string | null;
  openModal: (code: string) => void;
  closeModal: () => void;
}

const QRCodeModalContext = createContext<QRCodeModalContextType | undefined>(undefined);

export function QRCodeModalProvider({ children }: { children: ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);
  const [boxCode, setBoxCode] = useState<string | null>(null);
  const modalClearTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const openModal = (code: string) => {
    // Clear any pending timeout when opening new modal
    if (modalClearTimeoutRef.current) {
      clearTimeout(modalClearTimeoutRef.current);
      modalClearTimeoutRef.current = null;
    }
    setBoxCode(code);
    setIsOpen(true);
  };

  const closeModal = () => {
    setIsOpen(false);
    // Clear any existing timeout before scheduling new one
    if (modalClearTimeoutRef.current) {
      clearTimeout(modalClearTimeoutRef.current);
    }
    // Clear code after close animation
    modalClearTimeoutRef.current = setTimeout(() => {
      setBoxCode(null);
      modalClearTimeoutRef.current = null;
    }, 300);
  };

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (modalClearTimeoutRef.current) {
        clearTimeout(modalClearTimeoutRef.current);
      }
    };
  }, []);

  return (
    <QRCodeModalContext.Provider value={{ isOpen, boxCode, openModal, closeModal }}>
      {children}
    </QRCodeModalContext.Provider>
  );
}

export function useQRCodeModal() {
  const context = useContext(QRCodeModalContext);
  if (!context) {
    throw new Error("useQRCodeModal must be used within QRCodeModalProvider");
  }
  return context;
}
