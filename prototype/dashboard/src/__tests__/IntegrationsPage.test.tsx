import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import IntegrationsPage from '../components/IntegrationsPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({}) }))); });

describe('IntegrationsPage', () => {
  it('renders', () => { render(<IntegrationsPage />); expect(document.body.innerHTML).not.toBe(''); });
});
