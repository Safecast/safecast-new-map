# Swap web-chat LLM call to support North Mini Code (local test)

## Context
The Safecast AI assistant (`cmd/web-chat`) currently hardcodes the Anthropic
Messages API (`https://api.anthropic.com/v1/messages`, Anthropic request/response
structs, Anthropic tool-call format). The user wants to test **Cohere North Mini
Code** — a 30B-MoE/3B-active, Apache-2.0, agentic-coding model (256K ctx, trained
with 43% tool-use data) — as a drop-in replacement, running **locally first**.

North Mini Code is served through an **OpenAI-compatible** endpoint (self-hosted
via vLLM at `/v1/chat/completions`, or Cohere's Compatibility API). The MCP tool
layer is already model-agnostic (`model-adapter` has Claude/Kimi/Qwen hint
providers), so the only Anthropic-locked code is the single HTTP call + payload
translation in web-chat. Goal: add an OpenAI-compatible provider path behind an
env flag, leaving the agentic loop untouched.

## Approach
Add a provider switch at the `callAnthropic` boundary in
[cmd/web-chat/main.go](cmd/web-chat/main.go). The agentic loop
([main.go:398-466](cmd/web-chat/main.go#L398)) keeps operating on
`anthropicMessage`/`contentBlock` internally; a new OpenAI adapter translates
in and out so nothing downstream changes.

### 1. Config / wiring (`main()` ~[main.go:481-513](cmd/web-chat/main.go#L481))
Add env vars (defaults preserve current Anthropic behavior):
- `LLM_PROVIDER` = `anthropic` (default) | `openai`
- `LLM_BASE_URL` = e.g. `http://localhost:8000/v1` for local vLLM
- `LLM_API_KEY`  = falls back to `ANTHROPIC_API_KEY` if unset (vLLM accepts any)
- `LLM_MODEL`    = e.g. `CohereLabs/North-Mini-Code-1.0` (falls back to `CLAUDE_MODEL`)

Thread `provider`, `baseURL` into `handleChat` (extend its signature, or pass a
small `llmConfig` struct).

### 2. New `callOpenAI` (sibling of `callAnthropic`, [main.go:243](cmd/web-chat/main.go#L243))
Same signature shape, returns `*anthropicResponse` so the loop is unchanged.
- **Request translation** (`anthropicMessage[] + system + tools → OpenAI body`):
  - `system` → leading `{role:"system"}` message.
  - String-content messages → `{role, content}`.
  - Assistant messages whose content is `[]contentBlock`: emit `tool_use` blocks
    as `tool_calls[]` (`id`, `function.name`, `function.arguments` = JSON string
    of `Input`); any text block → `content`.
  - User messages whose content is `[]contentBlock` of `tool_result`: emit one
    `{role:"tool", tool_call_id, content}` per block.
  - `anthropicTool` → `{type:"function", function:{name, description, parameters:InputSchema}}`.
- **HTTP**: POST `LLM_BASE_URL + /chat/completions`, `Authorization: Bearer <key>`.
- **Response translation** (`choices[0].message → anthropicResponse`):
  - `message.content` → one `{Type:"text"}` block.
  - each `message.tool_calls[]` → `{Type:"tool_use", ID, Name, Input}` (Input =
    raw `function.arguments`).
  - `finish_reason`: `tool_calls` → leave `StopReason` empty (loop continues);
    `stop`/`length` → `"end_turn"`.

Add minimal OpenAI structs (`openAIRequest`, `openAIMessage`, `openAIToolCall`,
`openAIResponse`) in the types section near [main.go:167](cmd/web-chat/main.go#L167).

### 3. Dispatch
At [main.go:399](cmd/web-chat/main.go#L399) call `callAnthropic` or `callOpenAI`
based on `provider`. Keep `truncateHistory`/`systemPrompt`/`maxTokens` shared.

## Out of scope
- No streaming-token support for OpenAI path (current code already buffers
  per-block; one text block per turn is fine for local testing).
- No changes to MCP server, model-adapter hints, or production deploy.
- vLLM server setup (tool-call-parser + Cohere `melody` lib) is the user's local
  serving step, not a code change here.

## Run via OpenCode Zen (free, recommended for testing)
No serving needed — OpenCode Zen exposes `north-mini-code-free` over an
OpenAI-compatible endpoint, which the new `callOpenAI` path already supports.
Get an API key at https://opencode.ai/auth, then:
```
LLM_PROVIDER=openai \
LLM_BASE_URL=https://opencode.ai/zen/v1 \
LLM_MODEL=north-mini-code-free \
LLM_API_KEY=<opencode-zen-key> \
MCP_URL=http://localhost:3333/mcp-http ./web-chat
```
Note: during the free period OpenCode may use submitted data to improve the
model — do not send sensitive data.

## Verification (local self-host alternative)
1. Serve the model locally, OpenAI-compatible, e.g.:
   `vllm serve CohereLabs/North-Mini-Code-1.0 --tool-call-parser <cohere> --enable-auto-tool-choice`
   (per Cohere docs; needs `melody` for correct parsing).
2. Build web-chat: `/usr/local/go/bin/go build -o web-chat ./cmd/web-chat/`.
3. Run:
   ```
   LLM_PROVIDER=openai LLM_BASE_URL=http://localhost:8000/v1 \
   LLM_MODEL=CohereLabs/North-Mini-Code-1.0 LLM_API_KEY=local \
   MCP_URL=http://localhost:3333/... ./web-chat
   ```
4. POST to `/chat` with a question that forces a tool call (e.g. "radiation near
   Fukushima") and confirm: text streams back AND an MCP tool actually fires
   (check web-chat logs for `CallTool`). This validates the tool-call round-trip,
   the main risk with a small model.
5. Regression: unset `LLM_PROVIDER` (or `=anthropic`) → original Claude path still
   works unchanged.

## Key files
- [cmd/web-chat/main.go](cmd/web-chat/main.go) — only file changed.
