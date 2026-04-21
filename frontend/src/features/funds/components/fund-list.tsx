"use client";

import { useFunds } from "@/features/funds/hooks/use-funds";
import { Loader2, Plus, Target, PiggyBank } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function FundList() {
  const { data: funds, isLoading } = useFunds();

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
        <h2 className="text-2xl font-semibold">Savings Pockets</h2>
        <Button size="sm" className="rounded-xl">
          <Plus className="mr-2 h-4 w-4" /> New Fund
        </Button>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {funds?.length === 0 ? (
          <Card className="col-span-full border-dashed">
            <CardContent className="flex h-40 flex-col items-center justify-center text-muted-foreground">
              <p>No savings pockets created yet.</p>
              <Button variant="ghost">Start your first saving fund</Button>
            </CardContent>
          </Card>
        ) : (
          funds?.map((fund) => {
            const progress = fund.target_amount > 0 
              ? (fund.balance / fund.target_amount) * 100 
              : 0;

            return (
              <Card key={fund.id} className="overflow-hidden border-border/70 bg-card/95 transition-shadow hover:shadow-lg">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-lg flex items-center gap-2">
                      <PiggyBank className="h-5 w-5 text-primary" />
                      {fund.name}
                    </CardTitle>
                    <span className="rounded-full bg-primary/10 px-2 py-1 text-[10px] font-semibold text-primary uppercase">
                      {fund.currency}
                    </span>
                  </div>
                  <CardDescription>
                    {fund.description || "No description"}
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex items-end justify-between">
                    <div>
                      <p className="text-sm text-muted-foreground">Current Balance</p>
                      <p className="text-xl font-bold">{fund.balance.toLocaleString()} {fund.currency}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm text-muted-foreground">Target</p>
                      <p className="text-lg font-medium">{fund.target_amount > 0 ? `${fund.target_amount.toLocaleString()} ${fund.currency}` : "N/A"}</p>
                    </div>
                  </div>

                  {fund.target_amount > 0 && (
                    <div className="space-y-2">
                      <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
                        <div
                          className="h-full bg-primary transition-all duration-500"
                          style={{ width: `${Math.min(progress, 100)}%` }}
                        />
                      </div>
                      <div className="flex items-center justify-between text-xs font-medium">
                        <span className="flex items-center gap-1">
                          <Target className="h-3 w-3" />
                          {progress.toFixed(1)}% reached
                        </span>
                        <span className="text-muted-foreground">
                          {Math.max(0, fund.target_amount - fund.balance).toLocaleString()} {fund.currency} left
                        </span>
                      </div>
                    </div>
                  )}

                  <div className="flex gap-2 pt-2">
                    <Button variant="secondary" size="sm" className="flex-1 rounded-xl">Deposit</Button>
                    <Button variant="outline" size="sm" className="flex-1 rounded-xl">Withdraw</Button>
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
