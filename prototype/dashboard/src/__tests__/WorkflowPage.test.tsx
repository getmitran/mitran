import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import WorkflowBuilderPage from '../components/WorkflowBuilderPage';

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) })));
});

describe('WorkflowBuilderPage', () => {
  it('renders loading state initially', () => {
    render(<WorkflowBuilderPage />);
    expect(screen.getByText('Loading workflows...')).toBeDefined();
  });
});
