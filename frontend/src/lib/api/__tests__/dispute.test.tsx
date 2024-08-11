import {
  getDisputeList,
  getDisputeDetails,
  updateDisputeStatus,
  getStatusEnum,
} from '@/lib/api/dispute';
import { API_URL } from '@/lib/utils';
import { cookies } from 'next/headers';
const { expect, describe, it } = require('@jest/globals')

global.fetch = jest.fn();
jest.mock('next/headers', () => ({
  cookies: jest.fn(),
}));

beforeEach(() => {
  (fetch as jest.Mock).mockClear();
  (cookies as jest.Mock).mockClear();
});

describe('dispute API functions', () => {
  it('getDisputeList should return dispute list', async () => {
    const mockResponse = { data: 'some data' };
    (cookies as jest.Mock).mockReturnValue({
      get: jest.fn().mockReturnValue({ value: 'mock-jwt' }),
    });
    (fetch as jest.Mock).mockResolvedValue({
      json: jest.fn().mockResolvedValue(mockResponse),
    });

    const result = await getDisputeList();

    expect(fetch).toHaveBeenCalledWith(`${API_URL}/disputes`, {
      headers: {
        Authorization: 'Bearer mock-jwt',
      },
    });
    expect(result).toEqual(mockResponse);
  });

  it('getDisputeDetails should return dispute details', async () => {
    const mockResponse = { data: 'some details' };
    (cookies as jest.Mock).mockReturnValue({
      get: jest.fn().mockReturnValue({ value: 'mock-jwt' }),
    });
    (fetch as jest.Mock).mockResolvedValue({
      json: jest.fn().mockResolvedValue(mockResponse),
    });

    const id = '123';
    const result = await getDisputeDetails(id);

    expect(fetch).toHaveBeenCalledWith(`${API_URL}/disputes/${id}`, {
      headers: {
        Authorization: 'Bearer mock-jwt',
      },
    });
    expect(result).toEqual(mockResponse);
  });

  it('updateDisputeStatus should update dispute status', async () => {
    const mockResponse = { data: 'updated status' };
    (cookies as jest.Mock).mockReturnValue({
      get: jest.fn().mockReturnValue({ value: 'mock-jwt' }),
    });
    (fetch as jest.Mock).mockResolvedValue({
      json: jest.fn().mockResolvedValue(mockResponse),
    });

    const id = '123';
    const status = 'resolved';
    const result = await updateDisputeStatus(id, status);

    expect(fetch).toHaveBeenCalledWith(`${API_URL}/disputes/dispute/status`, {
      method: 'PUT',
      headers: {
        Authorization: 'Bearer mock-jwt',
      },
      body: JSON.stringify({ dispute_id: id, status }),
    });
    expect(result).toEqual(mockResponse);
  });

  it('getStatusEnum should return status enum', async () => {
    const mockResponse = { data: ['open', 'closed'] };
    (fetch as jest.Mock).mockResolvedValue({
      json: jest.fn().mockResolvedValue(mockResponse),
    });

    const result = await getStatusEnum();

    expect(fetch).toHaveBeenCalledWith(`${API_URL}/utils/dispute_statuses`, {
      method: 'GET',
    });
    expect(result).toEqual(mockResponse.data);
  });

  it('should handle fetch errors', async () => {
    const mockError = new Error('Fetch error');
    (fetch as jest.Mock).mockRejectedValue(mockError);
    (cookies as jest.Mock).mockReturnValue({
      get: jest.fn().mockReturnValue({ value: 'mock-jwt' }),
    });

    const result = await getDisputeList();

    expect(result).toEqual({ error: mockError.message });
  });

  // it('should handle unauthorized errors', async () => {
  //   (cookies as jest.Mock).mockReturnValue({
  //     get: jest.fn().mockReturnValue(undefined),
  //   });

  //   const result = await getDisputeList();

  //   expect(result).toEqual({ error: 'Unauthorized' });
  // });
});