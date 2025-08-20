import { Link } from '@tanstack/react-router';
import { Button } from './ui/button';
import { useQuery } from '@tanstack/react-query';
import agent from '@/apis/agents';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { LoaderCircle } from 'lucide-react';

export default function Header() {
  const userQuery = useQuery({ queryKey: ['user'], queryFn: agent.user.getUser, retry: 0 });
  const hasUser = userQuery.isSuccess;

  return (
    <header className=" py-4 flex items-center justify-between">
      <Link
        to="/"
        className="flex items-center"
      >
        <img
          src="/kwento-ko-logo.svg"
          className="h-10 w-10 mr-2"
        />
        <h1 className="uppercase text-xl">Kwento Ko</h1>
      </Link>

      {hasUser && (
        <div>
          <Avatar>
            <AvatarImage src={userQuery.data.data.profile_image} />
            <AvatarFallback>{userQuery.data.data.name[0].toUpperCase()}</AvatarFallback>
          </Avatar>
        </div>
      )}

      {!hasUser && !userQuery.isLoading && (
        <Button asChild>
          <a href="http://localhost:8080/api/v1/auth/google">Sign In</a>
        </Button>
      )}

      {userQuery.isLoading && (
        <div>
          <LoaderCircle className=" animate-spin" />
        </div>
      )}
    </header>
  );
}
