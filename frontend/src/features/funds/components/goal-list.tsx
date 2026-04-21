"use client";

import { useGoals } from "@/features/funds/hooks/use-goals";
import { Loader2, Plus, Flag, Calendar } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function GoalList() {
  const { data: goals, isLoading } = useGoals();

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="size-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between px-2">
        <h2 className="text-2xl font-semibold text-primary">Savings Goals</h2>
        <Button size="sm" variant="outline" className="rounded-xl">
          <Plus className="mr-2 h-4 w-4" /> New Goal
        </Button>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {goals?.length === 0 ? (
          <Card className="col-span-full border-dashed bg-muted/20">
            <CardContent className="flex h-40 flex-col items-center justify-center text-muted-foreground">
              <p>No savings goals set yet.</p>
              <Button variant="ghost">Set your first big goal</Button>
            </CardContent>
          </Card>
        ) : (
          goals?.map((goal) => (
            <Card key={goal.id} className="overflow-hidden border-border/70 bg-card/95 transition-shadow hover:shadow-lg">
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-lg flex items-center gap-2">
                    <Flag className={cn("h-5 w-5", goal.is_reached ? "text-emerald-500" : "text-primary")} />
                    {goal.name}
                  </CardTitle>
                  {goal.is_reached && (
                    <span className="rounded-full bg-emerald-500/10 px-2 py-1 text-[10px] font-semibold text-emerald-600 uppercase">
                      Completed
                    </span>
                  )}
                </div>
                {goal.deadline && (
                  <CardDescription className="flex items-center gap-1">
                    <Calendar className="h-3 w-3" />
                    Target: {new Date(goal.deadline).toLocaleDateString()}
                  </CardDescription>
                )}
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-end justify-between">
                  <div>
                    <p className="text-sm text-muted-foreground">Saved</p>
                    <p className="text-xl font-bold">{goal.saved_amount.toLocaleString()} VND</p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-muted-foreground">Target</p>
                    <p className="text-lg font-medium">{goal.target_amount.toLocaleString()} VND</p>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
                    <div
                      className={cn(
                        "h-full transition-all duration-500",
                        goal.is_reached ? "bg-emerald-500" : "bg-primary"
                      )}
                      style={{ width: `${Math.min(goal.progress_percent, 100)}%` }}
                    />
                  </div>
                  <div className="flex items-center justify-between text-xs font-medium">
                    <span>{goal.progress_percent.toFixed(1)}% complete</span>
                    {!goal.is_reached && (
                      <span className="text-muted-foreground">
                        {(goal.target_amount - goal.saved_amount).toLocaleString()} VND to go
                      </span>
                    )}
                  </div>
                </div>
                
                {!goal.is_reached && (
                  <Button className="w-full rounded-xl" size="sm">Add Funds</Button>
                )}
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  );
}
