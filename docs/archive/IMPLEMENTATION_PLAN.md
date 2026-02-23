# Plan: Integrating Claude AI Chat into Safecast Map

**Status:** 🗄️ ARCHIVED - Not implemented
**Date Archived:** 2026-02-22
**Alternative Approach:** Instead of embedding chat in the map, the [Safecast MCP server](https://github.com/Safecast/safecast-map-MCP) was created as a standalone service that can be integrated with Claude.ai's built-in MCP integration system. See [README.md](../../README.md) for current architecture.

---

This plan outlines the steps to integrate the Claude AI chat interface directly into the `safecast-new-map` application. This integration allows users to ask questions about radiation data while viewing the map, leveraging the underlying Safecast MCP tools.

## high-level Architecture

1.  **Frontend (`map.html`):** A new Chat UI (floating button + modal) triggers API requests to the Go backend.
2.  **Backend (`safecast-new-map.go`):** A new `/api/chat` endpoint receives messages, invokes the `claude` CLI, and returns the response.
3.  **Intelligence (`claude` CLI):** The backend executes the `claude` CLI, which is pre-configured to communicate with the Safecast MCP server to fetch data.

## Step 1: Backend Implementation (`safecast-new-map.go`)

We will add a new HTTP handler to the existing Go server.

*   **Task:** Create a `handleChat` function.
*   **Responsibility:**
    *   Accept `POST` requests with a JSON body: `{"message": "..."}`.
    *   Construct a System Prompt that guides Claude on how to use Safecast tools (Real-time vs Historical data keys, Unit conversions).
    *   Execute the `claude` CLI command:
        ```go
        exec.CommandContext(ctx, "claude", "-p", prompts, "--allowedTools", "mcp__claude_ai_Safecast_MCP__*")
        ```
    *   Capture standard output/error and return the text response to the frontend.
*   **Registration:** Register the handler in `main()` via `http.HandleFunc("/api/chat", handleChat)`.

## Step 2: Frontend Implementation (`public_html/map.html`)

We will add the chat interface elements to the main map page.

*   **Task:** Add HTML/CSS/JS for the chat widget.
*   **Components:**
    *   **Button:** A floating "Chat" icon button in the bottom-right corner.
    *   **Modal:** A chat window with message history and an input field.
    *   **Logic (JS):**
        *   Handle opening/closing the modal.
        *   Send `fetch` requests to `/api/chat`.
        *   Render Markdown responses (using `marked.js` library).
        *   Show "Thinking..." loading states.

## Step 3: Server Configuration & Prerequisties

Ensure the environment where `safecast-new-map` runs is ready.

*   **Claude CLI:** The `claude` command line tool must be installed and authenticated on the server.
*   **MCP Configuration:** The `claude` CLI must be configured to connect to the Safecast MCP server (likely via `claude mcp add` or a config file).
*   **Environment Variables:** Ensure necessary API keys (if any, separate from the CLI's internal auth) are available to the process.

## Step 4: Verification

1.  **Build:** Compile the modified `safecast-new-map`.
2.  **Run:** Start the server.
3.  **Test:**
    *   Open the map in a browser.
    *   Click the chat button.
    *   Ask: "What is the radiation level in Tokyo?"
    *   Verify that Claude answers and correctly uses the MCP tools (checking the logs can confirm tool usage).

## Questions / Assumptions

*   **Assume:** The `claude` CLI is already installed and configured on the machine running `safecast-new-map`.
*   **Assume:** The MCP server is running and accessible by the `claude` CLI.
