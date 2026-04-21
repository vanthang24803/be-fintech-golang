import Link from "next/link";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { buildApiUrl, getApiBaseUrl } from "@/lib/api/client";

export default function HomePage() {
  return (
    <main className="px-4 py-8 md:px-8 md:py-12">
      <section className="mx-auto max-w-6xl rounded-[2rem] border border-border/70 bg-card/80 p-6 shadow-2xl shadow-blue-950/5 backdrop-blur md:p-10">
        <div className="grid gap-8 lg:grid-cols-[1.4fr_0.8fr] lg:items-end">
          <div>
            <p className="text-xs uppercase tracking-[0.35em] text-primary">Expense Manager</p>
            <h1 className="mt-4 max-w-3xl text-4xl font-semibold leading-none md:text-6xl">
              A sharper web shell for your Go finance backend.
            </h1>
            <p className="mt-6 max-w-2xl text-base leading-7 text-muted-foreground">
              The frontend now uses a shadcn-style component foundation and a responsive app shell.
              Keep the API running, point <code>NEXT_PUBLIC_API_BASE_URL</code> at it, then move into the dashboard experience.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Button asChild size="lg">
                <Link href="/login">Login</Link>
              </Button>
              <Button asChild variant="outline" size="lg">
                <Link href="/register">Register</Link>
              </Button>
              <Button asChild variant="ghost" size="lg">
                <a href={buildApiUrl("/docs")} target="_blank" rel="noreferrer">
                  Open API docs
                </a>
              </Button>
            </div>
          </div>
          <Card className="border-border/80 bg-background/70">
            <CardHeader>
              <CardTitle>Frontend wiring</CardTitle>
              <CardDescription>Current environment and integration baseline.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4 text-sm text-muted-foreground">
              <div>
                <p className="font-medium text-foreground">API base URL</p>
                <code className="mt-1 block rounded-xl bg-muted px-3 py-2 text-xs">{getApiBaseUrl()}</code>
              </div>
              <div>
                <p className="font-medium text-foreground">Suggested next steps</p>
                <p className="mt-1">Hook auth flows, dashboard metrics, and CRUD pages into the new shell.</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </section>
    </main>
  );
}
