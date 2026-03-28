import dayjs from "dayjs";
import { emptyApi } from "./emptyApi";

export const calendarsApi = emptyApi.injectEndpoints({
  endpoints: (builder) => ({
    createCalendar: builder.mutation<{
      id: string;
      admin_token: string;
    }, {
      title: string;
      description?: string;
      location?: string;
      accept_responses_until?: string;
      password?: string;
    }>({
      query: (arg) => ({
        url: "/calendars",
        method: "POST",
        body: {
          title: arg.title,
          description: arg.description || undefined,
          location: arg.location || undefined,
          password: arg.password || undefined,
          accept_responses_until: arg.accept_responses_until
            ? dayjs(arg.accept_responses_until).toISOString()
            : undefined,
        }
      }),
    }),
    createCalendarTimeSlots: builder.mutation<undefined, {
      calendar_id: string;
      admin_token: string;
      time_slots: {
        start_date: string,
        end_date: string
      }[]
    }>({
      query: (arg) => ({
        url: `/calendars/${arg.calendar_id}/time-slots`,
        method: "POST",
        body: {
          admin_token: arg.admin_token,
          time_slots: arg.time_slots.map(slot => {
            const parseLocalDateTime = (dateTimeStr: string) => {
              const [datePart, timePart] = dateTimeStr.split("T");
              const [year, month, day] = datePart.split("-").map(Number);
              const [hour, minute, second] = timePart.split(":").map(Number);
              return new Date(year, month - 1, day, hour, minute, second || 0);
            };

            const startLocal = parseLocalDateTime(slot.start_date);
            const endLocal = parseLocalDateTime(slot.end_date);

            return {
              start_date: startLocal.toISOString(),
              end_date: endLocal.toISOString(),
            };
          }),
        }
      }),
    }),
    getCalendar: builder.query<{
      calendar: {
        id: string;
        title: string;
        description?: string;
        location?: string;
        acceptResponsesUntil?: string;
        createdAt: string;
        updatedAt: string;
      };
      time_slots: {
        id: string;
        start_date: string;
        end_date: string;
      }[];
    }, {
      calendar_id: string;
    }>({
      query: (arg) => ({
        url: `/calendars/${arg.calendar_id}`,
        method: "GET",
      }),
    }),
  })
});

export const {
  useCreateCalendarMutation,
  useCreateCalendarTimeSlotsMutation,
  useGetCalendarQuery,
} = calendarsApi;