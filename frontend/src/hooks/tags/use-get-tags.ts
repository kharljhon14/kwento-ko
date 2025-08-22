import agent from '@/apis/agents';
import { useQuery } from '@tanstack/react-query';

export function useGetTags() {
  const tags = useQuery({ queryKey: ['tags'], queryFn: agent.tags.getTags });

  return tags;
}
