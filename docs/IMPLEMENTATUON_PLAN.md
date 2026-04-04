# Đề xuất bổ sung Roadmap (Giai đoạn 2)

Sau khi hoàn thành các tính năng cơ bản, tôi đề xuất mở rộng hệ thống với các module sau để tăng tính cạnh tranh và tiện ích cho người dùng.

## 1. Các Module đề xuất bổ sung

### 🎯 Savings Goals (Mục tiêu tiết kiệm) [DONE - 2026-03-29]

- Giúp người dùng đạt được các mục tiêu tài chính cụ thể.
- Thuộc tính: Tên mục tiêu, Số tiền mục tiêu, Số tiền đã tích lũy, Thời hạn.
- Logic: Tích hợp với Notification để nhắc nhở khi hụt tiến độ.

### 🏦 Debt & Loan (Quản lý Nợ)

- Quản lý các khoản nợ phải trả (Mortgage, Credit Card) và nợ phải thu (Lent money).
- Theo dõi lịch sử trả nợ từng phần.

### 💱 Multi-currency (Đa tiền tệ)

- Hỗ trợ lưu trữ số dư bằng nhiều loại tiền tệ.
- Tự động quy đổi tổng tài sản (Net Worth) về tiền tệ mặc định.

### 👥 Shared Funds (Quỹ chung)

- Cho phép quản lý tài chính nhóm/gia đình.
- Permission system: Chủ quỹ (Admin) và Thành viên (Viewer/Editor).

### 📊 Data Export (Xuất dữ liệu)

- Cung cấp tính năng xuất báo cáo CSV/XLSX.
- Gửi báo cáo định kỳ qua Email.

## 2. Thay đổi trong ARCHITECTURE.md

- Cập nhật Gantt chart để bổ sung timeline cho Giai đoạn 2.
- Thêm định nghĩa các đối tượng mới (Goals, Debts) vào sơ đồ Domain Model.

## Câu hỏi cho bạn

- Bạn muốn ưu tiên triển khai tính năng nào đầu tiên trong danh sách trên?
- Bạn có muốn thêm tính năng **Split Bill** (Chia tiền hóa đơn) vào phần Quỹ chung không?
