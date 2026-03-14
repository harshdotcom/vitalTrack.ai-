# AI Health Journal – Final API Specification

BASE URL:
`/api/v1`

---

## AUTH APIs

### POST /auth/signup

Request:

```json
{
  "fullName": "Eren Yaegar",
  "email": "eren@gmail.com",
  "password": "StrongPassword123",
  "dob": "YYYY-MM-DD",
  "gender": "Male",
  "profile_pic": "File"
}
```

Response:

```json
{
    "users": {
        "user_id": 4,
        "email": "eren@gmail.com",
        "password": "$2a$12$XhgVWOwypzvYE1gz6CtAAeorysihE7Uzg9t0JExx16t/EDs2Luhlq",
        "GoogleId": null,
        "name": "Eren Yaegar",
        "dob": "1950-11-23T00:00:00+05:30",
        "gender": "Male",
        "profile_pic": "profile-pics-eren@gmail.com",
        "CreatedAt": "2026-03-14T02:38:17.7815419+05:30",
        "UpdatedAt": "2026-03-14T02:38:17.7815419+05:30"
    }
}
```

---

### POST /auth/login

Request:

```json
{
    "email": "eren@gmail.com",
    "password": "StrongPassword123"
}
```

Response:

```json
{
    "message": "Successfully login",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1pa2FzYUBnbWFpbC5jb20iLCJleHAiOjE3NzI2NjE2MzgsInVzZXJJZCI6MzZ9.R1M7xjJGezETqKNvQPKBLn5wZ3ejFLCLFXbuzpieh1Y",
    "users": {
        "user_id": 4,
        "email": "eren@gmail.com",
        "password": "$2a$12$XhgVWOwypzvYE1gz6CtAAeorysihE7Uzg9t0JExx16t/EDs2Luhlq",
        "GoogleId": null,
        "name": "Eren Yaegar",
        "dob": "1950-11-23T00:00:00+05:30",
        "gender": "Male",
        "profile_pic": "profile-pics-eren@gmail.com",
        "CreatedAt": "2026-03-14T02:38:17.7815419+05:30",
        "UpdatedAt": "2026-03-14T02:38:17.7815419+05:30"
    }
}
```

---

### POST /auth/google-login [Needs to be done]

Request:

```json
{
  "googleToken": "GOOGLE_ID_TOKEN"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "token": "JWT_TOKEN",
    "userId": "USER_ID",
    "loginType": "GOOGLE"
  }
}
```

---
## FILE / DOCUMENT APIs

### POST /files/upload

Request (FormData):

```
files: report.pdf

```

Response:

```json
{
    "files": [
        {
            "file_id": "a000c1d5-5bc7-40e2-8d32-962a99127832",
            "original_name": "Screenshot_22.png"
        }
    ]
}
```

### GET /files/:id

Request (Params):

```
id: 8044891e-e5cf-44fc-bdc4-2b3e8f7efcfb

```

Response:

```json
{
    "url": "https://vitatrack-documents-dev.s3.ap-south-1.amazonaws.com/ba4dcca3-05d9-410e-8122-a9783668bce0.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Checksum-Mode=ENABLED&X-Amz-Credential=AKIAWLX2NNJXXXTUF5PS%2F20260304%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20260304T205954Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=71a4d3d723eb10526f18aafe626197da12e9ad7c06c99d0178d387b00033c43e"
}
```

---

### POST /documents/calendar

Request:

```json
{
    "month": 3,
    "year": 2026
}
```

Response:

```json
{
    "days": {
        "2026-03-03": {
            "count": 2,
            "documents": [
                {
                    "id": "27d16371-64af-4514-aec7-8dfa16023919",
                    "user_id": 36,
                    "file_id": "ec037f8a-ba71-4c36-99ad-ff551b884317",
                    "category": "general",
                    "report_type": "Luffy",
                    "file_type": "lab_report",
                    "tags": "[\"fasting\",\"vitamin d\",\"abc\"]",
                    "status": "uploaded",
                    "report_date": "2026-03-03T05:30:00+05:30"
                },
                {
                    "id": "f9c85001-b718-4c66-acbf-0dccc56ddcdb",
                    "user_id": 36,
                    "file_id": "711dcdff-f233-45d6-8601-e151a7314b80",
                    "category": "general",
                    "report_type": "present",
                    "file_type": "lab_report",
                    "tags": "[\"fasting\"]",
                    "status": "uploaded",
                    "report_date": "2026-03-03T05:30:00+05:30"
                }
            ]
        }
    },
    "month": 3,
    "year": 2026
}
```

---

### GET /documents/:documentId

Request (Params):

```
id: 0e23fd6b-f6ad-42b4-a676-3682ab20b8cf

```

Response:

```json
{
    "id": "0e23fd6b-f6ad-42b4-a676-3682ab20b8cf",
    "user_id": 36,
    "file_id": "7d94725a-cba7-46f9-b04c-2fe7bf7f56bb",
    "category": "misc",
    "report_type": "nothing",
    "file_type": "pdf",
    "tags": "[\"no tag\"]",
    "status": "uploaded",
    "report_date": "2026-02-25T05:30:00+05:30"
}
```

---

### POST /documents

Request:

```json
   {
        "file_id":"7d94725a-cba7-46f9-b04c-2fe7bf7f56bb",
        "category": "misc",
        "report_type":"nothing",
        "file_type":"pdf",
        "tags": ["no tag"],
        "report_date": "2026-02-25"
    }
```

Response:

```json
{
    "document_id": "0e23fd6b-f6ad-42b4-a676-3682ab20b8cf",
    "status": "uploaded"
}
```

---

### DELETE /documents/:documentId

Request (Params):

```
id: 0e23fd6b-f6ad-42b4-a676-3682ab20b8cf

```

Response:

```json
{
    "message": "document deleted"
}
```

---

## AI APIs

### GET /files/ai/:fileId

Request (Params):

```
id: 8044891e-e5cf-44fc-bdc4-2b3e8f7efcfb

```

Response:

```json
{
    "json": {
        "report_metadata": {
            "report_date": "14-11-2025",
            "report_type": "Laboratory Report",
            "hospital_or_lab_name": "JEEVAN DIAGNOSTIC CENTRE"
        },
        "metrics": [
            {
                "test_name": "Serum Bilirubin Total",
                "value": "0.76",
                "unit": "mg/dl",
                "reference_range": "0.3-1.2 mg/dl",
                "status": "Normal"
            },
            {
                "test_name": "Serum Bilirubin Direct",
                "value": "0.24",
                "unit": "mg/dl",
                "reference_range": "0.1-0.4 mg/dl",
                "status": "Normal"
            },
            {
                "test_name": "Serum Bilirubin Indirect",
                "value": "0.52",
                "unit": "mg/dl",
                "reference_range": "0.2-0.8 mg/dl",
                "status": "Normal"
            },
            {
                "test_name": "ALT (SGPT)",
                "value": "84.0",
                "unit": "IU/L",
                "reference_range": "5-40 IU/L",
                "status": "High"
            },
            {
                "test_name": "AST (SGOT)",
                "value": "72.0",
                "unit": "IU/L",
                "reference_range": "5-40 IU/L",
                "status": "High"
            },
            {
                "test_name": "Alkaline Phosphatase",
                "value": "153.0",
                "unit": "IU/L",
                "reference_range": "25-130 IU/L",
                "status": "High"
            },
            {
                "test_name": "Serum Protein Total",
                "value": "9.41",
                "unit": "g/dl",
                "reference_range": "5.5-8.0 g/dl",
                "status": "High"
            },
            {
                "test_name": "Serum Albumin",
                "value": "4.23",
                "unit": "g/dl",
                "reference_range": "3.5-5.5 g/dl",
                "status": "Normal"
            },
            {
                "test_name": "Serum Globulin",
                "value": "5.10",
                "unit": "g/dl",
                "reference_range": "2.0-3.5 g/dl",
                "status": "High"
            },
            {
                "test_name": "A:G Ratio",
                "value": "0.82",
                "unit": "",
                "reference_range": "2:1",
                "status": "Low"
            },
            {
                "test_name": "Blood Urea",
                "value": "64.0",
                "unit": "mg/dl",
                "reference_range": "10-40 mg/dl",
                "status": "High"
            },
            {
                "test_name": "Blood Urea Nitrogen (BUN)",
                "value": "29.89",
                "unit": "mg/dl",
                "reference_range": "8-20 mg/dl",
                "status": "High"
            },
            {
                "test_name": "Serum Creatinine",
                "value": "0.14",
                "unit": "mg/dl",
                "reference_range": "0.5-1.5 mg/dl",
                "status": "Low"
            },
            {
                "test_name": "Serum Uric Acid",
                "value": "5.37",
                "unit": "mg/dl",
                "reference_range": "1.5-7.0 mg/dl",
                "status": "Normal"
            },
            {
                "test_name": "Serum Sodium",
                "value": "149.0",
                "unit": "mEq/L",
                "reference_range": "135-145 mEq/L",
                "status": "High"
            },
            {
                "test_name": "Serum Potassium",
                "value": "3.63",
                "unit": "mEq/L",
                "reference_range": "3.5-5.5 mEq/L",
                "status": "Normal"
            },
            {
                "test_name": "Serum Chloride",
                "value": "8.21",
                "unit": "mEq/L",
                "reference_range": "97-107 mEq/L",
                "status": "Low"
            },
            {
                "test_name": "WBC Count",
                "value": "8500",
                "unit": "cells/µL",
                "reference_range": "4000-11000 cells/µL",
                "status": "Normal"
            },
            {
                "test_name": "Platelet Count",
                "value": "324000",
                "unit": "cells/µL",
                "reference_range": "150000-450000 cells/µL",
                "status": "Normal"
            },
            {
                "test_name": "RBC Count",
                "value": "4.51",
                "unit": "million/µL",
                "reference_range": "3.5-5.5 million/µL",
                "status": "Normal"
            },
            {
                "test_name": "Hemoglobin",
                "value": "11.7",
                "unit": "g/dl",
                "reference_range": "13.5-17.5 g/dl (male)",
                "status": "Low"
            },
            {
                "test_name": "Hematocrit (PCV)",
                "value": "33.4",
                "unit": "%",
                "reference_range": "34-47 %",
                "status": "Low"
            },
            {
                "test_name": "MCV",
                "value": "74.06",
                "unit": "fL",
                "reference_range": "80-96 fL",
                "status": "Low"
            },
            {
                "test_name": "MCH",
                "value": "25.94",
                "unit": "pg",
                "reference_range": "27.5-33.2 pg",
                "status": "Low"
            },
            {
                "test_name": "MCHC",
                "value": "35.03",
                "unit": "%",
                "reference_range": "33.4-35.5 %",
                "status": "Normal"
            },
            {
                "test_name": "Neutrophils",
                "value": "76",
                "unit": "%",
                "reference_range": "40-70 %",
                "status": "High"
            },
            {
                "test_name": "Lymphocytes",
                "value": "17",
                "unit": "%",
                "reference_range": "20-45 %",
                "status": "Low"
            },
            {
                "test_name": "Eosinophils",
                "value": "5",
                "unit": "%",
                "reference_range": "2-6 %",
                "status": "Normal"
            },
            {
                "test_name": "Monocytes",
                "value": "2",
                "unit": "%",
                "reference_range": "1-6 %",
                "status": "Normal"
            },
            {
                "test_name": "Basophils",
                "value": "0",
                "unit": "%",
                "reference_range": "0-1 %",
                "status": "Normal"
            },
            {
                "test_name": "ESR (First Hour)",
                "value": "35",
                "unit": "mm",
                "reference_range": "0-10 mm",
                "status": "High"
            },
            {
                "test_name": "Malaria Parasite",
                "value": "Negative",
                "unit": "",
                "reference_range": "Negative",
                "status": "Normal"
            }
        ],
        "abnormal_findings": [
            "ALT (SGPT)",
            "AST (SGOT)",
            "Alkaline Phosphatase",
            "Serum Protein Total",
            "Serum Globulin",
            "A:G Ratio",
            "Blood Urea",
            "Blood Urea Nitrogen (BUN)",
            "Serum Creatinine",
            "Serum Sodium",
            "Serum Chloride",
            "Hemoglobin",
            "Hematocrit (PCV)",
            "MCV",
            "MCH",
            "Neutrophils",
            "Lymphocytes",
            "ESR (First Hour)"
        ],
        "simple_explanation": "Multiple liver enzymes, kidney function markers, and blood count parameters are outside their normal ranges, indicating hepatic stress, possible renal impairment, and anemia with a reactive neutrophil shift.",
        "overall_risk_level": "Moderate",
        "recommendations": {
            "diet": [
                "Limit high‑fat and fried foods",
                "Reduce sodium intake",
                "Increase intake of fruits, vegetables, and lean protein",
                "Avoid alcohol"
            ],
            "lifestyle": [
                "Stay well‑hydrated",
                "Engage in moderate aerobic exercise 3‑4 times per week",
                "Get adequate sleep (7‑8 hours)",
                "Follow up with a physician within 2 weeks"
            ]
        },
        "follow_up_suggestions": [
            "Repeat liver function tests in 4‑6 weeks",
            "Repeat renal panel and electrolytes",
            "Complete iron studies and vitamin B12 assessment for anemia",
            "Consult a gastroenterologist if liver enzymes remain elevated"
        ]
    }
}
```
---

## User Details API

### GET /user-details/usage

Response:

```json
{
    "data": {
        "UserID": 2,
        "TotalStorageUsed": 321344
    },
    "message": "User storage usage fetched successfully"
}
```

### PATCH /user-details/update

Request:

```json
   {
    "name": "Ereh Yaegar",
    "dob": "2000-11-24",
    "profile_pic": "file",
    "gender": "Female",
    "delete_profile_pic": "false"
    }
```

Response:

```json
{
    "message": "Update successful",
    "user": {
        "user_id": 4,
        "email": "eren@gmail.com",
        "password": "$2a$12$XhgVWOwypzvYE1gz6CtAAeorysihE7Uzg9t0JExx16t/EDs2Luhlq",
        "GoogleId": null,
        "name": "Ereh Yaegar",
        "dob": "2000-11-24T00:00:00+05:30",
        "gender": "Male",
        "profile_pic": null,
        "CreatedAt": "2026-03-14T02:38:17.781541+05:30",
        "UpdatedAt": "2026-03-14T02:38:17.781541+05:30"
    }
}
```
---

## HEALTH JOURNAL API

### POST /user/get-health-timeline

Request:

```json
{
  "userId": "USER_ID",
  "startDate": "2026-02-01",
  "endDate": "2026-02-28"
}
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "date": "2026-02-25",
      "documentId": "DOC_ID",
      "documentType": "Blood Test",
      "summary": "Hemoglobin slightly low"
    }
  ]
}
```

---

## ID FLOW (REFERENCE)

```
Login → USER_ID
Upload Document → DOC_ID
Analyze Document → ANALYSIS_ID
Generate Food Recommendation → USER_ID + ANALYSIS_ID
Timeline → USER_ID + Date Range
```

## Current Deploy Flow

1. **Create the build:**

   ```bash
   ng build --configuration production
   ```

2. **Navigate to the output folder:**

   `AI_Bharat\vitalTrack.ai-\frontend\dist\vitalTrackFrontend\browser`

3. **Copy all the generated files.**

4. **On the VM**, navigate to:

   `C:\inetpub\wwwroot`

5. **Remove all existing files** in that folder and **paste the new files** (from step 3).

6. The latest version of the frontend will now be live via the IIS-hosted server / public IP.
