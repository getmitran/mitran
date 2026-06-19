import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import AuditLogPage from '../components/AuditLogPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('AuditLogPage', () => {
  it('renders audit log heading', () => { render(<AuditLogPage />); expect(screen.getByText(/audit|log/i)).toBeDefined(); });
  it('fetches logs on mount', () => { render(<AuditLogPage />); expect(fetch).toHaveBeenCalled(); });
});
