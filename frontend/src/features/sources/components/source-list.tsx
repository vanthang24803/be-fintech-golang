"use client";

import { useSources, useDeleteSource } from "@/features/sources/hooks/use-sources";
import { Loader2, Plus, Trash2, Wallet } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

export function SourceList() {
  const { data: sources, isLoading } = useSources();
  const { mutate: deleteSource } = useDeleteSource();

  if (isLoading) {
    return (
      <div className="flex h-48 items-center justify-center">
        <Loader2 className="size-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle>Payment Sources</CardTitle>
          <CardDescription>Accounts and wallets used for transactions.</CardDescription>
        </div>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" /> Add Source
        </Button>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Balance</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sources?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center text-muted-foreground">
                  No payment sources found.
                </TableCell>
              </TableRow>
            ) : (
              sources?.map((source) => (
                <TableRow key={source.id}>
                  <TableCell className="font-medium flex items-center gap-2">
                    <Wallet className="h-4 w-4 text-primary" />
                    {source.name}
                  </TableCell>
                  <TableCell className="capitalize">{source.type.replace("_", " ")}</TableCell>
                  <TableCell className="font-semibold">
                    {source.balance.toLocaleString()} {source.currency}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button 
                      variant="ghost" 
                      size="icon" 
                      className="text-destructive hover:bg-destructive/10 hover:text-destructive"
                      onClick={() => {
                        if (confirm("Are you sure you want to delete this payment source?")) {
                          deleteSource(source.id);
                        }
                      }}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
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
