import React, { useEffect } from "react";

export default function Alert({ alert, onClose }) {
  useEffect(() => {
    if (alert) {
      const t = setTimeout(onClose, 5000);
      return () => clearTimeout(t);
    }
  }, [alert, onClose]);

  if (!alert) return null;
  
  const baseStyles = "fixed bottom-4 right-4 px-6 py-3 rounded-lg shadow-lg backdrop-blur-sm transition-all duration-300 border";
  const typeStyles = alert.type === "success"
    ? "bg-emerald-900/90 text-emerald-400 border-emerald-500/20"
    : "bg-red-900/90 text-red-400 border-red-500/20";
  
  return (
    <div className={`${baseStyles} ${typeStyles}`} role="alert">
      <div className="flex items-center gap-2">
        {alert.type === "success" ? (
          <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
          </svg>
        ) : (
          <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        )}
        <span className="font-medium">{alert.message}</span>
        <button
          onClick={onClose}
          className="ml-2 text-gray-400 hover:text-white transition-colors"
          aria-label="Cerrar alerta"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  );
}