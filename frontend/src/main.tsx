import { StrictMode } from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider, createRouter } from '@tanstack/react-router';
import './index.css';

// Import the generated route tree
import { routeTree } from './routeTree.gen';
import Providers from './providers/providers';
import { useAuth } from './hooks/use-auth';
import { useUserStore } from './stores/user';

// Create a new router instance
const router = createRouter({
  routeTree,
  defaultPreload: 'intent',
  scrollRestoration: true,
  context: {
    auth: { user: undefined, isAuthenticated: false }
  }
});

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

function InnerApp() {
  const auth = useAuth();

  if (auth.isLoading) {
    return undefined;
  }

  return (
    <Providers>
      <RouterProvider
        router={router}
        context={{ auth }}
      />
    </Providers>
  );
}

// Render the app
const rootElement = document.getElementById('root')!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <StrictMode>
      <Providers>
        <InnerApp />
      </Providers>
    </StrictMode>
  );
}
