import json, os
from datetime import date
from collections import defaultdict

DATA_FILE = os.path.join(os.path.dirname(__file__), '..', 'data', 'cost_tracking.json')
MODEL_PRICES = {
    "claude-sonnet": (3.0, 15.0),
    "claude-haiku": (0.25, 1.25),
    "gpt-4o": (5.0, 15.0),
    "gpt-4o-mini": (0.15, 0.6),
    "llama3": (0.0, 0.0),
}

class CostTracker:
    def __init__(self):
        os.makedirs(os.path.dirname(DATA_FILE), exist_ok=True)
        self.data = self._load()

    def _load(self):
        if os.path.exists(DATA_FILE):
            with open(DATA_FILE) as f:
                return json.load(f)
        return {"records": []}

    def _save(self):
        with open(DATA_FILE, 'w') as f:
            json.dump(self.data, f, indent=2)

    def track(self, agent_id, model, input_tokens, output_tokens, cost_usd=None):
        if cost_usd is None:
            inp_rate, out_rate = MODEL_PRICES.get(model, MODEL_PRICES["claude-sonnet"])
            cost_usd = input_tokens * (inp_rate / 1_000_000) + output_tokens * (out_rate / 1_000_000)
        self.data["records"].append({"agent_id": agent_id, "model": model,
            "input_tokens": input_tokens, "output_tokens": output_tokens,
            "cost_usd": round(cost_usd, 6), "date": str(date.today())})
        self._save()

    def get_summary(self, agent_id=None):
        recs = [r for r in self.data["records"] if not agent_id or r["agent_id"] == agent_id]
        return {"total_input": sum(r["input_tokens"] for r in recs),
                "total_output": sum(r["output_tokens"] for r in recs),
                "total_cost_usd": round(sum(r["cost_usd"] for r in recs), 6),
                "requests": len(recs)}

    def get_daily(self, agent_id):
        daily = defaultdict(lambda: {"input": 0, "output": 0, "cost": 0.0})
        for r in self.data["records"]:
            if r["agent_id"] == agent_id:
                daily[r["date"]]["input"] += r["input_tokens"]
                daily[r["date"]]["output"] += r["output_tokens"]
                daily[r["date"]]["cost"] += r["cost_usd"]
        return dict(daily)

    def reset(self):
        self.data = {"records": []}
        self._save()
