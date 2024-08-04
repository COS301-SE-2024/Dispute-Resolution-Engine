// import { render, screen, fireEvent } from '@testing-library/react';
const { render, screen, fireEvent } = require('@testing-library/react')
import Archive from "@/app/archive/page";
import { fetchArchiveHighlights } from '@/lib/api/archive';
import "@testing-library/jest-dom";
const { expect, describe, it } = require('@jest/globals')

jest.mock('@/lib/api/archive');
beforeEach(() => {
  (fetchArchiveHighlights as jest.Mock).mockResolvedValue({
    data: {
      archives: [
        {
          id: '1',
          title: 'Mock Dispute 1',
          date_resolved: '2024-07-31',
          category: ['Category 1', 'Category 2'],
          summary: 'This is a mock summary for dispute 1'
        },
      ]
    }
  });
})
describe('Archive Page', () => {
  it('renders without crashing', async () => {
    render(await Archive());
    expect(screen.getByText('Archive')).toBeInTheDocument();
  });

  it('can type into search bar', async () => {
    render(await Archive());
    const input = screen.getByPlaceholderText('Search the Archive...');
    fireEvent.change(input, { target: { value: 'test search' } });
    expect(input.value).toBe('test search');
  })

  it('renders archive highlights', async () => {
    render(await Archive());

    expect(screen.getByText('Archive')).toBeInTheDocument();
    expect(screen.getByText('Explore our previously handled cases')).toBeInTheDocument();
    expect(screen.getByText('Mock Dispute 1')).toBeInTheDocument();
    expect(screen.getByText('2024-07-31')).toBeInTheDocument();
    expect(screen.getByText('Category 1')).toBeInTheDocument();
    expect(screen.getByText('Category 2')).toBeInTheDocument();
    expect(screen.getByText('This is a mock summary for dispute 1')).toBeInTheDocument();
  })
});