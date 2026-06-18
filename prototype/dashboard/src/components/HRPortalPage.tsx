

const modules = [
  { icon: '👥', title: 'Directory', desc: 'Employee profiles and contact lookup', path: '/hr/directory' },
  { icon: '🏖️', title: 'Leave', desc: 'Request and manage time off', path: '/hr/leave' },
  { icon: '🚀', title: 'Onboarding', desc: 'New hire workflows and checklists', path: '/hr/onboarding' },
  { icon: '📊', title: 'Performance', desc: 'Reviews, goals, and feedback cycles', path: '/hr/performance' },
  { icon: '🏢', title: 'Org Chart', desc: 'Team structure and reporting lines', path: '/hr/org-chart' },
  { icon: '📢', title: 'Announcements', desc: 'Company-wide news and updates', path: '/hr/announcements' },
  { icon: '⏰', title: 'Attendance', desc: 'Timesheets and clock-in records', path: '/hr/attendance' },
  { icon: '💰', title: 'Payroll', desc: 'Salary, payslips, and deductions', path: '/hr/payroll' },
  { icon: '🎯', title: 'Recruitment', desc: 'Job postings and candidate pipeline', path: '/hr/recruitment' },
  { icon: '👋', title: 'Offboarding', desc: 'Exit workflows and asset returns', path: '/hr/offboarding' },
  { icon: '📚', title: 'Training', desc: 'Courses, certifications, and learning paths', path: '/hr/training' },
  { icon: '📈', title: 'Analytics', desc: 'HR metrics and workforce insights', path: '/hr/analytics' },
];

export default function HRPortalPage() {
  return (
    <div className="p-6 max-w-6xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">HR Portal</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {modules.map((m) => (
          <a
            key={m.title}
            href={m.path}
            className="block p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 hover:shadow-md transition-shadow"
          >
            <div className="text-3xl mb-2">{m.icon}</div>
            <h2 className="font-semibold text-lg">{m.title}</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{m.desc}</p>
          </a>
        ))}
      </div>
    </div>
  );
}
