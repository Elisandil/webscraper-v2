import React, { useEffect, useRef } from "react";

export default function ChatMessage({ message, isUser }) {
  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"} mb-4`}>
      <div
        className={`max-w-[80%] rounded-lg px-4 py-3 ${
          isUser
            ? "bg-gradient-to-r from-cyan-600 to-teal-600 text-white"
            : "bg-white/10 text-gray-100 border border-white/20"
        }`}
      >
        <p className="text-sm whitespace-pre-wrap">{message.text}</p>
        <span className="text-xs opacity-70 mt-1 block">
          {message.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
        </span>
      </div>
    </div>
  );
}
