import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { API_CONSTANTS } from '../constants/api.constants';

export interface AuthUser {
  user_id: number;
  email: string;
  password?: string;
  google_id?: string | null;
  name?: string;
  age?: number | null;
  gender?: string;
  profile_pic?: string | null;
  dob?: string | null;
  is_verified?: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface UpdateProfilePayload {
  name?: string;
  gender?: string;
  dob?: string | null;
  delete_profile_pic?: boolean;
  profile_pic?: File | null;
}

export interface AuthResponse {
  message: string;
  token: string;
  user: AuthUser;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private http = inject(HttpClient);
  private tokenKey = 'auth_token';
  private userKey = 'auth_user';

  login(credentials: { email: string; password: string }): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(API_CONSTANTS.LOGIN_URL, credentials).pipe(
      tap((res) => this.persistAuthSession(res))
    );
  }

  googleLogin(payload: { token: string }): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(API_CONSTANTS.GOOGLE_LOGIN_URL, payload).pipe(
      tap((res) => this.persistAuthSession(res))
    );
  }

  signup(userData: { name: string; email: string; password: string; dob?: string | null; gender?: string | null }): Observable<any> {
    const formData = new FormData();
    formData.append('name', userData.name);
    formData.append('email', userData.email);
    formData.append('password', userData.password);

    if (userData.dob) {
      formData.append('dob', userData.dob);
    }

    if (userData.gender) {
      formData.append('gender', userData.gender);
    }
    
    return this.http.post(API_CONSTANTS.SIGNUP_URL, formData);
  }

  verifyOTP(email: string, otp: string): Observable<any> {
    return this.http.post(API_CONSTANTS.VERIFY_OTP_URL, { email, otp });
  }

  resendOTP(email: string): Observable<any> {
    return this.http.post(API_CONSTANTS.RESEND_OTP_URL, { email });
  }

  forgotPassword(email: string): Observable<any> {
    return this.http.post(API_CONSTANTS.FORGOT_PASSWORD_URL, { email });
  }

  resetPassword(payload: { email: string; new_password: string; otp: string }): Observable<any> {
    return this.http.post(API_CONSTANTS.RESET_PASSWORD_URL, payload);
  }

  getUserUsage(): Observable<any> {
    return this.http.get(API_CONSTANTS.USER_USAGE_URL);
  }

  getAICredits(): Observable<any> {
    return this.http.get(API_CONSTANTS.AI_CREDITS_URL);
  }

  updateProfile(payload: UpdateProfilePayload): Observable<any> {
    const hasFile = payload.profile_pic instanceof File;

    if (hasFile) {
      const formData = new FormData();

      if (payload.name !== undefined) {
        formData.append('name', payload.name);
      }
      if (payload.gender !== undefined) {
        formData.append('gender', payload.gender);
      }
      if (payload.dob !== undefined) {
        formData.append('dob', payload.dob ?? '');
      }
      if (payload.delete_profile_pic !== undefined) {
        formData.append('delete_profile_pic', String(payload.delete_profile_pic));
      }
      if (payload.profile_pic) {
        formData.append('profile_pic', payload.profile_pic);
      }

      return this.http.patch(API_CONSTANTS.UPDATE_PROFILE_URL, formData).pipe(
        tap((res: any) => this.persistUserFromResponse(res))
      );
    }

    const body: Record<string, unknown> = {};
    if (payload.name !== undefined) {
      body['name'] = payload.name;
    }
    if (payload.gender !== undefined) {
      body['gender'] = payload.gender;
    }
    if (payload.dob !== undefined) {
      body['dob'] = payload.dob ?? '';
    }
    if (payload.delete_profile_pic !== undefined) {
      body['delete_profile_pic'] = payload.delete_profile_pic;
    }

    return this.http.patch(API_CONSTANTS.UPDATE_PROFILE_URL, body).pipe(
      tap((res: any) => this.persistUserFromResponse(res))
    );
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  getCurrentUser(): AuthUser | null {
    const storedUser = localStorage.getItem(this.userKey);
    if (!storedUser) {
      return null;
    }

    try {
      return JSON.parse(storedUser) as AuthUser;
    } catch {
      return null;
    }
  }

  logout(): void {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.userKey);
  }

  private persistUserFromResponse(res: any): void {
    if (res?.user) {
      const { password, ...safeUser } = res.user;
      localStorage.setItem(this.userKey, JSON.stringify(safeUser));
    }
  }

  private persistAuthSession(res: AuthResponse): void {
    if (res?.token) {
      localStorage.setItem(this.tokenKey, res.token);
    }

    this.persistUserFromResponse(res);
  }
}
