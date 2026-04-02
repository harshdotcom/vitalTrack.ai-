import { Component, OnInit, OnDestroy, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { DocumentService } from '../../../core/services/document';
import { ToastService } from '../../../core/services/toast';

@Component({
  selector: 'app-details-view',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './details-view.html',
  styleUrl: './details-view.css',
})
export class DetailsView implements OnInit, OnDestroy {
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private documentService = inject(DocumentService);
  private toastService = inject(ToastService);
  private sanitizer = inject(DomSanitizer);
  private cdr = inject(ChangeDetectorRef);

  docId: string | null = null;
  isLoading = true;
  error = '';
  docDetails: any = null;
  rawFileUrl = '';
  selectedFileUrl: SafeResourceUrl | null = null;
  pdfBlobUrl: string | null = null;
  isFullscreenImage = false;
  parsedTags: string[] = [];

  // AI Analysis state
  isAnalyzing = false;
  analyzeError = '';

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.docId = params['id'] || null;
      if (this.docId) {
        this.loadDocumentDetails(this.docId);
      } else {
        this.error = 'No document ID provided.';
        this.isLoading = false;
      }
    });
  }

  loadDocumentDetails(id: string) {
    this.isLoading = true;
    this.error = '';
    this.documentService.getDocumentDetails(id).subscribe({
      next: (response) => {
        let docData = response;
        if (response?.data) docData = response.data;
        else if (typeof response === 'string') {
          try { const p = JSON.parse(response); docData = p.data || p; } catch (e) {}
        }
        this.docDetails = docData;
        this.parsedTags = this.parseTags(docData?.tags);

        // Load file URL
        const fileId = docData?.file_id || docData?.id;
        if (fileId) {
          this.documentService.getFileUrl(fileId).subscribe({
            next: (fileRes) => {
              let fileData = fileRes;
              if (fileRes?.data) fileData = fileRes.data;
              if (fileData?.url) {
                this.rawFileUrl = fileData.url;
                if (this.isPdfFile()) {
                  fetch(this.rawFileUrl)
                    .then(res => res.blob())
                    .then(blob => {
                      const pdfBlob = new Blob([blob], { type: 'application/pdf' });
                      this.pdfBlobUrl = URL.createObjectURL(pdfBlob);
                      this.selectedFileUrl = this.sanitizer.bypassSecurityTrustResourceUrl(this.pdfBlobUrl);
                      this.isLoading = false;
                      this.cdr.detectChanges();
                    })
                    .catch(() => {
                      const googleProxy = `https://docs.google.com/viewer?url=${encodeURIComponent(this.rawFileUrl)}&embedded=true`;
                      this.selectedFileUrl = this.sanitizer.bypassSecurityTrustResourceUrl(googleProxy);
                      this.isLoading = false;
                      this.cdr.detectChanges();
                    });
                } else {
                  this.selectedFileUrl = this.sanitizer.bypassSecurityTrustResourceUrl(fileData.url);
                  this.isLoading = false;
                  this.cdr.detectChanges();
                }
              } else {
                this.isLoading = false;
                this.cdr.detectChanges();
              }
            },
            error: () => {
              this.isLoading = false;
              this.cdr.detectChanges();
            }
          });
        } else {
          this.isLoading = false;
          this.cdr.detectChanges();
        }
      },
      error: () => {
        this.error = 'Failed to load document details.';
        this.isLoading = false;
        this.cdr.detectChanges();
      }
    });
  }

  getAiAnalysis() {
    if (!this.docDetails) return;
    // Use file_id for the AI endpoint
    const fileId = this.docDetails?.file_id || this.docDetails?.id;
    if (!fileId) {
      this.toastService.showError('No file ID found for AI analysis.');
      return;
    }

    this.isAnalyzing = true;
    this.analyzeError = '';

    this.documentService.getAiAnalysis(fileId).subscribe({
      next: (response) => {
        this.isAnalyzing = false;
        // Extract json data from response
        let analysisData = response;
        if (response?.json) analysisData = response.json;
        else if (response?.data?.json) analysisData = response.data.json;
        else if (response?.data) analysisData = response.data;

        // Navigate to AI analysis page and pass data via router state
        this.router.navigate(['/analysis'], {
          state: { analysisData, docName: this.docDetails?.report_type || 'Report' }
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

  hasGeneratedAnalysis(): boolean {
    return !!this.docDetails?.analysis_generated;
  }

  get aiActionLabel(): string {
    return this.hasGeneratedAnalysis() ? 'View AI Analysis' : 'Get AI Analysis';
  }

  goBack() {
    this.router.navigate(['/dashboard']);
  }

  isImageFile(): boolean {
    if (!this.rawFileUrl) return false;
    const url = this.rawFileUrl.toLowerCase().split('?')[0];
    return url.endsWith('.png') || url.endsWith('.jpg') || url.endsWith('.jpeg');
  }

  isPdfFile(): boolean {
    if (!this.rawFileUrl) return false;
    return this.rawFileUrl.toLowerCase().split('?')[0].endsWith('.pdf');
  }

  parseTags(tags: string | string[]): string[] {
    if (!tags) return [];
    if (Array.isArray(tags)) return tags;
    try {
      const parsed = JSON.parse(tags);
      return Array.isArray(parsed) ? parsed : [];
    } catch { return []; }
  }

  ngOnDestroy() {
    if (this.pdfBlobUrl) URL.revokeObjectURL(this.pdfBlobUrl);
  }
}
