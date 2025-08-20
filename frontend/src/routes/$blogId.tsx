import agent from '@/apis/agents';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import { useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/$blogId')({
  component: RouteComponent
});

function RouteComponent() {
  const { blogId } = Route.useParams();

  const blogQuery = useQuery({
    queryKey: ['blogs', blogId],
    queryFn: () => agent.blogs.getBlog(blogId)
  });

  if (blogQuery.isLoading) {
    return (
      <>
        <Skeleton className="h-[44rem] rounded-lg" />
      </>
    );
  }

  if (blogQuery.isError) {
    return (
      <div className="border bg-card rounded-lg p-10">
        <h2 className="text-center text-5xl uppercase text-red-500 font-semibold font-orbiton">
          404 not found!
        </h2>
      </div>
    );
  }

  if (blogQuery.isSuccess) {
    return (
      <div className="border bg-card rounded-lg p-10">
        <div>
          <h2 className="text-center text-5xl font-semibold font-orbiton">
            {blogQuery.data.data.title}
          </h2>
          <div className="space-x-2 mt-4 flex items-center justify-center">
            {blogQuery.data.data.tags.map((tag) => (
              <Badge
                variant="outline"
                key={tag}
              >
                {tag}
              </Badge>
            ))}
          </div>
          <div className="mt-4 text-center text-sm">
            <span>By {blogQuery.data?.data.author}</span>
            <span className="mx-2">•</span>
            <span>
              {new Date(blogQuery.data?.data.created_at).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
              })}
            </span>
          </div>
        </div>
        <Separator className="my-6" />
        <div dangerouslySetInnerHTML={{ __html: blogQuery.data?.data.content }}></div>
      </div>
    );
  }
}
