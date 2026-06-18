import { useEffect, useState } from 'react';

interface GitHubStatus {
  connected: boolean;
  repos?: string[];
}

interface SlackStatus {
  connected: boolean;
  workspace?: string;
}

function Badge({ connected }: { connected: boolean }) {
  return (
    <span className={`px-2 py-0.5 rounded text-xs font-medium ${connected ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
      {connected ? 'Connected' : 'Disconnected'}
    </span>
  );
}

export default function IntegrationsPage() {
  const [github, setGithub] = useState<GitHubStatus | null>(null);
  const [slack, setSlack] = useState<SlackStatus | null>(null);

  useEffect(() => {
    fetch('http://localhost:7780/api/v1/integrations/github/status')
      .then(r => r.json())
      .then(setGithub)
      .catch(() => setGithub({ connected: false }));

    fetch('http://localhost:7780/api/v1/integrations/slack/status')
      .then(r => r.json())
      .then(setSlack)
      .catch(() => setSlack({ connected: false }));
  }, []);

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Integrations</h1>

      <section className="border rounded-lg p-4">
        <div className="flex items-center gap-3 mb-3">
          <h2 className="text-lg font-semibold">GitHub</h2>
          {github && <Badge connected={github.connected} />}
        </div>
        {github?.connected && github.repos && (
          <ul className="list-disc list-inside text-sm space-y-1">
            {github.repos.map(repo => <li key={repo}>{repo}</li>)}
          </ul>
        )}
        {!github && <p className="text-sm text-gray-500">Loading...</p>}
      </section>

      <section className="border rounded-lg p-4">
        <div className="flex items-center gap-3 mb-3">
          <h2 className="text-lg font-semibold">Slack</h2>
          {slack && <Badge connected={slack.connected} />}
        </div>
        {slack?.connected && slack.workspace && (
          <p className="text-sm">Workspace: {slack.workspace}</p>
        )}
        {!slack && <p className="text-sm text-gray-500">Loading...</p>}
      </section>
    </div>
  );
}
