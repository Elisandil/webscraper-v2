import React, { useEffect } from "react";

export default function Alert({ alert, onClose }) {
  useEffect(() => {
    if (alert) {
      const t = setTimeout(onClose, 5000);
      return () => clearTimeout(t);
    }
  }, [alert, onClose]);

  if (!alert) return null;
  const base =
    "fixed bottom-4 right-4 px-6 py-3 rounded-lg shadow-lg backdrop-blur-sm";
  const style =
    alert.type === "success"
      ? "bg-gray-800/90 text-emerald-400 border border-emerald-500/20"
      : "bg-gray-800/90 text-gray-400 border border-gray-700/50";
  return (
    <div className={`${base} ${style}`} role="alert">
      {alert.message}
    </div>
  );
}
