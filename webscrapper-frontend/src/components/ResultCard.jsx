import React from "react";

export default function ResultCard({ result, onView, onDelete }) {
  const statusColor =
    result.status_code === 200 ? "text-emerald-400" : "text-gray-400";
  const date = new Date(result.created_at).toLocaleString();
  return (
    <div className="bg-gray-800/30 backdrop-blur-sm rounded-lg border border-gray-700/50 p-6 hover:bg-gray-800/40 transition-all duration-300 group">
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1 min-w-0">
          <h3
            className="text-lg font-semibold text-white truncate mb-2"
            title={result.title || "No title"}
          >
            {result.title || "No title"}
          </h3>
          <p className="text-sm text-blue-400 truncate mb-2" title={result.url}>
            {result.url}
          </p>
          <p className="text-sm text-gray-400 mb-3 line-clamp-2">
            {result.description || "No description available"}
          </p>
        </div>
        <span
          className={`px-2 py-1 text-xs font-medium ${statusColor} bg-white/10 rounded-full ml-4`}
        >
          {result.status_code}
        </span>
      </div>
      <div className="flex items-center justify-between text-sm text-gray-400 mb-4">
        <div className="flex items-center space-x-4">
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4" /* icon */></svg>
            <span>{(result.links || []).length} links</span>
          </span>
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4" /* icon */></svg>
            <span>{(result.images || []).length} imágenes</span>
          </span>
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4" /* icon */></svg>
            <span>{result.word_count || 0} palabras</span>
          </span>
        </div>
        <span>{date}</span>
      </div>
      <div className="flex items-center justify-between">
        <div className="flex items-center text-sm text-gray-500">
          <svg className="w-4 h-4" /* icon */></svg>
          <span>Tiempo de carga: {result.load_time_ms}ms</span>
        </div>
        <div className="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            onClick={onView}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-100 text-sm rounded-lg transition-colors"
          >
            Ver detalles
          </button>
          <button
            onClick={() => {
              if (
                window.confirm(
                  "¿Estás seguro de que quieres eliminar este resultado?"
                )
              ) {
                onDelete(result.id);
              }
            }}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-100 text-sm rounded-lg transition-colors"
          >
            Eliminar
          </button>
        </div>
      </div>
    </div>
  );
}
