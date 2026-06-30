package engine

// --- Variant & Size Definitions ---

// Sizes maps size names to padding/font-size/height values
var sizeMap = map[string][3]string{
	"xs": {"4px 8px", "11px", "24px"},
	"sm": {"6px 12px", "12px", "28px"},
	"md": {"8px 16px", "13px", "36px"},
	"lg": {"12px 24px", "14px", "44px"},
}

// ValidVariants for buttons
var buttonVariants = map[string]bool{
	"solid":   true,
	"outline": true,
	"ghost":   true,
	"link":    true,
}

// ValidSizes for components
var validSizes = map[string]bool{
	"xs": true,
	"sm": true,
	"md": true,
	"lg": true,
}

// --- CSS Generation ---

// tokenCSS generates all design token CSS custom properties (MUI-derived).
func tokenCSS() string {
	return `/* --- ChefScript Design Tokens (MUI-derived) --- */
:root {
  /* Spacing — 8px base. Every margin, padding, gap must be a multiple. */
  --spacing-unit: 8px;
  --spacing-0: 0;
  --spacing-1: 4px;
  --spacing-2: 8px;
  --spacing-3: 12px;
  --spacing-4: 16px;
  --spacing-5: 20px;
  --spacing-6: 24px;
  --spacing-8: 32px;
  --spacing-10: 40px;
  --spacing-12: 48px;

  /* Typography — MUI scale, rem for responsive sizing. */
  --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  --font-h1: 300 6rem/1.167 var(--font-family);
  --font-h2: 300 3.75rem/1.2 var(--font-family);
  --font-h3: 400 3rem/1.167 var(--font-family);
  --font-h4: 400 2.125rem/1.235 var(--font-family);
  --font-h5: 400 1.5rem/1.334 var(--font-family);
  --font-h6: 500 1.25rem/1.6 var(--font-family);
  --font-subtitle1: 400 1rem/1.75 var(--font-family);
  --font-subtitle2: 500 0.875rem/1.57 var(--font-family);
  --font-body1: 400 1rem/1.5 var(--font-family);
  --font-body2: 400 0.875rem/1.43 var(--font-family);
  --font-caption: 400 0.75rem/1.66 var(--font-family);
  --font-overline: 400 0.75rem/2.66 var(--font-family);
  --font-button: 500 0.875rem/1.75 var(--font-family);

  /* Palette — MUI defaults */
  --color-primary: #1976d2;
  --color-primary-light: #42a5f5;
  --color-primary-dark: #1565c0;
  --color-primary-contrast: #fff;
  --color-secondary: #9c27b0;
  --color-secondary-light: #ba68c8;
  --color-secondary-dark: #7b1fa2;
  --color-error: #d32f2f;
  --color-error-light: #ef5350;
  --color-error-dark: #c62828;
  --color-warning: #ed6c02;
  --color-warning-light: #ff9800;
  --color-warning-dark: #e65100;
  --color-info: #0288d1;
  --color-info-light: #03a9f4;
  --color-info-dark: #01579b;
  --color-success: #2e7d32;
  --color-success-light: #4caf50;
  --color-success-dark: #1b5e20;

  /* Text */
  --text-primary: rgba(0,0,0,0.87);
  --text-secondary: rgba(0,0,0,0.6);
  --text-disabled: rgba(0,0,0,0.38);

  /* Surfaces */
  --color-divider: rgba(0,0,0,0.12);
  --color-background: #fff;
  --color-surface: #fff;
  --color-hover: rgba(0,0,0,0.04);
  --color-selected: rgba(0,0,0,0.08);
  --color-focus: rgba(0,0,0,0.12);
  --color-disabled-bg: rgba(0,0,0,0.12);

  /* Elevation — graduated shadow scale */
  --elevation-0: none;
  --elevation-1: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24);
  --elevation-2: 0 3px 6px rgba(0,0,0,0.15), 0 2px 4px rgba(0,0,0,0.12);
  --elevation-3: 0 10px 20px rgba(0,0,0,0.15), 0 3px 6px rgba(0,0,0,0.10);
  --elevation-4: 0 15px 25px rgba(0,0,0,0.15), 0 5px 10px rgba(0,0,0,0.05);
  --elevation-6: 0 20px 40px rgba(0,0,0,0.2);
  --elevation-8: 0 25px 50px rgba(0,0,0,0.25);
  --elevation-12: 0 30px 60px rgba(0,0,0,0.3);
  --elevation-24: 0 40px 80px rgba(0,0,0,0.35);

  /* Shape */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-full: 9999px;
  --radius-chip: 16px;

  /* Transitions */
  --ease-standard: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-decelerate: cubic-bezier(0.0, 0, 0.2, 1);
  --ease-accelerate: cubic-bezier(0.4, 0, 1, 1);
  --ease-sharp: cubic-bezier(0.4, 0, 0.6, 1);
  --duration-shortest: 150ms;
  --duration-shorter: 200ms;
  --duration-short: 250ms;
  --duration-standard: 300ms;
  --duration-complex: 375ms;

  /* Layout */
  --container-max-width: 1120px;

  /* Z-index */
  --z-drawer: 1200;
  --z-modal: 1300;
  --z-snackbar: 1400;
  --z-tooltip: 1500;

}

/* ── Dark mode — flip palette and surface tokens only ───────────── */
[data-theme="dark"] {
  /* Palette — MUI dark defaults (same hues, adjusted for dark surfaces) */
  --color-primary: #90caf9;
  --color-primary-light: #e3f2fd;
  --color-primary-dark: #42a5f5;
  --color-primary-contrast: rgba(0,0,0,0.87);
  --color-secondary: #ce93d8;
  --color-secondary-light: #f3e5f5;
  --color-secondary-dark: #ab47bc;
  --color-error: #f44336;
  --color-error-light: #e57373;
  --color-error-dark: #d32f2f;
  --color-warning: #ffa726;
  --color-warning-light: #ffb74d;
  --color-warning-dark: #f57c00;
  --color-info: #29b6f6;
  --color-info-light: #4fc3f7;
  --color-info-dark: #0288d1;
  --color-success: #66bb6a;
  --color-success-light: #81c784;
  --color-success-dark: #388e3c;

  /* Text */
  --text-primary: rgba(255,255,255,0.87);
  --text-secondary: rgba(255,255,255,0.6);
  --text-disabled: rgba(255,255,255,0.38);

  /* Surfaces */
  --color-divider: rgba(255,255,255,0.12);
  --color-background: #121212;
  --color-surface: #121212;
  --color-hover: rgba(255,255,255,0.08);
  --color-selected: rgba(255,255,255,0.16);
  --color-focus: rgba(255,255,255,0.12);
  --color-disabled-bg: rgba(255,255,255,0.12);

  /* Elevation — dark mode uses lighter overlays */
  --elevation-1: 0 1px 3px rgba(0,0,0,0.3), 0 1px 2px rgba(0,0,0,0.5);
  --elevation-2: 0 3px 6px rgba(0,0,0,0.35), 0 2px 4px rgba(0,0,0,0.3);
  --elevation-3: 0 10px 20px rgba(0,0,0,0.35), 0 3px 6px rgba(0,0,0,0.25);
  --elevation-4: 0 15px 25px rgba(0,0,0,0.35), 0 5px 10px rgba(0,0,0,0.2);
  --elevation-6: 0 20px 40px rgba(0,0,0,0.45);
  --elevation-8: 0 25px 50px rgba(0,0,0,0.5);
  --elevation-12: 0 30px 60px rgba(0,0,0,0.55);
  --elevation-24: 0 40px 80px rgba(0,0,0,0.6);
}

:root {
  /* ── Backward-compat aliases (old Tailwind tokens → MUI tokens) ──
     These exist so existing components don't break during Phase 2 migration.
     Remove once all components reference the MUI tokens directly. */

  /* Old spacing (4px base --space-N) */
  --space-0: 0;
  --space-px: 1px;
  --space-0\.5: 2px;
  --space-1: 4px;
  --space-1\.5: 6px;
  --space-2: 8px;
  --space-2\.5: 10px;
  --space-3: 12px;
  --space-3\.5: 14px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-7: 28px;
  --space-8: 32px;
  --space-9: 36px;
  --space-10: 40px;
  --space-12: 48px;
  --space-14: 56px;
  --space-16: 64px;
  --space-20: 80px;
  --space-24: 96px;

  /* Old radii */
  --radius-none: 0;
  --radius-xl: 8px;
  --radius-2xl: 12px;
  --radius-3xl: 16px;

  /* Old shadows */
  --shadow-none: none;
  --shadow-sm: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24);
  --shadow-md: 0 3px 6px rgba(0,0,0,0.15), 0 2px 4px rgba(0,0,0,0.12);
  --shadow-lg: 0 10px 20px rgba(0,0,0,0.15), 0 3px 6px rgba(0,0,0,0.10);
  --shadow-xl: 0 15px 25px rgba(0,0,0,0.15), 0 5px 10px rgba(0,0,0,0.05);
  --shadow-2xl: 0 20px 40px rgba(0,0,0,0.2);

  /* Old font sizes */
  --text-xs: 11px;
  --text-sm: 12px;
  --text-base: 14px;
  --text-lg: 16px;
  --text-xl: 18px;
  --text-2xl: 22px;
  --text-3xl: 28px;
  --text-4xl: 36px;
  --text-5xl: 48px;
  --leading-xs: 1.45;
  --leading-sm: 1.45;
  --leading-base: 1.5;
  --leading-lg: 1.5;
  --leading-xl: 1.4;
  --leading-2xl: 1.35;
  --leading-3xl: 1.3;
  --leading-4xl: 1.25;
  --leading-5xl: 1.2;

  /* Old font weights */
  --font-light: 300;
  --font-normal: 400;
  --font-medium: 500;
  --font-semibold: 600;
  --font-bold: 700;

  /* Old transitions */
  --duration-fast: 150ms;
  --duration-normal: 200ms;
  --duration-slow: 300ms;
  --ease-default: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-in: cubic-bezier(0.4, 0, 1, 1);
  --ease-out: cubic-bezier(0, 0, 0.2, 1);

  /* Old breakpoints */
  --bp-sm: 640px;
  --bp-md: 768px;
  --bp-lg: 1024px;
  --bp-xl: 1280px;
  --bp-2xl: 1536px;
}
`
}

// variantCSS generates the CSS for all component variants and sizes
func variantCSS() string {
	return `
/* --- Button Variants --- */
.cs-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-2);
  font-family: inherit;
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--duration-shorter) var(--ease-standard);
  text-decoration: none;
  white-space: nowrap;
  user-select: none;
  border: 1px solid transparent;
}
.cs-button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

/* Sizes */
.cs-button--xs { padding: var(--spacing-1) var(--spacing-2); font-size: 0.75rem; min-height: 24px; }
.cs-button--sm { padding: var(--spacing-2) var(--spacing-3); font-size: 0.75rem; min-height: 28px; }
.cs-button--md { padding: var(--spacing-2) var(--spacing-4); font-size: 0.75rem; min-height: 36px; letter-spacing: 0.01em; }
.cs-button--lg { padding: var(--spacing-3) var(--spacing-6); font-size: 0.875rem; min-height: 44px; }

/* Variant: solid */
.cs-button--solid {
  background: var(--accent);
  color: var(--accent-text);
  border-color: var(--accent);
}
.cs-button--solid:hover {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
  box-shadow: var(--elevation-1);
}

/* Variant: outline */
.cs-button--outline {
  background: transparent;
  color: var(--accent);
  border-color: var(--glass-border);
}
.cs-button--outline:hover {
  border-color: var(--accent);
  background: var(--color-primary-glass, rgba(102,126,234,0.08));
}

/* Variant: ghost */
.cs-button--ghost {
  background: transparent;
  color: var(--text-secondary);
  border-color: transparent;
}
.cs-button--ghost:hover {
  background: var(--glass-hover, rgba(0,0,0,0.04));
  color: var(--text-primary);
}

/* Variant: link */
.cs-button--link {
  background: transparent;
  color: var(--accent);
  border-color: transparent;
  padding-left: 0;
  padding-right: 0;
  min-height: auto;
}
.cs-button--link:hover {
  color: var(--accent-hover);
  text-decoration: underline;
}

/* Semantic color overrides */
.cs-button--color-primary.cs-button--solid { background: var(--color-primary); border-color: var(--color-primary); color: #fff; }
.cs-button--color-primary.cs-button--solid:hover { background: var(--color-primary-hover); border-color: var(--color-primary-hover); }
.cs-button--color-primary.cs-button--outline { color: var(--color-primary); border-color: var(--color-primary); }
.cs-button--color-primary.cs-button--outline:hover { background: var(--color-primary-glass); }

.cs-button--color-success.cs-button--solid { background: var(--color-success); border-color: var(--color-success); color: #fff; }
.cs-button--color-success.cs-button--solid:hover { background: var(--color-success-hover); border-color: var(--color-success-hover); }
.cs-button--color-success.cs-button--outline { color: var(--color-success); border-color: var(--color-success); }
.cs-button--color-success.cs-button--outline:hover { background: var(--color-success-glass); }

.cs-button--color-danger.cs-button--solid { background: var(--color-danger); border-color: var(--color-danger); color: #fff; }
.cs-button--color-danger.cs-button--solid:hover { background: var(--color-danger-hover); border-color: var(--color-danger-hover); }
.cs-button--color-danger.cs-button--outline { color: var(--color-danger); border-color: var(--color-danger); }
.cs-button--color-danger.cs-button--outline:hover { background: var(--color-danger-glass); }

.cs-button--color-warning.cs-button--solid { background: var(--color-warning); border-color: var(--color-warning); color: #fff; }
.cs-button--color-warning.cs-button--solid:hover { background: var(--color-warning-hover); border-color: var(--color-warning-hover); }
.cs-button--color-warning.cs-button--outline { color: var(--color-warning); border-color: var(--color-warning); }
.cs-button--color-warning.cs-button--outline:hover { background: var(--color-warning-glass); }

/* --- Button loading + disabled states --- */
.cs-button:disabled {
  opacity: 0.55;
  cursor: not-allowed;
  pointer-events: none;
}
.cs-button.cs-loading {
  opacity: 0.7;
  cursor: wait;
  pointer-events: none;
}

/* --- Field-level form errors --- */
.cs-input-error {
  border-color: var(--color-danger) !important;
  box-shadow: 0 0 0 2px var(--color-danger-glass) !important;
}
.cs-field-error {
  color: var(--color-danger);
  font-size: 0.75rem;
  margin-top: var(--spacing-1);
  padding-left: var(--spacing-1);
}

/* --- Icon --- */
.cs-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.cs-icon svg {
  display: block;
  width: 100%;
  height: 100%;
}
`
}
