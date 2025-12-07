import React, { useEffect, useState } from "react";
import { apiRequest } from "../../api/client";

export default function HealthIndicator() {
  const [online, setOnline] = useState(true);
  
  useEffect(() => {
    const check = () => {
      apiRequest("/health").then(({ ok, data }) => {
        setOnline(ok && data.data?.status === "ok");
      });
    };
    check();
    const iv = setInterval(check, 30000);
    return () => clearInterval(iv);
  }, []);
  
  return (
    <div className="fixed top-4 right-6 flex items-center space-x-2 text-sm bg-black/80 backdrop-blur-lg rounded-lg px-3 py-2 border border-white/20 shadow-lg">
      <span
        className={`inline-block w-2 h-2 rounded-full ${
          online ? "bg-cyan-400 animate-pulse shadow-lg shadow-cyan-500/50" : "bg-red-400"
        }`}
      ></span>
      <span className="text-gray-300 font-medium">
        {online ? "Servicio en línea" : "Servicio fuera"}
      </span>
    </div>
  );
}