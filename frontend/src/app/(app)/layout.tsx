import { AppShell } from "@/components/app/app-shell";
import { AuthGuard } from "@/features/auth/components/auth-guard";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <AppShell>{children}</AppShell>
    </AuthGuard>
  );
}

