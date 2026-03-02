import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { DocumentService } from '../../../core/services/document';
import { AuthService } from '../../../core/services/auth';

interface CalendarDay {
  date: Date;
  isCurrentMonth: boolean;
  documents: any[];
}

@Component({
  selector: 'app-calendar-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  templateUrl: './calendar-dashboard.html',
  styleUrl: './calendar-dashboard.css',
})
export class CalendarDashboard implements OnInit {
  private documentService = inject(DocumentService);
  private authService = inject(AuthService);
  private router = inject(Router);
  private fb = inject(FormBuilder);
  private cdr = inject(ChangeDetectorRef);
  private sanitizer = inject(DomSanitizer);

  currentDate = new Date();
  calendarGrid: CalendarDay[] = [];
  weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  
  // API Data mapping: date string (YYYY-MM-DD) -> array of documents
  reportsMap: { [key: string]: any[] } = {};

  isLoading = true;

  // Day Modal State (For "+X more" clicks)
  isDayModalOpen = false;
  selectedDayDate: Date | null = null;
  selectedDayDocuments: any[] = [];

  // Details Modal State
  isDetailsModalOpen = false;
  isDetailsLoading = false;
  selectedDocDetails: any = null;
  rawFileUrl = '';
  selectedFileUrl: SafeResourceUrl | null = null;

  // Upload Modal State
  isUploadModalOpen = false;
  isUploading = false;
  selectedFile: File | null = null;
  uploadError = '';
  
  uploadForm: FormGroup = this.fb.group({
    category: ['general', Validators.required],
    report_type: ['', Validators.required],
    file_type: ['lab_report', Validators.required],
    tags: [''], // user will input comma separated values
    report_date: ['', Validators.required]
  });

  ngOnInit() {
    // Default to today for new reports
    this.uploadForm.patchValue({
      report_date: this.formatDateForApi(new Date())
    });
    this.generateCalendar();
    this.fetchMonthData();
  }

  get currentMonthName(): string {
    return this.currentDate.toLocaleString('default', { month: 'long', year: 'numeric' });
  }

  get isNextMonthDisabled(): boolean {
    const today = new Date();
    return this.currentDate.getMonth() === today.getMonth() && 
           this.currentDate.getFullYear() === today.getFullYear();
  }

  get hasReportsForCurrentMonth(): boolean {
    return Object.keys(this.reportsMap).length > 0;
  }

  previousMonth() {
    this.currentDate = new Date(this.currentDate.getFullYear(), this.currentDate.getMonth() - 1, 1);
    this.reportsMap = {};
    this.generateCalendar();
    this.fetchMonthData();
  }

  nextMonth() {
    this.currentDate = new Date(this.currentDate.getFullYear(), this.currentDate.getMonth() + 1, 1);
    this.reportsMap = {};
    this.generateCalendar();
    this.fetchMonthData();
  }

  fetchMonthData() {
    this.isLoading = true;
    const month = this.currentDate.getMonth() + 1; // 1-12
    const year = this.currentDate.getFullYear();

    this.documentService.getMonthlyReports(month, year).subscribe({
      next: (response) => {
        try {
          this.reportsMap = {};
          let daysData: any = null;

          if (response && response.days) {
            daysData = response.days;
          } else if (response && response.data && response.data.days) {
            daysData = response.data.days;
          } else if (typeof response === 'string') {
            try {
              const parsed = JSON.parse(response);
              daysData = parsed.days || (parsed.data ? parsed.data.days : null);
            } catch (e) {
              console.warn('Could not parse response string', e);
            }
          }

          if (daysData) {
            Object.keys(daysData).forEach(dateStr => {
              const dayItem = daysData[dateStr];
              if (dayItem && Array.isArray(dayItem.documents)) {
                this.reportsMap[dateStr] = dayItem.documents;
              } else if (Array.isArray(dayItem)) {
                this.reportsMap[dateStr] = dayItem;
              } else {
                this.reportsMap[dateStr] = [];
              }
            });
          }
        } catch (e) {
          console.error('Error mapping calendar data:', e);
        } finally {
          this.isLoading = false;
          this.generateCalendar();
          this.cdr.detectChanges();
        }
      },
      error: (err) => {
        console.error('Error fetching calendar', err);
        this.isLoading = false;
        this.generateCalendar();
        this.cdr.detectChanges();
      }
    });
  }

  generateCalendar() {
    this.calendarGrid = [];
    const year = this.currentDate.getFullYear();
    const month = this.currentDate.getMonth();

    const firstDayIndex = new Date(year, month, 1).getDay();
    const lastDay = new Date(year, month + 1, 0).getDate();
    
    // Previous month filling
    const prevMonthLastDay = new Date(year, month, 0).getDate();
    for (let i = firstDayIndex; i > 0; i--) {
      const d = new Date(year, month - 1, prevMonthLastDay - i + 1);
      this.calendarGrid.push({
        date: d,
        isCurrentMonth: false,
        documents: []
      });
    }

    // Current month filing
    for (let i = 1; i <= lastDay; i++) {
      const d = new Date(year, month, i);
      const dateStr = this.formatDateForApi(d);
      this.calendarGrid.push({
        date: d,
        isCurrentMonth: true,
        documents: this.reportsMap[dateStr] || []
      });
    }

    // Next month filling
    const remainingSlots = 42 - this.calendarGrid.length; // 6 rows of 7
    for (let i = 1; i <= remainingSlots; i++) {
        const d = new Date(year, month + 1, i);
        this.calendarGrid.push({
          date: d,
          isCurrentMonth: false,
          documents: []
        });
    }
  }

  formatDateForApi(date: Date): string {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, '0');
    const d = String(date.getDate()).padStart(2, '0');
    return `${y}-${m}-${d}`;
  }

  logout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  // --- Upload Modal Methods ---

  openUploadModal() {
    this.isUploadModalOpen = true;
    document.body.style.overflow = 'hidden';
    this.uploadError = '';
    this.selectedFile = null;
    this.uploadForm.reset({
      category: 'general',
      file_type: 'lab_report',
      report_date: this.formatDateForApi(new Date())
    });
  }

  closeUploadModal() {
    this.isUploadModalOpen = false;
    document.body.style.overflow = '';
  }

  onFileSelected(event: any) {
    if (event.target.files && event.target.files.length > 0) {
      this.selectedFile = event.target.files[0];
    }
  }

  submitDocument() {
    if (this.uploadForm.invalid) {
      this.uploadForm.markAllAsTouched();
      return;
    }
    if (!this.selectedFile) {
      this.uploadError = 'Please select a file to upload.';
      return;
    }

    this.isUploading = true;
    this.uploadError = '';

    // Step 1: Upload File
    const fileType = this.uploadForm.get('file_type')?.value || 'lab_report';
    this.documentService.uploadFile(this.selectedFile, fileType).subscribe({
      next: (uploadRes) => {
        if (uploadRes && uploadRes.files && uploadRes.files.length > 0) {
          const fileId = uploadRes.files[0].file_id;
          
          // Step 2: Post Details
          const rawTags = this.uploadForm.get('tags')?.value || '';
          const tagsArray = rawTags.split(',').map((t: string) => t.trim()).filter((t: string) => t);

          const payload = {
            ...this.uploadForm.value,
            tags: tagsArray,
            file_id: fileId
          };

          this.documentService.submitDocument(payload).subscribe({
            next: () => {
              this.isUploading = false;
              this.closeUploadModal();
              this.fetchMonthData(); // Refresh calendar to show new document
            },
            error: (err) => {
              console.error('Submit Doc Error', err);
              this.uploadError = 'Failed to save document details.';
              this.isUploading = false;
            }
          });

        } else {
            this.uploadError = 'File upload successful but no ID returned.';
            this.isUploading = false;
        }
      },
      error: (err) => {
        console.error('File Upload Error', err);
        this.uploadError = 'Failed to upload the file.';
        this.isUploading = false;
      }
    });
  }

  // --- Day Modal Methods ---
  openDayModal(date: Date, documents: any[]) {
    this.selectedDayDate = date;
    this.selectedDayDocuments = documents;
    this.isDayModalOpen = true;
    document.body.style.overflow = 'hidden';
  }

  closeDayModal() {
    this.isDayModalOpen = false;
    document.body.style.overflow = '';
    this.selectedDayDate = null;
    this.selectedDayDocuments = [];
  }

  // --- Delete Modal Methods ---
  isDeleteModalOpen = false;
  documentToDelete: any = null;
  isDeleting = false;
  deleteError = '';

  openDeleteModal(doc: any) {
    this.documentToDelete = doc;
    this.isDeleteModalOpen = true;
    this.deleteError = '';
    document.body.style.overflow = 'hidden';
  }

  closeDeleteModal() {
    this.isDeleteModalOpen = false;
    this.documentToDelete = null;
    this.isDeleting = false;
    this.deleteError = '';
    
    // Only clear overflow if no other modals are stubbornly open
    if (!this.isUploadModalOpen && !this.isDetailsModalOpen && !this.isDayModalOpen) {
        document.body.style.overflow = '';
    }
  }

  confirmDelete() {
    if (!this.documentToDelete || !this.documentToDelete.id) return;
    
    this.isDeleting = true;
    this.deleteError = '';

    this.documentService.deleteDocument(this.documentToDelete.id).subscribe({
      next: () => {
        this.isDeleting = false;
        this.closeDeleteModal();
        // Since the data changed, force a clean refresh of the month
        this.reportsMap = {};
        this.generateCalendar();
        this.fetchMonthData();
        
        // If the day modal was open, securely close it to prevent orphaned data
        if (this.isDayModalOpen) {
            this.closeDayModal();
        }
      },
      error: (err) => {
        console.error('Failed to delete document', err);
        this.deleteError = 'Failed to delete the document. Please try again.';
        this.isDeleting = false;
        this.cdr.detectChanges();
      }
    });
  }

  detailsError: string = '';
  isFullscreenImage: boolean = false;
  pdfBlobUrl: string | null = null;

  // --- Details Modal Methods ---
  openDocumentDetails(docId: string) {
    this.isDetailsModalOpen = true;
    document.body.style.overflow = 'hidden';
    this.isDetailsLoading = true;
    this.detailsError = '';
    this.selectedDocDetails = null;
    this.selectedFileUrl = null;
    this.rawFileUrl = '';
    this.isFullscreenImage = false;

    if (!docId) {
      this.detailsError = 'Invalid document ID.';
      this.isDetailsLoading = false;
      return;
    }

    this.documentService.getDocumentDetails(docId).subscribe({
      next: (response) => {
        try {
          let docData = response;
          // Unwind potential wrappers
          if (response && response.data) docData = response.data;
          else if (typeof response === 'string') {
            try { 
              const parsed = JSON.parse(response); 
              docData = parsed.data || parsed; 
            } catch(e) { console.warn('JSON string parse fail', e); }
          }

          this.selectedDocDetails = {
            ...docData,
            parsedTags: this.parseTags(docData.tags)
          };

          if (docData && docData.file_id) {
            this.documentService.getFileUrl(docData.file_id).subscribe({
              next: (fileRes) => {
                try {
                  let fileData = fileRes;
                  if (fileRes && fileRes.data) fileData = fileRes.data;
                  else if (typeof fileRes === 'string') {
                    try { 
                      const parsed = JSON.parse(fileRes); 
                      fileData = parsed.data || parsed; 
                    } catch(e) {}
                  }

                  if (fileData && fileData.url) {
                    this.rawFileUrl = fileData.url;
                    
                    if (this.isPdfFile()) {
                      // Fetch the PDF directly and construct a Blob URL so the browser renders it inline
                      // and bypasses S3's rigid Content-Disposition: attachment headers
                      fetch(this.rawFileUrl)
                        .then(res => res.blob())
                        .then(blob => {
                            const pdfBlob = new Blob([blob], { type: 'application/pdf' });
                            this.pdfBlobUrl = URL.createObjectURL(pdfBlob);
                            this.selectedFileUrl = this.sanitizer.bypassSecurityTrustResourceUrl(this.pdfBlobUrl);
                            this.isDetailsLoading = false;
                            this.cdr.detectChanges();
                        })
                        .catch(e => {
                            console.warn('CORS prevented inline PDF blob. Falling back to explicit Google Viewer Proxy.', e);
                            const googleProxyUrl = `https://docs.google.com/viewer?url=${encodeURIComponent(this.rawFileUrl)}&embedded=true`;
                            this.selectedFileUrl = this.sanitizer.bypassSecurityTrustResourceUrl(googleProxyUrl);
                            this.isDetailsLoading = false;
                            this.cdr.detectChanges();
                        });
                    } else {
                      this.selectedFileUrl = this.sanitizer.bypassSecurityTrustResourceUrl(fileData.url);
                      this.isDetailsLoading = false;
                      this.cdr.detectChanges();
                    }
                  } else {
                    this.detailsError = 'No file URL returned from server.';
                    this.isDetailsLoading = false;
                    this.cdr.detectChanges();
                  }
                } catch(e) {
                  this.detailsError = 'Failed to map file URL data.';
                } finally {
                  this.isDetailsLoading = false;
                  this.cdr.detectChanges();
                }
              },
              error: (err) => {
                console.error('File URL fetch error', err);
                this.detailsError = 'Failed to load document file link.';
                this.isDetailsLoading = false;
                this.cdr.detectChanges();
              }
            });
          } else {
            // Missing file_id but doc exists
            this.isDetailsLoading = false;
            this.cdr.detectChanges();
          }
        } catch(e) {
          console.error('Doc payload mapping error', e);
          this.detailsError = 'Unexpected error rendering document.';
          this.isDetailsLoading = false;
          this.cdr.detectChanges();
        }
      },
      error: (err) => {
        console.error('Doc details fetch error', err);
        this.detailsError = 'Failed to load document information.';
        this.isDetailsLoading = false;
        this.cdr.detectChanges();
      }
    });
  }

  closeDocumentDetails() {
    this.isDetailsModalOpen = false;
    document.body.style.overflow = '';
    this.selectedDocDetails = null;
    this.selectedFileUrl = null;
    this.rawFileUrl = '';
    this.isFullscreenImage = false;
    
    // Cleanup memory from our temporary blob URLs
    if (this.pdfBlobUrl) {
      URL.revokeObjectURL(this.pdfBlobUrl);
      this.pdfBlobUrl = null;
    }
  }

  isImageFile(): boolean {
    if (!this.rawFileUrl) return false;
    const lowerUrl = this.rawFileUrl.toLowerCase();
    const urlWithoutParams = lowerUrl.split('?')[0];
    return urlWithoutParams.endsWith('.png') || urlWithoutParams.endsWith('.jpg') || urlWithoutParams.endsWith('.jpeg');
  }

  isPdfFile(): boolean {
    if (!this.rawFileUrl) return false;
    const lowerUrl = this.rawFileUrl.toLowerCase();
    const urlWithoutParams = lowerUrl.split('?')[0];
    return urlWithoutParams.endsWith('.pdf');
  }

  parseTags(tags: string | string[]): string[] {
    if (!tags) return [];
    if (Array.isArray(tags)) return tags;
    try {
      const parsed = JSON.parse(tags);
      if (Array.isArray(parsed)) return parsed;
      return [];
    } catch (e) {
      return [];
    }
  }
}
