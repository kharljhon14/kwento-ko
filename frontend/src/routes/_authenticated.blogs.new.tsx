import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_authenticated/blogs/new')({
  component: RouteComponent,
  beforeLoad: ({ context, location }) => {
    console.log(context.auth);
    if (!context.auth.isAuthenticated) {
      throw redirect({
        to: '/',
        search: {
          redirect: location.href
        }
      });
    }
  }
});

function RouteComponent() {
  return <div>Hello "/_authenticated/blogs/new"!</div>;
}
