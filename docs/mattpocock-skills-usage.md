# mattpocock-skills: 按处境上手

> 目标：不了解这套技能也能选对入口，复制第一句话开干。
> 迷茫时打：`/ask-matt` + 你现在卡在哪。

## 30 秒先记住

1. **先选处境，再打技能** (多数技能不会自己开，要你先 `/技能名`)。
2. **Matt 管流程** (想清楚 -> 规格 -> 工单 -> 实现)。

---

---

## 先看来源 (入口分流)

多数「新需求」先问来源，比背技能名更快：

| 来源 | 入口 | 什么时候能开写 |
|------|------|----------------|
| 别人丢来的原文 | `/triage` | 标到 `ready-for-agent` (或给人 `ready-for-human`) |
| 开发者自己发起 | `/grill-with-docs` | 小：共识后同窗写；大：`to-spec` + 票后再 `/implement` |

(下图用纯文本，避免 Markdown 预览不渲染 mermaid。)

```text
先看来源
|
+-- [外来原文：策划/QA/运营/issue/外部PR]
|   |
|   v
|   /triage  ----落地唯一状态 (五选一) ----
|       |       needs-triage ----继续----> /triage
|       |       needs-info --对方回复--> /triage
|       |       ready-for-agent -------> /implement  << 可开写
|       |       ready-for-human -------> 人做         << 可开写
|       |       wontfix --------------> 结束
|       |
|       +-- (旁路,不是第6态) 要拍板/大雾
|               -> grill / wayfinder -> 拍板后再实现
|
+-- [开发者自己发起：已知要做,要拍口径]
    |
    v
    /grill-with-docs
    |
    +-- 同窗口能做完? --是--> 确认共识后同窗口实现 << 可开写
    |
    +-- 否 --> /to-spec -> /to-tickets -> /implement << 可开写
    |
    +--(旁路) 边界都画不出 -> /wayfinder 先 Chart
               -> 雾散后再 to-spec / implement
```

---

口诀：外来走蓝 (triage 五态)；自己走紫 (grill)。两边最后都可能到 implement。
**别把外来原文直接 grill，也别把自己的设计稿丢进 triage。**
坏了等其它处境，仍用下一节对照表（不只这两种来源）。

---

## 易混淆

三个名字都带 grill，本仓只记这一条：

```text
grilling        = 引擎 (怎么问)
grill-me        = 引擎，不写文件 (无仓库/纯想法)
grill-with-docs = 引擎 + 对着本仓写术语/ADR  <- 本仓用这个
```

| 技能 | 你打不打 | 本仓 |
|------|---------|------|
| `/grilling` | 一般不打 | 内核，会被其它技能拉起 |
| `/grill-me` | 不要打 | 无落盘；有工作目录时不用 |
| `/grill-with-docs` | **打这个** | 开发者自己发起需求时的入口 |

一次会话装不下 / 边界画不出 -> `/wayfinder` (里面仍用 grilling，不用再换 grill-me)。

## 我现在是哪种情况？

从上到下对号，**第一个匹配就停**。右栏整段可直接粘到 Cursor。

| # | 你现在的情况 | 打什么 | 第一句话（照抄改括号） |
|---|---|---|---|
| 1 | 线上/测试坏了、偶发、回归、原因说不清 | `/diagnosing-bugs` | `/diagnosing-bugs` 现象：(一句话)。复现：(步骤)。日志/栈：(贴关键几行) |
| 2 | 策划/QA/运营丢来一段原文，还没整理 | `/triage` | `/triage` 下面是对方原文：(粘贴)。请标状态并说明下一步 |
| 3 | 系统级大改，连边界/权威数据都画不出 | `/wayfinder` | `/wayfinder` 目的地：(一句话)。先 Chart map，这轮不要解题 |
| 4 | 已有 `.scratch/.../map.md`，继续解雾 | `/wayfinder` | `/wayfinder` map: (map.md 路径)。解下一张 frontier 票 |
| 5 | 代码已改完，只要审 | `/code-review` | `/code-review` 固定点：(commit/branch/tag)。规格：(spec 路径或刚才共识) |
| 6 | 新功能大，一次对话装不下 | 完整链路 | `/grill-with-docs` 需求：(一句话)。有 repo，请落 CONTEXT/ADR |
| 7 | 有几处产品决策，但实现同窗口能做完 | 小功能 | `/grill-with-docs` 需求：()。问清后同窗口实现，不要拆票 |
| 8 | 一句话说清「改哪、改成啥、别动什么」 | 直接改 | (普通对话) 改：(文件/函数)。改成：()。不要动：() |
| 9 | grill/wayfinder 卡住，答案在别人脑子里 | `/to-questionnaire` | `/to-questionnaire` 收件人：(角色)。要收回：(缺口列表) |
| 10 | Agent 上一句看不懂 | `/wait-what` | `/wait-what` |

策划 docx / 外来案：**不能**直接 `/implement`，从 `/triage` 或 `/wayfinder` 进。

## 对上号之后怎么走

### 1) 坏了 -> diagnosing-bugs

先锁一条能红的反馈 (复现/断言/日志)，再修。
**不要**先 grill 新产品设计。修完需要审再用 `/code-review`。

### 2) 外来原文 -> triage

见上节入口图。triage **状态只有五选一** (不是后面步骤)：`needs-triage` / `needs-info` / `ready-for-agent` / `ready-for-human` / `wontfix`。

| 标成 | 你下一步 |
|------|---------|
| `needs-triage` | 继续 triage，未完成 |
| `needs-info` | 问对方，或同对话 `/to-questionnaire`；对方回复后再 triage |
| `ready-for-agent` | **可开启后面**：`/implement` + 说明范围 |
| `ready-for-human` | **可开启后面**：人做 |
| `wontfix` | 结束 |
| 要拍板 / 大雾 (旁路，不是第 6 态) | `/grill-with-docs` 或 `/wayfinder`，拍板后再实现 |

triage **不写代码、不写 spec**。

### 3) 大雾 -> wayfinder

- **第一轮**: 只要 Chart (map + 决策票)，不要解题。
- **之后每轮**: 带 `map.md` 路径，一次解一张 (research 可并行)。
- research 票关掉 ≠ 需求明确；grilling 才拍口径。
- 雾散后再 `/to-spec`，**禁止**从 map 直接 implement。

落盘：`.scratch/<effort>/map.md` + `issues/` (research 长文在 `research/`)。

### 4) 完整链路

**同一对话**做完前三段，再按票新开对话实现：

```text
/grill-with-docs -> /to-spec -> /to-tickets -> (/clear 后) /implement -> /code-review
```

| 步骤 | 你做什么 | 看到什么算过 |
|------|---------|-------------|
| grill | 只答决策；可用「按推荐」 | 共识复述你点头；CONTEXT/ADR 该写已写 |
| to-spec | 确认测试 seam (怎么验收) | `.scratch/<feature>/spec.md` 且 `ready-for-agent` |
| to-tickets | 批准垂直切片与依赖 | `issues/01-*.md`... (不是按文件类型拆票) |
| implement | `/implement` + 票路径 | 按 seam 手工/测试验收 |
| code-review | 固定点 (commit/branch/tag) + spec 路径 | Standards / Spec 两轴报告 |

中途纸面谈不清：`/prototype`，结论带回原对话。
缺别人口径：同对话 `/to-questionnaire`，收回再继续。

实现示例：

```text
/implement 按 .scratch/<feature>/issues/01-xxx.md 实现
```

### 5) 小功能

```text
/grill-with-docs -> 你确认共识 -> 同窗口「按刚才共识改」
```

砍掉：`to-spec` / `to-tickets` / 为拆票而拆票。

### 6) 直接改

不打 Matt 技能。说清范围即可。有产品语义变化时改走小功能（先 `/grill-with-docs`）。

---

## 会话怎么切 (少踩坑)

| 阶段 | 做法 |
|------|------|
| grill -> to-spec -> to-tickets | **同一对话**，中途不要 `/clear` |
| 每张 implement 票 | `/clear` 后新开，只带票路径 |
| 窗口快满且未到 tickets | 在阶段边界 `/compact`，别在 grill 半中间压 |

换会话优先：继续原对话 -> `/clear` -> `/handoff` -> subagent -> `/compact`。

---

## 产物在哪

```text
CONTEXT.md                        术语
docs/adr/NNNN-*.md                难回退决策
.scratch/<feature>/spec.md        规格
.scratch/<feature>/issues/NN-*.md  实现票
.scratch/<effort>/map.md          wayfinder 地图
.scratch/<effort>/issues/NN-*.md   决策票
```

约定详见：`docs/agents/issue-tracker.md`、`docs/agents/triage-labels.md`、`AGENTS.md`。

---

## 常见踩坑

| 别这样 | 该这样 |
|--------|--------|
| 空指针/纯日志 bug 去 grill | `/diagnosing-bugs` 或直接改 |
| 策划原文直接 implement | `/triage` |
| 大雾直接 to-spec / implement | `/wayfinder` 先 Chart |
| 小改拆成按文件类型的多票 | 小功能或直接改 |
| 把问卷当 spec | 问卷收回后继续 grill，再 to-spec |

---

## 开干前

1. 先约定验收 seam；无单测则手工验收也算。
2. 文档与实现冲突：行为对则以代码为准，改 spec/ADR 再审。

---

## 最小检查

- [ ] 用上面表格对上了处境 (没误走完整链路)
- [ ] 第一句话已带齐：需求/原文/复现/票路径/固定点
- [ ] 完整链路才拆票，且是垂直可验收票
- [ ] wayfinder: Chart 与解题分开；雾散再 to-spec
