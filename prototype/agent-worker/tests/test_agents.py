import pytest
from worker.agents import AGENTS

def test_all_agents_registered():
    expected = {'dev', 'docs', 'ops', 'review', 'hr', 'cicd', 'tickets', 'wiki'}
    assert expected.issubset(set(AGENTS.keys()))

def test_agents_have_system_prompt():
    for name, agent in AGENTS.items():
        assert hasattr(agent, 'system_prompt'), f'{name} missing system_prompt'
        assert len(agent.system_prompt) > 10, f'{name} system_prompt too short'

def test_agents_have_execute():
    for name, agent in AGENTS.items():
        assert hasattr(agent, 'execute'), f'{name} missing execute method'
        assert callable(agent.execute), f'{name}.execute not callable'
