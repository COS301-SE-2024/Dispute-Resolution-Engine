import { render, screen, fireEvent } from '@testing-library/react';
import Archive from "@/app/archive/page";
import "@testing-library/jest-dom";

describe('Archive Page', () => {
  it('renders without crashing', () => {
    render(<Archive />);
    expect(screen.getByText('Archive')).toBeInTheDocument();
  });

  it('can type into search bar', () => {
    render(<Archive />);
    const input = screen.getByPlaceholderText('Search the Archive...');
    fireEvent.change(input, { target: { value: 'test search' } });
    expect(input.value).toBe('test search');
  });
});