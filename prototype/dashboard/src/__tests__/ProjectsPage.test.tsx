import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import ProjectsPage from '../pages/ProjectsPage';

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve([]) })));
});

describe('ProjectsPage', () => {
  it('renders New Project button', () => {
    render(<BrowserRouter><ProjectsPage /></BrowserRouter>);
    expect(screen.getByText('New Project')).toBeDefined();
  });
});
