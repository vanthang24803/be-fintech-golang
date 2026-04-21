import { FundList } from "@/features/funds/components/fund-list";
import { GoalList } from "@/features/funds/components/goal-list";

export default function FundsPage() {
  return (
    <main className="p-4 md:p-6 lg:p-8 space-y-12">
      <FundList />
      <GoalList />
    </main>
  );
}
