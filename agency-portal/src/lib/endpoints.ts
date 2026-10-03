// Typed wrappers for every API endpoint the frontend uses.
import { api } from "./api";
import type {
  Agency,
  Booking,
  BookingStatus,
  HomeData,
  Member,
  MemberProfile,
  Page,
  PaymentMethod,
  Permission,
  SearchResult,
  Tour,
  TourStatus,
  TourStatusChange,
} from "./types";

function query(params: Record<string, string | number | undefined>): string {
  const q = new URLSearchParams();
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== "") q.set(k, String(v));
  }
  const s = q.toString();
  return s ? `?${s}` : "";
}

// ---- public
export const getHome = () => api.get<{ success: boolean; data: HomeData }>("/home").then((r) => r.data);

export interface SearchParams {
  q?: string;
  min_price?: string;
  max_price?: string;
  start_date?: string;
  end_date?: string;
}
export const searchTours = (p: SearchParams) => api.get<SearchResult>(`/search${query({ ...p })}`);

export const getTour = (id: string) => api.get<Tour>(`/tours/${id}`);
export const getAgency = (id: string) => api.get<Agency>(`/agency/${id}`);
export const listAgencyTours = (agencyId: string, page = 1, limit = 12) =>
  api.get<Page<Tour>>(`/agency/${agencyId}/tours/list${query({ page, limit })}`);

// ---- accounts
export interface RegisterInput {
  name: string;
  email: string;
  password: string;
  phone: string;
}
export const register = (input: RegisterInput) => api.post<string>("/users", input);
export const loginUser = (email: string, password: string) => api.post<string>("/users/login", { email, password });
export const loginMember = (email: string, password: string) => api.post<string>("/members/login", { email, password });

// ---- bookings
export interface BookingInput {
  number_of_people: number;
  total_price: number;
  method: PaymentMethod;
  transaction_id: string;
}
export interface GuestBookingInput extends BookingInput {
  customer_name: string;
  customer_email: string;
  customer_phone: string;
}
export const createBooking = (tourId: string, input: BookingInput) => api.post<Booking>(`/bookings/${tourId}`, input);
export const createGuestBooking = (tourId: string, input: GuestBookingInput) =>
  api.post<Booking>(`/admin/bookings/${tourId}`, input);

export interface BookingFilter {
  status?: BookingStatus | "";
  tour_id?: string;
  page?: number;
  limit?: number;
}
export const myBookings = (f: BookingFilter) => api.get<Page<Booking>>(`/me/bookings${query({ ...f })}`);
export const cancelMyBooking = (id: string) => api.post<Booking>(`/me/bookings/${id}/cancel`);

export const agencyBookings = (agencyId: string, f: BookingFilter) =>
  api.get<Page<Booking>>(`/agency/${agencyId}/bookings${query({ ...f })}`);
export const setBookingStatus = (agencyId: string, bookingId: string, status: "confirmed" | "cancelled" | "completed") =>
  api.patch<Booking>(`/agency/${agencyId}/bookings/${bookingId}/status`, { status });

// ---- agency tours
export const createTour = (agencyId: string, form: FormData) => api.post<unknown>(`/agency/${agencyId}/tours`, form);

export interface TourUpdateInput {
  name: string;
  description: string;
  start_date: string;
  end_date: string;
  last_enrollment_date: string;
  total_seat: number;
  price: number;
  discount: number;
}
export const updateTour = (agencyId: string, tourId: string, input: TourUpdateInput) =>
  api.put<{ message: string; total_seat: number; available_seat: number }>(`/agency/${agencyId}/tours/${tourId}`, input);

export interface ImageUpdateResult {
  image_path: string;
  image_url: string;
  message: string;
}
export const updateTourImage = (agencyId: string, tourId: string, form: FormData) =>
  api.put<ImageUpdateResult>(`/agency/${agencyId}/tours/${tourId}/image`, form);
export const updateAgencyImage = (agencyId: string, form: FormData) =>
  api.put<ImageUpdateResult>(`/agency/${agencyId}/image`, form);
export const setTourStatus = (tourId: string, status: TourStatus) =>
  api.patch<TourStatusChange>(`/tours/${tourId}/tour-status`, status);
export const deleteTour = (tourId: string) => api.delete<string>(`/tours/${tourId}`);

// ---- agency members
export interface MemberInput {
  name: string;
  email: string;
  phone: string;
  password: string;
  role_name: string;
  permissions: number[];
}
export const listMembers = (agencyId: string) => api.get<Member[]>(`/members/${agencyId}`);
export const getMyProfile = () => api.get<MemberProfile>("/members/me");
export const createMember = (agencyId: string, input: MemberInput) => api.post<string>(`/members/${agencyId}`, input);
export const updateMemberPermissions = (memberId: string, permissions: number[]) =>
  api.put<string>(`/members/${memberId}/permissions`, { permissions });
export const deleteMember = (memberId: string) => api.delete<string>(`/members/${memberId}`);
export const listPermissions = () => api.get<Permission[]>("/permissions");

// ---- platform admin
export const createAgency = (form: FormData) =>
  api.post<{ message: string; agency_id: string; image_url: string }>("/agency", form);
