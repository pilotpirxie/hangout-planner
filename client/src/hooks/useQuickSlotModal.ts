import { useState } from "react";
import type { TimeSlot } from "../types";
import { generateTimeSlots } from "../utils/generateTimeSlots";

export const useQuickSlotModal = () => {
  const [isQuickModalOpen, setIsQuickModalOpen] = useState(false);
  const [quickData, setQuickData] = useState({
    startDate: "",
    endDate: "",
    dailyStartTime: "",
    dailyEndTime: "",
    duration: "",
    isOverlapping: true,
    isWholeDay: false,
  });

  const handleOpenQuickModal = () => {
    setQuickData({
      startDate: "",
      endDate: "",
      dailyStartTime: "",
      dailyEndTime: "",
      duration: "",
      isOverlapping: true,
      isWholeDay: false,
    });
    setIsQuickModalOpen(true);
  };

  const handleCloseQuickModal = () => {
    setIsQuickModalOpen(false);
    setQuickData({
      startDate: "",
      endDate: "",
      dailyStartTime: "",
      dailyEndTime: "",
      duration: "",
      isOverlapping: true,
      isWholeDay: false,
    });
  };

  const handleGenerateQuickSlots = (onGenerate: (slots: TimeSlot[]) => void) => {
    const slots = generateTimeSlots(quickData);
    onGenerate(slots);
    handleCloseQuickModal();
  };

  const isQuickFormValid = quickData.isWholeDay
    ? !!(quickData.startDate && quickData.endDate)
    : !!(quickData.startDate &&
    quickData.endDate &&
    quickData.dailyStartTime &&
    quickData.dailyEndTime &&
    quickData.duration);

  return {
    isQuickModalOpen,
    quickData,
    setQuickData,
    handleOpenQuickModal,
    handleCloseQuickModal,
    handleGenerateQuickSlots,
    isQuickFormValid,
  };
};
