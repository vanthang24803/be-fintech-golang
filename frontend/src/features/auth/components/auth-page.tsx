"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/features/auth/context/auth-context";
import { apiRequest } from "@/lib/api/client";
import { TokenPair } from "@/lib/api/types";
import { Loader2, Lock, LogIn, Mail, User, UserPlus } from "lucide-react";
import Link from "next/link";
import { AuthLeftPanel } from "./auth-left-panel";

/* ─── schemas ─── */
const loginSchema = z.object({
  email: z.string().email({ message: "Email không hợp lệ" }),
  password: z.string().min(8, { message: "Mật khẩu ít nhất 8 ký tự" }),
});
const registerSchema = z
  .object({
    username: z.string().min(3, { message: "Tên người dùng ít nhất 3 ký tự" }),
    email: z.string().email({ message: "Email không hợp lệ" }),
    password: z.string().min(8, { message: "Mật khẩu ít nhất 8 ký tự" }),
    confirmPassword: z.string(),
  })
  .refine((d) => d.password === d.confirmPassword, {
    message: "Mật khẩu không khớp",
    path: ["confirmPassword"],
  });

type Tab = "login" | "register";
type LoginValues = z.infer<typeof loginSchema>;
type RegisterValues = z.infer<typeof registerSchema>;

/* ─── shared helpers ─── */
function FieldWrapper({ label, extra, children }: { label: string; extra?: React.ReactNode; children: React.ReactNode }) {
  return (
    <div>
      <div className="flex items-center justify-between mb-1.5">
        <label className="text-sm font-medium text-gray-700">{label}</label>
        {extra}
      </div>
      {children}
    </div>
  );
}

function IconInput({ icon, ...props }: { icon: React.ReactNode } & React.ComponentProps<typeof Input>) {
  return (
    <div className="relative">
      <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none">{icon}</span>
      <Input className="pl-9 h-11 border-gray-200 rounded-xl bg-white focus-visible:ring-2 focus-visible:ring-[#001BB7]/30 focus-visible:border-[#001BB7]" {...props} />
    </div>
  );
}

function SubmitBtn({ loading, children }: { loading: boolean; children: React.ReactNode }) {
  return (
    <button
      type="submit"
      disabled={loading}
      className="w-full h-11 rounded-xl text-white text-sm font-semibold flex items-center justify-center gap-2 transition-opacity hover:opacity-90 disabled:opacity-60 mt-1"
      style={{ background: "linear-gradient(135deg,#00137a 0%,#001BB7 100%)" }}
    >
      {loading && <Loader2 className="w-4 h-4 animate-spin" />}
      {children}
    </button>
  );
}

function Divider() {
  return (
    <div className="flex items-center gap-3 my-5">
      <div className="flex-1 h-px bg-gray-200" />
      <span className="text-xs font-medium text-gray-400 tracking-widest">HOẶC SỬ DỤNG</span>
      <div className="flex-1 h-px bg-gray-200" />
    </div>
  );
}

function GoogleBtn() {
  return (
    <button
      type="button"
      className="w-full h-11 rounded-xl border border-gray-200 bg-white flex items-center justify-center gap-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
    >
      <GoogleIcon />
      Kết nối với Google
    </button>
  );
}

function Terms() {
  return (
    <p className="text-center text-xs text-gray-400 mt-6">
      Bằng cách tham gia, bạn đồng ý với{" "}
      <Link href="#" className="font-semibold text-gray-600 underline-offset-2 hover:underline">
        Điều khoản dịch vụ
      </Link>{" "}
      và{" "}
      <Link href="#" className="font-semibold text-gray-600 underline-offset-2 hover:underline">
        Chính sách bảo mật
      </Link>{" "}
      của chúng tôi.
    </p>
  );
}

/* ─── login inner form ─── */
function LoginInner({ onSwitch }: { onSwitch: () => void }) {
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const form = useForm<LoginValues>({ resolver: zodResolver(loginSchema), defaultValues: { email: "", password: "" } });

  async function onSubmit(v: LoginValues) {
    setLoading(true);
    setError(null);
    try {
      const data = await apiRequest<{ tokens: TokenPair }>("auth/login", {
        body: JSON.stringify({ identifier: v.email, password: v.password }),
      });
      await login(data.tokens);
    } catch (e: any) {
      setError(e.message || "Thông tin đăng nhập không hợp lệ");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FieldWrapper label="Email">
                <FormControl>
                  <IconInput icon={<Mail className="w-4 h-4" />} placeholder="name@company.com" {...field} />
                </FormControl>
              </FieldWrapper>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FieldWrapper
                label="Mật khẩu"
                extra={
                  <button type="button" className="text-xs font-medium" style={{ color: "#001BB7" }}>
                    Quên mật khẩu?
                  </button>
                }
              >
                <FormControl>
                  <IconInput icon={<Lock className="w-4 h-4" />} type="password" placeholder="••••••••" {...field} />
                </FormControl>
              </FieldWrapper>
              <FormMessage />
            </FormItem>
          )}
        />
        {error && <p className="text-sm text-red-500">{error}</p>}
        <SubmitBtn loading={loading}>Tiếp tục đăng nhập →</SubmitBtn>
      </form>
    </Form>
  );
}

/* ─── register inner form ─── */
function RegisterInner({ onSwitch }: { onSwitch: () => void }) {
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const form = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { username: "", email: "", password: "", confirmPassword: "" },
  });

  async function onSubmit(v: RegisterValues) {
    setLoading(true);
    setError(null);
    try {
      await apiRequest("auth/register", {
        body: JSON.stringify({ username: v.username, email: v.email, password: v.password }),
      });
      const data = await apiRequest<{ tokens: TokenPair }>("auth/login", {
        body: JSON.stringify({ identifier: v.email, password: v.password }),
      });
      await login(data.tokens);
    } catch (e: any) {
      setError(e.message || "Đăng ký thất bại");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FieldWrapper label="Tên người dùng">
                <FormControl>
                  <IconInput icon={<User className="w-4 h-4" />} placeholder="e.g. Nguyễn Văn A" {...field} />
                </FormControl>
              </FieldWrapper>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FieldWrapper label="Email">
                <FormControl>
                  <IconInput icon={<Mail className="w-4 h-4" />} placeholder="name@company.com" {...field} />
                </FormControl>
              </FieldWrapper>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FieldWrapper label="Mật khẩu">
                <FormControl>
                  <IconInput icon={<Lock className="w-4 h-4" />} type="password" placeholder="••••••••" {...field} />
                </FormControl>
              </FieldWrapper>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="confirmPassword"
          render={({ field }) => (
            <FormItem>
              <FieldWrapper label="Nhập lại mật khẩu">
                <FormControl>
                  <IconInput icon={<Lock className="w-4 h-4" />} type="password" placeholder="••••••••" {...field} />
                </FormControl>
              </FieldWrapper>
              <FormMessage />
            </FormItem>
          )}
        />
        {error && <p className="text-sm text-red-500">{error}</p>}
        <SubmitBtn loading={loading}>Hoàn tất đăng ký →</SubmitBtn>
      </form>
    </Form>
  );
}

/* ─── main export ─── */
export function AuthPage({ defaultTab = "login" }: { defaultTab?: Tab }) {
  const [tab, setTab] = useState<Tab>(defaultTab);
  const isLogin = tab === "login";

  return (
    <div className="flex min-h-screen w-full">
      <AuthLeftPanel />

      {/* Right panel */}
      <div className="flex flex-1 items-center justify-center bg-gray-50/60 px-6 py-12">
        {/* Card */}
        <div className="w-full max-w-sm bg-white rounded-2xl shadow-[0_4px_32px_rgba(0,0,0,0.08)] border border-gray-100 px-8 py-9">
          {/* Title — animates when tab changes */}
          <div key={tab} className="mb-7 animate-fade-in">
            <h2 className="text-2xl font-bold text-gray-900">
              {isLogin ? "Chào mừng trở lại" : "Bắt đầu hành trình"}
            </h2>
            <p className="text-sm text-gray-500 mt-1">
              {isLogin
                ? "Vui lòng đăng nhập để tiếp tục quản lý tài chính."
                : "Đăng ký tài khoản để bắt đầu làm chủ thu nhập của bạn."}
            </p>
          </div>

          {/* Tabs */}
          <div className="flex rounded-xl border border-gray-200 p-1 mb-6 bg-gray-50">
            {(["login", "register"] as Tab[]).map((t) => {
              const active = tab === t;
              return (
                <button
                  key={t}
                  type="button"
                  onClick={() => setTab(t)}
                  className={[
                    "flex-1 flex items-center justify-center gap-1.5 rounded-lg py-2 text-sm font-medium transition-all duration-200",
                    active
                      ? "bg-white shadow-sm text-gray-900"
                      : "text-gray-400 hover:text-gray-600",
                  ].join(" ")}
                >
                  {t === "login" ? <LogIn className="w-3.5 h-3.5" /> : <UserPlus className="w-3.5 h-3.5" />}
                  {t === "login" ? "Đăng nhập" : "Đăng ký"}
                </button>
              );
            })}
          </div>

          {/* Form content — fades when tab changes */}
          <div key={tab + "-form"} className="animate-fade-in">
            {isLogin ? <LoginInner onSwitch={() => setTab("register")} /> : <RegisterInner onSwitch={() => setTab("login")} />}
          </div>

          <Divider />
          <GoogleBtn />
          <Terms />
        </div>
      </div>
    </div>
  );
}

function GoogleIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none">
      <path d="M17.64 9.2c0-.637-.057-1.251-.164-1.84H9v3.481h4.844a4.14 4.14 0 01-1.796 2.716v2.259h2.908c1.702-1.567 2.684-3.875 2.684-6.615z" fill="#4285F4" />
      <path d="M9 18c2.43 0 4.467-.806 5.956-2.18l-2.908-2.259c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.711H.957v2.332A8.997 8.997 0 009 18z" fill="#34A853" />
      <path d="M3.964 10.71A5.41 5.41 0 013.682 9c0-.593.102-1.17.282-1.71V4.958H.957A8.996 8.996 0 000 9c0 1.452.348 2.827.957 4.042l3.007-2.332z" fill="#FBBC05" />
      <path d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0A8.997 8.997 0 00.957 4.958L3.964 7.29C4.672 5.163 6.656 3.58 9 3.58z" fill="#EA4335" />
    </svg>
  );
}
