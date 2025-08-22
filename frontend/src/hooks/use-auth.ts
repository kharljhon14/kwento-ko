import agent from '@/apis/agents';
import { useQuery } from '@tanstack/react-query';

export function useAuth() {
  const query = useQuery({ queryKey: ['user'], queryFn: agent.user.getUser, retry: 0 });

  return query;
}
