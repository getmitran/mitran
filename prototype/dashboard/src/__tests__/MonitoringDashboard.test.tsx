import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import MonitoringPage from '../components/MonitoringPage';

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({}) })));
});

describe('MonitoringPage', () => {
  it('renders page heading', () => {
    render(<MonitoringPage />);
    expect(screen.getByText('Monitoring')).toBeDefined();
  });
});
