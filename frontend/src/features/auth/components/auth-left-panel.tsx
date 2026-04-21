"use client";

import { BarChart2, Shield, TrendingUp, Wallet } from "lucide-react";

export function AuthLeftPanel() {
  return (
    <div className="hidden lg:flex lg:w-1/2 flex-col justify-between p-12 text-white relative overflow-hidden"
      style={{ background: "linear-gradient(135deg, #00137a 0%, #001BB7 100%)" }}>
      {/* Decorative blobs */}
      <div className="absolute top-0 right-0 w-72 h-72 rounded-full opacity-10"
        style={{ background: "radial-gradient(circle, #fff 0%, transparent 70%)", transform: "translate(30%, -30%)" }} />
      <div className="absolute bottom-0 left-0 w-96 h-96 rounded-full opacity-10"
        style={{ background: "radial-gradient(circle, #fff 0%, transparent 70%)", transform: "translate(-40%, 40%)" }} />

      {/* Logo */}
      <div className="flex items-center gap-3 z-10">
        <div className="w-10 h-10 bg-white/20 rounded-xl flex items-center justify-center backdrop-blur-sm">
          <Wallet className="w-5 h-5 text-white" />
        </div>
        <span className="font-bold text-lg tracking-wide">FINANSMART</span>
      </div>

      {/* Hero text */}
      <div className="z-10 space-y-8">
        <div>
          <h1 className="text-4xl font-extrabold leading-tight mb-1">Làm chủ tài chính</h1>
          <h1 className="text-4xl font-extrabold leading-tight" style={{ color: "#a5b4fc" }}>
            Kiến tạo tương lai.
          </h1>
        </div>

        <div className="space-y-5">
          <FeatureItem
            icon={<BarChart2 className="w-4 h-4 text-white" />}
            title="Theo dõi thông minh"
            desc="Tự động phân loại chi tiêu và thu nhập một cách chính xác."
          />
          <FeatureItem
            icon={<TrendingUp className="w-4 h-4 text-white" />}
            title="Báo cáo trực quan"
            desc="Biểu đồ phân tích chuyên sâu giúp bạn thấu hiểu thói quen chi tiêu."
          />
          <FeatureItem
            icon={<Shield className="w-4 h-4 text-white" />}
            title="Bảo mật tuyệt đối"
            desc="Dữ liệu của bạn được mã hóa và bảo vệ an toàn 24/7."
          />
        </div>
      </div>

      {/* Bottom savings card */}
      <div className="z-10 self-end animate-float">
        <div className="flex items-center gap-3 bg-white/15 backdrop-blur-md border border-white/20 rounded-2xl px-5 py-3 shadow-lg">
          <div className="w-9 h-9 rounded-full flex items-center justify-center" style={{ background: "#22c55e" }}>
            <TrendingUp className="w-4 h-4 text-white" />
          </div>
          <div>
            <p className="text-sm font-semibold">Tiết kiệm tháng này</p>
            <p className="text-xs text-white/70">+15.2% so với tháng trước</p>
          </div>
        </div>
      </div>
    </div>
  );
}

function FeatureItem({ icon, title, desc }: { icon: React.ReactNode; title: string; desc: string }) {
  return (
    <div className="flex items-start gap-3">
      <div className="w-8 h-8 rounded-full bg-white/15 flex items-center justify-center flex-shrink-0 mt-0.5">
        {icon}
      </div>
      <div>
        <p className="font-semibold text-sm">{title}</p>
        <p className="text-xs text-white/70 mt-0.5">{desc}</p>
      </div>
    </div>
  );
}
