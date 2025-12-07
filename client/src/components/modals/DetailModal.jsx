import React from "react";

export default function DetailModal({ result, onClose }) {
  if (!result) return null;

  const {
    url,
    status_code,
    content_type,
    load_time_ms,
    word_count,
    created_at,
    title,
    description,
    keywords,
    author,
    language,
    site_name,
    headers,
    links,
    images,
    content,
  } = result;
  const safeHeaders = headers ?? [];
  const safeLinks = links ?? [];
  const safeImages = images ?? [];
  const safeContent = content ?? "";

  const date = new Date(created_at).toLocaleString();
  const statusColor = status_code === 200 ? "text-teal-400" : "text-red-400";

  return (
    <div
      id="detailModal"
      className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      onClick={(e) => e.target.id === "detailModal" && onClose()}
    >
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 max-w-4xl w-full max-h-[90vh] overflow-y-auto relative p-6">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 hover:text-white"
          aria-label="Cerrar"
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>

        <h3 className="text-2xl font-bold text-white mb-4">
          Detalles del Scraping
        </h3>

        {/* Información Básica */}
        <div className="bg-white/5 rounded-lg p-4 mb-6">
          <h4 className="text-lg font-semibold text-white mb-3">
            Información Básica
          </h4>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-300">
            <div>
              <span className="text-gray-400">URL:</span>
              <p className="text-cyan-400 break-all">{url}</p>
            </div>
            <div>
              <span className="text-gray-400">Código de estado:</span>
              <p className={`${statusColor} font-medium`}>{status_code}</p>
            </div>
            <div>
              <span className="text-gray-400">Tipo de contenido:</span>
              <p className="text-white">{content_type || "Desconocido"}</p>
            </div>
            <div>
              <span className="text-gray-400">Tiempo de carga:</span>
              <p className="text-white">{load_time_ms}ms</p>
            </div>
            <div>
              <span className="text-gray-400">Número de palabras:</span>
              <p className="text-white">{word_count || 0}</p>
            </div>
            <div>
              <span className="text-gray-400">Scraped:</span>
              <p className="text-white">{date}</p>
            </div>
          </div>
        </div>

        {/* Metadatos */}
        <div className="bg-white/5 rounded-lg p-4 mb-6">
          <h4 className="text-lg font-semibold text-white mb-3">Metadatos</h4>
          <div className="space-y-2 text-sm text-gray-300">
            <div>
              <span className="text-gray-400">Título:</span>
              <p className="text-white">{title || "Sin título"}</p>
            </div>
            <div>
              <span className="text-gray-400">Descripción:</span>
              <p className="text-white">{description || "Sin descripción"}</p>
            </div>
            <div>
              <span className="text-gray-400">Palabras clave:</span>
              <p className="text-white">{keywords || "Sin palabras clave"}</p>
            </div>
            <div>
              <span className="text-gray-400">Autor:</span>
              <p className="text-white">{author || "Desconocido"}</p>
            </div>
            <div>
              <span className="text-gray-400">Idioma:</span>
              <p className="text-white">{language || "Desconocido"}</p>
            </div>
            <div>
              <span className="text-gray-400">Nombre del sitio:</span>
              <p className="text-white">{site_name || "Desconocido"}</p>
            </div>
          </div>
        </div>

        {/* Cabeceras */}
        {safeHeaders.length > 0 && (
          <div className="bg-white/5 rounded-lg p-4 mb-6">
            <h4 className="text-lg font-semibold text-white mb-3">
              Cabeceras ({safeHeaders.length})
            </h4>
            <div className="max-h-40 overflow-y-auto text-sm space-y-2">
              {safeHeaders.map((h, i) => (
                <div key={i}>
                  <span className="text-teal-400 font-mono">H{h.level}:</span>{" "}
                  <span className="text-white">{h.text}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Links */}
        {safeLinks.length > 0 && (
          <div className="bg-white/5 rounded-lg p-4 mb-6">
            <h4 className="text-lg font-semibold text-white mb-3">
              Links ({safeLinks.length})
            </h4>
            <div className="max-h-40 overflow-y-auto text-sm space-y-2">
              {safeLinks.map((link, i) => (
                <div key={i}>
                  <a
                    href={link}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-cyan-400 hover:text-cyan-300 truncate block"
                  >
                    {link}
                  </a>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Imágenes */}
        {safeImages.length > 0 && (
          <div className="bg-white/5 rounded-lg p-4 mb-6">
            <h4 className="text-lg font-semibold text-white mb-3">
              Imágenes ({safeImages.length})
            </h4>
            <div className="max-h-40 overflow-y-auto text-sm space-y-2">
              {safeImages.map((img, i) => (
                <div key={i}>
                  <a
                    href={img}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-teal-400 hover:text-teal-300 truncate block"
                  >
                    {img}
                  </a>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Vista previa de contenido */}
        {safeContent && (
          <div className="bg-white/5 rounded-lg p-4 mb-6">
            <h4 className="text-lg font-semibold text-white mb-3">
              Vista previa del contenido
            </h4>
            <pre className="whitespace-pre-wrap break-words text-sm text-gray-300 max-h-60 overflow-y-auto">
              {safeContent.length > 2000
                ? `${safeContent.substring(0, 2000)}...\n\n[Content truncated]`
                : safeContent}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}
