import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
  component: Homepage
});

function Homepage() {
  return (
    <div>
      <h3>Homepage</h3>
    </div>
  );
}
