import { useState } from 'react';

const STEPS = ['Organization', 'Agents', 'Integrations'];

export default function OnboardingPage() {
  const [step, setStep] = useState(0);
  const [config, setConfig] = useState({ org_name: '', agents: [] as string[], integrations: [] as string[] });

  const next = () => {
    fetch('/api/v1/onboarding/step', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ tenant_id: config.org_name, step_name: STEPS[step], config }),
    });
    setStep(s => s + 1);
  };

  return (
    <div className="p-6 max-w-lg mx-auto">
      <h1 className="text-2xl font-bold mb-4">Setup Wizard</h1>
      <div className="flex gap-1 mb-6">
        {STEPS.map((s, i) => (
          <div key={s} className={`flex-1 h-2 rounded ${i <= step ? 'bg-blue-500' : 'bg-gray-300'}`} />
        ))}
      </div>

      {step === 0 && (
        <div>
          <label className="block text-sm mb-1">Organization Name</label>
          <input value={config.org_name} onChange={e => setConfig({ ...config, org_name: e.target.value })}
            className="w-full px-3 py-2 border rounded" placeholder="Acme Corp" />
        </div>
      )}
      {step === 1 && (
        <div>
          <p className="text-sm mb-2">Select agents to enable:</p>
          {['DevOps', 'HR', 'Finance', 'Support', 'Analytics'].map(a => (
            <label key={a} className="flex items-center gap-2 mb-1">
              <input type="checkbox" checked={config.agents.includes(a)}
                onChange={e => setConfig({ ...config, agents: e.target.checked ? [...config.agents, a] : config.agents.filter(x => x !== a) })} />
              {a}
            </label>
          ))}
        </div>
      )}
      {step === 2 && (
        <div>
          <p className="text-sm mb-2">Connect integrations:</p>
          {['GitHub', 'Slack', 'Jira'].map(i => (
            <label key={i} className="flex items-center gap-2 mb-1">
              <input type="checkbox" checked={config.integrations.includes(i)}
                onChange={e => setConfig({ ...config, integrations: e.target.checked ? [...config.integrations, i] : config.integrations.filter(x => x !== i) })} />
              {i}
            </label>
          ))}
        </div>
      )}
      {step >= 3 && <p className="text-green-600 font-semibold">✅ Setup complete!</p>}

      {step < 3 && (
        <button onClick={next} className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
          {step === 2 ? 'Finish' : 'Next'}
        </button>
      )}
    </div>
  );
}
