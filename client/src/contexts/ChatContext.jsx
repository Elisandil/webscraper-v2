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
  // Stores context when the bot asks a follow-up question (missing URL or frequency).
  // Shape: { awaitingUrl: bool, awaitingFrequency: bool, partialIntent: ChatIntent }
  const [conversationState, setConversationState] = useState(null);
  const { showSuccess, showError } = useAlert();
  const { refreshSchedules } = useSchedule();

  const addMessage = (message, isUser = false) => {
    setMessages((prev) => [...prev, { text: message, isUser, timestamp: new Date() }]);
  };

  // Builds a full natural-language message when the bot previously asked for
  // a missing piece (URL or frequency), so the backend receives complete context.
  const buildEnrichedMessage = (userMessage, state) => {
    if (!state) return userMessage;
    const pi = state.partialIntent || {};
    if (state.awaitingUrl) {
      if (pi.action === "create_schedule") {
        const freq = pi.frequency || "diariamente";
        return `programa ${userMessage} ${freq}`;
      }
      return `scrapea ${userMessage} ahora`;
    }
    if (state.awaitingFrequency && pi.url) {
      return `programa ${pi.url} ${userMessage}`;
    }
    return userMessage;
  };

  const sendMessage = async (message) => {
    if (!message.trim()) return;

    const enriched = buildEnrichedMessage(message, conversationState);
    setConversationState(null);

    addMessage(message, true);
    setIsLoading(true);

    try {
      const { ok, data } = await apiRequest("/chat/parse", {
        method: "POST",
        body: JSON.stringify({ message: enriched }),
      });

      if (ok) {
        addMessage(data.data.message, false);

        if (data.data.needs_confirm && data.data.intent) {
          setPendingIntent(data.data.intent);
        } else if (data.data.action === "ask_url") {
          setConversationState({ awaitingUrl: true, partialIntent: data.data.intent });
          setPendingIntent(null);
        } else if (data.data.action === "ask_frequency") {
          setConversationState({ awaitingFrequency: true, partialIntent: data.data.intent });
          setPendingIntent(null);
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
        
        if (pendingIntent.action === "create_schedule") {
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
    setConversationState(null);
  };

  const clearChat = () => {
    setMessages([]);
    setPendingIntent(null);
    setConversationState(null);
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
    conversationState,
    sendMessage,
    confirmAction,
    cancelAction,
    clearChat,
    toggleChat,
  };

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
};
