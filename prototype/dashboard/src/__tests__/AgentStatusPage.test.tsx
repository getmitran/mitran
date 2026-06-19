import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import AgentStatusPage from '../components/AgentStatusPage';

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) })));
});

describe('AgentStatusPage', () => {
  it('renders loading state initially', () => {
    render(<AgentStatusPage />);
    expect(screen.getByText('Loading agents...')).toBeDefined();
  });
});
