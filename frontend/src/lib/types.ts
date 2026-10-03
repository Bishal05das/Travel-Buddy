// Shapes returned by the Travel Buddy API.

export type Role = "super" | "user" | "member";

export type TourStatus = "open" | "closed" | "cancelled";
export type BookingStatus = "pending" | "confirmed" | "cancelled" | "completed";
export type PaymentMethod = "bkash" | "nagad" | "bank";

export interface Tour {
  tour_id: string;
  agency_id: string;
  name: string;
  start_date: string;
  end_date: string;
  total_seat: number;
  available_seat: number;
  description: string;
  last_enrollment_date: string;
  price: number;
  discount: number;
  status: TourStatus;
  image_path: string;
}

export interface Agency {
  agency_id: string;
  name: string;
  address: string;
  reg_id: string;
  rating: number;
  is_active: boolean;
  image_path?: string;
}

export interface HomeTour {
  tour_id: string;
  name: string;
  agency_id: string;
  agency_name: string;
  agency_rating: number;
  start_date: string;
  end_date: string;
  price: number;
  discount: number;
  final_price: number;
  available_seat: number;
  last_enrollment_date: string;
  description: string;
  status: TourStatus;
  total_bookings: number;
  image_path: string;
}

export interface HomeAgency {
  agency_id: string;
  name: string;
  address: string;
  rating: number;
  total_tours: number;
  image_path: string;
}

export interface HomeData {
  top_tours: HomeTour[];
  top_agencies: HomeAgency[];
}

export interface SearchTour {
  tour_id: string;
  agency_id: string;
  agency_name: string;
  name: string;
  start_date: string;
  end_date: string;
  available_seat: number;
  price: number;
  discount: number;
  status: TourStatus;
}

export interface SearchResult {
  Tours: SearchTour[];
  Agencies: Agency[];
}

export interface Page<T> {
  Data: T[];
  Meta: { Page: number; Limit: number; TotalCount: number; TotalPage: number };
}

export interface Booking {
  booking_id: string;
  status: BookingStatus;
  number_of_people: number;
  total_price: number;
  booking_date: string;
  created_by: "user" | "agency_member";
  cancellation_reason?: "customer" | "agency" | "tour_cancelled";
  customer_name: string;
  customer_email: string;
  customer_phone: string;
  tour_id: string;
  tour_name: string;
  tour_start_date: string;
  tour_status: TourStatus;
  agency_id: string;
  agency_name: string;
  payment_method: string;
  transaction_id: string;
  payment_status: "pending" | "success" | "failed" | "";
}

export interface Member {
  member_id: string;
  name: string;
  email: string;
  phone: string;
  permissions: number[];
}

export interface Permission {
  permission_id: number;
  name: string;
  resource: string;
  action: string;
}

export interface TourStatusChange {
  tour_id: string;
  status: TourStatus;
  cancelled_bookings: number;
}
