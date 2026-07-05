# Memory Index

- [Bedrock Fable env var](bedrock-fable-env-var.md) — Claude Code 支持 ANTHROPIC_DEFAULT_FABLE_MODEL 来 pin Fable 档位(Bedrock/Vertex/Foundry)
- [ith5:ai subagent 团队](ith5-ai-subagent-team.md) — /yd:ai 并行执行靠 ~/.claude/agents/ 下 4 个角色化 subagent(frontend/backend/database/qa),各薄封装同名 ith5-* skill
- [Goal+codex review 工作流](goal-codex-review-workflow.md) — /goal 写完代码先停(改动留工作区不提交),Stop hook 跑 codex review;BLOCK 就改完再停循环,过了才提交。提前 commit 会让 hook 静默放行=跳过审查
