import { Link } from '@tanstack/react-router';
import { Button } from './ui/button';

import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';

import { useAuth } from '@/hooks/use-auth';

export default function Header() {
  const { user } = useAuth();

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

      {user && (
        <div className="flex items-center gap-2">
          <Avatar>
            <AvatarImage src={user.profile_image} />
            <AvatarFallback>{user.name[0].toUpperCase()}</AvatarFallback>
          </Avatar>
          <p>{user.name}</p>
        </div>
      )}

      {!user && (
        <Button asChild>
          <a href="http://localhost:8080/api/v1/auth/google">Sign In</a>
        </Button>
      )}
    </header>
  );
}
