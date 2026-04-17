import { Component, inject, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, ActivatedRoute, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth';
import { ToastService } from '../../../core/services/toast';
import { environment } from '../../../../environments/environment';
import { ThemeToggleComponent } from '../../../core/components/theme-toggle/theme-toggle';

@Component({
  selector: 'app-verify-otp',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink, ThemeToggleComponent],
  templateUrl: './verify-otp.html',
  styleUrl: './verify-otp.css',
})
export class VerifyOtp implements OnInit, OnDestroy {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private toastService = inject(ToastService);
  private cdr = inject(ChangeDetectorRef);

  otpForm: FormGroup = this.fb.nonNullable.group({
    otp: ['', [Validators.required, Validators.minLength(6), Validators.maxLength(6), Validators.pattern(/^\d{6}$/)]]
  });

  email = '';
  password = '';
  isLoading = false;
  isResending = false;
  errorMessage = '';
  isVerified = false;

  // Retry timer
  retryCountdown = 0;           // seconds remaining
  private timerInterval: ReturnType<typeof setInterval> | null = null;

  ngOnInit(): void {
    if (!environment.emailVerificationEnabled) {
      this.router.navigate(['/login']);
      return;
    }
    this.email = this.route.snapshot.queryParamMap.get('email') ?? '';
    this.password = history.state?.password || '';
    // Start the initial 60-second cooldown immediately (OTP was just sent on signup)
    this.startRetryTimer();
  }

  ngOnDestroy(): void {
    this.clearTimer();
  }

  private startRetryTimer(): void {
    this.retryCountdown = 60;
    this.clearTimer();
    this.timerInterval = setInterval(() => {
      this.retryCountdown--;
      this.cdr.detectChanges();
      if (this.retryCountdown <= 0) {
        this.clearTimer();
      }
    }, 1000);
  }

  private clearTimer(): void {
    if (this.timerInterval) {
      clearInterval(this.timerInterval);
      this.timerInterval = null;
    }
  }

  get canResend(): boolean {
    return this.retryCountdown === 0 && !this.isResending;
  }

  onSubmit(): void {
    if (this.otpForm.invalid) {
      this.otpForm.markAllAsTouched();
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    const { otp } = this.otpForm.getRawValue();

    this.authService.verifyOTP(this.email, otp).subscribe({
      next: () => {
        this.isLoading = false;
        this.isVerified = true;
        this.toastService.showSuccess('Email verified successfully! You can now log in.');
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.isLoading = false;
        const errMsg = err.error?.message || 'Invalid or expired OTP. Please try again.';
        this.errorMessage = errMsg;
        this.toastService.showError(errMsg);
        this.cdr.detectChanges();
      }
    });
  }

  resendOTP(): void {
    if (!this.canResend) return;

    this.isResending = true;
    this.errorMessage = '';

    this.authService.resendOTP(this.email).subscribe({
      next: () => {
        this.isResending = false;
        this.toastService.showSuccess('A new OTP has been sent to your email.');
        this.startRetryTimer();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.isResending = false;
        const errMsg = err.error?.message || 'Failed to resend OTP. Please try again.';
        this.errorMessage = errMsg;
        this.toastService.showError(errMsg);
        this.cdr.detectChanges();
      }
    });
  }

  goToLogin(): void {
    this.router.navigate(['/login'], { 
      state: { email: this.email, password: this.password } 
    });
  }
}
