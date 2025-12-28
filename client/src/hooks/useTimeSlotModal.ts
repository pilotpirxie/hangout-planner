import { useState } from "react";
import type { TimeSlot } from "../types";

export const useTimeSlotModal = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [modalData, setModalData] = useState({
    date: "",
    startTime: "",
    endTime: "",
  });

  const handleOpenModal = () => {
    setEditingId(null);
    setModalData({ date: "", startTime: "", endTime: "" });
    setIsModalOpen(true);
  };

  const handleOpenEditModal = (slot: TimeSlot) => {
    setEditingId(slot.id);
    setModalData({
      date: slot.date,
      startTime: slot.startTime,
      endTime: slot.endTime,
    });
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingId(null);
    setModalData({ date: "", startTime: "", endTime: "" });
  };

  const handleSaveTimeSlot = (
    onAdd: (slot: TimeSlot) => void,
    onUpdate: (id: string, slot: TimeSlot) => void
  ) => {
    if (editingId) {
      const updatedSlot: TimeSlot = {
        id: editingId,
        date: modalData.date,
        startTime: modalData.startTime,
        endTime: modalData.endTime,
      };
      onUpdate(editingId, updatedSlot);
    } else {
      const newSlot: TimeSlot = {
        id: crypto.randomUUID(),
        date: modalData.date,
        startTime: modalData.startTime,
        endTime: modalData.endTime,
      };
      onAdd(newSlot);
    }
    handleCloseModal();
  };

  const isFormValid = modalData.date && modalData.startTime && modalData.endTime;

  return {
    isModalOpen,
    modalData,
    setModalData,
    editingId,
    handleOpenModal,
    handleOpenEditModal,
    handleCloseModal,
    handleSaveTimeSlot,
    isFormValid,
  };
};
