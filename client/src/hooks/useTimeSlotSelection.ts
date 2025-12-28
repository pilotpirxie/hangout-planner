import { useMemo, useState } from "react";
import type { TimeSlot } from "../types";

export const useTimeSlotSelection = (timeSlots: TimeSlot[]) => {
  const [selectedTimeSlotId, setSelectedTimeSlotId] = useState<string | null>(null);
  const [nickname, setNickname] = useState<string>("");

  const selectedTimeSlot = useMemo(() => {
    return timeSlots.find(slot => slot.id === selectedTimeSlotId) ?? null;
  }, [selectedTimeSlotId, timeSlots]);

  const handleClickTimeSlot = (timeSlotId: string) => {
    setSelectedTimeSlotId(timeSlotId);
  };

  const handleCloseModal = () => {
    setSelectedTimeSlotId(null);
  };

  const handleConfirm = () => {
    setSelectedTimeSlotId(null);
  };

  return {
    selectedTimeSlotId,
    selectedTimeSlot,
    nickname,
    setNickname,
    handleClickTimeSlot,
    handleCloseModal,
    handleConfirm,
  };
};
