import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import WebhooksPage from '../components/WebhooksPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('WebhooksPage', () => {
  it('renders', () => { render(<WebhooksPage />); expect(document.body.innerHTML).not.toBe(''); });
});
