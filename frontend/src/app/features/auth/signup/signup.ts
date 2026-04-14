import { Component, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth';
import { ToastService } from '../../../core/services/toast';
import { environment } from '../../../../environments/environment';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './signup.html',
  styleUrl: './signup.css',
})
export class Signup {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private toastService = inject(ToastService);
  private cdr = inject(ChangeDetectorRef);

  // Form group for Name, Email, Password
  signupForm: FormGroup = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]]
  });

  isLoading = false;
  errorMessage = '';
  showPassword = false;

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

    const userData = this.signupForm.getRawValue();

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
}
