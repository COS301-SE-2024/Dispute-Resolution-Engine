import { API_URL } from '@/lib/utils';
import { searchArchive, fetchArchiveHighlights, fetchArchivedDispute } from '@/lib/api/archive';
import "@testing-library/jest-dom";
const { expect, describe, it } = require('@jest/globals')

global.fetch = jest.fn();

beforeEach(() => {
  (fetch as jest.Mock).mockClear();
});

describe('archive API functions', () => {
  it('searchArchive should return search results', async () => {
    const mockResponse = { data: 'some data' };
    (fetch as jest.Mock).mockResolvedValue({
      json: jest.fn().mockResolvedValue(mockResponse),
    });

    const params = { query: 'test' };
    const result = await searchArchive(params);

    expect(fetch).toHaveBeenCalledWith(`${API_URL}/archive/search`, {
      cache: 'no-store',
      method: 'POST',
      body: JSON.stringify(params),
    });
    expect(result).toEqual(mockResponse);
  });

  it('fetchArchiveHighlights should return highlights', async () => {
    const mockResponse = { data: 'some highlights' };
    (fetch as jest.Mock).mockResolvedValue({
      json: jest.fn().mockResolvedValue(mockResponse),
    });

    const limit = 10;
    const result = await fetchArchiveHighlights(limit);

    expect(fetch).toHaveBeenCalledWith(`${API_URL}/archive/highlights?limit=${limit}`, {
      cache: 'no-store',
    });
    expect(result).toEqual(mockResponse);
  });

  it('fetchArchivedDispute should return a specific dispute', async () => {
    const mockResponse = { data: 'specific dispute' };
    (fetch as jest.Mock).mockResolvedValue({
      json: jest.fn().mockResolvedValue(mockResponse),
    });

    const id = '123';
    const result = await fetchArchivedDispute(id);

    expect(fetch).toHaveBeenCalledWith(`${API_URL}/archive/${id}`);
    expect(result).toEqual(mockResponse);
  });

});