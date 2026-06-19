import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

describe('Sidebar', () => {
  it('renders nav links', () => {
    const { container } = render(
      <MemoryRouter>
        <Sidebar active="chat" onNavigate={vi.fn()} onInit={vi.fn()} />
      </MemoryRouter>
    );
    expect(container.innerHTML).not.toBe('');
  });
});
