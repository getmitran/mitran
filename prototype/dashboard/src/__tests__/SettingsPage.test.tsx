import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import SettingsPage from '../components/SettingsPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({}) }))); });

describe('SettingsPage', () => {
  it('renders settings heading', () => { render(<SettingsPage />); expect(screen.getByText(/settings/i)).toBeDefined(); });
  it('fetches config on mount', () => { render(<SettingsPage />); expect(fetch).toHaveBeenCalled(); });
});
