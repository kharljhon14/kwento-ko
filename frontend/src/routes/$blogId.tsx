import agent from '@/apis/agents';
import { Separator } from '@/components/ui/separator';
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

  return (
    <div className="border bg-card rounded-lg p-10">
      <div>
        <h2 className="text-center text-5xl font-semibold font-orbiton">
          {blogQuery.data?.data.title}
        </h2>
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
