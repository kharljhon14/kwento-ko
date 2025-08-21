import Footer from '@/components/footer';
import Header from '@/components/header';
import type { User } from '@/types/user';
import { createRootRouteWithContext, Outlet } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

interface RouterContext {
  auth: {
    user: User | undefined;
    isAuthenticated: boolean;
  };
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: () => (
    <>
      <div className="container mx-auto min-h-screen flex flex-col justify-between space-y-10">
        <Header />
        <main className="mb-auto">
          <Outlet />
        </main>
        <TanStackRouterDevtools />
        <Footer />
      </div>
    </>
  )
});
