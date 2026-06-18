import { useEffect, useRef, useState } from 'react';

interface Metrics {
  tasks_per_hour: number[];
  agent_invocations: Record<string, number>;
  token_usage: { prompt: number; completion: number; total: number };
  system_health: { cpu: number; memory: number; uptime: number; status: string };
}

const DEFAULT_METRICS: Metrics = {
  tasks_per_hour: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  agent_invocations: { Dev: 0, Docs: 0, Ops: 0, Review: 0, HR: 0, 'CI/CD': 0, Tickets: 0, Wiki: 0 },
  token_usage: { prompt: 0, completion: 0, total: 0 },
  system_health: { cpu: 0, memory: 0, uptime: 0, status: 'unknown' },
};

declare const Chart: any;

export default function MonitoringPage() {
  const [metrics, setMetrics] = useState<Metrics>(DEFAULT_METRICS);
  const [error, setError] = useState('');
  const tasksChartRef = useRef<HTMLCanvasElement>(null);
  const agentChartRef = useRef<HTMLCanvasElement>(null);
  const tokenChartRef = useRef<HTMLCanvasElement>(null);
  const chartsRef = useRef<any[]>([]);

  useEffect(() => {
    const script = document.createElement('script');
    script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
    script.onload = () => fetchAndRender();
    document.head.appendChild(script);
    const interval = setInterval(fetchAndRender, 30000);
    return () => {
      clearInterval(interval);
      chartsRef.current.forEach(c => c?.destroy());
    };
  }, []);

  async function fetchAndRender() {
    try {
      const res = await fetch('http://localhost:7780/api/v1/metrics');
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      setMetrics(data);
      setError('');
      renderCharts(data);
    } catch (e: any) {
      setError(e.message);
      renderCharts(DEFAULT_METRICS);
    }
  }

  function renderCharts(data: Metrics) {
    if (typeof Chart === 'undefined') return;
    chartsRef.current.forEach(c => c?.destroy());
    chartsRef.current = [];

    if (tasksChartRef.current) {
      chartsRef.current.push(new Chart(tasksChartRef.current, {
        type: 'line',
        data: {
          labels: data.tasks_per_hour.map((_, i) => `${i}h`),
          datasets: [{ label: 'Tasks/Hour', data: data.tasks_per_hour, borderColor: '#3b82f6', tension: 0.3, fill: true, backgroundColor: 'rgba(59,130,246,0.1)' }],
        },
        options: { responsive: true, plugins: { legend: { display: false } } },
      }));
    }

    if (agentChartRef.current) {
      const labels = Object.keys(data.agent_invocations);
      chartsRef.current.push(new Chart(agentChartRef.current, {
        type: 'bar',
        data: {
          labels,
          datasets: [{ label: 'Invocations', data: labels.map(k => data.agent_invocations[k]), backgroundColor: '#8b5cf6' }],
        },
        options: { responsive: true, plugins: { legend: { display: false } } },
      }));
    }

    if (tokenChartRef.current) {
      chartsRef.current.push(new Chart(tokenChartRef.current, {
        type: 'doughnut',
        data: {
          labels: ['Prompt', 'Completion'],
          datasets: [{ data: [data.token_usage.prompt, data.token_usage.completion], backgroundColor: ['#f59e0b', '#10b981'] }],
        },
        options: { responsive: true },
      }));
    }
  }

  const { system_health: health } = metrics;

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Monitoring</h1>
      {error && <div className="text-sm text-red-500">⚠ Metrics unavailable: {error}</div>}

      <div className="grid grid-cols-4 gap-4">
        <StatCard label="Status" value={health.status} color={health.status === 'healthy' ? 'text-green-500' : 'text-red-500'} />
        <StatCard label="CPU" value={`${health.cpu}%`} />
        <StatCard label="Memory" value={`${health.memory}%`} />
        <StatCard label="Uptime" value={`${Math.floor(health.uptime / 3600)}h`} />
      </div>

      <div className="grid grid-cols-2 gap-6">
        <ChartCard title="Tasks / Hour"><canvas ref={tasksChartRef} /></ChartCard>
        <ChartCard title="Agent Invocations"><canvas ref={agentChartRef} /></ChartCard>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <ChartCard title="LLM Token Usage">
          <canvas ref={tokenChartRef} />
          <div className="mt-2 text-xs text-gray-500 text-center">Total: {metrics.token_usage.total.toLocaleString()} tokens</div>
        </ChartCard>
        <ChartCard title="Token Breakdown">
          <div className="space-y-3 pt-4">
            <TokenBar label="Prompt" value={metrics.token_usage.prompt} total={metrics.token_usage.total} color="bg-amber-500" />
            <TokenBar label="Completion" value={metrics.token_usage.completion} total={metrics.token_usage.total} color="bg-emerald-500" />
          </div>
        </ChartCard>
      </div>
    </div>
  );
}

function StatCard({ label, value, color }: { label: string; value: string; color?: string }) {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
      <div className="text-xs text-gray-500">{label}</div>
      <div className={`text-xl font-semibold ${color || ''}`}>{value}</div>
    </div>
  );
}

function ChartCard({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
      <h3 className="text-sm font-medium mb-3">{title}</h3>
      {children}
    </div>
  );
}

function TokenBar({ label, value, total, color }: { label: string; value: number; total: number; color: string }) {
  const pct = total > 0 ? (value / total) * 100 : 0;
  return (
    <div>
      <div className="flex justify-between text-xs mb-1">
        <span>{label}</span>
        <span>{value.toLocaleString()}</span>
      </div>
      <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full">
        <div className={`h-2 ${color} rounded-full`} style={{ width: `${pct}%` }} />
      </div>
    </div>
  );
}
