# AGENTS.md - Agentic Coding Guidelines

This document provides guidelines for agents working in this codebase.

## Project Overview

- **Type**: React 19 + TypeScript + Vite application
- **Package Manager**: pnpm
- **Styling**: Tailwind CSS v4 with shadcn/ui patterns
- **Routing**: React Router 7
- **Forms**: React Hook Form + Zod validation
- **Data Fetching**: TanStack Query (React Query)
- **UI Components**: Base UI primitives with class-variance-authority

---

## Commands

### Development
```bash
pnpm dev          # Start Vite dev server
pnpm preview      # Preview production build
```

### Build & Lint
```bash
pnpm build        # Type-check (tsc -b) + Vite build
pnpm lint         # Run ESLint on all files
```

### Running a Single Test
> **Note**: Currently no test framework is configured. If tests are added (e.g., Vitest), run them with:
```bash
pnpm test         # Run all tests
pnpm test run <file>  # Run single test file (if using Vitest)
```

---

## Code Style Guidelines

### TypeScript
- **Strict mode enabled** in `tsconfig.json`
- Use `type` for type aliases, `interface` for object shapes
- Use `as const` for literal types where appropriate
- Avoid `any`, use `unknown` when type is uncertain
- Path alias: use `@/` prefix (e.g., `@/components/ui/button`)

### Imports
- Order: external libs → internal aliases → relative
- Use path alias `@/` for internal imports (configured in tsconfig + vite.config)
- Example:
  ```typescript
  import { useState } from "react"
  import { Button } from "@/components/ui/button"
  import { cn } from "@/lib/utils"
  import { SomeComponent } from "./SomeComponent"
  ```

### React Patterns
- Use functional components with explicit prop typing
- Prefer composition over inheritance
- Use React 19 features (use of `use` hook if applicable)
- Memoize expensive computations with `useMemo`/`useCallback`

### Error Handling
- Use Zod for runtime validation with React Hook Form
- Handle async errors with try/catch + user feedback
- Define error types explicitly rather than using generic errors

### CSS & Styling
- Use Tailwind CSS utility classes
- Define custom theme variables in `src/index.css` using `@theme inline`
- Use `cn()` utility (clsx + tailwind-merge) for conditional classes
- Use CSS custom properties (e.g., `var(--radius)`) for consistent values
- Use shadcn/ui color tokens: `primary`, `secondary`, `muted`, `destructive`, `border`, `input`, `ring`, `foreground`, `background`

### Component Patterns
- Use `cva` (class-variance-authority) for variant components
- Export both component and variants:
  ```typescript
  const buttonVariants = cva("...", { variants: {...} })
  function Button({ className, variant, ...props }) { ... }
  export { Button, buttonVariants }
  ```
- Use `data-slot` attribute for polymorphic components
- Follow shadcn/ui conventions for component structure

### Naming Conventions
- **Files**: kebab-case (e.g., `button.tsx`, `some-component.tsx`)
- **Components**: PascalCase (e.g., `Button`, `SomeComponent`)
- **Hooks**: camelCase with `use` prefix (e.g., `useAuth`, `useCounter`)
- **Constants**: UPPER_SNAKE_CASE for config, camelCase for others
- **Props**: camelCase, prefix event handlers with `on` (e.g., `onClick`, `onSubmit`)

### File Organization
```
src/
├── components/
│   └── ui/          # Reusable UI components
├── lib/
│   └── utils.ts     # Utilities (cn, etc.)
├── App.tsx          # Root component
└── main.tsx         # Entry point
```

---

## Routing

Routes are defined in `src/routes/index.tsx` using React Router 7's `createBrowserRouter`. The file exports a `router` object that is used by `RouterProvider` in `App.tsx`.

### Protected Routes
All pages requiring authentication are wrapped inside the `ProtectedLayout` component (`src/components/protected-layout.tsx`). This layout includes the sidebar, header, and an `<Outlet />` for nested routes. Protected routes use the `requireAuth` loader to enforce authentication.

Public routes (login, register) use the `requireGuest` loader to redirect authenticated users away.

---

## Linting

- ESLint is configured with:
  - `@eslint/js` (recommended)
  - `typescript-eslint` (recommended)
  - `eslint-plugin-react-hooks` (recommended)
  - `eslint-plugin-react-refresh`
  - `@tanstack/eslint-plugin-query`
- Run `pnpm lint` before committing
- Fix auto-fixable issues with `pnpm lint --fix`


---

## Adding Dependencies
> **Note**: Only user can adding new dependencies
