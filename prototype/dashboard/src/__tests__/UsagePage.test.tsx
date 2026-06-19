import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import UsagePage from '../pages/UsagePage';

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({
    ok: true,
    json: () => Promise.resolve({ cpu: 45, memory: 62, storage: 30 })
  })));
});

describe('UsagePage', () => {
  it('renders Usage & Billing heading', () => {
    render(<BrowserRouter><UsagePage /></BrowserRouter>);
    expect(screen.getByText('Usage & Billing')).toBeDefined();
  });
});
