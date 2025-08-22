import BlogForm from '@/features/blogs/blog-form';
import { useAuth } from '@/hooks/use-auth';
import { createFileRoute, useRouter } from '@tanstack/react-router';

export const Route = createFileRoute('/_authenticated/blogs/new')({
  component: RouteComponent
});

function RouteComponent() {
  const auth = useAuth();

  const router = useRouter();
  if (auth.isError) router.navigate({ to: '/' });

  return (
    <div>
      <BlogForm />
    </div>
  );
}
