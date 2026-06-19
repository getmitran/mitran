import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import SettingsPage from '../components/SettingsPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({}) }))); });

describe('SettingsPage', () => {
  it('renders', () => { render(<SettingsPage />); expect(document.body.innerHTML).not.toBe(''); });
});
