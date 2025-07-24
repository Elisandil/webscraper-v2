import React, { useState } from "react";
import { apiRequest } from "../api/client";

export default function ScrapeForm({ onSuccess, onError }) {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);

  const submit = async () => {
    if (!/^https?:\/\//.test(url)) {
      return onError(
        "Por favor, ingresa una URL válida con http:// o https://"
      );
    }
    setLoading(true);
    const { ok, data } = await apiRequest("/scrape", {
      method: "POST",
      body: JSON.stringify({ url }),
    });
    setLoading(false);
    if (ok) {
      onSuccess(`Se ha scrapeado el sitio web con éxito: ${url}`);
      setUrl("");
    } else {
      onError(data.error || "Fallo al scrapear la URL");
    }
  };

  return (
    <div className="mb-8">
      <div className="bg-gray-800/50 backdrop-blur-lg rounded-lg shadow-lg border border-gray-700/50 p-8">
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-100">
          Extraer contenido web
        </h2>
        <div className="space-y-6">
          <div className="relative">
            <input
              type="url"
              placeholder="https://example.com"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              className="w-full px-6 py-4 bg-gray-700/50 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-transparent transition-all duration-300"
            />
          </div>
          <button
            onClick={submit}
            disabled={loading}
            className="w-full py-4 bg-gray-700 hover:bg-gray-600 rounded-lg font-semibold text-gray-100 transition-all transform hover:scale-[1.02] disabled:opacity-50 disabled:hover:scale-100"
          >
            {loading ? "Scraping..." : "Scrapear sitio web"}
          </button>
        </div>
      </div>
    </div>
  );
}
