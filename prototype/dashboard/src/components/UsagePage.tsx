import { useEffect, useState } from 'react';

interface ModelUsage {
  model: string;
  total_input_tokens: number;
  total_output_tokens: number;
  total_cost: number;
  request_count: number;
}

interface UsageSummary {
  models: ModelUsage[];
  total_cost: number;
}

export default function UsagePage() {
  const [data, setData] = useState<UsageSummary | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('http://localhost:7780/api/v1/usage/summary')
      .then(r => { if (!r.ok) throw new Error(r.statusText); return r.json(); })
      .then(setData)
      .catch(e => setError(e.message));
  }, []);

  if (error) return <div className="p-6 text-red-500">Error: {error}</div>;
  if (!data) return <div className="p-6">Loading usage data...</div>;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-2">Usage & Billing</h1>
      <div className="mb-4 text-lg font-semibold">
        Total Cost: <span className="text-green-400">${data.total_cost.toFixed(4)}</span>
      </div>
      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b border-gray-700">
            <th className="text-left py-2 px-3">Model</th>
            <th className="text-right py-2 px-3">Input Tokens</th>
            <th className="text-right py-2 px-3">Output Tokens</th>
            <th className="text-right py-2 px-3">Requests</th>
            <th className="text-right py-2 px-3">Cost</th>
          </tr>
        </thead>
        <tbody>
          {data.models.map(m => (
            <tr key={m.model} className="border-b border-gray-800">
              <td className="py-2 px-3 font-mono">{m.model}</td>
              <td className="text-right py-2 px-3">{m.total_input_tokens.toLocaleString()}</td>
              <td className="text-right py-2 px-3">{m.total_output_tokens.toLocaleString()}</td>
              <td className="text-right py-2 px-3">{m.request_count}</td>
              <td className="text-right py-2 px-3">${m.total_cost.toFixed(4)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
