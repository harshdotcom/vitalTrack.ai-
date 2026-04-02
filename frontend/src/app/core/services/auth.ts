import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { API_CONSTANTS } from '../constants/api.constants';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private http = inject(HttpClient);
  private tokenKey = 'auth_token';

  login(credentials: any): Observable<any> {
    return this.http.post(API_CONSTANTS.LOGIN_URL, credentials).pipe(
      tap((res: any) => {
        if (res && res.token) {
          localStorage.setItem(this.tokenKey, res.token);
        }
      })
    );
  }

  signup(userData: {name: string, email: string, password: string}): Observable<any> {
    const formData = new FormData();
    formData.append('name', userData.name);
    formData.append('email', userData.email);
    formData.append('password', userData.password);
    
    return this.http.post(API_CONSTANTS.SIGNUP_URL, formData);
  }

  verifyOTP(email: string, otp: string): Observable<any> {
    return this.http.post(API_CONSTANTS.VERIFY_OTP_URL, { email, otp });
  }

  getUserUsage(): Observable<any> {
    return this.http.get(API_CONSTANTS.USER_USAGE_URL);
  }

  getAICredits(): Observable<any> {
    return this.http.get(API_CONSTANTS.AI_CREDITS_URL);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  logout(): void {
    localStorage.removeItem(this.tokenKey);
  }
}
