import React, { createContext, useContext, useState } from "react";
import { apiRequest } from "../api/client";
import { useAlert } from "./AlertContext";
import { useSchedule } from "./ScheduleContext";

const ChatContext = createContext();

export const useChat = () => {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error("useChat must be used within ChatProvider");
  }
  return context;
};

export const ChatProvider = ({ children }) => {
  const [messages, setMessages] = useState([]);
  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [pendingIntent, setPendingIntent] = useState(null);
  const { showSuccess, showError } = useAlert();
    const { refreshSchedules } = useSchedule();

  const addMessage = (message, isUser = false) => {
    setMessages((prev) => [...prev, { text: message, isUser, timestamp: new Date() }]);
  };

  const sendMessage = async (message) => {
    if (!message.trim()) return;

    addMessage(message, true);
    setIsLoading(true);

    try {
      const { ok, data } = await apiRequest("/chat/parse", {
        method: "POST",
        body: JSON.stringify({ message }),
      });

      if (ok) {
        addMessage(data.data.message, false);

        if (data.data.needs_confirm && data.data.intent) {
          setPendingIntent(data.data.intent);
        } else {
          setPendingIntent(null);
        }
      } else {
        addMessage("Lo siento, hubo un error al procesar tu mensaje.", false);
        showError(data.error || "Error al procesar el mensaje");
      }
    } catch (error) {
      addMessage("Error de conexión. Por favor, intenta de nuevo.", false);
      showError("Error de conexión");
    } finally {
      setIsLoading(false);
    }
  };

  const confirmAction = async () => {
    if (!pendingIntent) return;

    setIsLoading(true);
    addMessage("Confirmado, ejecutando...", true);

    try {
      const { ok, data } = await apiRequest("/chat/execute", {
        method: "POST",
        body: JSON.stringify({ intent: pendingIntent }),
      });

      if (ok) {
        addMessage(data.data.message, false);
        showSuccess(data.message);
        
        if (pendingIntent.action === "scrape_now") {
          window.dispatchEvent(new Event("reload-results"));
        } else if (pendingIntent.action === "create_schedule") {
          await refreshSchedules();
        }
      } else {
        addMessage("Error al ejecutar la acción: " + (data.error || "Error desconocido"), false);
        showError(data.error || "Error al ejecutar la acción");
      }
    } catch (error) {
      addMessage("Error de conexión al ejecutar la acción.", false);
      showError("Error de conexión");
    } finally {
      setIsLoading(false);
      setPendingIntent(null);
    }
  };

  const cancelAction = () => {
    addMessage("Acción cancelada", true);
    addMessage("De acuerdo, ¿en qué más puedo ayudarte?", false);
    setPendingIntent(null);
  };

  const clearChat = () => {
    setMessages([]);
    setPendingIntent(null);
  };

  const toggleChat = () => {
    setIsOpen((prev) => !prev);
    
    if (!isOpen && messages.length === 0) {
      addMessage("¡Hola! Puedo ayudarte a scrapear URLs o programar scraping automático. ¿Qué necesitas?", false);
    }
  };

  const value = {
    messages,
    isOpen,
    isLoading,
    pendingIntent,
    sendMessage,
    confirmAction,
    cancelAction,
    clearChat,
    toggleChat,
  };

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
};
