import { render, screen, fireEvent } from "@testing-library/react";
import ExpertItem, { ExpertItemProps } from "../negotiator";
import { expect, describe, it } from "@jest/globals";

describe("Expert Item", () => {
  const mockDispute: ExpertItemProps = {
    dispute_id: "dispute",
    id: "expert",
    role: "role",
    full_name: "john smith",
    email: "email",
    phone: "phone",
  };

  it("should not crash", () => {
    render(<ExpertItem {...mockDispute} />);
  });

  it("should show expert details", () => {
    render(<ExpertItem {...mockDispute} />);

    expect(screen.getByText(mockDispute.role)).toBeDefined();
    expect(screen.getByText(mockDispute.full_name)).toBeDefined();
    expect(screen.getByText(mockDispute.email)).toBeDefined();
    expect(screen.getByText(mockDispute.phone)).toBeDefined();
  });

  it("should have reject button", () => {
    render(<ExpertItem {...mockDispute} />);
    expect(screen.getByRole("button")).toBeDefined();
  });
});
