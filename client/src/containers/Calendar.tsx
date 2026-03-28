import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { CalendarHeader } from "../components/CalendarHeader";
import { DaySlotsModal } from "../components/DaySlotsModal";
import { MonthGridView } from "../components/MonthGridView";
import { TimeSlotConfirmationModal } from "../components/TimeSlotConfirmationModal";
import { WeekView } from "../components/WeekView";
import { calendarsApi } from "../data/calendarsApi";
import { useTimeSlotSelection } from "../hooks/useTimeSlotSelection";
import type { TimeSlot } from "../types";

export const Calendar = () => {
  const [viewMode, setViewMode] = useState<"week" | "month">("week");
  const [currentWeek, setCurrentWeek] = useState(dayjs().toDate());
  const [currentMonth, setCurrentMonth] = useState(dayjs().toDate());
  const [selectedDaySlots, setSelectedDaySlots] = useState<TimeSlot[] | null>(null);
  const [selectedDate, setSelectedDate] = useState<string>("");

  const params = useParams<{ id: string }>();
  const [getCalendar, { data: wholeCalendarDetails }] = calendarsApi.useLazyGetCalendarQuery();
  useEffect(() => {
    if (!params.id) return;
    void getCalendar({ calendar_id: params.id });
  }, [params.id, getCalendar]);

  const mappedTimeSlots = wholeCalendarDetails?.time_slots.map(slot => ({
    id: slot.id,
    startDate: dayjs(slot.start_date).toDate(),
    endDate: dayjs(slot.end_date).toDate(),
  })) ?? [];

  const {
    selectedTimeSlotId,
    selectedTimeSlot,
    nickname,
    setNickname,
    handleClickTimeSlot,
    handleCloseModal,
    handleConfirm,
  } = useTimeSlotSelection(mappedTimeSlots);

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

  if (!wholeCalendarDetails) return <div>Loading...</div>;

  return (
    <div className="bg-success vh-100 overflow-auto">
      <div
        className="container py-5">
        <div className="card">
          <CalendarHeader
            eventTitle={wholeCalendarDetails.calendar.title}
            eventDescription={wholeCalendarDetails.calendar.description}
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
                timeSlots={mappedTimeSlots}
                currentWeek={currentWeek}
                onTimeSlotClick={handleClickTimeSlot}
              />
            ) : (
              <MonthGridView
                timeSlots={mappedTimeSlots}
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