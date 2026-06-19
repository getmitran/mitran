import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import AuditLogPage from '../components/AuditLogPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('AuditLogPage', () => {
  it('renders', () => { render(<AuditLogPage />); expect(document.body.innerHTML).not.toBe(''); });
});
