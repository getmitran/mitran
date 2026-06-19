import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import IntegrationsPage from '../components/IntegrationsPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('IntegrationsPage', () => {
  it('renders page heading', () => { render(<IntegrationsPage />); expect(screen.getByText(/integration/i)).toBeDefined(); });
  it('fetches integrations on mount', () => { render(<IntegrationsPage />); expect(fetch).toHaveBeenCalled(); });
});
