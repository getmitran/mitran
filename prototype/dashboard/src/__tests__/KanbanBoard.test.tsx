import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import KanbanBoard from '../components/KanbanBoard';

beforeEach(() => { vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({ columns: [] }) }))); });

describe('KanbanBoard', () => {
  it('renders board heading', () => { render(<KanbanBoard />); expect(screen.getByText(/kanban|board|tasks/i)).toBeDefined(); });
  it('fetches board data on mount', () => { render(<KanbanBoard />); expect(fetch).toHaveBeenCalled(); });
});
