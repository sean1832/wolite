# Wolite Project Instructions

## Engineering Philosophy

_"Clay becomes pottery through craft, but it's the emptiness that makes a pot useful."_ — Laozi

Our engineering culture is built on **Locality**, **Simplicity**, and **Explicitness**. We prioritize code that is easy to reason about, easy to delete, and easy to maintain over code that is "clever" or "dry" at the cost of cognitive load.

### 1. Locality First

**Cognitive Proximity**
Code logic must reside as close as possible to the data it manipulates. Minimize the physical distance between variable declaration and usage.

- **The "Jump" Test:** If understanding a function requires opening more than two other files, the code is too fragmented.
- **Minimizing Context Switching:** Engineers should be able to interact with documentation and code using the same tools.

**Data Locality**
Prioritize contiguous memory layouts (arrays/structs) over pointer-chasing structures (linked lists/trees) to maximize CPU cache hits.

### 2. Radical Simplicity

**The "Rule of Three"**
Strictly prohibit abstraction until a pattern emerges three times. Duplication is preferable to the wrong abstraction.

**Usage-Based Extraction**
Functions or classes exist solely to reduce complexity, not for semantic categorization. If a block of logic is used once, it remains inline.

**Shallow Hierarchies**
Enforce flat dependency structures. Deep inheritance trees and multi-layer wrappers ("lasagna code") are anti-patterns. Composition is strictly preferred over inheritance.

### 3. Explicit Execution

**Linearity**
Prefer explicit control flow over implicit behaviors (e.g., magic methods, AOP, hidden middleware). Code should be readable linearly from top to bottom.

- **Early Returns:** Minimize nesting depth. Use guard clauses to exit functions early.
- **Switch/Match Completeness:** Exhaustive pattern matching is mandatory.

**Type Strictness**
Use strong, static typing to enforce constraints at compile time. Avoid `any` or `void` pointers unless interfacing with low-level boundaries.

### 4. Documentation & Maintenance

**Minimum Viable Documentation**
Docs thrive when they're treated like tests: a necessary chore. Brief and utilitarian is better than long and exhaustive.

- **Readable Source Text:** Content and presentation should not mingle. Plain text suffices.
- **Freshness:** Static content is better than dynamic, but fresh is better than stale.

**Better is Better than Best**
Incremental improvement is better than prolonged debate. Patience and tolerance of imperfection allow projects to evolve organically.

**Deprecation Policy**
Dead code is deleted immediately, not commented out. Version control is the archive; the codebase is the current state.

## Project Purpose

**Wolite** is a Wake-on-LAN (WoL) service that enables remote machine power control over the network. The application provides a simple, intuitive interface for users to remotely turn on their computers.

**Security & Simplicity**

- **JSON Storage**: For simplicity and robustness, the application uses a JSON-only database. This is sufficient for the intended scale (handful of machines) and minimizes operational complexity.
- **Backend Responsibility**: The frontend **must not** contain any sensitive backend logic (e.g., password hashing, OTP generation, session management). All security-critical operations must be handled by the Go backend.
- **Security Precautions**: Given the sensitive nature of remote power control, security is a critical design requirement. Prevent unauthorized actors from controlling machines and audit all wake requests.

## Architecture Overview

Wolite is a monorepo consisting of three main components:

1.  **Backend (`backend/`)**: A Go application that provides the API, manages the JSON database, and handles authentication/security.
2.  **Frontend (`frontend/`)**: A SvelteKit application for the UI. It is built independently and then its output (`build/`) is copied to `backend/internal/ui/dist` to be embedded into the Go binary.
3.  **Client (`client/`)**: (Planned) A lightweight background application for target machines to provide additional information (status) and commands (sleep, shutdown). Designed for minimal system impact.

### Embedding Flow

The frontend is built into the `frontend/build` folder, which is then mirrored to `backend/internal/ui/dist`. The Go backend embeds these assets using `go:embed`.

## Tech Stack

- **Backend**: Go (API, JSON Database)
- **Frontend**: Svelte 5 + SvelteKit with TypeScript
- **Styling**: Tailwind CSS v4 + tailwind-variants (`tv()`)
- **UI Library**: shadcn-svelte (bits-ui primitives)
- **Icons**: @lucide/svelte
- **Build System**: [Taskfile](https://taskfile.dev) (`taskfile.yml`)
- **Client**: (Planned) C/C++ (Preferred for minimal footprint) or Go (Backup)

## Development Commands

All development tasks are managed via `task`:

```bash
task build        # Build both frontend and backend
task frontend     # Build frontend and copy to backend dist
task backend      # Build Go backend binary
task clean        # Remove build artifacts
```

## Component Architecture

Frontend components follow **Atomic Design** in `frontend/src/lib/components/`:

- `atoms/` - Basic building blocks (buttons, inputs)
- `molecules/` - Combinations of atoms
- `organisms/` - Complex UI sections
- `ui/` - shadcn-svelte generated components (DO NOT manually edit)

### Component Patterns

```svelte
<!-- Use Svelte 5 runes -->
<script lang="ts">
  import { cn } from "$lib/utils.js";

  let { class: className, variant = "default", ...restProps } = $props();
</script>
```

- Use `$props()` for props, `$state()` for reactive state, `$bindable()` for two-way binding.
- Use `cn()` from `$lib/utils.js` for class merging.
- Use `tailwind-variants` (`tv()`) for component variants.

### Path Aliases (Frontend)

```typescript
$lib = frontend/src/lib
$lib/components/ui = UI components
$lib/utils = Utility functions
```

### Icons

- Use `@lucide/svelte` for icons
  e.g. `import { UserIcon } from "@lucide/svelte";`

## Key Conventions

1.  **Security First**: Never leak backend logic into the frontend.
2.  **Explicit Control**: Prefer linear, readable code over magic.
3.  **Locality**: Keep logic close to the data it uses.
4.  **Minimalism**: Maintain a small footprint, especially for the planned client.
5.  **New UI Components**: Run `npx shadcn-svelte@latest add <component>` from `frontend/`.

---

## Svelte MCP Tools

When asked about Svelte/SvelteKit topics, use these tools for documentation:

### 1. list-sections

Use FIRST to discover available documentation sections.

### 2. get-documentation

Fetch full documentation for relevant sections.

### 3. svelte-autofixer

MUST use when writing Svelte code.

### 4. playground-link

Only after user confirms they want a link.
