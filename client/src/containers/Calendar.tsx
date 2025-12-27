import { useEffect, useRef, useState } from "react";
import { CalendarHeader } from "../components/CalendarHeader";
import { DayColumn } from "../components/DayColumn";
import { useGroupedTimeSlots } from "../hooks/useGroupedTimeSlots";
import { useIsMobile } from "../hooks/useIsMobile";
import { useMockTimeSlots } from "../hooks/useMockTimeSlots";

const DESKTOP_SPEED = 550;
const MOBILE_SPEED = 250;

export const Calendar = () => {
  const timeSlots = useMockTimeSlots();
  const groupedSlots = useGroupedTimeSlots(timeSlots);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);

  const isMobile = useIsMobile();
  const scrollSpeed = isMobile ? MOBILE_SPEED : DESKTOP_SPEED;

  const checkScrollPosition = () => {
    const container = scrollContainerRef.current;
    if (!container) return;
    setCanScrollLeft(container.scrollLeft > 0);
    setCanScrollRight(container.scrollLeft < container.scrollWidth - container.clientWidth - 1);
  };

  const handleScroll = (direction: "left" | "right") => {
    const container = scrollContainerRef.current;
    if (!container) return;
    const newScrollLeft = container.scrollLeft + (direction === "left" ? -scrollSpeed : scrollSpeed);
    container.scrollTo({ left: newScrollLeft, behavior: "smooth" });
  };

  useEffect(() => {
    checkScrollPosition();
    const container = scrollContainerRef.current;
    if (!container) return;
    container.addEventListener("scroll", checkScrollPosition);

    return () => {
      container.removeEventListener("scroll", checkScrollPosition);
    };
  }, [groupedSlots]);

  return (
    <div className="bg-success vh-100 overflow-auto">
      <div
        className="container py-5">
        <div className="card">
          <CalendarHeader eventName="Event name" />
          <div className="card-body p-4 position-relative">
            {canScrollLeft ? <button
              className="btn btn-primary position-absolute top-0 start-0 translate-middle-y ms-5 z-1"
              onClick={() => { handleScroll("left"); }}>
              <i className="ri-arrow-left-long-line" />
            </button> : null}

            <div
              className="d-flex gap-3 overflow-x-scroll pb-3"
              style={{ scrollbarWidth: "thin" }}
              ref={scrollContainerRef}>
              {groupedSlots.map(day => (
                <DayColumn
                  key={day.date}
                  date={day.date}
                  slots={day.slots}
                />
              ))}
            </div>

            {canScrollRight ? <button
              className="btn btn-primary position-absolute top-0 end-0 translate-middle-y me-5"
              onClick={() => { handleScroll("right"); }}>
              <i className="ri-arrow-right-long-line" />
            </button> : null}
          </div>
        </div>
      </div>
    </div>
  );
};