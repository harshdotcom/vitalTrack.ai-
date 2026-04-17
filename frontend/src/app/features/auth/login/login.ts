import { Component, inject, ChangeDetectorRef, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth';
import { ToastService } from '../../../core/services/toast';
import { finalize } from 'rxjs';
import { environment } from '../../../../environments/environment';
import { ThemeToggleComponent } from '../../../core/components/theme-toggle/theme-toggle';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink, ThemeToggleComponent],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login implements OnInit {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private toastService = inject(ToastService);
  private cdr = inject(ChangeDetectorRef);

  loginForm: FormGroup = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]]
  });
  forgotPasswordForm: FormGroup = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    otp: ['', [Validators.required]],
    newPassword: ['', [Validators.required, Validators.minLength(6)]]
  });

  readonly emailVerificationEnabled = environment.emailVerificationEnabled;

  isLoading = false;
  isForgotPasswordOpen = false;
  isOtpRequested = false;
  isForgotPasswordLoading = false;
  forgotPasswordError = '';
  forgotPasswordInfo = '';
  showNewPassword = false;

  ngOnInit(): void {
    const state = history.state;
    if (state?.email) {
      this.loginForm.patchValue({ email: state.email });
      this.forgotPasswordForm.patchValue({ email: state.email });
    }
    if (state?.password) {
      this.loginForm.patchValue({ password: state.password });
    }
  }
  errorMessage = '';
  showPassword = false;

  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  toggleForgotPassword(): void {
    this.isForgotPasswordOpen = !this.isForgotPasswordOpen;
    this.forgotPasswordError = '';
    this.forgotPasswordInfo = '';

    const loginEmail = this.loginForm.controls['email'].value;
    if (loginEmail && !this.forgotPasswordForm.controls['email'].value) {
      this.forgotPasswordForm.patchValue({ email: loginEmail });
    }

    if (!this.isForgotPasswordOpen) {
      this.resetForgotPasswordState();
    }
  }

  requestPasswordResetOtp(): void {
    const emailControl = this.forgotPasswordForm.controls['email'];
    emailControl.markAsTouched();

    if (emailControl.invalid) {
      return;
    }

    this.isForgotPasswordLoading = true;
    this.forgotPasswordError = '';
    this.forgotPasswordInfo = '';

    this.authService.forgotPassword(emailControl.value).pipe(
      finalize(() => {
        this.isForgotPasswordLoading = false;
        this.cdr.detectChanges();
      })
    ).subscribe({
      next: (response) => {
        this.isOtpRequested = true;
        this.forgotPasswordInfo = response?.message || 'OTP sent to your email address.';
        this.loginForm.patchValue({ email: emailControl.value });
        this.toastService.showSuccess(this.forgotPasswordInfo);
      },
      error: (err) => {
        const errMsg = err.error?.message || 'Unable to send OTP. Please try again.';
        this.forgotPasswordError = errMsg;
        this.toastService.showError(errMsg);
      }
    });
  }

  submitPasswordReset(): void {
    this.forgotPasswordForm.markAllAsTouched();
    if (this.forgotPasswordForm.invalid) {
      return;
    }

    this.isForgotPasswordLoading = true;
    this.forgotPasswordError = '';
    this.forgotPasswordInfo = '';

    const raw = this.forgotPasswordForm.getRawValue();
    this.authService.resetPassword({
      email: raw.email,
      otp: raw.otp,
      new_password: raw.newPassword
    }).pipe(
      finalize(() => {
        this.isForgotPasswordLoading = false;
        this.cdr.detectChanges();
      })
    ).subscribe({
      next: (response) => {
        const successMessage = response?.message || 'Password updated successfully.';
        this.toastService.showSuccess(successMessage);
        this.forgotPasswordInfo = 'Password reset complete. You can log in with your new password.';
        this.loginForm.patchValue({
          email: raw.email,
          password: raw.newPassword
        });
        this.resetForgotPasswordState(raw.email);
        this.isForgotPasswordOpen = false;
        this.showNewPassword = false;
      },
      error: (err) => {
        const errMsg = err.error?.message || 'Unable to reset password. Please verify your OTP and try again.';
        this.forgotPasswordError = errMsg;
        this.toastService.showError(errMsg);
      }
    });
  }

  private resetForgotPasswordState(email: string = this.forgotPasswordForm.controls['email'].value): void {
    this.isOtpRequested = false;
    this.forgotPasswordForm.reset({
      email,
      otp: '',
      newPassword: ''
    });
    this.forgotPasswordForm.markAsPristine();
    this.forgotPasswordForm.markAsUntouched();
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    this.authService.login(this.loginForm.getRawValue()).subscribe({
      next: (response) => {
        // Handle success (e.g., store token, navigate)
        console.log('Login successful', response);
        this.router.navigate(['/dashboard']);
        this.isLoading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Login Error:', err);
        const errMsg = err.error?.message || 'Invalid email or password. Please try again.';
        this.errorMessage = errMsg;
        this.toastService.showError(errMsg);
        this.isLoading = false;
        this.cdr.detectChanges();
      }
    });
  }
}
