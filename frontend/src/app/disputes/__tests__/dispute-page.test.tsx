import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import Disputes from "@/app/disputes/page";

describe('Disputes', () => {
  it('renders without crashing', () => {
    render(<Disputes />);
    expect(screen.getByText('Select a dispute to view it')).toBeInTheDocument();
  });
});