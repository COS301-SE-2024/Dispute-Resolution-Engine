import { render, screen, fireEvent } from "@testing-library/react";
import { expect, describe, it } from "@jest/globals";
import ExpertRejectForm from "../expert-reject-form";

describe("Expert reject form", () => {
  const props = {
    disputeId: "dispute",
    expertId: "expert",
    name: "name",
  };

  it("should render", () => {
    render(<ExpertRejectForm {...props} />);
  });

  it("shows reject button", () => {
    render(<ExpertRejectForm {...props} />);
    expect(screen.getByRole("button")).toBeDefined();
    expect(screen.getByText("Reject")).toBeDefined();
  });

  it("pops up dialog when button clicked", () => {
    render(<ExpertRejectForm {...props} />);
    fireEvent.click(screen.getByText("Reject"));

    expect(screen.getByRole("dialog")).toBeDefined();
  });
});
