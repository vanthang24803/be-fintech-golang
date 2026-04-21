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
import { Loader2, LogIn, Lock, Mail, User, UserPlus } from "lucide-react";
import Link from "next/link";
import { AuthLeftPanel } from "./auth-left-panel";

const registerSchema = z.object({
  username: z.string().min(3, { message: "Tên người dùng ít nhất 3 ký tự" }),
  email: z.string().email({ message: "Email không hợp lệ" }),
  password: z.string().min(8, { message: "Mật khẩu ít nhất 8 ký tự" }),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Mật khẩu không khớp",
  path: ["confirmPassword"],
});

type RegisterFormValues = z.infer<typeof registerSchema>;

export function RegisterForm() {
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuth();
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { username: "", email: "", password: "", confirmPassword: "" },
  });

  async function onSubmit(values: RegisterFormValues) {
    setIsLoading(true);
    setError(null);
    try {
      await apiRequest("auth/register", {
        body: JSON.stringify({ username: values.username, email: values.email, password: values.password }),
      });
      const loginData = await apiRequest<{ tokens: TokenPair }>("auth/login", {
        body: JSON.stringify({ identifier: values.email, password: values.password }),
      });
      await login(loginData.tokens);
    } catch (err: any) {
      setError(err.message || "Đăng ký thất bại");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen w-full">
      <AuthLeftPanel />

      {/* Right panel */}
      <div className="flex flex-1 items-center justify-center bg-white px-8 py-12">
        <div className="w-full max-w-sm">
          {/* Title */}
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-gray-900">Bắt đầu hành trình</h2>
            <p className="text-sm text-gray-500 mt-1">Đăng ký tài khoản để bắt đầu làm chủ thu nhập của bạn.</p>
          </div>

          {/* Tabs */}
          <div className="flex rounded-xl border border-gray-200 p-1 mb-7 bg-gray-50">
            <Link href="/login"
              className="flex-1 flex items-center justify-center gap-1.5 rounded-lg py-2 text-sm font-medium text-gray-400 hover:text-gray-600 transition-colors">
              <LogIn className="w-3.5 h-3.5" />
              Đăng nhập
            </Link>
            <div className="flex-1 flex items-center justify-center gap-1.5 rounded-lg bg-white shadow-sm py-2 text-sm font-medium text-gray-900 cursor-default">
              <UserPlus className="w-3.5 h-3.5" />
              Đăng ký
            </div>
          </div>

          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              {/* Username */}
              <FormField control={form.control} name="username"
                render={({ field }) => (
                  <FormItem>
                    <label className="text-sm font-medium text-gray-700 block mb-1.5">Tên người dùng</label>
                    <FormControl>
                      <div className="relative">
                        <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input className="pl-9 h-11 border-gray-200 rounded-xl" placeholder="e.g. Nguyễn Văn A" {...field} />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Email */}
              <FormField control={form.control} name="email"
                render={({ field }) => (
                  <FormItem>
                    <label className="text-sm font-medium text-gray-700 block mb-1.5">Email</label>
                    <FormControl>
                      <div className="relative">
                        <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input className="pl-9 h-11 border-gray-200 rounded-xl" placeholder="name@company.com" {...field} />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Password */}
              <FormField control={form.control} name="password"
                render={({ field }) => (
                  <FormItem>
                    <label className="text-sm font-medium text-gray-700 block mb-1.5">Mật khẩu</label>
                    <FormControl>
                      <div className="relative">
                        <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input type="password" className="pl-9 h-11 border-gray-200 rounded-xl" placeholder="••••••••" {...field} />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Confirm Password */}
              <FormField control={form.control} name="confirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <label className="text-sm font-medium text-gray-700 block mb-1.5">Nhập lại mật khẩu</label>
                    <FormControl>
                      <div className="relative">
                        <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input type="password" className="pl-9 h-11 border-gray-200 rounded-xl" placeholder="••••••••" {...field} />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {error && <p className="text-sm text-red-500">{error}</p>}

              <button type="submit" disabled={isLoading}
                className="w-full h-11 rounded-xl text-white font-semibold flex items-center justify-center gap-2 mt-2 transition-opacity hover:opacity-90 disabled:opacity-60"
                style={{ background: "linear-gradient(135deg, #3547e8 0%, #4f63f0 100%)" }}>
                {isLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                Hoàn tất đăng ký →
              </button>
            </form>
          </Form>

          {/* Divider */}
          <div className="flex items-center gap-3 my-5">
            <div className="flex-1 h-px bg-gray-200" />
            <span className="text-xs font-medium text-gray-400 tracking-widest">HOẶC SỬ DỤNG</span>
            <div className="flex-1 h-px bg-gray-200" />
          </div>

          {/* Google */}
          <button className="w-full h-11 rounded-xl border border-gray-200 bg-white flex items-center justify-center gap-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
            <GoogleIcon />
            Kết nối với Google
          </button>

          {/* Terms */}
          <p className="text-center text-xs text-gray-400 mt-7">
            Bằng cách tham gia, bạn đồng ý với{" "}
            <Link href="#" className="font-semibold text-gray-600 underline-offset-2 hover:underline">Điều khoản dịch vụ</Link>
            {" "}và{" "}
            <Link href="#" className="font-semibold text-gray-600 underline-offset-2 hover:underline">Chính sách bảo mật</Link>
            {" "}của chúng tôi.
          </p>
        </div>
      </div>
    </div>
  );
}

function GoogleIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M17.64 9.2c0-.637-.057-1.251-.164-1.84H9v3.481h4.844a4.14 4.14 0 01-1.796 2.716v2.259h2.908c1.702-1.567 2.684-3.875 2.684-6.615z" fill="#4285F4"/>
      <path d="M9 18c2.43 0 4.467-.806 5.956-2.18l-2.908-2.259c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.711H.957v2.332A8.997 8.997 0 009 18z" fill="#34A853"/>
      <path d="M3.964 10.71A5.41 5.41 0 013.682 9c0-.593.102-1.17.282-1.71V4.958H.957A8.996 8.996 0 000 9c0 1.452.348 2.827.957 4.042l3.007-2.332z" fill="#FBBC05"/>
      <path d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0A8.997 8.997 0 00.957 4.958L3.964 7.29C4.672 5.163 6.656 3.58 9 3.58z" fill="#EA4335"/>
    </svg>
  );
}
