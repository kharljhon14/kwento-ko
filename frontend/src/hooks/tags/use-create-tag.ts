import agent from '@/apis/agents';
import { useMutation, useQueryClient } from '@tanstack/react-query';

export function useCreateTag() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: agent.tags.createTag,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['tags'] });
    }
  });

  return mutation;
}
