"""Skill loader: reads SKILL.md files and injects into agent prompts."""

import os
from dataclasses import dataclass, field
from pathlib import Path
from typing import Dict, List, Optional


@dataclass
class Skill:
    name: str
    content: str
    path: str


class SkillLoader:
    def __init__(self, skill_dirs: Optional[List[str]] = None):
        self._skill_dirs: List[str] = skill_dirs or []
        self._skills: Dict[str, Skill] = {}
        self._assignments: Dict[str, List[str]] = {}  # agent_name -> [skill_names]
        self.reload()

    def reload(self) -> None:
        """Reload all skills from configured directories."""
        self._skills.clear()
        for d in self._skill_dirs:
            p = Path(d)
            if not p.is_dir():
                continue
            for skill_dir in p.iterdir():
                if not skill_dir.is_dir():
                    continue
                skill_file = skill_dir / "SKILL.md"
                if skill_file.exists():
                    name = skill_dir.name
                    self._skills[name] = Skill(
                        name=name,
                        content=skill_file.read_text(encoding="utf-8"),
                        path=str(skill_file),
                    )

    def list_skills(self) -> List[Skill]:
        """Return all loaded skills."""
        return list(self._skills.values())

    def get(self, name: str) -> Optional[Skill]:
        """Get a skill by name."""
        return self._skills.get(name)

    def assign(self, agent_name: str, skill_names: List[str]) -> None:
        """Assign skills to an agent by name."""
        self._assignments[agent_name] = skill_names

    def get_agent_skills(self, agent_name: str) -> List[Skill]:
        """Get skills assigned to an agent."""
        names = self._assignments.get(agent_name, [])
        return [self._skills[n] for n in names if n in self._skills]

    def inject_into_prompt(self, agent_name: str, base_prompt: str) -> str:
        """Inject assigned skills into an agent's base prompt."""
        skills = self.get_agent_skills(agent_name)
        if not skills:
            return base_prompt
        sections = [f"## Skill: {s.name}\n\n{s.content}" for s in skills]
        skill_block = "\n\n---\n\n".join(sections)
        return f"{base_prompt}\n\n# Loaded Skills\n\n{skill_block}"
