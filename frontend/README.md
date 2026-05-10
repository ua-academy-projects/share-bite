# ShareBite Frontend

This is the React frontend for ShareBite, a food discovery and sharing application.

## Prerequisites
- Node.js (v18+)
- npm or yarn

## Setup and Development

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

## Environment Variables

Currently, the application uses Vite proxy configuration (see `vite.config.ts`) to route API calls to the local Go backend microservices:
- `/api/auth` -> `http://localhost:3850`
- `/api/guest` -> `http://localhost:3800`
- `/api/business` -> `http://localhost:3950`

## Authentication Flow

ShareBite uses JWT for authentication.
- **Login/Register**: Handled via `/api/auth/login` and `/api/auth/register`. The JWT token is saved in `localStorage`.
- **Protected Routes**: The `/profile`, `/user/:id`, and `/post/create` routes require authentication and will redirect unauthenticated users to the `/auth` page.
- **API Requests**: The Axios client in `src/api/client.ts` automatically attaches the JWT from `localStorage` as a Bearer token in the `Authorization` header for all requests to backend endpoints.

## Project Structure
- `src/components/`: Reusable UI components
- `src/pages/`: Application pages (Home, Auth, Create Post, Profile, etc.)
- `src/api/`: Axios client and API integration logic
- `src/types/`: TypeScript interfaces and DTOs

## Testing & Linting
Run `npm run tsc` to verify TypeScript typings.
