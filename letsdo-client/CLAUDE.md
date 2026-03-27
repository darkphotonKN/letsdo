# LetsDo Frontend - Development Guide

This is a Next.js 15 frontend application for {{.ProjectDescription}}.

**IMPORTANT**: This document contains implementation details (HOW we build). For project requirements and specifications (WHAT we're building), refer to `SPECIFICATION.md`. For testing support, use `AGENTS.md` with the `/test` command.

## Tech Stack

- **Next.js 15** - React framework with App Router
- **TypeScript** - Type-safe development
- **Tailwind CSS** - Utility-first styling
- **shadcn/ui** - Copy-paste component library
- **Zustand** - State management
- **TanStack Query v5** - Data fetching & caching
- **Axios** - HTTP client with interceptors
- **React Hook Form + Zod** - Forms & validation

## API Connection

This frontend is configured to connect to the Go API on port 8028.
API URL is configured in `.env` file as `NEXT_PUBLIC_API_URL`.

## Project Structure

```
src/
├── app/              # Next.js App Router pages
├── components/       # Reusable UI components
│   └── ui/          # shadcn/ui components
├── features/        # Feature-based modules
│   └── item/        # Item feature
│       ├── components/  # Item-specific components
│       ├── hooks/      # React Query hooks
│       ├── services/   # API service layer
│       └── types/      # TypeScript types
├── lib/             # Utilities and configurations
│   └── api/         # Axios client setup
├── stores/          # Zustand state stores
└── types/           # Global TypeScript types
```

## Available Scripts

```bash
npm run dev        # Start development server on port 3028
npm run build      # Build for production
npm run start      # Start production server
npm run lint       # Run ESLint
npm run format     # Format code with Prettier
npm run type-check # Check TypeScript types
```

## Reference Docs

Before implementing features that touch APIs, state management, or UI patterns, check the relevant doc in `docs/`:

| Task | Read First |
|------|------------|
| Component patterns | `docs/ui/*.md` |
| API integration | `docs/api/*.md` |
| State management patterns | `docs/state/*.md` |
| Authentication flow | `docs/auth/*.md` |
| File upload handling | `docs/uploads/*.md` |
| Form validation rules | `docs/validation/*.md` |
| Design system usage | `docs/design-system/*.md` |
| Generated plans | `docs/plans/*.md` |

**Do NOT guess component APIs, state shapes, or validation rules.** The docs have the exact specs from the design document.

The docs/schema/ folder contains one file per table with exact field specs, types, constraints, indexes, and business rules extracted from the design document.

If a doc doesn't exist for what you're implementing, ask before proceeding. If there are very sensible defaults for a component or API integration then you can use that while asking for confirmation.

## Development Workflow

### 1. Entity Management

The main entity is **{{.EntityCapitalized}}** with the following operations:
- List all {{.EntityPlural}}
- Create new {{.PrimaryEntity}}
- Update existing {{.PrimaryEntity}}
- Delete {{.PrimaryEntity}}

### 2. API Integration

All API calls go through Axios with:
- Automatic token injection (if auth enabled)
- Error handling and toast notifications
- Request/response interceptors

### 3. State Management

- **Zustand** for client state (auth, UI)
- **TanStack Query** for server state (API data)
- Automatic cache invalidation on mutations

### 4. Type Safety

All API responses and requests are typed using TypeScript interfaces.
The main entity type is defined in `src/features/{{.PrimaryEntity}}/types/index.ts`.

## Adding New Features

### 1. Create Feature Module

```bash
mkdir -p src/features/new-feature/{components,hooks,services,types}
```

### 2. Define Types

```typescript
// src/features/new-feature/types/index.ts
export interface NewFeature {
  id: string;
  // ... fields
}
```

### 3. Create Service

```typescript
// src/features/new-feature/services/api.ts
import { apiClient } from "@/lib/api/client";

export const newFeatureService = {
  getAll: async () => {
    const { data } = await apiClient.get("/new-features");
    return data;
  },
  // ... other methods
};
```

### 4. Create Hooks

```typescript
// src/features/new-feature/hooks/use-new-feature.ts
import { useQuery } from "@tanstack/react-query";
import { newFeatureService } from "../services/api";

export const useNewFeatureList = () => {
  return useQuery({
    queryKey: ["new-features"],
    queryFn: newFeatureService.getAll,
  });
};
```

### 5. Create Components

Build UI components using shadcn/ui primitives and Tailwind CSS.

## Authentication Flow{{if .IncludeAuth}}

Authentication is handled through:
1. Login form sends credentials to `/api/auth/login`
2. JWT token stored in Zustand persist store
3. Axios interceptor adds token to all requests
4. 401 responses trigger automatic logout{{else}}

Authentication is not enabled for this project.{{end}}

## Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:8028
NEXT_PUBLIC_APP_NAME=LetsDo{{if .IncludeAuth}}
NEXT_PUBLIC_ENABLE_AUTH=true{{end}}{{if .IncludeS3}}
NEXT_PUBLIC_ENABLE_UPLOADS=true{{end}}
```

## Common Tasks

### Add a New shadcn/ui Component

Components are already included in `src/components/ui/`.
To add more, copy from [ui.shadcn.com](https://ui.shadcn.com).

### Update API Endpoint

Edit `src/lib/api/endpoints.ts` to add new endpoints.

### Handle Loading States

Use TanStack Query's built-in loading states:
```typescript
const { data, isLoading, error } = useQuery(...);
```

### Show Toast Notifications

```typescript
import { useToast } from "@/components/ui/use-toast";

const { toast } = useToast();
toast({
  title: "Success",
  description: "Operation completed",
});
```

## Troubleshooting

### CORS Issues
Ensure the Go API has proper CORS configuration for `http://localhost:3028`.

### Port Conflicts
If port 3028 is in use, update the port in `package.json` scripts.

### API Connection Failed
1. Check if the API is running on port 8028
2. Verify `NEXT_PUBLIC_API_URL` in `.env`
3. Check browser console for errors

## Production Build

```bash
npm run build
npm run start
```

The production build will be optimized and ready for deployment.