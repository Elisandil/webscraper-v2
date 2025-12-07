import React, { useState } from "react";
import { publicApiRequest } from "../api/client";

export default function Landing({ onGetStarted }) {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  const handleTryFree = async (e) => {
    e.preventDefault();
    if (!url.trim()) {
      setError("Por favor, ingresa una URL válida");
      return;
    }

    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await publicApiRequest("/scrape", {
        method: "POST",
        body: JSON.stringify({ url }),
      });

      if (response.status === 429) {
        // Rate limit exceeded
        setError("Has alcanzado el límite de pruebas gratuitas. ¡Regístrate para continuar!");
        setTimeout(() => {
          onGetStarted();
        }, 3000);
      } else if (!response.ok) {
        setError(response.data.message || "Error al procesar la URL");
      } else {
        setResult(response.data.data);
        setUrl("");
      }
    } catch (err) {
      setError("Error de conexión. Inténtalo de nuevo.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
      <nav className="container mx-auto px-6 py-6">
        <div className="flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-cyan-500 to-teal-500 rounded-lg flex items-center justify-center shadow-lg shadow-cyan-500/40">
              <span className="text-white font-bold text-xl">W</span>
            </div>
            <span className="text-2xl font-bold bg-gradient-to-r from-cyan-400 to-teal-400 text-transparent bg-clip-text">WebScraper</span>
          </div>
          <button
            onClick={onGetStarted}
            className="px-6 py-2 bg-cyan-600 hover:bg-cyan-500 rounded-lg font-semibold transition-colors duration-200 shadow-lg shadow-cyan-500/20"
          >
            Iniciar Sesión
          </button>
        </div>
      </nav>

      <section className="container mx-auto px-6 py-20 text-center">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-r from-cyan-400 via-teal-400 to-green-400 text-transparent bg-clip-text">
            Extrae Datos Web de Forma Inteligente
          </h1>
          <p className="text-xl md:text-2xl text-gray-300 mb-10">
            Automatiza la extracción de datos de cualquier sitio web con nuestra
            plataforma potente y fácil de usar.
          </p>

          <div className="max-w-2xl mx-auto mb-8">
            <form onSubmit={handleTryFree} className="flex flex-col sm:flex-row gap-4">
              <input
                type="url"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="https://ejemplo.com"
                className="flex-1 px-6 py-4 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent"
                disabled={loading}
              />
              <button
                type="submit"
                disabled={loading}
                className="px-8 py-4 bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 disabled:from-gray-600 disabled:to-gray-600 rounded-lg font-bold text-lg transition-all duration-200 transform hover:scale-105 shadow-lg shadow-cyan-500/30 disabled:transform-none disabled:shadow-none"
              >
                {loading ? "Procesando..." : "Probar Gratis"}
              </button>
            </form>

            {error && (
              <div className="mt-4 p-4 bg-red-500 bg-opacity-20 border border-red-500 rounded-lg text-red-200">
                {error}
              </div>
            )}

            {result && (
              <div className="mt-6 p-6 bg-gray-800 bg-opacity-50 backdrop-blur-sm rounded-xl border border-gray-700 text-left">
                <h3 className="text-xl font-bold mb-4 text-cyan-400">Resultado del Scraping</h3>
                <div className="space-y-2 text-gray-300">
                  <p><span className="font-semibold text-white">URL:</span> {result.url}</p>
                  <p><span className="font-semibold text-white">Código de Estado:</span> {result.status_code}</p>
                  <p><span className="font-semibold text-white">Título:</span> {result.title || "N/A"}</p>
                  <p><span className="font-semibold text-white">Palabras:</span> {result.word_count}</p>
                  <p className="mt-4 text-sm text-gray-400">
                    ¡Regístrate para guardar tus resultados y acceder a más funciones!
                  </p>
                  <button
                    onClick={onGetStarted}
                    className="mt-3 px-6 py-2 bg-cyan-600 hover:bg-cyan-500 rounded-lg font-semibold transition-colors"
                  >
                    Registrarse Ahora
                  </button>
                </div>
              </div>
            )}
          </div>

          <p className="text-sm text-gray-400 mb-6">
            Prueba gratis hasta 3 consultas sin registro
          </p>
        </div>
      </section>

      <section className="container mx-auto px-6 py-20">
        <h2 className="text-4xl font-bold text-center mb-16">
          Características Principales
        </h2>
        <div className="grid md:grid-cols-3 gap-8">
          <div className="bg-gray-800 bg-opacity-50 backdrop-blur-sm rounded-xl p-8 border border-gray-700 hover:border-cyan-500 transition-all duration-300 hover:shadow-lg hover:shadow-cyan-500/20">
            <div className="w-14 h-14 bg-cyan-600 rounded-lg flex items-center justify-center mb-6 shadow-lg shadow-cyan-500/30">
              <svg
                className="w-8 h-8"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                />
              </svg>
            </div>
            <h3 className="text-2xl font-bold mb-4">Rápido y Eficiente</h3>
            <p className="text-gray-400">
              Extrae datos en segundos con nuestro motor de scraping optimizado
              y de alto rendimiento.
            </p>
          </div>

          <div className="bg-gray-800 bg-opacity-50 backdrop-blur-sm rounded-xl p-8 border border-gray-700 hover:border-teal-500 transition-all duration-300 hover:shadow-lg hover:shadow-teal-500/20">
            <div className="w-14 h-14 bg-teal-600 rounded-lg flex items-center justify-center mb-6 shadow-lg shadow-teal-500/30">
              <svg
                className="w-8 h-8"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h3 className="text-2xl font-bold mb-4">Programación Flexible</h3>
            <p className="text-gray-400">
              Programa tus scraping para ejecutarse automáticamente en los
              intervalos que desees.
            </p>
          </div>

          <div className="bg-gray-800 bg-opacity-50 backdrop-blur-sm rounded-xl p-8 border border-gray-700 hover:border-green-500 transition-all duration-300 hover:shadow-lg hover:shadow-green-500/20">
            <div className="w-14 h-14 bg-green-600 rounded-lg flex items-center justify-center mb-6 shadow-lg shadow-green-500/30">
              <svg
                className="w-8 h-8"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                />
              </svg>
            </div>
            <h3 className="text-2xl font-bold mb-4">Seguro y Confiable</h3>
            <p className="text-gray-400">
              Tus datos están seguros con autenticación robusta y almacenamiento
              cifrado.
            </p>
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="container mx-auto px-6 py-20 bg-gray-800 bg-opacity-30 rounded-3xl my-10">
        <h2 className="text-4xl font-bold text-center mb-16">
          Cómo Funciona
        </h2>
        <div className="grid md:grid-cols-3 gap-12 max-w-5xl mx-auto">
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-br from-cyan-500 to-cyan-600 rounded-full flex items-center justify-center mx-auto mb-6 text-2xl font-bold shadow-lg shadow-cyan-500/40">
              1
            </div>
            <h3 className="text-xl font-bold mb-3">Ingresa la URL</h3>
            <p className="text-gray-400">
              Proporciona la URL del sitio web del que deseas extraer datos
            </p>
          </div>
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-br from-teal-500 to-teal-600 rounded-full flex items-center justify-center mx-auto mb-6 text-2xl font-bold shadow-lg shadow-teal-500/40">
              2
            </div>
            <h3 className="text-xl font-bold mb-3">Configura tu Scraping</h3>
            <p className="text-gray-400">
              Define qué datos necesitas y programa la frecuencia si lo deseas
            </p>
          </div>
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-br from-green-500 to-green-600 rounded-full flex items-center justify-center mx-auto mb-6 text-2xl font-bold shadow-lg shadow-green-500/40">
              3
            </div>
            <h3 className="text-xl font-bold mb-3">Obtén tus Datos</h3>
            <p className="text-gray-400">
              Recibe los datos extraídos listos para usar y exportar
            </p>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="container mx-auto px-6 py-20 text-center">
        <div className="max-w-3xl mx-auto">
          <h2 className="text-4xl font-bold mb-6">
            ¿Listo para Empezar?
          </h2>
          <p className="text-xl text-gray-300 mb-10">
            Únete a nosotros y comienza a extraer datos de forma profesional hoy mismo.
          </p>
          <button
            onClick={onGetStarted}
            className="px-10 py-4 bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 rounded-lg font-bold text-lg transition-all duration-200 transform hover:scale-105 shadow-lg shadow-cyan-500/30"
          >
            Comenzar Gratis →
          </button>
        </div>
      </section>

      {/* Footer */}
      <footer className="container mx-auto px-6 py-8 border-t border-gray-800">
        <div className="text-center text-gray-400">
          <p>&copy; 2025 WebScraper. Todos los derechos reservados.</p>
        </div>
      </footer>
    </div>
  );
}
