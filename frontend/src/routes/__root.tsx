import Header from '@/components/header';
import { createRootRoute, Outlet } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="container mx-auto">
        <Header />

        <Outlet />
        <TanStackRouterDevtools />
      </div>
    </>
  )
});
