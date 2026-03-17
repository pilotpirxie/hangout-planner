import { useState } from "react";
import type { TimeSlot } from "../types";
import { generateTimeSlots } from "../utils/generateTimeSlots";

export const useQuickSlotModal = () => {
  const [isQuickSlotModalOpen, setIsQuickSlotModalOpen] = useState(false);
  const [quickSlotModalData, setQuickSlotModalData] = useState({
    startDate: "",
    endDate: "",
    dailyStartTime: "",
    dailyEndTime: "",
    duration: "",
    isOverlapping: true,
    isWholeDay: false,
  });

  const handleOpenQuickSlotModal = () => {
    setQuickSlotModalData({
      startDate: "",
      endDate: "",
      dailyStartTime: "",
      dailyEndTime: "",
      duration: "",
      isOverlapping: true,
      isWholeDay: false,
    });
    setIsQuickSlotModalOpen(true);
  };

  const handleCloseQuickSlotModal = () => {
    setIsQuickSlotModalOpen(false);
  };

  const handleGenerateQuickSlots = (onGenerate: (slots: TimeSlot[]) => void) => {
    const slots = generateTimeSlots(quickSlotModalData);
    onGenerate(slots);
    handleCloseQuickSlotModal();
  };

  const isQuickSlotFormValid = quickSlotModalData.isWholeDay
    ? !!(quickSlotModalData.startDate && quickSlotModalData.endDate)
    : !!(quickSlotModalData.startDate &&
    quickSlotModalData.endDate &&
    quickSlotModalData.dailyStartTime &&
    quickSlotModalData.dailyEndTime &&
    quickSlotModalData.duration);

  return {
    isQuickSlotModalOpen,
    quickSlotModalData,
    setQuickSlotModalData,
    handleOpenQuickSlotModal,
    handleCloseQuickSlotModal,
    handleGenerateQuickSlots,
    isQuickSlotFormValid,
  };
};
