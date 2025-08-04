import React from "react";
import { usePagination } from "../hooks/usePagination";
import Pagination from "./Pagination";
import LoadingSpinner from "./LoadingSpinner";

export default function PaginatedResultsList({ onView, onAlert }) {
  const {
    results,
    pagination,
    loading,
    error,
    currentPage,
    perPage,
    handlePageChange,
    handlePerPageChange,
    handleDelete,
    totalItems,
    totalPages,
    startItem,
    endItem
  } = usePagination(1, 10);

  const handleDeleteResult = async (id, e) => {
    e?.stopPropagation();
    
    if (!window.confirm("¿Estás seguro de que quieres eliminar este resultado?")) {
      return;
    }

    const success = await handleDelete(id);
    if (success) {
      onAlert?.({ 
        type: "success", 
        message: "Resultado eliminado correctamente" 
      });
    } else {
      onAlert?.({ 
        type: "error", 
        message: "Error al eliminar el resultado" 
      });
    }
  };

  const getStatusColor = (statusCode) => {
    if (statusCode >= 200 && statusCode < 300) return "text-green-400";
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

  if (error) {
    return (
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-red-500/20 overflow-hidden">
        <div className="p-8 text-center">
          <div className="w-16 h-16 mx-auto mb-4 bg-red-500/10 rounded-full flex items-center justify-center">
            <svg className="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-red-300 mb-2">Error al cargar resultados</h3>
          <p className="text-gray-400 text-sm">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 overflow-hidden">
        <div className="px-6 py-4 border-b border-white/20">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-2xl font-bold text-white">Resultados de Scraping</h2>
              <p className="text-gray-400 text-sm mt-1">
                {loading ? "Cargando..." : `${totalItems} resultado${totalItems !== 1 ? "s" : ""} encontrado${totalItems !== 1 ? "s" : ""}`}
              </p>
            </div>
            
            {totalItems > 0 && !loading && (
              <div className="text-sm text-gray-400">
                Mostrando {startItem} - {endItem} de {totalItems}
              </div>
            )}
          </div>
        </div>

        {loading ? (
          <div className="p-8">
            <LoadingSpinner />
          </div>
        ) : results.length === 0 ? (
          <div className="p-8 text-center">
            <div className="w-16 h-16 mx-auto mb-4 bg-white/10 rounded-full flex items-center justify-center">
              <svg
                className="w-8 h-8 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-300 mb-2">
              No hay resultados aún
            </h3>
            <p className="text-gray-400 text-sm">
              Realiza tu primer web scraping para ver los resultados aquí
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
                        <p className="text-blue-400 text-sm mb-3 truncate">
                          {result.url}
                        </p>
                        
                        {result.description && (
                          <p className="text-gray-300 text-sm mb-3 line-clamp-2">
                            {result.description}
                          </p>
                        )}
                        
                        <div className="flex flex-wrap items-center gap-4 text-sm text-gray-400">
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
                      onClick={(e) => handleDeleteResult(result.id, e)}
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

      {/* Componente de paginación */}
      {totalPages > 1 && (
        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          totalItems={totalItems}
          perPage={perPage}
          onPageChange={handlePageChange}
          onPerPageChange={handlePerPageChange}
        />
      )}
    </div>
  );
}