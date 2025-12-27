interface CalendarHeaderProps {
  eventName: string;
}

export const CalendarHeader = ({ eventName }: CalendarHeaderProps) => {
  return (
    <div className="rounded-top p-3 d-flex flex-column align-items-center justify-content-center">
      <h1>{eventName}</h1>
      <h3 className="mb-0 text-center">Choose the date that works best for you</h3>
    </div>
  );
};
