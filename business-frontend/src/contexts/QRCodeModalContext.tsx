import { createContext, useContext, useState, ReactNode } from "react";

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

  const openModal = (code: string) => {
    setBoxCode(code);
    setIsOpen(true);
  };

  const closeModal = () => {
    setIsOpen(false);
    // Очистити код після закриття анімації
    setTimeout(() => setBoxCode(null), 300);
  };

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
