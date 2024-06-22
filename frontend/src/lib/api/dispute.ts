"use server";

import { Result } from "@/lib/types";
import { DisputeListResponse, DisputeResponse } from "../interfaces/dispute";
import { cookies } from "next/headers";
import { JWT_KEY } from "../constants";

export async function getDisputeList(): Promise<Result<DisputeListResponse>> {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return {
      error: "Unauthorized",
    };
  }

  return {
    data: [...Array(10).keys()].map((i) => ({
      id: i.toString(),
      title: `Dispute ${i}`,
      description: "Lorem ipsum",
      status: "active",
    })),
  };

  // TODO: Uncomment once API is working
  // return fetch(`${API_URL}/disputes`, {
  //   headers: {
  //     Authorization: `Bearer ${jwt}`,
  //   },
  // })
  //   .then((res) => res.json())
  //   .catch((e: Error) => ({
  //     error: e.message,
  //   }));
}

export async function getDisputeDetails(id: string): Promise<Result<DisputeResponse>> {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return {
      error: "Unauthorized",
    };
  }

  return {
    data: {
      id: id,
      title: `Dispute ${id}`,
      description: "Dispute description",
      status: "status",
      date_created: "",
      evidence: [...Array(5).keys()].map((i) => ({
        label: `Image ${i}`,
        url: "https://picsum.photos/200",
        date_submitted: "today",
      })),
      experts: [...Array(3).keys()].map((i) => ({
        full_name: `Name ${i}`,
        email: `coolguy${i}@example.com`,
        phone: "phone number yes",
      })),
    },
  };

  // TODO: Uncomment once API is working
  // return fetch(`${API_URL}/disputes/${id}`, {
  //   headers: {
  //     Authorization: `Bearer ${jwt}`,
  //   },
  // })
  //   .then((res) => res.json())
  //   .catch((e: Error) => ({
  //     error: e.message,
  //   }));
}
