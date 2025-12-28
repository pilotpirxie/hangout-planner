export const TimeSlotForm = ({
  date,
  startTime,
  endTime,
  onDateChange,
  onStartTimeChange,
  onEndTimeChange,
}: {
  date: string;
  startTime: string;
  endTime: string;
  onDateChange: (date: string) => void;
  onStartTimeChange: (time: string) => void;
  onEndTimeChange: (time: string) => void;
}) => {
  return (
    <>
      <div className="mt-3">
        <label htmlFor="date">Date</label>
        <input
          id="date"
          type="date"
          className="form-control"
          value={date}
          onChange={(e) => {
            onDateChange(e.target.value);
          }}
        />
      </div>

      <div className="mt-3">
        <label htmlFor="start-time">Start time</label>
        <input
          id="start-time"
          type="time"
          className="form-control"
          value={startTime}
          onChange={(e) => {
            onStartTimeChange(e.target.value);
          }}
        />
      </div>

      <div className="mt-3">
        <label htmlFor="end-time">End time</label>
        <input
          id="end-time"
          type="time"
          className="form-control"
          value={endTime}
          onChange={(e) => {
            onEndTimeChange(e.target.value);
          }}
        />
      </div>
    </>
  );
};
