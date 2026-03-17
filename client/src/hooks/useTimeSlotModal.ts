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
    const startDate = new Date(slot.startDate);
    const endDate = new Date(slot.endDate);
    setModalData({
      date: startDate.toISOString().split("T")[0],
      startTime: startDate.toTimeString().split(" ")[0].substring(0, 5),
      endTime: endDate.toTimeString().split(" ")[0].substring(0, 5),
    });
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingId(null);
    setModalData({ date: "", startTime: "", endTime: "" });
  };

  const handleSaveTimeSlot = (
    onSave: (slot: TimeSlot) => void,
  ) => {
    onSave({
      id: editingId || crypto.randomUUID(),
      startDate: `${modalData.date}T${modalData.startTime}:00`,
      endDate: `${modalData.date}T${modalData.endTime}:00`,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });

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
