import React from "react";
import ResultCard from "./ResultCard";

export default function ResultsList({ results, onView, onDelete, onLogout }) {
  return (
    <div className="mb-8">
      <div className="flex justify-end mb-4">
        <button
          onClick={onLogout}
          className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg"
        >
          Cerrar sesión
        </button>
      </div>
      <div className="bg-black/20 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/10 p-8">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold bg-gradient-to-r from-green-400 to-blue-400 bg-clip-text text-transparent">
            Resultados
          </h2>
        </div>
        <div className="space-y-4">
          {results.length === 0 ? (
            <div className="text-center text-gray-400 py-8">
              Sin resultados. Introduce una URL arriba.
            </div>
          ) : (
            results.map((r) => (
              <ResultCard
                key={r.id}
                result={r}
                onView={() => onView(r)}
                onDelete={() => onDelete(r.id)}
              />
            ))
          )}
        </div>
      </div>
    </div>
  );
}
