import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import WebhooksPage from '../components/WebhooksPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('WebhooksPage', () => {
  it('renders webhooks heading', () => { render(<WebhooksPage />); expect(screen.getByText(/webhook/i)).toBeDefined(); });
  it('fetches webhooks on mount', () => { render(<WebhooksPage />); expect(fetch).toHaveBeenCalled(); });
});
