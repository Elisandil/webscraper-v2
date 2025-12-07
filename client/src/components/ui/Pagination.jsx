import React from "react";

export default function Pagination({ 
  currentPage, 
  totalPages, 
  totalItems, 
  perPage, 
  onPageChange, 
  onPerPageChange 
}) {
  const perPageOptions = [10, 30, 50];

  const handleFirstPage = () => {
    if (currentPage > 1) {
      onPageChange(1);
    }
  };

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  const handleLastPage = () => {
    if (currentPage < totalPages) {
      onPageChange(totalPages);
    }
  };
  if (totalItems === 0) {
    return null;
  }

  const startItem = totalItems === 0 ? 0 : Math.min((currentPage - 1) * perPage + 1, totalItems);
  const endItem = totalItems === 0 ? 0 : Math.min(currentPage * perPage, totalItems);
  const hasNavigation = totalPages > 1;

  return (
    <div className="bg-black/60 backdrop-blur-lg rounded-xl border border-white/20 px-6 py-4 mt-6">
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
        {/* Información de resultados */}
        <div className="flex items-center gap-4">
          <p className="text-sm text-gray-300">
            Mostrando <span className="font-medium text-white">{startItem}</span> a{" "}
            <span className="font-medium text-white">{endItem}</span> de{" "}
            <span className="font-medium text-white">{totalItems}</span> resultados
          </p>
          
          {/* Selector de elementos por página */}
          <div className="flex items-center gap-2">
            <label htmlFor="perPage" className="text-sm text-gray-400">
              Por página:
            </label>
            <select
              id="perPage"
              value={perPage}
              onChange={(e) => onPerPageChange(parseInt(e.target.value))}
              className="bg-gray-800/90 border border-white/20 text-white text-sm rounded-md px-3 py-1 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent"
            >
              {perPageOptions.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Controles de paginación - Solo mostrar si hay más de una página */}
        {hasNavigation && (
          <div className="flex items-center gap-2">
            {/* Información de página */}
            <span className="text-sm text-gray-300 mr-4">
              Página <span className="font-medium text-white">{currentPage}</span> de{" "}
              <span className="font-medium text-white">{totalPages}</span>
            </span>

            {/* Botones de navegación */}
            <div className="flex items-center gap-1">
              {/* Primera página */}
              <button
                onClick={handleFirstPage}
                disabled={currentPage === 1}
                className="p-2 rounded-md text-gray-400 hover:text-white hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                title="Primera página"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
                </svg>
              </button>

              {/* Página anterior */}
              <button
                onClick={handlePreviousPage}
                disabled={currentPage === 1}
                className="p-2 rounded-md text-gray-400 hover:text-white hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                title="Página anterior"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 19l-7-7 7-7" />
                </svg>
              </button>

              {/* Separador visual */}
              <div className="w-px h-6 bg-white/20 mx-2"></div>

              {/* Página siguiente */}
              <button
                onClick={handleNextPage}
                disabled={currentPage === totalPages}
                className="p-2 rounded-md text-gray-400 hover:text-white hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                title="Página siguiente"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7" />
                </svg>
              </button>

              {/* Última página */}
              <button
                onClick={handleLastPage}
                disabled={currentPage === totalPages}
                className="p-2 rounded-md text-gray-400 hover:text-white hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                title="Última página"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 5l7 7-7 7M5 5l7 7-7 7" />
                </svg>
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}