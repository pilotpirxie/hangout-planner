import dayjs from "dayjs";
import { useState } from "react";
import { CalendarHeader } from "../components/CalendarHeader";
import { DaySlotsModal } from "../components/DaySlotsModal";
import { MonthGridView } from "../components/MonthGridView";
import { TimeSlotConfirmationModal } from "../components/TimeSlotConfirmationModal";
import { WeekView } from "../components/WeekView";
import { useMockTimeSlots } from "../hooks/useMockTimeSlots";
import { useTimeSlotSelection } from "../hooks/useTimeSlotSelection";
import type { TimeSlot } from "../types";

type ViewMode = "week" | "month";

export const Calendar = () => {
  const [viewMode, setViewMode] = useState<ViewMode>("week");
  const [currentWeek, setCurrentWeek] = useState(dayjs().toDate());
  const [currentMonth, setCurrentMonth] = useState(dayjs().toDate());
  const [selectedDaySlots, setSelectedDaySlots] = useState<TimeSlot[] | null>(null);
  const [selectedDate, setSelectedDate] = useState<string>("");

  const timeSlots = useMockTimeSlots();

  const {
    selectedTimeSlotId,
    selectedTimeSlot,
    nickname,
    setNickname,
    handleClickTimeSlot,
    handleCloseModal,
    handleConfirm,
  } = useTimeSlotSelection(timeSlots);

  const handleWeekChange = (direction: "prev" | "next") => {
    setCurrentWeek((prev) => {
      const newWeek = dayjs(prev);
      if (direction === "prev") {
        return newWeek.subtract(7, "days").toDate();
      } else {
        return newWeek.add(7, "days").toDate();
      }
    });
  };

  const handleMonthChange = (direction: "prev" | "next") => {
    setCurrentMonth((prev) => {
      const newMonth = dayjs(prev);
      if (direction === "prev") {
        return newMonth.subtract(1, "month").toDate();
      } else {
        return newMonth.add(1, "month").toDate();
      }
    });
  };

  const handleDayClick = (slots: TimeSlot[], date: string) => {
    if (slots.length === 1) {
      handleClickTimeSlot(slots[0].id);
    } else {
      setSelectedDaySlots(slots);
      setSelectedDate(date);
    }
  };

  const handleCloseDaySlotsModal = () => {
    setSelectedDaySlots(null);
    setSelectedDate("");
  };

  const handleGoToToday = () => {
    const today = dayjs().toDate();
    setCurrentWeek(today);
    setCurrentMonth(today);
  };

  return (
    <div className="bg-success vh-100 overflow-auto">
      <div
        className="container py-5">
        <div className="card">
          <CalendarHeader
            eventName="Event name"
            viewMode={viewMode}
            onViewModeChange={setViewMode}
            currentWeek={currentWeek}
            currentMonth={currentMonth}
            onWeekChange={handleWeekChange}
            onMonthChange={handleMonthChange}
            onGoToToday={handleGoToToday}
          />
          <div className="card-body px-4">
            {viewMode === "week" ? (
              <WeekView
                timeSlots={timeSlots}
                currentWeek={currentWeek}
                onTimeSlotClick={handleClickTimeSlot}
              />
            ) : (
              <MonthGridView
                timeSlots={timeSlots}
                currentMonth={currentMonth}
                onDayClick={handleDayClick}
              />
            )}
          </div>
        </div>
      </div>

      <DaySlotsModal
        isVisible={!!selectedDaySlots}
        date={selectedDate}
        slots={selectedDaySlots ?? []}
        onClose={handleCloseDaySlotsModal}
        onSelectSlot={handleClickTimeSlot}
      />

      <TimeSlotConfirmationModal
        isVisible={!!selectedTimeSlotId}
        selectedTimeSlot={selectedTimeSlot}
        nickname={nickname}
        onNicknameChange={setNickname}
        onClose={handleCloseModal}
        onConfirm={handleConfirm}
      />
    </div>
  );
};