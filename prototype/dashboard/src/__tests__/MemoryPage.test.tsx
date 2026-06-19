import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import MemoryPage from '../components/MemoryPage';

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({ lessons: [], history: [] }) })));
});

describe('MemoryPage', () => {
  it('renders Memory heading', () => {
    render(<MemoryPage />);
    expect(screen.getByText('Memory')).toBeDefined();
  });
});
