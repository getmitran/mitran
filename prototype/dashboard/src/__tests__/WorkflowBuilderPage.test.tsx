import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import WorkflowBuilderPage from '../components/WorkflowBuilderPage';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({}) }))); });

describe('WorkflowBuilderPage', () => {
  it('renders', () => { render(<WorkflowBuilderPage />); expect(document.body.innerHTML).not.toBe(''); });
});
