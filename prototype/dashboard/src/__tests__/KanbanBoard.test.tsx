import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import KanbanBoard from '../components/KanbanBoard';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) }))); });

describe('KanbanBoard', () => {
  it('renders', () => { render(<KanbanBoard tasks={[]} refetch={() => {}} />); expect(document.body.innerHTML).not.toBe(''); });
});
