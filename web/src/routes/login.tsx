import { Button } from '#/components/ui/button'
import { Card, CardContent } from '#/components/ui/card'
import { Field, FieldGroup, FieldLabel, FieldSet } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs'
import { useLogin } from '#/lib/mutations/login'
import { useSignup } from '#/lib/mutations/signup'
import { cn } from '#/lib/utils'
import { useForm } from '@tanstack/react-form'
import { createFileRoute, redirect } from '@tanstack/react-router'
import { GiSharpShuriken } from 'react-icons/gi'
import { toast } from 'sonner'
import z from 'zod'

const authSchema = z.object({
  form: z.enum(['login', 'signup']),
  email: z.email(),
  password: z.string().min(4),
  secret: z.string(),
})

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => {
    if (context.auth.user) {
      throw redirect({ to: '/team-builder' })
    }
  },
  component: Login,
})

export function Login() {
  const signup = useSignup()
  const login = useLogin()
  const navigate = Route.useNavigate()
  const form = useForm({
    defaultValues: {
      form: 'login',
      email: '',
      password: '',
      secret: '',
    },
    validators: {
      onChange: authSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        if (value.form === 'signup') {
          await signup.mutateAsync(value)
          await navigate({ to: '/team-builder' })
        }
        if (value.form === 'login') {
          await login.mutateAsync(value)
          await navigate({ to: '/team-builder' })
        }
      } catch (e) {
        console.error(e)
        toast.error(String(e))
      }
    },
  })
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <a className="flex items-center justify-center gap-2 self-center font-medium p-4 mb-3">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <GiSharpShuriken className="size-4" />
          </div>
        </a>
        <form
          className={cn('flex flex-col gap-3')}
          onSubmit={form.handleSubmit}
        >
          <form.Field name="form">
            {(field) => (
              <Card>
                <CardContent>
                  <Tabs
                    value={field.state.value}
                    onValueChange={field.handleChange}
                    className="flex flex-col gap-6"
                  >
                    <TabsList className="flex flex-row gap-2 self-center">
                      <TabsTrigger value="login">Log In</TabsTrigger>
                      <TabsTrigger value="signup">Sign Up</TabsTrigger>
                    </TabsList>
                    <TabsContent value="signup">
                      <FieldSet className="w-full max-w-xs">
                        <FieldGroup>
                          <form.Field name="email">
                            {(field) => (
                              <Field>
                                <FieldLabel htmlFor="email">Email</FieldLabel>
                                <Input
                                  id="email"
                                  type="text"
                                  placeholder="you@email.com"
                                  value={field.state.value}
                                  onChange={(e) =>
                                    field.handleChange(e.target.value)
                                  }
                                />
                              </Field>
                            )}
                          </form.Field>
                          <form.Field name="password">
                            {(field) => (
                              <Field>
                                <FieldLabel htmlFor="password">
                                  Password
                                </FieldLabel>
                                <Input
                                  id="password"
                                  type="password"
                                  placeholder="••••••••"
                                  value={field.state.value}
                                  onChange={(e) =>
                                    field.handleChange(e.target.value)
                                  }
                                />
                              </Field>
                            )}
                          </form.Field>
                          <form.Field name="secret">
                            {(field) => (
                              <Field>
                                <FieldLabel htmlFor="secret">Secret</FieldLabel>
                                <Input
                                  id="secret"
                                  type="secret"
                                  placeholder=""
                                  value={field.state.value}
                                  onChange={(e) =>
                                    field.handleChange(e.target.value)
                                  }
                                />
                              </Field>
                            )}
                          </form.Field>
                          <form.Subscribe>
                            {({ canSubmit }) => (
                              <Field>
                                <Button
                                  type="submit"
                                  disabled={!canSubmit}
                                  onClick={() => form.handleSubmit()}
                                >
                                  Sign Up
                                </Button>
                              </Field>
                            )}
                          </form.Subscribe>
                        </FieldGroup>
                      </FieldSet>
                    </TabsContent>
                    <TabsContent value="login">
                      <FieldSet className="w-full max-w-xs">
                        <FieldGroup>
                          <form.Field name="email">
                            {(field) => (
                              <Field>
                                <FieldLabel htmlFor="email">Email</FieldLabel>
                                <Input
                                  id="email"
                                  type="email"
                                  placeholder="you@email.com"
                                  value={field.state.value}
                                  onChange={(e) =>
                                    field.handleChange(e.target.value)
                                  }
                                />
                              </Field>
                            )}
                          </form.Field>
                          <form.Field name="password">
                            {(field) => (
                              <Field>
                                <FieldLabel htmlFor="password">
                                  Password
                                </FieldLabel>
                                <Input
                                  id="password"
                                  type="password"
                                  placeholder="••••••••"
                                  value={field.state.value}
                                  onChange={(e) =>
                                    field.handleChange(e.target.value)
                                  }
                                />
                              </Field>
                            )}
                          </form.Field>
                          <form.Subscribe>
                            {({ canSubmit }) => (
                              <Field>
                                <Button
                                  type="submit"
                                  disabled={!canSubmit}
                                  onClick={() => form.handleSubmit()}
                                >
                                  Log In
                                </Button>
                              </Field>
                            )}
                          </form.Subscribe>
                        </FieldGroup>
                      </FieldSet>
                    </TabsContent>
                  </Tabs>
                </CardContent>
              </Card>
            )}
          </form.Field>
          <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
            Copyright © 2026 snaxm games. All rights reserved.
          </div>
        </form>
      </div>
    </div>
  )
}
