import agent from '@/apis/agents';
import { useMutation, useQueryClient } from '@tanstack/react-query';

export function useCreateBlog() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: agent.blogs.createBlog,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['blogs'] });
    }
  });

  return mutation;
}
