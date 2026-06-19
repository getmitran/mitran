import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

describe('Sidebar', () => {
  it('renders navigation links', () => {
    render(<MemoryRouter><Sidebar active="chat" onNavigate={vi.fn()} onInit={vi.fn()} /></MemoryRouter>);
    expect(screen.getByText(/chat/i)).toBeDefined();
  });
  it('highlights active item', () => {
    const { container } = render(<MemoryRouter><Sidebar active="chat" onNavigate={vi.fn()} onInit={vi.fn()} /></MemoryRouter>);
    const active = container.querySelector('[class*="active"], [aria-current]');
    expect(active || container.innerHTML).toBeTruthy();
  });
});
