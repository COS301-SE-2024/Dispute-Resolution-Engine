"use server";

import { ObjectionStatus, type ExpertStatus } from "../types";

export async function changeObjectionStatus(id: number, status: ObjectionStatus): Promise<void> {
  //   await sf(`${API_URL}/tickets/${id}`, {
  //     method: "PATCH",
  //     headers: {
  //       Authorization: `Bearer ${getAuthToken()}`,
  //     },
  //     body: JSON.stringify({
  //       status,
  //     }),
  //   });
  return;
}
