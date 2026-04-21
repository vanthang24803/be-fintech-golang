"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useTransactions } from "@/features/transactions/hooks/use-transactions";
import { Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";

import { QuickCaptureDialog } from "./quick-capture-dialog";

export function RecentActivityTable() {
  const { data: transactions, isLoading } = useTransactions();

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Recent activity</CardTitle>
          <CardDescription>Loading recent transactions...</CardDescription>
        </CardHeader>
        <CardContent className="flex h-48 items-center justify-center">
          <Loader2 className="size-8 animate-spin text-primary" />
        </CardContent>
      </Card>
    );
  }

  const recentTransactions = transactions?.slice(0, 5) || [];

  return (
    <Card>
      <CardHeader className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <CardTitle>Recent activity</CardTitle>
          <CardDescription>Real-time view of your latest financial moves.</CardDescription>
        </div>
        <QuickCaptureDialog />
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Category</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Source</TableHead>
              <TableHead>Date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {recentTransactions.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center text-muted-foreground">
                  No recent transactions.
                </TableCell>
              </TableRow>
            ) : (
              recentTransactions.map((row) => (
                <TableRow key={row.id}>
                  <TableCell className="font-medium">{row.category_name}</TableCell>
                  <TableCell className={cn(
                    "font-semibold",
                    row.type === "income" ? "text-primary" : "text-foreground"
                  )}>
                    {row.type === "income" ? "+" : "-"}{row.amount.toLocaleString()} VND
                  </TableCell>
                  <TableCell>{row.source_name}</TableCell>
                  <TableCell className="text-muted-foreground">
                    {new Date(row.transaction_date).toLocaleDateString()}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
