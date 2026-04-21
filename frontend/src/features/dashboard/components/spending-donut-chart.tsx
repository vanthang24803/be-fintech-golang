"use client";

import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } from "recharts";
import { useCategorySummary } from "@/features/dashboard/hooks/use-reports";

const PALETTE = [
  "#001BB7", "#22c55e", "#f59e0b", "#ef4444", "#a855f7",
  "#ec4899", "#14b8a6", "#f97316", "#6366f1", "#84cc16",
];

function fmt(n: number) {
  return n.toLocaleString("vi-VN") + " đ";
}

export function SpendingDonutChart() {
  const { data: summary, isLoading } = useCategorySummary();

  const items = (summary ?? [])
    .filter((c) => c.total_amount > 0)
    .sort((a, b) => b.total_amount - a.total_amount)
    .slice(0, 8);

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
      <h3 className="text-sm font-semibold text-gray-800 mb-4">Phân tích chi tiêu</h3>

      {isLoading ? (
        <div className="h-52 flex items-center justify-center">
          <div className="w-8 h-8 border-2 border-[#001BB7] border-t-transparent rounded-full animate-spin" />
        </div>
      ) : items.length === 0 ? (
        <div className="h-52 flex items-center justify-center text-sm text-gray-400">
          Chưa có dữ liệu
        </div>
      ) : (
        <>
          {/* Donut */}
          <ResponsiveContainer width="100%" height={180}>
            <PieChart>
              <Pie
                data={items}
                dataKey="total_amount"
                nameKey="category_name"
                cx="50%"
                cy="50%"
                innerRadius={52}
                outerRadius={80}
                paddingAngle={2}
              >
                {items.map((_, i) => (
                  <Cell key={i} fill={PALETTE[i % PALETTE.length]} />
                ))}
              </Pie>
              <Tooltip
                formatter={(v: unknown) => [fmt(Number(v)), "Chi tiêu"]}
                contentStyle={{ borderRadius: 10, border: "none", boxShadow: "0 4px 20px rgba(0,0,0,.08)", fontSize: 12 }}
              />
            </PieChart>
          </ResponsiveContainer>

          {/* Legend */}
          <div className="flex flex-wrap gap-x-3 gap-y-1 justify-center mb-4">
            {items.map((item, i) => (
              <div key={item.category_id} className="flex items-center gap-1">
                <span className="w-2.5 h-2.5 rounded-sm flex-shrink-0" style={{ background: PALETTE[i % PALETTE.length] }} />
                <span className="text-xs text-gray-500">{item.category_name}</span>
              </div>
            ))}
          </div>

          {/* Breakdown list */}
          <div>
            <p className="text-xs font-semibold text-gray-400 tracking-widest uppercase mb-3">
              Chi tiêu theo hạng mục
            </p>
            <div className="space-y-3">
              {items.slice(0, 4).map((item, i) => (
                <div key={item.category_id}>
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-sm text-gray-700">{item.category_name}</span>
                    <span className="text-sm font-medium text-gray-700">{fmt(item.total_amount)}</span>
                  </div>
                  <div className="h-1.5 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full"
                      style={{
                        width: `${item.percentage}%`,
                        background: PALETTE[i % PALETTE.length],
                      }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
