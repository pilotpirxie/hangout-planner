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
  const [isPasswordProtected, setIsPasswordProtected] = useState<boolean | null>(null);
  const [password, setPassword] = useState<string>("");
  const params = useParams<{ id: string }>();

  const [checkIfPasswordProtected] = calendarsApi.useLazyCheckIfCalendarPasswordProtectedQuery();
  useEffect(() => {
    if (!params.id) return;
    checkIfPasswordProtected({ calendar_id: params.id }).then(({ data }) => {
      setIsPasswordProtected(data?.is_password_protected ?? false);
    }).catch(() => {
      console.error("Failed to check if calendar is password protected");
    });
  }, [params.id, checkIfPasswordProtected]);

  const [getCalendar, { data: wholeCalendarDetails }] = calendarsApi.useLazyGetCalendarQuery();
  useEffect(() => {
    if (!params.id || isPasswordProtected === null || isPasswordProtected) return;
    void getCalendar({ calendar_id: params.id });
  }, [params.id, getCalendar, password, isPasswordProtected]);

  const handleGetCalendarManually = () => {
    if (!params.id) return;
    void getCalendar({ calendar_id: params.id, password });
  };

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
  } = useTimeSlotSelection(mappedTimeSlots, params.id);

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

  if (isPasswordProtected && !wholeCalendarDetails) {
    return (
      <div className="bg-success vh-100 d-flex flex-column justify-content-center align-items-center">
        <h2 className="mb-4">This calendar is password protected</h2>
        <div
          className="input-group mb-3"
          style={{ maxWidth: "400px" }}>
          <input
            type="password"
            className="form-control"
            placeholder="Enter password"
            value={password}
            onChange={(e) => { setPassword(e.target.value); }}
          />
          <button
            className="btn btn-primary"
            onClick={() => { handleGetCalendarManually(); }}>
            Submit
          </button>
        </div>
      </div>
    );
  }

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