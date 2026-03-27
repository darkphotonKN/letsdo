# SPECIFICATION.md - Frontend

## Project Overview
**Name:** LetsDo Client
**Type:** Next.js Frontend Application
**One-liner:** Frontend UI for Simple todo app LetsDo
**API Backend:** Connects to letsdo API on port 8028

## Domain Terms
| Term | Definition |
|------|------------|
| Todo | The primary entity displayed and managed in the UI |
| Todo List | Table/grid view showing all todos |
| Todo Form | UI for creating/editing todos |
{{if .IncludeAuth}}| Dashboard | Main authenticated user interface |
| Session | User's authenticated state in the application |{{end}}

## Features

### Core (MVP)
- [ ] **Todo List View**: Display all todos in a responsive table
- [ ] **Create Todo Form**: Modal/page for adding new todos
- [ ] **Edit Todo**: In-place or modal editing of existing todos
- [ ] **Delete Todo**: Confirmation dialog and deletion
- [ ] **Todo Details**: Detailed view of single todo
- [ ] **Search & Filter**: Find todos by name or properties
- [ ] **Pagination**: Navigate through large datasets
{{if .IncludeAuth}}- [ ] **User Authentication**: Login/register forms with JWT management
- [ ] **Protected Routes**: Secure pages requiring authentication
- [ ] **User Profile**: View and edit user information{{end}}
{{if .IncludeS3}}- [ ] **File Upload**: Drag-and-drop or click to upload images
- [ ] **Image Preview**: Display uploaded images in todo views{{end}}

### UI/UX Requirements
- [ ] **Responsive Design**: Mobile, tablet, and desktop layouts
- [ ] **Loading States**: Skeleton screens and spinners
- [ ] **Error Handling**: User-friendly error messages
- [ ] **Success Feedback**: Toast notifications for actions
- [ ] **Dark Mode**: Optional theme switching
- [ ] **Accessibility**: WCAG 2.1 AA compliance

### Future Enhancements
- [ ] **Bulk Actions**: Select and act on multiple todos
- [ ] **Export**: Download todos data as CSV/Excel
- [ ] **Real-time Updates**: WebSocket connection for live data
- [ ] **Offline Support**: PWA with service workers
- [ ] **Advanced Filters**: Complex query builder UI

## User Interface Structure

### Page Hierarchy
```
/
├── / (Home/Landing)
{{if .IncludeAuth}}├── /login (Authentication)
├── /register (User Registration)
├── /dashboard (Main App){{end}}
├── /todos (List View)
├── /todos/new (Create Form)
├── /todos/[id] (Detail View)
└── /todos/[id]/edit (Edit Form)
```

### Component Architecture
```
components/
├── ui/                    # Reusable UI primitives
├── layout/                # App layout components
│   ├── Header
│   ├── Sidebar
│   └── Footer
└── features/
    └── todo/
        ├── TodoList
        ├── TodoForm
        ├── TodoCard
        └── TodoDetails
```

## State Management

### Client State (Zustand)
{{if .IncludeAuth}}- **Auth Store**: User session, tokens, login state{{end}}
- **UI Store**: Theme, sidebar state, modals
- **Filter Store**: Active filters, search queries

### Server State (TanStack Query)
- **Todo Data**: List, individual items, mutations
{{if .IncludeAuth}}- **User Data**: Profile, preferences{{end}}
- **Cache Management**: Invalidation strategies

## API Integration

### Data Fetching Patterns
| Operation | Method | Cache Strategy |
|-----------|--------|---------------|
| List todos | GET with pagination | Cache 5 minutes |
| Get Todo | GET by ID | Cache until update |
| Create Todo | POST + invalidate list | Optimistic update |
| Update Todo | PUT + invalidate | Optimistic update |
| Delete Todo | DELETE + remove from cache | Pessimistic update |

{{if .IncludeAuth}}### Authentication Flow
1. User submits login form
2. API returns JWT token
3. Store token in Zustand persist
4. Add token to Axios interceptor
5. Redirect to dashboard
6. Handle 401 responses with re-login{{end}}

## External Dependencies

### Core Libraries
| Library | Purpose | Version |
|---------|---------|---------|
| Next.js | React framework | 15.x |
| React | UI library | 18.x |
| TypeScript | Type safety | 5.x |
| Tailwind CSS | Styling | 3.x |
| shadcn/ui | Component library | Latest |

### Data Management
| Library | Purpose | Version |
|---------|---------|---------|
| TanStack Query | Server state | 5.x |
| Zustand | Client state | 4.x |
| Axios | HTTP client | 1.x |
| React Hook Form | Forms | 7.x |
| Zod | Validation | 3.x |

## Business Rules (Frontend)

### Form Validation
1. **Required Fields**: Name and description are mandatory
2. **Character Limits**: Name max 255, description max 1000
3. **Status Values**: Only allow predefined status options
{{if .IncludeS3}}4. **File Validation**: Max 10MB, images only{{end}}

### User Experience
1. **Confirmation Dialogs**: Require confirmation for destructive actions
2. **Unsaved Changes**: Warn before navigation with unsaved forms
3. **Auto-save**: Save draft every 30 seconds (if implemented)
{{if .IncludeAuth}}4. **Session Timeout**: Re-authenticate after 24 hours{{end}}

## Edge Cases

### Network & API
- **Offline Mode**: Show cached data with offline indicator
- **Slow Network**: Show loading states after 300ms
- **API Errors**: Display user-friendly error messages
- **Rate Limiting**: Handle 429 responses with retry

### Data Handling
- **Empty States**: Show helpful messages when no data
- **Large Lists**: Virtual scrolling for 1000+ items
- **Concurrent Edits**: Handle version conflicts
{{if .IncludeS3}}- **Upload Failures**: Retry mechanism with progress{{end}}

### Browser Compatibility
- **Local Storage**: Fallback if unavailable
- **JavaScript Disabled**: Basic functionality message
- **Old Browsers**: Polyfills for modern features

## Performance Requirements

### Loading Performance
- First Contentful Paint < 1.2s
- Time to Interactive < 2.5s
- Core Web Vitals passing scores

### Runtime Performance
- 60fps scrolling and animations
- Debounced search inputs
- Lazy loading for images and components
- Code splitting per route

### Bundle Size
- Initial JS bundle < 200KB
- Lazy load heavy dependencies
- Tree shaking for unused code

## Accessibility Requirements

### WCAG 2.1 Level AA
- Keyboard navigation for all interactions
- Screen reader announcements
- Color contrast ratios (4.5:1 minimum)
- Focus indicators visible
- ARIA labels and roles

### User Preferences
- Respect prefers-reduced-motion
- Support browser zoom to 200%
- High contrast mode support

## Testing Strategy

### Unit Tests
- Utility functions
- Custom hooks
- Form validation logic

### Component Tests
- Render testing with React Testing Library
- User interaction flows
- Error state handling
{{if .IncludeAuth}}- Authentication flow{{end}}

### E2E Tests (Future)
- Critical user journeys
- Cross-browser testing
- Mobile device testing

## Development Workflow

### Local Development
```bash
npm run dev       # Start on port 3028
npm run build     # Production build
npm run lint      # ESLint checks
npm run test      # Run test suite
```

### Environment Variables
```env
NEXT_PUBLIC_API_URL=http://localhost:8028
NEXT_PUBLIC_APP_NAME=LetsDo
{{if .IncludeAuth}}NEXT_PUBLIC_ENABLE_AUTH=true{{end}}
{{if .IncludeS3}}NEXT_PUBLIC_ENABLE_UPLOADS=true{{end}}
```

---

*This specification defines WHAT the frontend should do. For HOW to implement it, see CLAUDE.md. For testing support, use AGENTS.md with the /test command.*