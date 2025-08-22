import { Link } from '@tanstack/react-router';
import { Button } from './ui/button';

import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';

import { useAuth } from '@/hooks/use-auth';

import { LoaderCircle } from 'lucide-react';

export default function Header() {
  const userQuery = useAuth();
  const hasUser = userQuery.isSuccess;

  return (
    <header className="bg-card py-4 flex items-center justify-between border rounded-lg px-12 mt-8">
      <Link
        to="/"
        className="flex items-center"
      >
        <img
          src="/kwento-ko-logo.svg"
          className="h-8 w-8 mr-2"
        />
        <h1 className="uppercase text-lg">Kwento Ko</h1>
      </Link>

      {hasUser && !userQuery.isLoading && (
        <div className="flex items-center gap-2">
          <Avatar>
            <AvatarImage src={userQuery.data.data.profile_image} />
            <AvatarFallback>{userQuery.data.data.name[0].toUpperCase()}</AvatarFallback>
          </Avatar>
          <p>{userQuery.data.data.name}</p>
        </div>
      )}

      {!hasUser && !userQuery.isLoading && (
        <Button asChild>
          <a href="http://localhost:8080/api/v1/auth/google">Sign In</a>
        </Button>
      )}

      {userQuery.isLoading && (
        <div>
          <LoaderCircle className="animate-spin " />
        </div>
      )}
    </header>
  );
}
