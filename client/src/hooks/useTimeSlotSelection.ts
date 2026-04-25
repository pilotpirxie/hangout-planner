import { useMemo, useState } from "react";
import { useVoteOnTimeSlotMutation } from "../data/calendarsApi";
import type { TimeSlot } from "../types";

export const useTimeSlotSelection = (timeSlots: TimeSlot[], calendarId: string | undefined, password: string | undefined) => {
  const [selectedTimeSlotId, setSelectedTimeSlotId] = useState<string | null>(null);
  const [username, setUsername] = useState<string>("");
  const [vote] = useVoteOnTimeSlotMutation();

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
    if (!selectedTimeSlotId || !username || !calendarId) return;
    void vote({
      calendar_id: calendarId,
      time_slot_id: selectedTimeSlotId,
      username,
      password
    });
    setSelectedTimeSlotId(null);
    setUsername("");
  };

  return {
    selectedTimeSlotId,
    selectedTimeSlot,
    username,
    setUsername,
    handleClickTimeSlot,
    handleCloseModal,
    handleConfirm,
  };
};
