export interface Calendar {
  id: string;
  title: string;
  description?: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface TimeSlot {
  id: string;
  startDate: Date;
  endDate: Date;
}
