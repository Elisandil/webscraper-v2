import React, { useState, useRef, useEffect } from "react";
import { useChat } from "../../contexts/ChatContext";
import ChatMessage from "./ChatMessage";

export default function ChatWidget() {
  const {
    messages,
    isOpen,
    isLoading,
    pendingIntent,
    sendMessage,
    confirmAction,
    cancelAction,
    clearChat,
    toggleChat,
  } = useChat();

  const [input, setInput] = useState("");
  const messagesEndRef = useRef(null);
  const inputRef = useRef(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  useEffect(() => {
    if (isOpen && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isOpen]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (input.trim() && !isLoading) {
      sendMessage(input);
      setInput("");
    }
  };

  const handleConfirm = () => {
    confirmAction();
  };

  const handleCancel = () => {
    cancelAction();
  };

  return (
    <>
      {!isOpen && (
        <button
          onClick={toggleChat}
          className="fixed bottom-6 right-6 z-50 bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 text-white p-4 rounded-full shadow-2xl shadow-cyan-500/50 transition-all duration-300 hover:scale-110"
          aria-label="Abrir chat asistente"
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
              d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
            />
          </svg>
          <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center animate-pulse">
            AI
          </span>
        </button>
      )}

      {isOpen && (
        <div className="fixed bottom-6 right-6 z-50 w-96 h-[600px] bg-gray-900/95 backdrop-blur-xl border border-white/20 rounded-2xl shadow-2xl flex flex-col overflow-hidden">
          <div className="bg-gradient-to-r from-cyan-600 to-teal-600 px-4 py-3 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
              <h3 className="text-white font-semibold">Asistente de Scraping</h3>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={clearChat}
                className="text-white/80 hover:text-white transition-colors"
                title="Limpiar chat"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
              <button
                onClick={toggleChat}
                className="text-white/80 hover:text-white transition-colors"
                title="Cerrar chat"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-2">
            {messages.map((message, index) => (
              <ChatMessage key={index} message={message} isUser={message.isUser} />
            ))}
            
            {isLoading && (
              <div className="flex justify-start mb-4">
                <div className="bg-white/10 rounded-lg px-4 py-3 border border-white/20">
                  <div className="flex space-x-2">
                    <div className="w-2 h-2 bg-cyan-500 rounded-full animate-bounce"></div>
                    <div className="w-2 h-2 bg-cyan-500 rounded-full animate-bounce" style={{ animationDelay: "0.1s" }}></div>
                    <div className="w-2 h-2 bg-cyan-500 rounded-full animate-bounce" style={{ animationDelay: "0.2s" }}></div>
                  </div>
                </div>
              </div>
            )}
            
            <div ref={messagesEndRef} />
          </div>

          {pendingIntent && !isLoading && (
            <div className="px-4 py-3 bg-teal-500/10 border-t border-teal-500/20 flex gap-2">
              <button
                onClick={handleConfirm}
                className="flex-1 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-500 hover:to-green-600 text-white py-2 px-4 rounded-lg transition-all duration-200 font-medium"
              >
                ✓ Confirmar
              </button>
              <button
                onClick={handleCancel}
                className="flex-1 bg-red-600/80 hover:bg-red-600 text-white py-2 px-4 rounded-lg transition-all duration-200 font-medium"
              >
                ✕ Cancelar
              </button>
            </div>
          )}

          <form onSubmit={handleSubmit} className="p-4 border-t border-white/10">
            <div className="flex gap-2">
              <input
                ref={inputRef}
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Escribe tu mensaje..."
                className="flex-1 bg-white/10 border border-white/20 rounded-lg px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent transition-all duration-200"
                disabled={isLoading}
              />
              <button
                type="submit"
                disabled={isLoading || !input.trim()}
                className="bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 text-white p-2 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                </svg>
              </button>
            </div>
            
            {messages.length === 1 && (
              <div className="mt-3 text-xs text-gray-400">
                <p className="mb-1">💡 Ejemplos:</p>
                <ul className="list-disc list-inside space-y-1">
                  <li>"Scrapea https://example.com ahora"</li>
                  <li>"Programa reddit.com cada 2 horas"</li>
                </ul>
              </div>
            )}
          </form>
        </div>
      )}
    </>
  );
}
