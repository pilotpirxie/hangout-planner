import { emptyApi } from "./emptyApi";

export const calendarsApi = emptyApi.injectEndpoints({
  endpoints: (builder) => ({
    createCalendar: builder.mutation<{
      id: string;
      admin_token: string;
    }, {
      title: string;
      description?: string;
      password?: string;
    }>({
      query: (arg) => ({
        url: "/calendars",
        method: "POST",
        body: {
          title: arg.title,
          description: arg.description || undefined,
          password: arg.password || undefined,
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
      password?: string;
    }>({
      query: (arg) => {
        const queryParams = new URLSearchParams();
        if (arg.password) {
          queryParams.append("password", arg.password);
        }

        return {
          url: `/calendars/${arg.calendar_id}?${queryParams.toString()}`,
          method: "GET",
        };
      },
    }),
    checkIfCalendarPasswordProtected: builder.query<{
      is_password_protected: boolean
    }, {
      calendar_id: string
    }>({
      query: (arg) => ({
        url: `/calendars/${arg.calendar_id}/password-protected`,
        method: "GET",
      }),
    }),
    voteOnTimeSlot: builder.mutation<undefined, {
      calendar_id: string;
      time_slot_id: string;
      username: string;
      password?: string;
    }>({
      query: (arg) => ({
        url: `/calendars/${arg.calendar_id}/time-slots/${arg.time_slot_id}/votes`,
        method: "POST",
        body: {
          username: arg.username,
          password: arg.password || undefined,
        },
      }),
    }),
    getCalendarVotes: builder.query<
      {
        id: string;
        calendar_id: string;
        username: string;
        time_slot: {
          id: string;
          start_date: string;
          end_date: string;
        };
      }[],
    {
      calendar_id: string;
      password?: string;
    }>({
      query: (arg) => {
        const queryParams = new URLSearchParams();
        if (arg.password) {
          queryParams.append("password", arg.password);
        }

        return {
          url: `/calendars/${arg.calendar_id}/votes?${queryParams.toString()}`,
          method: "GET",
        };
      },
    }),
  }),
});

export const {
  useCreateCalendarMutation,
  useCreateCalendarTimeSlotsMutation,
  useGetCalendarQuery,
  useLazyCheckIfCalendarPasswordProtectedQuery,
  useCheckIfCalendarPasswordProtectedQuery,
  useVoteOnTimeSlotMutation,
  useLazyGetCalendarVotesQuery,
} = calendarsApi;