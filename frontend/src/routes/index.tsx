import agent from '@/apis/agents';
import { Skeleton } from '@/components/ui/skeleton';

import BlogCard from '@/features/blogs/blog-card';
import { useQuery } from '@tanstack/react-query';

import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
  component: Homepage
});

function Homepage() {
  const blogsQuery = useQuery({ queryKey: ['blogs'], queryFn: agent.blogs.getBlogs });

  if (blogsQuery.isLoading) {
    return (
      <div className="grid grid-cols-4 gap-4">
        <Skeleton className="h-[20rem] rounded-lg" />
        <Skeleton className="h-[20rem] rounded-lg" />
        <Skeleton className="h-[20rem] rounded-lg" />
        <Skeleton className="h-[20rem] rounded-lg" />
        <Skeleton className="h-[20rem] rounded-lg" />
        <Skeleton className="h-[20rem] rounded-lg" />
        <Skeleton className="h-[20rem] rounded-lg" />
        <Skeleton className="h-[20rem] rounded-lg" />
      </div>
    );
  }

  const hasBlogs = blogsQuery.data && blogsQuery.isSuccess;

  return (
    <div>
      {hasBlogs && (
        <div className="grid grid-cols-4 gap-4">
          {blogsQuery.data.data.map((blog) => (
            <BlogCard
              key={blog.id}
              blog={blog}
            />
          ))}
        </div>
      )}
    </div>
  );
}
