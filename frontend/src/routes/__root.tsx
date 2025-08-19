import Header from '@/components/header';
import { createRootRoute, Outlet } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

export const Route = createRootRoute({
  component: () => (
    <>
      <div>
        <Header />
        <hr />
        <Outlet />
        <TanStackRouterDevtools />
      </div>
    </>
  )
});
