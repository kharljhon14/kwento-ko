import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import type { Blog } from '@/types/blog';
import { Link } from '@tanstack/react-router';

type Props = {
  blog: Blog;
};

export default function BlogCard({ blog }: Props) {
  return (
    <Link
      to="/$blogId"
      params={{ blogId: blog.id }}
    >
      <Card className="h-[320px] bg-card text-text">
        <CardHeader>
          <CardTitle className="line-clamp-2 font-orbiton text-primary">{blog.title}</CardTitle>
          <div className="space-x-2 mt-1">
            {blog.tags.map((tag) => (
              <Badge
                variant="outline"
                key={tag}
              >
                {tag}
              </Badge>
            ))}
          </div>
        </CardHeader>
        <Separator />
        <CardContent
          className="line-clamp-7 text-sm h-full"
          dangerouslySetInnerHTML={{ __html: blog.content }}
        ></CardContent>
        <Separator />
        <CardFooter className="mt-auto text-secondary">
          <div className="flex justify-between w-full">
            <small>{blog.author}</small>
            <small>
              {new Intl.DateTimeFormat('en-US', {
                year: 'numeric',
                month: 'long',
                day: '2-digit'
              }).format(new Date(blog.created_at))}
            </small>
          </div>
        </CardFooter>
      </Card>
    </Link>
  );
}
