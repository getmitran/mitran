import { useState } from 'react'
import { ListOrdered, Bot, ShieldCheck, Ticket, BookOpen, Zap, LayoutGrid, Settings, MessageSquare, History, FolderOpen, GitBranch, Puzzle, AppWindow, Link2, Webhook, DollarSign, Rocket, FileSearch, Users, CalendarDays, ChevronDown, BarChart3 } from 'lucide-react'
import ThemeToggle from './ThemeToggle'

type Section = 'Core' | 'Platform' | 'Operations' | 'Content' | 'System'

const nav = [
  { id: 'queue', label: 'Queue', icon: ListOrdered, section: 'Core' as Section },
  { id: 'kanban', label: 'Kanban', icon: LayoutGrid, section: 'Core' as Section },
  { id: 'chat', label: 'Chat', icon: MessageSquare, section: 'Core' as Section },
  { id: 'sessions', label: 'Sessions', icon: History, section: 'Core' as Section },
  { id: 'projects', label: 'Projects', icon: FolderOpen, section: 'Platform' as Section },
  { id: 'workflows', label: 'Workflows', icon: GitBranch, section: 'Platform' as Section },
  { id: 'agents', label: 'Agents', icon: Bot, section: 'Platform' as Section },
  { id: 'plugins', label: 'Plugins', icon: Puzzle, section: 'Platform' as Section },
  { id: 'apps', label: 'Apps', icon: AppWindow, section: 'Platform' as Section },
  { id: 'monitoring', label: 'Monitoring', icon: BarChart3, section: 'Operations' as Section },
  { id: 'integrations', label: 'Integrations', icon: Link2, section: 'Operations' as Section },
  { id: 'webhooks', label: 'Webhooks', icon: Webhook, section: 'Operations' as Section },
  { id: 'usage', label: 'Usage', icon: DollarSign, section: 'Operations' as Section },
  { id: 'sprints', label: 'Sprints', icon: CalendarDays, section: 'Operations' as Section },
  { id: 'tickets', label: 'Tickets', icon: Ticket, section: 'Content' as Section },
  { id: 'wiki', label: 'Wiki', icon: BookOpen, section: 'Content' as Section },
  { id: 'hr', label: 'HR Portal', icon: Users, section: 'Content' as Section },
  { id: 'audit', label: 'Audit Log', icon: FileSearch, section: 'Content' as Section },
  { id: 'onboarding', label: 'Onboarding', icon: Rocket, section: 'System' as Section },
  { id: 'settings', label: 'Settings', icon: Settings, section: 'System' as Section },
  { id: 'checkpoints', label: 'Checkpoints', icon: ShieldCheck, section: 'System' as Section },
] as const

const sections: Section[] = ['Core', 'Platform', 'Operations', 'Content', 'System']

type View = (typeof nav)[number]['id']

interface Props {
  active: View
  onNavigate: (v: View) => void
  onInit: () => void
}

export default function Sidebar({ active, onNavigate, onInit }: Props) {
  const [collapsed, setCollapsed] = useState<Record<Section, boolean>>({
    Core: false, Platform: false, Operations: false, Content: false, System: false,
  })

  const toggle = (s: Section) => setCollapsed(prev => ({ ...prev, [s]: !prev[s] }))

  return (
    <aside className="w-56 bg-surface border-r border-gray-800 flex flex-col">
      <div className="px-5 py-5 border-b border-gray-800">
        <h1 className="text-xl font-bold text-accent tracking-tight">mitran</h1>
        <p className="text-[11px] text-gray-500 mt-0.5">agent orchestration</p>
      </div>
      <div className="px-4 py-3 border-b border-gray-800">
        <button onClick={onInit}
          className="w-full flex items-center justify-center gap-2 bg-accent hover:bg-accent/80 text-black font-medium text-xs py-2 rounded-md transition-colors">
          <Zap size={14} /> Init Team
        </button>
      </div>
      <nav className="flex-1 py-2 overflow-y-auto">
        {sections.map(section => (
          <div key={section}>
            <button
              onClick={() => toggle(section)}
              className="w-full flex items-center justify-between px-5 py-1.5 mt-1 text-[10px] font-semibold uppercase tracking-wider text-gray-500 hover:text-gray-400 transition-colors"
            >
              {section}
              <ChevronDown size={12} className={`transition-transform ${collapsed[section] ? '-rotate-90' : ''}`} />
            </button>
            {!collapsed[section] && nav.filter(n => n.section === section).map(({ id, label, icon: Icon }) => (
              <button
                key={id}
                onClick={() => onNavigate(id)}
                className={`w-full flex items-center gap-3 px-5 py-2 text-sm transition-colors ${
                  active === id
                    ? 'text-accent bg-accent/10 border-r-2 border-accent'
                    : 'text-gray-400 hover:text-gray-200 hover:bg-card'
                }`}
              >
                <Icon size={15} />
                {label}
              </button>
            ))}
          </div>
        ))}
      </nav>
      <div className="px-5 py-4 border-t border-gray-800 flex items-center justify-between">
        <span className="text-[11px] text-gray-600">v0.1.0</span>
        <ThemeToggle />
      </div>
    </aside>
  )
}
