Let me create both an architecture diagram and Go pseudocode for this.Now let me build the architecture diagram and the Go pseudocode.The flow works like this: every incoming question gets embedded and checked against a semantic cache of previously well-answered questions. Cache hits return instantly; cache misses go to Claude with RAG context injected from past good answers and curated location knowledge.

Now here's the Go pseudocode showing how this plugs into the existing MCP server:Here's the key idea in the architecture:

**The system never modifies the LLM itself** — instead it builds a smarter layer around it. When someone asks "why is there high radiation near Eubank and Southern in Albuquerque?", the flow is:

1. The question gets embedded into a vector
2. DuckDB checks for semantically similar past questions (using `list_cosine_similarity`, which DuckDB supports natively)
3. If there's a high-confidence match with positive user feedback → return the cached answer instantly, no LLM call needed
4. If no match, the top-k most similar past Q&A pairs get injected into Claude's system prompt as RAG context, along with any curated location knowledge for the area
5. After the response, the user can thumbs-up/down it
6. Positively-rated answers enter the cache and also auto-populate the location knowledge store if coordinates are detected

The nice thing for Safecast specifically is the **location knowledge store** — once someone confirms that the Nuclear Museum explains the readings at that spot, every future question about that area gets that context automatically. Over time you'd build up a geographic knowledge base of "why this area reads high/low" — naturally radioactive geology, nuclear facilities, museums, medical centers, etc.

The Go code plugs directly into the existing `web-chat` handler. The main additions are the DuckDB tables, an embedding service call, and a `/api/feedback` endpoint. DuckDB is already in your stack, so no new dependencies there. For embeddings you could start with OpenAI's `text-embedding-3-small` (very cheap) and later swap in a local model to keep everything self-hosted.
