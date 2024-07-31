// import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import Disputes from "@/app/disputes/page";
import DisputeRootLayout from "../layout";
import { getDisputeList } from "@/lib/api/dispute";
const { expect, describe, it } = require('@jest/globals')
const { render, screen, fireEvent } = require('@testing-library/react')

jest.mock("@/lib/api/dispute")
beforeEach(() => {
  (getDisputeList as jest.Mock).mockResolvedValue({
    data: [
        {
          id: "1",
          title: "Mock Title",
          description: "Mock Description lorem ipsum",
          status: "Mock Status",
          role: "Complainant",
        },
      ]
  });
})

describe('Disputes', () => {
  it('renders without crashing', () => {
    render(<Disputes />);
    expect(screen.getByText('Select a dispute to view it')).toBeInTheDocument();
  })
  it('layout renders correctly', () => {
    render(<DisputeRootLayout><div/></DisputeRootLayout>)
    expect(screen.getByText('Select a dispute to view it')).toBeInTheDocument();
  })
});