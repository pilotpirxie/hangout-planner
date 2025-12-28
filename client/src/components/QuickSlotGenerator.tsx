export const QuickSlotGenerator = ({
  startDate,
  endDate,
  dailyStartTime,
  dailyEndTime,
  duration,
  isOverlapping,
  isWholeDay,
  onStartDateChange,
  onEndDateChange,
  onDailyStartTimeChange,
  onDailyEndTimeChange,
  onDurationChange,
  onOverlappingChange,
  onWholeDayChange,
}: {
  startDate: string;
  endDate: string;
  dailyStartTime: string;
  dailyEndTime: string;
  duration: string;
  isOverlapping: boolean;
  isWholeDay: boolean;
  onStartDateChange: (date: string) => void;
  onEndDateChange: (date: string) => void;
  onDailyStartTimeChange: (time: string) => void;
  onDailyEndTimeChange: (time: string) => void;
  onDurationChange: (duration: string) => void;
  onOverlappingChange: (isOverlapping: boolean) => void;
  onWholeDayChange: (isWholeDay: boolean) => void;
}) => {
  return (
    <>
      <div className="mt-3">
        <label htmlFor="start-date">Start date</label>
        <input
          id="start-date"
          type="date"
          className="form-control"
          value={startDate}
          onChange={(e) => {
            onStartDateChange(e.target.value);
          }}
        />
      </div>

      <div className="mt-3">
        <label htmlFor="end-date">End date</label>
        <input
          id="end-date"
          type="date"
          className="form-control"
          value={endDate}
          onChange={(e) => {
            onEndDateChange(e.target.value);
          }}
        />
      </div>

      <div className="mt-3">
        <div className="form-check">
          <input
            id="whole-day"
            type="checkbox"
            className="form-check-input"
            checked={isWholeDay}
            onChange={(e) => {
              const checked = e.target.checked;
              onWholeDayChange(checked);
              if (checked) {
                onDailyStartTimeChange("00:00");
                onDailyEndTimeChange("23:59");
              } else {
                onDailyStartTimeChange("");
                onDailyEndTimeChange("");
              }
            }}
          />
          <label
            className="form-check-label"
            htmlFor="whole-day">
            Whole day
          </label>
        </div>
        <small className="text-muted">
          {isWholeDay
            ? "Slots will cover the entire day (00:00 - 23:59)"
            : "Specify custom time range for daily slots"}
        </small>
      </div>

      {!isWholeDay ? <>
        <div className="mt-3">
          <label htmlFor="daily-start-time">Daily start time</label>
          <input
            id="daily-start-time"
            type="time"
            className="form-control"
            value={dailyStartTime}
            onChange={(e) => {
              onDailyStartTimeChange(e.target.value);
            }}
          />
        </div>

        <div className="mt-3">
          <label htmlFor="daily-end-time">Daily end time</label>
          <input
            id="daily-end-time"
            type="time"
            className="form-control"
            value={dailyEndTime}
            onChange={(e) => {
              onDailyEndTimeChange(e.target.value);
            }}
          />
        </div>
      </> : null}

      {!isWholeDay ? <>
        <div className="mt-3">
          <label htmlFor="duration">Slot duration (hours)</label>
          <select
            id="duration"
            className="form-control"
            value={duration}
            onChange={(e) => {
              onDurationChange(e.target.value);
            }}>
            <option value="">Select duration</option>
            <option value="0.5">30 minutes</option>
            <option value="1">1 hour</option>
            <option value="1.5">1.5 hours</option>
            <option value="2">2 hours</option>
            <option value="2.5">2.5 hours</option>
            <option value="3">3 hours</option>
            <option value="4">4 hours</option>
          </select>
        </div>

        <div className="mt-3">
          <div className="form-check">
            <input
              id="overlapping"
              type="checkbox"
              className="form-check-input"
              checked={isOverlapping}
              onChange={(e) => {
                onOverlappingChange(e.target.checked);
              }}
            />
            <label
              className="form-check-label"
              htmlFor="overlapping">
              Allow overlapping time slots
            </label>
          </div>
          <small className="text-muted">
            {isOverlapping
              ? "Slots will overlap (e.g., 16-18, 17-19, 18-20) for maximum flexibility"
              : "Slots will be consecutive (e.g., 16-18, 18-20, 20-22) with no overlap"}
          </small>
        </div>
      </> : null}

      <div className="alert alert-info mt-3">
        <small>
          This will generate time slots for each day in the date range within the
          specified time window.
        </small>
      </div>
    </>
  );
};
