import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { DocumentService } from '../../../core/services/document';
import { ToastService } from '../../../core/services/toast';
import { AuthService } from '../../../core/services/auth';
import { ThemeToggleComponent } from '../../../core/components/theme-toggle/theme-toggle';

interface CalendarEntry {
  id: string;
  entry_type: 'document' | 'direct_entry';
  category: string;
  document_name: string;
  status: string;
  document_date: string;
  analysis_generated: boolean;
  metric_type?: string;
  metric_label?: string;
  metric_summary?: string;
  timestamp?: string;
  tags?: string | string[];
}

interface CalendarDay {
  date: Date;
  isCurrentMonth: boolean;
  documents: CalendarEntry[];
}

@Component({
  selector: 'app-calendar-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, ThemeToggleComponent],
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
  private toastService = inject(ToastService);

  currentDate = new Date();
  calendarGrid: CalendarDay[] = [];
  weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  
  // API Data mapping: date string (YYYY-MM-DD) -> array of documents
  reportsMap: { [key: string]: CalendarEntry[] } = {};

  isLoading = true;

  // Day Modal State (For "+X more" clicks)
  isDayModalOpen = false;
  selectedDayDate: Date | null = null;
  selectedDayDocuments: CalendarEntry[] = [];

  // Details Modal State
  isDetailsModalOpen = false;
  isDetailsLoading = false;
  selectedDocDetails: any = null;
  rawFileUrl = '';
  selectedFileUrl: SafeResourceUrl | null = null;

  isEditingDetails = false;
  isUpdatingDocument = false;
  updateDocError = '';

  editForm: FormGroup = this.fb.group({
    document_name: ['', Validators.required],
    category: ['Medical Report', Validators.required],
    document_date: ['', Validators.required],
    tags: ['']
  });

  // Entry Origin Modal State
  isEntryOriginModalOpen = false;

  // Vitals Modal State
  isVitalsModalOpen = false;
  isSavingVitals = false;
  vitalsError = '';
  
  vitalsForm: FormGroup = this.fb.group({
    vital_type: ['', Validators.required],
    date: ['', Validators.required],
    time: ['', Validators.required],
    bp_systolic: [null],
    bp_diastolic: [null],
    unit: ['mg/dL'],
    value: [null],
    notes: ['']
  });

  // Upload Modal State
  isUploadModalOpen = false;
  isUploading = false;
  selectedFile: File | null = null;
  uploadError = '';
  
  uploadForm: FormGroup = this.fb.group({
    category: ['general', Validators.required],
    document_name: ['', Validators.required],
    file_type: ['lab_report', Validators.required],
    tags: [''], // user will input comma separated values
    document_date: ['', Validators.required]
  });

  // Storage State
  totalStorageUsedBytes = 0;
  MAX_STORAGE_BYTES = 100 * 1024 * 1024; // 100 MB
  
  get formattedStorage(): string {
    const mb = this.totalStorageUsedBytes / (1024 * 1024);
    return mb.toFixed(1) + ' MB';
  }

  get storagePercentage(): number {
    const pct = (this.totalStorageUsedBytes / this.MAX_STORAGE_BYTES) * 100;
    return Math.min(pct, 100);
  }

  get storageExceeded(): boolean {
    return this.totalStorageUsedBytes >= this.MAX_STORAGE_BYTES;
  }

  // AI Credit State
  aiLeftCredit = 0;
  aiTotalCredit = 0;
  aiRenewDate: string = '';

  get aiUsedCredit(): number {
    return this.aiTotalCredit - this.aiLeftCredit;
  }

  get aiCreditPercentage(): number {
    if (this.aiTotalCredit === 0) return 0;
    return Math.min((this.aiUsedCredit / this.aiTotalCredit) * 100, 100);
  }

  get aiCreditsExceeded(): boolean {
    return this.aiLeftCredit === 0;
  }

  get formattedRenewDate(): string {
    if (!this.aiRenewDate) return '';
    return new Date(this.aiRenewDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  }

  ngOnInit() {
    this.uploadForm.patchValue({
      document_date: this.formatDateForApi(new Date())
    });
    this.generateCalendar();
    this.fetchMonthData();
    this.fetchStorageUsage();
    this.fetchAICredits();
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
                this.reportsMap[dateStr] = dayItem.documents.map((entry: any) => this.normalizeCalendarEntry(entry));
              } else if (Array.isArray(dayItem)) {
                this.reportsMap[dateStr] = dayItem.map((entry: any) => this.normalizeCalendarEntry(entry));
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
        this.toastService.showError('Failed to load calendar data.');
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

  formatTimeForInput(date: Date): string {
    const hh = String(date.getHours()).padStart(2, '0');
    const mm = String(date.getMinutes()).padStart(2, '0');
    return `${hh}:${mm}`;
  }

  buildLocalTimestamp(date: string, time: string): string {
    if (!date) {
      return new Date().toISOString();
    }

    const timeValue = time || '00:00';
    const localDate = new Date(`${date}T${timeValue}:00`);
    const offsetMinutes = -localDate.getTimezoneOffset();
    const sign = offsetMinutes >= 0 ? '+' : '-';
    const absOffset = Math.abs(offsetMinutes);
    const offsetHours = String(Math.floor(absOffset / 60)).padStart(2, '0');
    const offsetMins = String(absOffset % 60).padStart(2, '0');

    return `${date}T${timeValue}:00${sign}${offsetHours}:${offsetMins}`;
  }

  logout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  goToProfile() {
    this.router.navigate(['/profile']);
  }

  fetchStorageUsage() {
    this.authService.getUserUsage().subscribe({
      next: (res) => {
        if (res) {
          const usedBytes = res?.usage?.TotalStorageUsed ?? res?.data?.TotalStorageUsed;
          if (usedBytes !== undefined) {
            this.totalStorageUsedBytes = usedBytes;
            this.cdr.detectChanges();
          }
        }
      },
      error: (err) => console.error('Failed to fetch storage usage', err)
    });
  }

  fetchAICredits() {
    this.authService.getAICredits().subscribe({
      next: (res) => {
        if (res?.usage) {
          this.aiLeftCredit = res.usage.leftCredit ?? 0;
          this.aiTotalCredit = res.usage.totalCredit ?? 0;
          this.aiRenewDate = res.usage.renewDate ?? '';
          this.cdr.detectChanges();
        }
      },
      error: (err) => console.error('Failed to fetch AI credits', err)
    });
  }

  openEntryOriginModal() {
    this.isEntryOriginModalOpen = true;
    document.body.style.overflow = 'hidden';
  }

  closeEntryOriginModal() {
    this.isEntryOriginModalOpen = false;
    document.body.style.overflow = '';
  }

  selectEntryOrigin(type: 'vitals' | 'report') {
    this.closeEntryOriginModal();
    if (type === 'report') {
      this.openUploadModal();
    } else {
      this.openVitalsModal();
    }
  }

  // --- Vitals Modal Methods ---
  openVitalsModal() {
    this.isVitalsModalOpen = true;
    document.body.style.overflow = 'hidden';
    this.vitalsError = '';
    this.vitalsForm.reset({
      date: this.formatDateForApi(new Date()),
      time: this.formatTimeForInput(new Date()),
      unit: 'mg/dL'
    });
  }

  closeVitalsModal() {
    this.isVitalsModalOpen = false;
    document.body.style.overflow = '';
  }

  onVitalTypeChange() {
    const type = this.vitalsForm.get('vital_type')?.value;
    if (type === 'blood_sugar') {
      this.vitalsForm.patchValue({ unit: 'mg/dL' });
    } else if (type === 'weight') {
      this.vitalsForm.patchValue({ unit: 'kg' });
    }
  }

  submitVitals() {
    if (this.vitalsForm.invalid) {
      this.vitalsForm.markAllAsTouched();
      return;
    }
    
    const formVals = this.vitalsForm.value;
    const type = formVals.vital_type;
    
    let payload: any = {};
    
    if (type === 'blood_pressure') {
      payload['blood_pressure'] = {
        systolic: Number(formVals.bp_systolic),
        diastolic: Number(formVals.bp_diastolic)
      };
    } else if (type === 'blood_sugar') {
      payload['blood_sugar'] = { unit: formVals.unit, value: Number(formVals.value) };
    } else if (['calories', 'heart_rate', 'oxygen_level', 'sleep_hours', 'steps'].includes(type)) {
      payload[type] = Number(formVals.value);
    } else if (type === 'weight') {
      payload['weight'] = { unit: formVals.unit, value: Number(formVals.value) };
    } else if (type === 'notes') {
      payload['notes'] = formVals.notes;
    }

    payload['timestamp'] = this.buildLocalTimestamp(formVals.date, formVals.time);

    this.isSavingVitals = true;
    this.vitalsError = '';
    
    this.documentService.saveHealthMetric(payload).subscribe({
      next: () => {
        this.isSavingVitals = false;
        this.closeVitalsModal();
        this.toastService.showSuccess('Vital entry saved successfully');
        this.fetchMonthData();
      },
      error: (err) => {
        this.isSavingVitals = false;
        const msg = err.error?.message || 'Failed to save vitals.';
        this.vitalsError = msg;
        this.toastService.showError(msg);
      }
    });
  }

  // --- Upload Modal Methods ---

  openUploadModal() {
    if (this.storageExceeded) {
      this.toastService.showError('Storage limit reached (100 MB). Please email harshjha92002@gmail.com to upgrade your plan.');
      return;
    }
    this.isUploadModalOpen = true;
    document.body.style.overflow = 'hidden';
    this.uploadError = '';
    this.selectedFile = null;
    this.uploadForm.reset({
      category: 'general',
      file_type: 'lab_report',
      document_date: this.formatDateForApi(new Date())
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
        // Automatically fetch latest usage since the file has hit the backend
        this.fetchStorageUsage();

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
              this.fetchStorageUsage(); // Refresh storage limit
            },
            error: (err) => {
              console.error('Submit Doc Error', err);
              const errMsg = err.error?.message || 'Failed to save document details.';
              this.uploadError = errMsg;
              this.toastService.showError(errMsg);
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
        const errMsg = err.error?.message || 'Failed to upload the file.';
        this.uploadError = errMsg;
        this.toastService.showError(errMsg);
        this.isUploading = false;
      }
    });
  }

  // --- Day Modal Methods ---
  openDayModal(date: Date, documents: CalendarEntry[]) {
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

  isDocumentEntry(entry: CalendarEntry): boolean {
    return entry.entry_type !== 'direct_entry';
  }

  isDirectEntry(entry: CalendarEntry): boolean {
    return entry.entry_type === 'direct_entry';
  }

  getEntryTitle(entry: CalendarEntry): string {
    if (this.isDirectEntry(entry)) {
      return entry.metric_label || entry.document_name || 'Direct Entry';
    }

    return entry.document_name || 'Document';
  }

  getEntrySummary(entry: CalendarEntry): string {
    if (this.isDirectEntry(entry)) {
      return entry.metric_summary || 'Logged directly in VitaTrack';
    }

    return entry.category || 'Document';
  }

  getEntryTime(entry: CalendarEntry): string {
    const raw = entry.timestamp || entry.document_date;
    if (!raw) {
      return '';
    }

    return new Date(raw).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  getEntryClass(entry: CalendarEntry): string {
    return this.isDirectEntry(entry) ? 'doc-badge--direct-entry' : 'doc-badge--document';
  }

  getEntryCountLabel(entries: CalendarEntry[]): string {
    const directCount = entries.filter(entry => this.isDirectEntry(entry)).length;
    const documentCount = entries.length - directCount;

    if (directCount > 0 && documentCount > 0) {
      return `${entries.length} entries`;
    }

    if (directCount > 0) {
      return `${directCount} vital${directCount > 1 ? 's' : ''}`;
    }

    return `${documentCount} doc${documentCount > 1 ? 's' : ''}`;
  }

  handleCalendarEntryClick(entry: CalendarEntry, date: Date, dayEntries: CalendarEntry[], event?: Event) {
    event?.stopPropagation();

    if (this.isDocumentEntry(entry)) {
      this.openDocumentDetails(entry.id);
      return;
    }

    this.openDayModal(date, dayEntries);
  }

  handleDayModalEntryClick(entry: CalendarEntry) {
    if (this.isDocumentEntry(entry)) {
      this.openDocumentDetails(entry.id);
      this.closeDayModal();
    }
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
    const deleteRequest = this.isDirectEntry(this.documentToDelete)
      ? this.documentService.deleteHealthMetric(this.documentToDelete.id)
      : this.documentService.deleteDocument(this.documentToDelete.id);

    deleteRequest.subscribe({
      next: () => {
        this.isDeleting = false;
        this.closeDeleteModal();
        // Since the data changed, force a clean refresh of the month
        this.reportsMap = {};
        this.generateCalendar();
        this.fetchMonthData();
        this.fetchStorageUsage(); // Refresh storage limit
        
        // If the day modal was open, securely close it to prevent orphaned data
        if (this.isDayModalOpen) {
            this.closeDayModal();
        }
      },
      error: (err) => {
        console.error('Failed to delete entry', err);
        const errMsg = err.error?.message || 'Failed to delete the entry. Please try again.';
        this.deleteError = errMsg;
        this.toastService.showError(errMsg);
        this.isDeleting = false;
        this.cdr.detectChanges();
      }
    });
  }

  detailsError: string = '';

  // AI Analysis state
  isAnalyzing = false;
  analyzeError = '';
  isFullscreenImage: boolean = false;
  pdfBlobUrl: string | null = null;
  showCreditsSupportPopup = false;

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

          if (docData && docData.id) {
            this.documentService.getFileUrl(docData.id).subscribe({
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
                const errMsg = err.error?.message || 'Failed to load document file link.';
                this.detailsError = errMsg;
                this.toastService.showError(errMsg);
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
    this.isAnalyzing = false;
    this.analyzeError = '';
    this.showCreditsSupportPopup = false;
    
    // Cleanup memory from our temporary blob URLs
    if (this.pdfBlobUrl) {
      URL.revokeObjectURL(this.pdfBlobUrl);
      this.pdfBlobUrl = null;
    }
  }

  requestCreditsUpgrade() {
    this.showCreditsSupportPopup = false;
    window.open(
      'https://mail.google.com/mail/?view=cm&fs=1&to=support.vitatrack@gmail.com&subject=Request%20for%20AI%20Credit%20Upgrade&body=Hi,%20I%20would%20like%20more%20credits',
      '_blank'
    );
  }

  toggleEditDoc() {
    if (!this.selectedDocDetails) return;
    this.isEditingDetails = true;
    this.updateDocError = '';
    this.editForm.patchValue({
      document_name: this.selectedDocDetails.document_name,
      category: this.selectedDocDetails.category || 'Medical Report',
      document_date: this.selectedDocDetails.document_date ? this.selectedDocDetails.document_date.split('T')[0] : '',
      tags: this.selectedDocDetails.parsedTags ? this.selectedDocDetails.parsedTags.join(', ') : ''
    });
  }

  cancelEditDoc() {
    this.isEditingDetails = false;
  }

  saveEditDoc() {
    if (this.editForm.invalid || !this.selectedDocDetails) return;
    this.isUpdatingDocument = true;
    this.updateDocError = '';

    const formValues = this.editForm.value;
    const payload = {
      document_name: formValues.document_name,
      category: formValues.category,
      document_date: formValues.document_date,
      Tags: formValues.tags || ''
    };

    this.documentService.updateDocument(this.selectedDocDetails.id, payload).subscribe({
      next: () => {
        this.isUpdatingDocument = false;
        this.isEditingDetails = false;
        this.toastService.showSuccess('Document updated successfully');
        this.openDocumentDetails(this.selectedDocDetails.id);
        this.fetchMonthData();
      },
      error: (err: any) => {
        this.isUpdatingDocument = false;
        this.updateDocError = err?.error?.message || 'Failed to update document.';
        this.toastService.showError(this.updateDocError);
      }
    });
  }
  toggleCreditsSupportPopup() {
    this.showCreditsSupportPopup = !this.showCreditsSupportPopup;
  }

  hasGeneratedAnalysis(doc: any): boolean {
    return !!doc?.analysis_generated;
  }

  getAiAnalysis() {
    if (!this.selectedDocDetails) return;
    const fileId = this.selectedDocDetails?.file_id || this.selectedDocDetails?.id;
    const docName = this.selectedDocDetails?.document_name || 'Report';
    if (!fileId) {
      this.toastService.showError('No file ID found for AI analysis.');
      return;
    }
    this.isAnalyzing = true;
    this.analyzeError = '';
    this.documentService.getAiAnalysis(fileId).subscribe({
      next: (response) => {
        this.isAnalyzing = false;
        let analysisData = response;
        if (response?.json) analysisData = response.json;
        else if (response?.data?.json) analysisData = response.data.json;
        else if (response?.data) analysisData = response.data;
        this.closeDocumentDetails();
        this.router.navigate(['/analysis'], {
          state: { analysisData, docName }
        });
      },
      error: (err) => {
        this.isAnalyzing = false;
        const msg = err?.error?.message || 'Failed to get AI analysis. Please try again.';
        this.analyzeError = msg;
        this.toastService.showError(msg);
        this.cdr.detectChanges();
      }
    });
  }

  getAiActionLabel(doc: any): string {
    return this.hasGeneratedAnalysis(doc) ? 'View AI Analysis' : 'Get AI Analysis';
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

  normalizeCalendarEntry(entry: any): CalendarEntry {
    const entryType = entry?.entry_type === 'direct_entry' ? 'direct_entry' : 'document';
    return {
      id: entry?.id || '',
      entry_type: entryType,
      category: entry?.category || (entryType === 'direct_entry' ? 'Direct Entry' : 'Document'),
      document_name: entry?.document_name || entry?.metric_label || 'Untitled',
      status: entry?.status || (entryType === 'direct_entry' ? 'logged' : 'uploaded'),
      document_date: entry?.document_date || entry?.timestamp || '',
      analysis_generated: !!entry?.analysis_generated,
      metric_type: entry?.metric_type,
      metric_label: entry?.metric_label,
      metric_summary: entry?.metric_summary,
      timestamp: entry?.timestamp,
      tags: entry?.tags
    };
  }
}
