import React, { Suspense } from "react";
const { render, screen, act } = require("@testing-library/react");
import DisputeRootLayout from "@/app/disputes/layout";
import { getDisputeList } from "@/lib/api/dispute";
import "@testing-library/jest-dom";
const {describe, it } = require("@jest/globals");
jest.mock("@/lib/api/dispute", () => ({
  getDisputeList: jest.fn(),
}));

describe("DisputeRootLayout", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders the search input", async () => {
    const mockData = {
      data: [
        {
          id: "1",
          title: "Dispute 1",
          description: "Mock Desctip",
          role: "Complainant",
          status: "Awaiting respondant",
        },
      ],
    };
    (getDisputeList as jest.Mock).mockResolvedValue(Promise.resolve(mockData));
    render(await <DisputeRootLayout>Test 2</DisputeRootLayout>);
    // expect(await screen.findByText('Dispute 1', {}, {timeout: 10000})).toBeInTheDocument();
}, 11000);
});
