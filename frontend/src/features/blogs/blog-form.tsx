import { Form, FormControl, FormField, FormItem } from '@/components/ui/form';
import { useForm, type SubmitHandler } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import SimpleMDE from 'react-simplemde-editor';

import 'easymde/dist/easymde.min.css';
import { Button } from '@/components/ui/button';
import { useGetTags } from '@/hooks/tags/use-get-tags';
import { MultiSelect } from '@/components/multi-select';
import { useCreateTag } from '@/hooks/tags/use-create-tag';
import { useCreateBlog } from '@/hooks/blogs/use-create-blog';

const blogFormSchema = z.object({
  title: z.string().min(1).max(60),
  tags: z.string().array().min(1).max(5),
  content: z.string()
});

export type BlogFormSchema = z.infer<typeof blogFormSchema>;

export default function BlogForm() {
  // Tags
  const tagMutation = useCreateTag();
  const onCreateTag = (name: string) => {
    tagMutation.mutate({ name });
  };

  const tagsQuery = useGetTags();
  const tagOptions = (tagsQuery.data?.data ?? []).map((tag) => ({
    label: tag.name,
    value: tag.id
  }));

  const createBlogMutation = useCreateBlog();

  const form = useForm<BlogFormSchema>({
    resolver: zodResolver(blogFormSchema),
    defaultValues: { title: '', content: '' }
  });

  const onSubmit: SubmitHandler<BlogFormSchema> = (values) => {
    createBlogMutation.mutate(values);
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="space-y-4"
      >
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <div>
                  <Label
                    className="mb-2"
                    htmlFor="title"
                  >
                    Title
                  </Label>
                  <Input
                    placeholder="Input title"
                    id="title"
                    {...field}
                  />
                </div>
              </FormControl>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="tags"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <div>
                  <Label
                    className="mb-2"
                    htmlFor="title"
                  >
                    Tags
                  </Label>
                  <MultiSelect
                    placeholder="Select tags"
                    onChange={field.onChange}
                    onCreate={onCreateTag}
                    value={field.value}
                    options={tagOptions}
                  />
                </div>
              </FormControl>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="content"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <SimpleMDE {...field} />
              </FormControl>
            </FormItem>
          )}
        />

        <Button>Create Blog</Button>
      </form>
    </Form>
  );
}
