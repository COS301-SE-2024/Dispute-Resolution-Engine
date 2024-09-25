"use server";

import { TicketListRequest, type Ticket, type TicketSummary } from "../types/tickets";

export async function getTicketSummaries(req: TicketListRequest): Promise<{
  total: number;
  tickets: TicketSummary[];
}> {
  let data = MOCK_TICKETS;
  if (req.search) {
    data = data.filter((t) => {
      const combined = [t.subject, t.user.full_name].join(" ").toLowerCase();
      return combined.includes(req.search!.toLowerCase());
    });
  }

  const total = data.length;
  data = data.slice(req.offset ?? 0, req.limit ? (req.offset ?? 0) + req.limit : undefined);

  return {
    total: total,
    tickets: data,
  };
}

const MOCK_TICKETS: Ticket[] = [
  {
    id: "0",
    user: { id: "0", full_name: "John Doe" },
    date_created: "2023-09-01T10:30:00Z",
    subject: "Unable to access account",
    status: "Open",
    body: "I cannot log in to my account. It keeps saying incorrect password.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "John Doe" },
        date_sent: "2023-09-01T10:31:00Z",
        message: "I cannot log in to my account. It keeps saying incorrect password.",
      },
    ],
  },
  {
    id: "1",
    user: { id: "0", full_name: "Jane Smith" },
    date_created: "2023-08-22T08:45:00Z",
    subject: "Billing issue on my account",
    status: "Solved",
    body: "I was overcharged for my subscription last month. Please assist.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Jane Smith" },
        date_sent: "2023-08-22T08:46:00Z",
        message: "I was overcharged for my subscription last month. Please assist.",
      },
      {
        id: "0",
        user: { id: "0", full_name: "Support Agent" },
        date_sent: "2023-08-23T09:00:00Z",
        message: "Your refund has been processed. Apologies for the inconvenience.",
      },
    ],
  },
  {
    id: "2",
    user: { id: "0", full_name: "Bob Johnson" },
    date_created: "2023-09-15T13:20:00Z",
    subject: "Feature request: Dark mode",
    status: "On Hold",
    body: "It would be great to have a dark mode option in the settings.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Bob Johnson" },
        date_sent: "2023-09-15T13:21:00Z",
        message: "It would be great to have a dark mode option in the settings.",
      },
      {
        id: "0",
        user: { id: "0", full_name: "Support Agent" },
        date_sent: "2023-09-16T11:30:00Z",
        message: "Thanks for the suggestion! We'll consider it for future updates.",
      },
    ],
  },
  {
    id: "3",
    user: { id: "0", full_name: "Alice Brown" },
    date_created: "2023-07-30T16:00:00Z",
    subject: "App crashes on startup",
    status: "Closed",
    body: "The app crashes every time I try to open it on my Android phone.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Alice Brown" },
        date_sent: "2023-07-30T16:02:00Z",
        message: "The app crashes every time I try to open it on my Android phone.",
      },
      {
        id: "0",
        user: { id: "0", full_name: "Support Agent" },
        date_sent: "2023-07-31T09:00:00Z",
        message: "We have released an update that should fix the issue. Please try again.",
      },
    ],
  },
  {
    id: "4",
    user: { id: "0", full_name: "Eve Davis" },
    date_created: "2023-08-10T14:50:00Z",
    subject: "Password reset email not received",
    status: "Open",
    body: "I requested a password reset, but I haven’t received the email yet.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Eve Davis" },
        date_sent: "2023-08-10T14:51:00Z",
        message: "I requested a password reset, but I haven’t received the email yet.",
      },
    ],
  },
  {
    id: "5",
    user: { id: "0", full_name: "Charlie Wilson" },
    date_created: "2023-09-05T12:10:00Z",
    subject: "Unable to upload profile picture",
    status: "Solved",
    body: "Whenever I try to upload a profile picture, it fails with an error.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Charlie Wilson" },
        date_sent: "2023-09-05T12:12:00Z",
        message: "Whenever I try to upload a profile picture, it fails with an error.",
      },
      {
        id: "0",
        user: { id: "0", full_name: "Support Agent" },
        date_sent: "2023-09-06T09:00:00Z",
        message: "The issue has been resolved. You should now be able to upload your picture.",
      },
    ],
  },
  {
    id: "6",
    user: { id: "0", full_name: "Daniel Martinez" },
    date_created: "2023-09-18T09:30:00Z",
    subject: "Cannot download the report",
    status: "On Hold",
    body: "I am unable to download my monthly report from the dashboard.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Daniel Martinez" },
        date_sent: "2023-09-18T09:32:00Z",
        message: "I am unable to download my monthly report from the dashboard.",
      },
    ],
  },
  {
    id: "7",
    user: { id: "0", full_name: "Fiona Clark" },
    date_created: "2023-08-29T11:25:00Z",
    subject: "Email notifications not working",
    status: "Closed",
    body: "I am not receiving email notifications for new messages.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Fiona Clark" },
        date_sent: "2023-08-29T11:26:00Z",
        message: "I am not receiving email notifications for new messages.",
      },
      {
        id: "0",
        user: { id: "0", full_name: "Support Agent" },
        date_sent: "2023-08-30T10:00:00Z",
        message: "We have fixed the email notification settings. Please check if it works now.",
      },
    ],
  },
  {
    id: "8",
    user: { id: "0", full_name: "George White" },
    date_created: "2023-07-18T13:15:00Z",
    subject: "System outage",
    status: "Closed",
    body: "The system was down for an hour. Was this expected?",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "George White" },
        date_sent: "2023-07-18T13:16:00Z",
        message: "The system was down for an hour. Was this expected?",
      },
      {
        id: "0",
        user: { id: "0", full_name: "Support Agent" },
        date_sent: "2023-07-18T14:00:00Z",
        message: "We had an unexpected outage due to server issues. Everything is back up now.",
      },
    ],
  },
  {
    id: "9",
    user: { id: "0", full_name: "Hannah Lee" },
    date_created: "2023-09-12T10:00:00Z",
    subject: "Account suspension",
    status: "Open",
    body: "My account was suspended without explanation. Please help me resolve this.",
    messages: [
      {
        id: "0",
        user: { id: "0", full_name: "Hannah Lee" },
        date_sent: "2023-09-12T10:01:00Z",
        message: "My account was suspended without explanation. Please help me resolve this.",
      },
    ],
  },
];
