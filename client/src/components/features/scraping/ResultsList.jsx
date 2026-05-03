import React from "react";
import { apiRequest } from "../../../api/client";

function SkeletonRow() {
  return (
    <div className="p-6 border-b border-white/10">
      <div className="h-5 w-2/3 bg-white/10 rounded animate-pulse mb-2" />
      <div className="h-4 w-1/2 bg-white/10 rounded animate-pulse mb-4" />
      <div className="flex gap-4">
        <div className="h-4 w-16 bg-white/10 rounded animate-pulse" />
        <div className="h-4 w-16 bg-white/10 rounded animate-pulse" />
        <div className="h-4 w-24 bg-white/10 rounded animate-pulse" />
        <div className="h-4 w-32 bg-white/10 rounded animate-pulse" />
      </div>
    </div>
  );
}

export default function ResultsList({ results, onView, onDelete, isLoading }) {
  const handleDelete = async (id, e) => {
    e.stopPropagation();
    if (!window.confirm("¿Estás seguro de que quieres eliminar este resultado?")) {
      return;
    }

    const { ok } = await apiRequest(`/results/${id}`, { method: "DELETE" });
    if (ok) {
      onDelete(id);
    }
  };

  const getStatusColor = (statusCode) => {
    if (statusCode >= 200 && statusCode < 300) return "text-emerald-400";
    if (statusCode >= 400 && statusCode < 500) return "text-orange-400";
    if (statusCode >= 500) return "text-red-400";
    return "text-gray-400";
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString("es-ES", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (isLoading && results.length === 0) {
    return (
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 overflow-hidden">
        <div className="px-6 py-4 border-b border-white/20">
          <div className="h-7 w-52 bg-white/10 rounded animate-pulse mb-2" />
          <div className="h-4 w-32 bg-white/10 rounded animate-pulse" />
        </div>
        <SkeletonRow />
        <SkeletonRow />
        <SkeletonRow />
      </div>
    );
  }

  return (
    <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 overflow-hidden">
      <div className="px-6 py-4 border-b border-white/20">
        <h2 className="text-2xl font-bold text-white">Resultados de Scraping</h2>
        <p className="text-gray-400 text-sm mt-1">
          {results.length} resultado{results.length !== 1 ? "s" : ""} encontrado{results.length !== 1 ? "s" : ""}
        </p>
      </div>

      {results.length === 0 ? (
        <div className="py-16 px-8 flex flex-col items-center text-center">
          <div className="w-20 h-20 mb-5 rounded-2xl bg-gradient-to-br from-violet-500/10 to-indigo-500/10 border border-white/10 flex items-center justify-center">
            <svg className="w-10 h-10 text-violet-500/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.5"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-white mb-2">Aún no hay resultados</h3>
          <p className="text-gray-400 text-sm max-w-xs">
            Ingresa una URL en el formulario de arriba para extraer su contenido y verlo aquí.
          </p>
        </div>
      ) : (
        <div className="divide-y divide-white/20">
          {results.map((result) => (
            <div
              key={result.id}
              className="p-6 hover:bg-white/5 transition-all duration-200 cursor-pointer"
              onClick={() => onView(result)}
            >
              <div className="flex items-start justify-between">
                <div className="flex-1 min-w-0">
                  <div className="flex items-start gap-4">
                    <div className="flex-1 min-w-0">
                      <h3 className="text-lg font-semibold text-white mb-1 truncate">
                        {result.title || "Sin título"}
                      </h3>
                      <p className="text-violet-400 text-sm mb-3 truncate">
                        {result.url}
                      </p>
                      
                      {result.description && (
                        <p className="text-gray-300 text-sm mb-3 line-clamp-2">
                          {result.description}
                        </p>
                      )}
                      
                      <div className="flex flex-wrap items-center gap-4 text-sm text-gray-400">
                        {result.seo_score !== undefined && (
                          <span className={`text-xs px-2 py-0.5 rounded font-semibold ${
                            result.seo_score >= 70
                              ? "bg-teal-500/20 text-teal-400"
                              : result.seo_score >= 40
                              ? "bg-orange-500/20 text-orange-400"
                              : "bg-red-500/20 text-red-400"
                          }`}>
                            SEO {result.seo_score}
                          </span>
                        )}
                        {result.robots_directive?.toLowerCase().includes("noindex") && (
                          <span className="text-xs px-2 py-0.5 rounded font-semibold bg-red-500/20 text-red-400">
                            noindex
                          </span>
                        )}
                        <div className="flex items-center gap-1">
                          <span>Estado:</span>
                          <span
                            className={`font-medium ${getStatusColor(
                              result.status_code
                            )}`}
                          >
                            {result.status_code}
                          </span>
                        </div>
                        
                        <div className="flex items-center gap-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span>{result.load_time_ms}ms</span>
                        </div>
                        
                        <div className="flex items-center gap-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                          </svg>
                          <span>{result.word_count || 0} palabras</span>
                        </div>
                        
                        <div className="flex items-center gap-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 7V3a2 2 0 012-2h4a2 2 0 012 2v4m-6 4v10a2 2 0 002 2h4a2 2 0 002-2V11M9 7h6" />
                          </svg>
                          <span>{formatDate(result.created_at)}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                
                <div className="ml-4 flex-shrink-0 flex items-center gap-2">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      onView(result);
                    }}
                    className="p-2 text-gray-400 hover:text-blue-400 hover:bg-blue-400/10 rounded-lg transition-colors"
                    title="Ver detalles"
                  >
                    <svg
                      className="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      />
                    </svg>
                  </button>
                  
                  <button
                    onClick={(e) => handleDelete(result.id, e)}
                    className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-colors"
                    title="Eliminar resultado"
                  >
                    <svg
                      className="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}