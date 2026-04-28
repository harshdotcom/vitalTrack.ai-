import { Component, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth';
import { ToastService } from '../../../core/services/toast';
import { finalize } from 'rxjs';
import { environment } from '../../../../environments/environment';
import { ThemeToggleComponent } from '../../../core/components/theme-toggle/theme-toggle';
import { GoogleSigninButtonComponent } from '../../../core/components/google-signin-button/google-signin-button';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink, ThemeToggleComponent, GoogleSigninButtonComponent],
  templateUrl: './signup.html',
  styleUrl: './signup.css',
})
export class Signup {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private toastService = inject(ToastService);
  private cdr = inject(ChangeDetectorRef);
  readonly googleClientId = environment.googleClientId;

  // Form group for Name, Email, Password
  signupForm: FormGroup = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    dob: [''],
    gender: ['']
  });

  isLoading = false;
  isGoogleLoading = false;
  errorMessage = '';
  showPassword = false;
  readonly genderOptions = [
    { value: 'Male', label: 'Male' },
    { value: 'Female', label: 'Female' },
    { value: 'Other', label: 'Other' },
    { value: 'Prefer not to say', label: 'Prefer not to say' },
  ];

  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  onSubmit() {
    if (this.signupForm.invalid) {
      this.signupForm.markAllAsTouched();
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    const rawValue = this.signupForm.getRawValue();
    const userData = {
      name: rawValue.name,
      email: rawValue.email,
      password: rawValue.password,
      dob: rawValue.dob || null,
      gender: rawValue.gender || null,
    };

    this.authService.signup(userData).subscribe({
      next: (_response) => {
        this.isLoading = false;
        this.cdr.detectChanges();
        if (environment.emailVerificationEnabled) {
          this.toastService.showSuccess('Account created! Please check your email for an OTP.');
          this.router.navigate(['/verify-otp'], {
            queryParams: { email: userData.email },
            state: { password: userData.password }
          });
        } else {
          this.toastService.showSuccess('Account created! You can now log in.');
          this.router.navigate(['/login'], {
            state: { email: userData.email, password: userData.password }
          });
        }
      },
      error: (err) => {
        console.error('Signup Error:', err);
        const errMsg = err.error?.message || 'Failed to create an account. Please try again later.';
        this.errorMessage = errMsg;
        this.toastService.showError(errMsg);
        this.isLoading = false;
        this.cdr.detectChanges();
      }
    });
  }

  onGoogleCredential(token: string): void {
    this.isGoogleLoading = true;
    this.errorMessage = '';

    this.authService.googleLogin({ token }).pipe(
      finalize(() => {
        this.isGoogleLoading = false;
        this.cdr.detectChanges();
      })
    ).subscribe({
      next: () => {
        this.toastService.showSuccess('Signed in with Google.');
        void this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        const errMsg = err.error?.message || 'Google sign-up failed. Please try again.';
        this.errorMessage = errMsg;
        this.toastService.showError(errMsg);
      }
    });
  }

  onGoogleInitializationFailed(message: string): void {
    this.errorMessage = message;
    this.cdr.detectChanges();
  }
}
