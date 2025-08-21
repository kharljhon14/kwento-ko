import agent from '@/apis/agents';
import { useQuery } from '@tanstack/react-query';

export function useAuth() {
  const userQuery = useQuery({ queryKey: ['user'], queryFn: agent.user.getUser, retry: 0 });
  const isAuthenticated = userQuery.isSuccess;

  if (userQuery.isError) {
    return { user: undefined, isAuthenticated, isLoading: userQuery.isLoading };
  }

  return { user: userQuery.data?.data, isAuthenticated };
}
