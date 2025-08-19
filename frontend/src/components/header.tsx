import { Link } from '@tanstack/react-router';
import { Button } from './ui/button';

export default function Header() {
  return (
    <header>
      <div>
        <Link to="/">Kwento Ko</Link>
        <Button asChild>
          <a href="http://localhost:8080/api/v1/auth/google">Sign In</a>
        </Button>
      </div>
    </header>
  );
}
