---
name: yolo-sahi
description: Slicing Aided Hyper Inference (SAHI) patterns, dynamic ROI slicing, overlap ratios, batch GPU inference, and Non-Maximum Merging (NMM) for small object detection in HydraVMS.
---

# 🔍 Slicing Aided Hyper Inference (SAHI) Skill (HydraVMS)

This skill guides the implementation and tuning of the SAHI pipeline for high-resolution cameras (1080p / 4K) to detect small objects (distant persons, license plates, small PPE).

---

## 1. SAHI Slicing Architecture

```text
Full High-Res Frame (e.g., 3840x2160 / 1920x1080)
  ├── 1. Motion Gate: Extract active bounding ROIs
  ├── 2. Slice Generation: 640x640 window with 20% overlap ratio
  ├── 3. Batch GPU Execution: Run all slices in 1 batched forward pass on RTX 5090
  └── 4. NMM (Non-Maximum Merging): Re-map slice coordinates to global frame space
```

---

## 2. Configuration Parameters
- `slice_height`: 640
- `slice_width`: 640
- `overlap_height_ratio`: 0.20
- `overlap_width_ratio`: 0.20
- `postprocess_type`: `NMM` (Non-Maximum Merging)
- `postprocess_match_threshold`: 0.50
- `adaptive_motion_slicing`: `true` (only slice regions with pixel variance)
