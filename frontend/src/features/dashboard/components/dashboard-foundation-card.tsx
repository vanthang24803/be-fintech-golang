import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function DashboardFoundationCard() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>What this shell establishes</CardTitle>
        <CardDescription>Stable structure before deeper product work.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4 text-sm leading-6 text-muted-foreground">
        <p>
          The sidebar is data-driven and route-aware, the header supports mobile navigation, and
          the component foundation now covers cards, forms, dialogs, inputs, and tables.
        </p>
        <p>
          That means the next screens can stay focused on business flows instead of rebuilding
          layout, interaction primitives, or styling tokens each time.
        </p>
      </CardContent>
    </Card>
  );
}

