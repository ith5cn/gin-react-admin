---
name: goal-codex-review-workflow
description: "/goal 写完代码先停不提交,Stop hook 跑 codex review;任一轮 CR 给 ALLOW 就在那一轮立刻提交,BLOCK 就改完再 CR 循环。铁律:先过 CR 再 commit"
metadata: 
  node_type: memory
  type: feedback
  originSessionId: 4f0152b8-bebe-4f59-a60d-525d41cbb03b
---

固定工作流(用户 2026-06-20 明确要求固定下来):

1. **写代码阶段绝不提前提交** —— 改动写完先停,**保持留在工作区(未提交)**,让 Stop hook 触发 codex review(`~/.claude/hooks/codex-review-on-stop.js`,审 working-tree 的未提交改动)。
2. **每次某一轮停下时 CR 给 ALLOW,就在那一轮立刻提交这一版**(不攒到最后):
   - 这一轮有改动、CR **通过 → 提交**。
   - 这一轮 CR **BLOCK → 改 → 再停 → 再 CR**,循环直到 **过 → 提交**。
   - 这一轮停下时 CR **直接通过(没让我改)→ 直接提交**。
3. 提交完继续后续工作。

铁律:**先过 CR 再 commit,绝不在 CR 之前 commit。**

**Why:** review gate 只审「未提交的 working-tree 改动」。一旦提前 commit,工作区变干净,`hasChanges()` 返回 false,hook 静默放行 → 等于跳过审查(这正是 yd-admin-login-cc2 那次「hook 没参与进来」的根因)。

**How to apply:** 任何 /goal 任务,写完代码后直接停,不要提交;等审查过了再提交。若已误提交,用 `git reset --soft HEAD~1` 把改动退回工作区再让 hook 审。
