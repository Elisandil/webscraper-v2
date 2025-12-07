import React from "react";
import ScheduleList from "./ScheduleList";
import ScheduleModal from "../../modals/ScheduleModal";
import { useSchedule } from "../../../contexts/ScheduleContext";

export default function ScheduleSection() {
  const { openScheduleModal, editSchedule } = useSchedule();

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
          onClick={() => openScheduleModal()}
          className="bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 text-white font-medium py-3 px-6 rounded-lg transition-all duration-200 flex items-center gap-2 shadow-lg shadow-cyan-500/30 backdrop-blur-sm border border-cyan-500/20"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Nuevo Schedule
        </button>
      </div>

      <ScheduleList onEdit={editSchedule} />

      <ScheduleModal />
    </div>
  );
}