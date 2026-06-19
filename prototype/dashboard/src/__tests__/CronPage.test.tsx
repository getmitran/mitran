import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import CronPage from '../components/CronPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('CronPage', () => {
  it('renders', () => { render(<CronPage />); expect(document.body.innerHTML).not.toBe(''); });
});
