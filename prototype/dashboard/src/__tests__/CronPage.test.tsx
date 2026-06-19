import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import CronPage from '../components/CronPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('CronPage', () => {
  it('renders page heading', () => { render(<CronPage />); expect(screen.getByText(/cron/i)).toBeDefined(); });
  it('calls fetch on mount', () => { render(<CronPage />); expect(fetch).toHaveBeenCalled(); });
});
