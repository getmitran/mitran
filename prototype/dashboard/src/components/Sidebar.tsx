import { ListOrdered, Bot, ShieldCheck, Ticket, BookOpen, Zap } from 'lucide-react'

const nav = [
  { id: 'queue', label: 'Queue', icon: ListOrdered },
  { id: 'agents', label: 'Agents', icon: Bot },
  { id: 'checkpoints', label: 'Checkpoints', icon: ShieldCheck },
  { id: 'tickets', label: 'Tickets', icon: Ticket },
  { id: 'wiki', label: 'Wiki', icon: BookOpen },
] as const

type View = (typeof nav)[number]['id']

interface Props {
  active: View
  onNavigate: (v: View) => void
  onInit: () => void
}

export default function Sidebar({ active, onNavigate, onInit }: Props) {
  return (
    <aside className="w-56 bg-surface border-r border-gray-800 flex flex-col">
      <div className="px-5 py-5 border-b border-gray-800">
        <h1 className="text-xl font-bold text-accent tracking-tight">mitran</h1>
        <p className="text-[11px] text-gray-500 mt-0.5">agent orchestration</p>
      </div>
      <div className="px-4 py-3 border-b border-gray-800">
        <button onClick={onInit}
          className="w-full flex items-center justify-center gap-2 bg-accent hover:bg-accent/80 text-black font-medium text-xs py-2 rounded-md transition-colors">
          <Zap size={14} /> Init Company
        </button>
      </div>
      <nav className="flex-1 py-3">
        {nav.map(({ id, label, icon: Icon }) => (
          <button
            key={id}
            onClick={() => onNavigate(id)}
            className={`w-full flex items-center gap-3 px-5 py-2.5 text-sm transition-colors ${
              active === id
                ? 'text-accent bg-accent/10 border-r-2 border-accent'
                : 'text-gray-400 hover:text-gray-200 hover:bg-card'
            }`}
          >
            <Icon size={16} />
            {label}
          </button>
        ))}
      </nav>
      <div className="px-5 py-4 border-t border-gray-800 text-[11px] text-gray-600">
        v0.1.0 • prototype
      </div>
    </aside>
  )
}
