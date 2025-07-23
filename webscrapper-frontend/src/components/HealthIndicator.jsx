import React, { useEffect, useState } from 'react';
import { apiRequest } from '../api/client';

export default function HealthIndicator() {
  const [online, setOnline] = useState(true);
  useEffect(() => {
    const check = () => {
      apiRequest('/health').then(({ ok, data }) => {
        setOnline(ok && data.data?.status === 'ok');
      });
    };
    check();
    const iv = setInterval(check, 30000);
    return () => clearInterval(iv);
  }, []);
  return (
    <div className="fixed top-4 right-6 flex items-center space-x-2 text-sm">
      <span className={`inline-block w-2 h-2 rounded-full ${online ? 'bg-green-400 animate-pulse' : 'bg-red-400'}`}></span>
      <span className="text-gray-400">{online ? 'Servicio en línea' : 'Servicio fuera'}</span>
    </div>
  );
}
