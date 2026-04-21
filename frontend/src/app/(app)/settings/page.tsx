"use client";

import { CategoryList } from "@/features/categories/components/category-list";
import { SourceList } from "@/features/sources/components/source-list";
import { useAuth } from "@/features/auth/context/auth-context";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

export default function SettingsPage() {
  const { user: userData } = useAuth();
  const user = userData?.user;
  const profile = userData?.profile;

  const displayName = profile?.full_name || user?.username || "Guest User";
  const initials = displayName
    ? displayName
        .split(" ")
        .filter(Boolean)
        .map((n) => n[0])
        .join("")
        .toUpperCase()
    : "??";

  return (
    <main className="space-y-8 p-4 md:p-6 lg:p-8 max-w-6xl mx-auto">
      <div className="flex items-center gap-6">
        <Avatar className="size-20">
          <AvatarImage src={profile?.avatar_url || undefined} />
          <AvatarFallback className="text-2xl">{initials}</AvatarFallback>
        </Avatar>
        <div>
          <h1 className="text-3xl font-bold">{displayName}</h1>
          <p className="text-muted-foreground">{profile?.phone_number || user?.email || "Personal workspace"}</p>
        </div>
      </div>

      <div className="grid gap-8">
        <section>
          <h2 className="text-xl font-semibold mb-4">Financial Configuration</h2>
          <div className="grid gap-6 lg:grid-cols-2 lg:items-start">
            <CategoryList />
            <SourceList />
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-4">Account Security</h2>
          <Card>
            <CardHeader>
              <CardTitle>Security Preferences</CardTitle>
              <CardDescription>Manage your password and biometric authentication.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between p-4 border rounded-2xl bg-muted/30">
                <div>
                  <p className="font-medium">FIDO2 Biometric Auth</p>
                  <p className="text-sm text-muted-foreground">Register this device for secure step-up authentication.</p>
                </div>
                <Button variant="outline" disabled>Not supported in this version</Button>
              </div>
            </CardContent>
          </Card>
        </section>
      </div>
    </main>
  );
}
