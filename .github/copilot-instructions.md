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

**Security Considerations**
Given the sensitive nature of remote power control, security is a critical design requirement:

- Prevent unauthorized actors from controlling machines on the network
- Implement authentication and authorization mechanisms
- Consider network-level security (LAN access, VPN requirements)
- Audit and log all wake requests for accountability

## Architecture Overview

Monorepo with two main directories:

- `server/` - SvelteKit frontend application (main active codebase)
- `client/` - Reserved for future client-side code

## Tech Stack

- **Framework**: Svelte 5 + SvelteKit with TypeScript
- **Styling**: Tailwind CSS v4 + tailwind-variants (`tv()`) for component variants
- **UI Library**: shadcn-svelte (bits-ui primitives)
- **Icons**: @lucide/svelte
- **Additional**: vaul-svelte (drawer), svelte-sonner (toasts), mode-watcher (dark mode)

## Component Architecture

Components follow **Atomic Design** in `server/src/lib/components/`:

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

- Use `$props()` for props, `$state()` for reactive state, `$bindable()` for two-way binding
- Use `cn()` from `$lib/utils.js` for class merging (never raw template literals)
- Use `tailwind-variants` (`tv()`) for component variants (see `button.svelte` for example)
- Each component folder has an `index.ts` barrel file exporting named components

### Icons

- Use `@lucide/svelte` for icons
  e.g. `import { UserIcon } from "@lucide/svelte";`

### Path Aliases

```typescript
$lib = server/src/lib
$lib/components/ui = UI components
$lib/utils = Utility functions (cn, type helpers)
```

## Development Commands

```bash
cd server
npm run dev          # Start dev server
npm run build        # Production build
npm run check        # TypeScript + Svelte checks
npm run lint         # Prettier + ESLint
npm run format       # Auto-format code
```

## Key Conventions

1. **Props**: Always use TypeScript types, prefer `WithElementRef<T>` for ref forwarding
2. **Styling**: Use Tailwind utilities; semantic color tokens like `bg-primary`, `text-muted-foreground`
3. **Dark Mode**: Uses `.dark` class on root; CSS variables in `layout.css` handle theming
4. **New UI Components**: Run `npx shadcn-svelte@latest add <component>` from `server/`

---

## Svelte MCP Tools

When asked about Svelte/SvelteKit topics, use these tools for documentation:

### 1. list-sections

Use FIRST to discover available documentation sections.

### 2. get-documentation

Fetch full documentation for relevant sections after analyzing list-sections results.

### 3. svelte-autofixer

MUST use when writing Svelte code. Keep calling until no issues returned.

### 4. playground-link

Only after user confirms they want a link (never for code written to project files).
