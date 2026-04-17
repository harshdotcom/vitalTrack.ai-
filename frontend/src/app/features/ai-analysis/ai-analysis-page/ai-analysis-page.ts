import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { ThemeToggleComponent } from '../../../core/components/theme-toggle/theme-toggle';

interface Metric {
  test_name: string;
  value: string;
  unit: string;
  reference_range: string;
  status: 'Normal' | 'High' | 'Low';
}

interface AnalysisData {
  report_metadata?: {
    document_date?: string;
    document_name?: string;
    hospital_or_lab_name?: string;
  };
  metrics?: Metric[];
  abnormal_findings?: string[];
  simple_explanation?: string;
  overall_risk_level?: string;
  recommendations?: {
    diet?: string[];
    lifestyle?: string[];
  };
  follow_up_suggestions?: string[];
}

@Component({
  selector: 'app-ai-analysis-page',
  standalone: true,
  imports: [CommonModule, ThemeToggleComponent],
  templateUrl: './ai-analysis-page.html',
  styleUrl: './ai-analysis-page.css',
})
export class AiAnalysisPage implements OnInit {
  private router = inject(Router);

  analysisData: AnalysisData | null = null;
  docName = 'Medical Report';
  hasData = false;

  // For metric table grouping
  metricGroups: { label: string; metrics: Metric[] }[] = [];

  ngOnInit(): void {
    // Angular router passes state via history.state when using router.navigate with state
    const state = window.history.state as any;
    if (state?.analysisData) {
      this.analysisData = state.analysisData;
      this.docName = state.docName || 'Medical Report';
      this.hasData = true;
      this.buildMetricGroups();
    } else {
      this.hasData = false;
    }
  }

  buildMetricGroups() {
    if (!this.analysisData?.metrics) return;
    const metrics = this.analysisData.metrics;

    // Group by test category based on known medical categories
    const liverTests = ['Serum Bilirubin Total','Serum Bilirubin Direct','Serum Bilirubin Indirect','ALT (SGPT)','AST (SGOT)','Alkaline Phosphatase','Serum Protein Total','Serum Albumin','Serum Globulin','A:G Ratio'];
    const kidneyTests = ['Blood Urea','Blood Urea Nitrogen (BUN)','Serum Creatinine','Serum Uric Acid'];
    const electrolytes = ['Serum Sodium','Serum Potassium','Serum Chloride'];
    const cbc = ['WBC Count','Platelet Count','RBC Count','Hemoglobin','Hematocrit (PCV)','MCV','MCH','MCHC','Neutrophils','Lymphocytes','Eosinophils','Monocytes','Basophils','ESR (First Hour)','Malaria Parasite'];

    const groups: { label: string; keys: string[] }[] = [
      { label: 'Liver Function Tests', keys: liverTests },
      { label: 'Kidney Function Tests', keys: kidneyTests },
      { label: 'Electrolytes', keys: electrolytes },
      { label: 'Complete Blood Count (CBC)', keys: cbc },
    ];

    const matched = new Set<string>();
    this.metricGroups = groups.map(g => {
      const groupMetrics = metrics.filter(m => g.keys.includes(m.test_name));
      groupMetrics.forEach(m => matched.add(m.test_name));
      return { label: g.label, metrics: groupMetrics };
    }).filter(g => g.metrics.length > 0);

    // Any unmatched goes to "Other Tests"
    const others = metrics.filter(m => !matched.has(m.test_name));
    if (others.length > 0) {
      this.metricGroups.push({ label: 'Other Tests', metrics: others });
    }
  }

  get riskLevel(): string {
    return this.analysisData?.overall_risk_level || 'Unknown';
  }

  get riskClass(): string {
    const r = this.riskLevel.toLowerCase();
    if (r === 'low') return 'risk-low';
    if (r === 'moderate') return 'risk-moderate';
    if (r === 'high') return 'risk-high';
    if (r === 'critical') return 'risk-critical';
    return 'risk-unknown';
  }

  get abnormalCount(): number {
    return this.analysisData?.abnormal_findings?.length || 0;
  }

  get normalCount(): number {
    const total = this.analysisData?.metrics?.length || 0;
    return total - this.abnormalCount;
  }

  get totalTests(): number {
    return this.analysisData?.metrics?.length || 0;
  }

  isAbnormal(testName: string): boolean {
    return (this.analysisData?.abnormal_findings || []).includes(testName);
  }

  statusClass(status: string): string {
    switch (status?.toLowerCase()) {
      case 'normal': return 'status-normal';
      case 'high': return 'status-high';
      case 'low': return 'status-low';
      default: return '';
    }
  }

  goBack() {
    this.router.navigate(['/dashboard']);
  }
}
