import React, { useState, useEffect } from "react";
import ScheduleList from "./ScheduleList";
import ScheduleModal from "./ScheduleModal";

export default function ScheduleSection({ onAlert }) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedSchedule, setSelectedSchedule] = useState(null);

  const handleNewSchedule = () => {
    setSelectedSchedule(null);
    setIsModalOpen(true);
  };

  const handleEditSchedule = (schedule) => {
    setSelectedSchedule(schedule);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedSchedule(null);
  };

  const handleSuccess = (message) => {
    onAlert({ type: "success", message });
  };

  const handleError = (message) => {
    onAlert({ type: "error", message });
  };

  useEffect(() => {
    const handleOpenScheduleModal = () => {
      setSelectedSchedule(null);
      setIsModalOpen(true);
      const event = new CustomEvent('switchToSchedules');
      window.dispatchEvent(event);
    };

    window.addEventListener('openScheduleModal', handleOpenScheduleModal);
    return () => {
      window.removeEventListener('openScheduleModal', handleOpenScheduleModal);
    };
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold text-white mb-2">Schedules Automatizados</h2>
          <p className="text-gray-400">
            Programa scraping automático de sitios web con expresiones cron
          </p>
        </div>
        
        <button
          onClick={handleNewSchedule}
          className="bg-blue-600/90 hover:bg-blue-600 text-white font-medium py-3 px-6 rounded-lg transition-all duration-200 flex items-center gap-2 shadow-lg backdrop-blur-sm border border-blue-500/20"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Nuevo Schedule
        </button>
      </div>

      <ScheduleList 
        onEdit={handleEditSchedule}
        onAlert={onAlert}
      />

      <ScheduleModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        schedule={selectedSchedule}
        onSuccess={handleSuccess}
        onError={handleError}
      />
    </div>
  );
}