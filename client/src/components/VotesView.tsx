import dayjs from "dayjs";
import { useState } from "react";

export const VotesView = ({
  votes
}: {
  votes: {
    id: string;
    calendar_id: string;
    username: string;
    time_slot: {
      id: string;
      start_date: string;
      end_date: string;
    };
  }[];
}) => {
  const [sortBy, setSortBy] = useState<"date" | "votes">("votes");

  const groupedTimeSlots = new Map<string, {
    startDate: string;
    endDate: string;
    votes: {
      id: string;
      username: string;
    }[];
  }>();

  votes.forEach(vote => {
    const slotId = vote.time_slot.id;
    if (!groupedTimeSlots.has(slotId)) {
      groupedTimeSlots.set(slotId, {
        startDate: vote.time_slot.start_date,
        endDate: vote.time_slot.end_date,
        votes: [],
      });
    }
    groupedTimeSlots.get(slotId)?.votes.push({
      id: vote.id,
      username: vote.username,
    });
  });

  const sortedVotes = Array.from(groupedTimeSlots.values()).sort((a, b) => {
    if (sortBy === "date") {
      return dayjs(a.startDate).diff(dayjs(b.startDate));
    } else {
      return b.votes.length - a.votes.length;
    }
  });

  return (
    <div>
      <div className="d-flex gap-2 mb-3 align-items-center">
        <span>Sort by:</span>
        <div className="btn-group">
          <button
            className={`btn btn-sm ${sortBy === "date" ? "btn-primary" : "btn-outline-primary"}`}
            onClick={() => { setSortBy("date"); }}>
            Date
          </button>
          <button
            className={`btn btn-sm ${sortBy === "votes" ? "btn-primary" : "btn-outline-primary"}`}
            onClick={() => { setSortBy("votes"); }}>
            Number of votes
          </button>
        </div>
      </div>

      {sortedVotes.map((vote, index) => (
        <div
          key={vote.startDate + vote.endDate}
          className="card mb-2">
          <div className="card-header d-flex align-items-center">
            <b className="me-1">#{index + 1}</b>
            <div className="d-flex align-items-center">
              {dayjs(vote.startDate).format("DD-MM-YYYY")} {dayjs(vote.startDate).format("HH:mm")} - {dayjs(vote.endDate).format("HH:mm")}
            </div>
          </div>

          <div className="card-body">
            <div className="gap-1 d-flex flex-column">
              {vote.votes.map(v => (
                <div
                  key={v.id}
                  className="d-flex align-items-center">
                  <div className="bg-info avatar me-1" />
                  <div>
                    {v.username}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};