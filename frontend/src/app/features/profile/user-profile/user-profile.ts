import { CommonModule } from '@angular/common';
import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';
import { AuthService, AuthUser, UpdateProfilePayload } from '../../../core/services/auth';

@Component({
  selector: 'app-user-profile',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './user-profile.html',
  styleUrl: './user-profile.css',
})
export class UserProfile implements OnInit {
  private authService = inject(AuthService);
  private router = inject(Router);
  private fb = inject(FormBuilder);

  protected readonly user = signal<AuthUser | null>(null);
  protected readonly isEditing = signal(false);
  protected readonly isSubmitting = signal(false);
  protected readonly selectedFileName = signal<string | null>(null);
  protected readonly previewImage = signal<string | null>(null);
  protected readonly removeProfileImage = signal(false);
  protected readonly updateError = signal<string | null>(null);
  protected readonly updateSuccess = signal<string | null>(null);

  protected readonly profileForm = this.fb.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    gender: [''],
    dob: [''],
  });

  protected readonly initials = computed(() => {
    const name = this.user()?.name?.trim();
    if (!name) {
      return 'VT';
    }

    return name
      .split(/\s+/)
      .slice(0, 2)
      .map((part) => part.charAt(0).toUpperCase())
      .join('');
  });

  ngOnInit(): void {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) {
      this.router.navigate(['/login']);
      return;
    }

    this.syncUserState(currentUser);
  }

  protected get joinedDate(): string {
    const createdAt = this.user()?.created_at;
    if (!createdAt) {
      return 'Not available';
    }

    return new Date(createdAt).toLocaleDateString('en-US', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });
  }

  protected get updatedDate(): string {
    const updatedAt = this.user()?.updated_at;
    if (!updatedAt) {
      return 'Not available';
    }

    return new Date(updatedAt).toLocaleString('en-US', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    });
  }

  protected get displayDob(): string {
    const dob = this.user()?.dob;
    if (!dob) {
      return 'Not added';
    }

    return new Date(dob).toLocaleDateString('en-US', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });
  }

  protected get displayGender(): string {
    const gender = this.user()?.gender?.trim();
    return gender ? gender : 'Not added';
  }

  protected get profileImage(): string | null {
    return this.previewImage() || this.user()?.profile_pic || null;
  }

  protected onProfilePicSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0] ?? null;

    this.updateError.set(null);
    this.updateSuccess.set(null);

    if (!file) {
      this.selectedFileName.set(null);
      this.previewImage.set(null);
      return;
    }

    this.selectedFileName.set(file.name);
    this.removeProfileImage.set(false);

    const reader = new FileReader();
    reader.onload = () => {
      this.previewImage.set(typeof reader.result === 'string' ? reader.result : null);
    };
    reader.readAsDataURL(file);
  }

  protected clearSelectedImage(fileInput: HTMLInputElement): void {
    fileInput.value = '';
    this.selectedFileName.set(null);
    this.previewImage.set(null);
  }

  protected markProfileImageForRemoval(fileInput: HTMLInputElement): void {
    this.removeProfileImage.set(true);
    this.clearSelectedImage(fileInput);
    this.updateError.set(null);
    this.updateSuccess.set(null);
  }

  protected saveProfile(fileInput: HTMLInputElement): void {
    if (this.profileForm.invalid) {
      this.profileForm.markAllAsTouched();
      return;
    }

    const currentUser = this.user();
    if (!currentUser) {
      return;
    }

    const formValue = this.profileForm.getRawValue();
    const trimmedName = (formValue.name ?? '').trim();
    const trimmedGender = (formValue.gender ?? '').trim();
    const dob = formValue.dob ?? '';
    const selectedFile = fileInput.files?.[0] ?? null;

    const payload: UpdateProfilePayload = {};

    if (trimmedName !== (currentUser.name ?? '')) {
      payload.name = trimmedName;
    }

    if (trimmedGender !== (currentUser.gender ?? '')) {
      payload.gender = trimmedGender;
    }

    if (dob !== this.toDateInputValue(currentUser.dob ?? null)) {
      payload.dob = dob;
    }

    if (this.removeProfileImage()) {
      payload.delete_profile_pic = true;
    }

    if (selectedFile) {
      payload.profile_pic = selectedFile;
    }

    if (Object.keys(payload).length === 0) {
      this.updateSuccess.set('No profile changes to save.');
      this.updateError.set(null);
      return;
    }

    this.isSubmitting.set(true);
    this.updateError.set(null);
    this.updateSuccess.set(null);

    this.authService.updateProfile(payload).pipe(
      finalize(() => this.isSubmitting.set(false))
    ).subscribe({
      next: (response) => {
        if (response?.user) {
          this.syncUserState(response.user as AuthUser);
        }

        if (this.removeProfileImage()) {
          fileInput.value = '';
        }

        this.selectedFileName.set(null);
        this.previewImage.set(null);
        this.removeProfileImage.set(false);
        this.updateSuccess.set(response?.message || 'Profile updated successfully.');
        this.isEditing.set(false);
      },
      error: (error) => {
        this.updateError.set(error?.error?.error || error?.error?.message || 'Unable to update profile.');
      }
    });
  }

  protected logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  protected openEditor(): void {
    this.updateError.set(null);
    this.updateSuccess.set(null);
    this.isEditing.set(true);
  }

  protected closeEditor(fileInput?: HTMLInputElement): void {
    const currentUser = this.user();
    if (currentUser) {
      this.syncUserState(currentUser);
    }

    if (fileInput) {
      fileInput.value = '';
    }

    this.updateError.set(null);
    this.updateSuccess.set(null);
    this.isEditing.set(false);
  }

  private syncUserState(user: AuthUser): void {
    this.user.set(user);
    this.profileForm.patchValue({
      name: user.name ?? '',
      gender: user.gender ?? '',
      dob: this.toDateInputValue(user.dob ?? null),
    });
    this.selectedFileName.set(null);
    this.previewImage.set(null);
    this.removeProfileImage.set(false);
  }

  private toDateInputValue(value: string | null): string {
    if (!value) {
      return '';
    }

    return value.includes('T') ? value.split('T')[0] : value;
  }
}
