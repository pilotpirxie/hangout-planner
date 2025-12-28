import { useState } from "react";
import type { TimeSlot } from "../types";

export const useHangoutForm = () => {
  const [timeSlots, setTimeSlots] = useState<TimeSlot[]>([]);

  const addTimeSlot = (slot: TimeSlot) => {
    setTimeSlots(prev => [...prev, slot]);
  };

  const updateTimeSlot = (id: string, updatedSlot: TimeSlot) => {
    setTimeSlots(prev => prev.map(slot => slot.id === id ? updatedSlot : slot));
  };

  const addTimeSlots = (slots: TimeSlot[]) => {
    setTimeSlots(prev => [...prev, ...slots]);
  };

  const handleDeleteTimeSlot = (id: string) => {
    setTimeSlots(prev => prev.filter(slot => slot.id !== id));
  };

  const handleClearAllTimeSlots = () => {
    setTimeSlots([]);
  };

  return {
    timeSlots,
    addTimeSlot,
    updateTimeSlot,
    addTimeSlots,
    handleDeleteTimeSlot,
    handleClearAllTimeSlots,
  };
};
