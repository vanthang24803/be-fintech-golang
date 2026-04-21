"use client";

import { useBudgets } from "@/features/budgets/hooks/use-budgets";
import { useCategories } from "@/features/categories/hooks/use-categories";
import { Loader2, Plus } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function BudgetList() {
  const { data: budgets, isLoading: isBudgetsLoading } = useBudgets();
  const { data: categories, isLoading: isCategoriesLoading } = useCategories();

  if (isBudgetsLoading || isCategoriesLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="size-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between px-2">
        <h2 className="text-2xl font-semibold">Your Budgets</h2>
        <Button size="sm" className="rounded-xl">
          <Plus className="mr-2 h-4 w-4" /> New Budget
        </Button>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {budgets?.length === 0 ? (
          <Card className="col-span-full border-dashed">
            <CardContent className="flex h-40 flex-col items-center justify-center text-muted-foreground">
              <p>No budgets created yet.</p>
              <Button variant="ghost">Create your first budget</Button>
            </CardContent>
          </Card>
        ) : (
          budgets?.map((budget) => {
            const category = categories?.find((c) => c.id === budget.category_id);
            const isExceeded = budget.progress_percent > 100;

            return (
              <Card key={budget.id} className="overflow-hidden border-border/70 bg-card/95 transition-shadow hover:shadow-lg">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-lg">
                      {category ? `${category.icon} ${category.name}` : "Overall Budget"}
                    </CardTitle>
                    <span className={cn(
                      "rounded-full px-2 py-1 text-[10px] font-semibold uppercase tracking-wider",
                      budget.is_active ? "bg-emerald-500/10 text-emerald-600" : "bg-muted text-muted-foreground"
                    )}>
                      {budget.is_active ? "Active" : "Inactive"}
                    </span>
                  </div>
                  <CardDescription>
                    {budget.period} budget
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex items-end justify-between">
                    <div>
                      <p className="text-sm text-muted-foreground">Spent</p>
                      <p className="text-xl font-bold">{budget.current_spending.toLocaleString()} VND</p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm text-muted-foreground">Limit</p>
                      <p className="text-lg font-medium">{budget.amount.toLocaleString()} VND</p>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
                      <div
                        className={cn(
                          "h-full transition-all duration-500",
                          isExceeded ? "bg-destructive" : "bg-primary"
                        )}
                        style={{ width: `${Math.min(budget.progress_percent, 100)}%` }}
                      />
                    </div>
                    <div className="flex items-center justify-between text-xs font-medium">
                      <span className={cn(isExceeded && "text-destructive")}>
                        {budget.progress_percent.toFixed(1)}% used
                      </span>
                      <span className="text-muted-foreground">
                        {Math.max(0, budget.amount - budget.current_spending).toLocaleString()} VND left
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })
        )}
      </div>
    </div>
  );
}
