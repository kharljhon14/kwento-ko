import agent from '@/apis/agents';
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

  console.log(blogQuery.data?.data.title);
  return <div>{blogQuery.data?.data.name} Hello</div>;
}
