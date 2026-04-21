import { TransactionTable } from "@/features/transactions/components/transaction-table";

export default function TransactionsPage() {
  return (
    <main className="p-4 md:p-6 lg:p-8">
      <TransactionTable />
    </main>
  );
}
