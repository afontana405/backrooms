package engine

// componentLibCSS contains CSS for all extended components.
// Injected after base componentCSS.
const componentLibCSS = `
@keyframes cs-ripple {
  to { transform: scale(2.5); opacity: 0; }
}
@keyframes cs-shimmer {
  0% { background-position: -400px 0; }
  100% { background-position: 400px 0; }
}
@keyframes cs-spin {
  to { transform: rotate(360deg); }
}
@keyframes cs-snackbar-in {
  from { transform: translateY(16px); opacity: 0; }
  to   { transform: translateY(0); opacity: 1; }
}

/* ── Input ───────────────────────────────────────────────────────────────── */
.cs-input {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.cs-input__wrap {
  position: relative;
  display: flex;
  flex-direction: column;
}
.cs-input__field {
  width: 100%;
  background: var(--glass);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-family: inherit;
  font-size: 0.875rem;
  padding: 20px 14px 8px;
  outline: none;
  transition: border-color var(--duration-shortest) var(--ease-standard), box-shadow var(--duration-shortest) var(--ease-standard);
  line-height: 1.4;
  -webkit-appearance: none;
}
.cs-input__field:focus {
  border-color: var(--glass-border-hover);
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.06);
}
.cs-input--error .cs-input__field {
  border-color: var(--color-danger);
}
.cs-input__label {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.875rem;
  color: var(--text-tertiary);
  pointer-events: none;
  transition: all var(--duration-shortest) var(--ease-standard);
  transform-origin: left top;
  white-space: nowrap;
}
.cs-input__label--float,
.cs-input__field:focus ~ .cs-input__label,
.cs-input__field:not(:placeholder-shown) ~ .cs-input__label {
  top: 10px;
  transform: translateY(0) scale(0.78);
  color: var(--accent);
  font-weight: var(--font-medium);
}
.cs-input--error .cs-input__label--float,
.cs-input--error .cs-input__field:focus ~ .cs-input__label {
  color: var(--color-danger);
}
.cs-input__hint {
  font-size: 0.75rem;
  color: var(--text-tertiary);
  padding-left: 14px;
}
.cs-input__hint--error {
  color: var(--color-danger);
}

/* Textarea */
.cs-textarea .cs-input__field {
  resize: vertical;
  min-height: 96px;
  padding-top: 24px;
}
.cs-textarea .cs-input__label {
  top: 18px;
  transform: none;
}
.cs-textarea .cs-input__label--float,
.cs-textarea .cs-input__field:focus ~ .cs-input__label,
.cs-textarea .cs-input__field:not(:placeholder-shown) ~ .cs-input__label {
  top: 8px;
  transform: scale(0.78);
  color: var(--accent);
}

/* ── Select ─────────────────────────────────────────────────────────────── */
.cs-select {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.cs-select__label {
  font-size: 0.75rem;
  color: var(--text-tertiary);
  font-weight: var(--font-medium);
}
.cs-select__trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--glass);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  padding: 10px 14px;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--text-primary);
  transition: border-color var(--duration-shortest) var(--ease-standard);
  user-select: none;
}
.cs-select__trigger:hover { border-color: var(--glass-border-hover); }
.cs-select__arrow { color: var(--text-tertiary); font-size: 10px; }
.cs-select__dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0; right: 0;
  background: rgba(15, 15, 20, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  box-shadow: var(--elevation-4);
  z-index: 100;
  overflow-y: auto;
  max-height: 240px;
}
.cs-select__option {
  padding: 14px 16px;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--text-primary);
  transition: background var(--duration-shortest) var(--ease-standard);
}
.cs-select__option:hover,
.cs-select__option--active { background: var(--glass-hover); }

/* ── Autocomplete ───────────────────────────────────────────────────────── */
.cs-autocomplete {
  position: relative;
}
.cs-autocomplete__dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0; right: 0;
  background: rgba(15, 15, 20, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  box-shadow: var(--elevation-4);
  z-index: 100;
  max-height: 240px;
  overflow-y: auto;
}
.cs-autocomplete__item {
  padding: 10px 14px;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--text-primary);
  transition: background var(--duration-shortest) var(--ease-standard);
}
.cs-autocomplete__item:hover,
.cs-autocomplete__item--active { background: var(--glass-hover); }

/* ── Checkbox ───────────────────────────────────────────────────────────── */
.cs-checkbox {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  user-select: none;
}
.cs-checkbox__input { display: none; }
.cs-checkbox__box {
  width: 18px; height: 18px;
  border: 2px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--bg-surface);
  display: flex; align-items: center; justify-content: center;
  transition: all var(--duration-shortest) var(--ease-standard);
  flex-shrink: 0;
}
.cs-checkbox__input:checked ~ .cs-checkbox__box {
  background: var(--accent);
  border-color: var(--accent);
}
.cs-checkbox__input:checked ~ .cs-checkbox__box::after {
  content: '';
  width: 5px; height: 9px;
  border: 2px solid var(--accent-text, #000);
  border-top: none; border-left: none;
  transform: rotate(45deg) translateY(-1px);
  display: block;
}
.cs-checkbox__label { font-size: 0.875rem; color: var(--text-primary); }

/* ── Radio ──────────────────────────────────────────────────────────────── */
.cs-radio {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  user-select: none;
}
.cs-radio__input { display: none; }
.cs-radio__dot {
  width: 18px; height: 18px;
  border: 2px solid var(--border);
  border-radius: 50%;
  background: var(--bg-surface);
  display: flex; align-items: center; justify-content: center;
  transition: all var(--duration-shortest) var(--ease-standard);
  flex-shrink: 0;
}
.cs-radio__input:checked ~ .cs-radio__dot {
  border-color: var(--accent);
}
.cs-radio__input:checked ~ .cs-radio__dot::after {
  content: '';
  width: 8px; height: 8px;
  border-radius: 50%;
  background: var(--accent);
  display: block;
}
.cs-radio__label { font-size: 0.875rem; color: var(--text-primary); }

/* ── Switch ─────────────────────────────────────────────────────────────── */
.cs-switch {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  user-select: none;
}
.cs-switch__input { display: none; }
.cs-switch__track {
  width: 40px; height: 22px;
  border-radius: var(--radius-full);
  background: var(--border);
  position: relative;
  transition: background var(--duration-shorter) var(--ease-standard);
  flex-shrink: 0;
}
.cs-switch__thumb {
  position: absolute;
  top: 3px; left: 3px;
  width: 16px; height: 16px;
  border-radius: 50%;
  background: #fff;
  transition: transform var(--duration-shorter) var(--ease-standard);
  box-shadow: var(--elevation-1);
}
.cs-switch__input:checked ~ .cs-switch__track { background: var(--accent); }
.cs-switch__input:checked ~ .cs-switch__track .cs-switch__thumb { transform: translateX(18px); }
.cs-switch__label { font-size: 0.875rem; color: var(--text-primary); }

/* ── Alert ──────────────────────────────────────────────────────────────── */
.cs-alert {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  border-radius: var(--radius-sm);
  border: 1px solid transparent;
  font-size: 0.875rem;
}
.cs-alert__icon { font-size: 16px; flex-shrink: 0; margin-top: 1px; }
.cs-alert__body { flex: 1; }
.cs-alert__title { display: block; font-weight: var(--font-semibold); margin-bottom: 2px; }
.cs-alert__dismiss {
  background: none; border: none; cursor: pointer;
  color: currentColor; opacity: 0.6; font-size: 14px; padding: 0; flex-shrink: 0;
}
.cs-alert__dismiss:hover { opacity: 1; }
.cs-alert--info    { background: rgba(65,105,225,0.1);  border-color: rgba(65,105,225,0.25);  color: var(--color-info);    backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px); }
.cs-alert--success { background: rgba(52,211,153,0.1); border-color: rgba(52,211,153,0.25); color: var(--color-success); backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px); }
.cs-alert--warning { background: rgba(234,179,8,0.1);  border-color: rgba(234,179,8,0.25);  color: var(--color-warning); backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px); }
.cs-alert--error   { background: rgba(220,38,38,0.1);  border-color: rgba(220,38,38,0.25);  color: var(--color-danger);  backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px); }

/* ── Badge ──────────────────────────────────────────────────────────────── */
.cs-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: var(--font-semibold);
  letter-spacing: 0.03em;
  line-height: 1.5;
}
.cs-badge--default { background: var(--bg-surface); color: var(--text-secondary); border: 1px solid var(--border); }
.cs-badge--primary { background: var(--accent); color: var(--accent-text); }
.cs-badge--success { background: var(--color-success); color: #fff; }
.cs-badge--warning { background: var(--color-warning); color: #fff; }
.cs-badge--error   { background: var(--color-danger);  color: #fff; }
.cs-badge--info    { background: var(--color-info);    color: #fff; }

/* ── Chip ───────────────────────────────────────────────────────────────── */
.cs-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: var(--font-medium);
  background: var(--bg-surface);
  border: 1px solid var(--border);
  color: var(--text-primary);
}
.cs-chip__dismiss {
  background: none; border: none; cursor: pointer;
  color: var(--text-tertiary); font-size: 11px; padding: 0;
  display: flex; align-items: center;
}
.cs-chip__dismiss:hover { color: var(--text-primary); }

/* ── Spinner ────────────────────────────────────────────────────────────── */
.cs-spinner {
  display: inline-block;
  border-radius: 50%;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  animation: cs-spin 0.7s linear infinite;
}
.cs-spinner--xs { width: 14px; height: 14px; }
.cs-spinner--sm { width: 18px; height: 18px; }
.cs-spinner--md { width: 24px; height: 24px; }
.cs-spinner--lg { width: 36px; height: 36px; }

/* ── Skeleton ───────────────────────────────────────────────────────────── */
.cs-skeleton {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}
.cs-skeleton__lines {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.cs-skeleton__block {
  height: 14px;
  border-radius: var(--radius-sm);
  background: linear-gradient(90deg, var(--bg-surface) 25%, var(--bg-elevated) 50%, var(--bg-surface) 75%);
  background-size: 400px 100%;
  animation: cs-shimmer 1.4s ease infinite;
}
.cs-skeleton__avatar {
  width: 40px; height: 40px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* ── Progress ───────────────────────────────────────────────────────────── */
.cs-progress { display: flex; flex-direction: column; gap: 6px; }
.cs-progress__label {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--text-secondary);
}
.cs-progress__track {
  height: 6px;
  background: var(--bg-surface);
  border-radius: var(--radius-full);
  overflow: hidden;
}
.cs-progress__fill {
  height: 100%;
  background: var(--accent);
  border-radius: var(--radius-full);
  transition: width var(--duration-standard) var(--ease-decelerate);
}

/* ── Tooltip ────────────────────────────────────────────────────────────── */
.cs-tooltip {
  position: relative;
  display: inline-flex;
}
.cs-tooltip__text {
  position: absolute;
  background: var(--text-primary);
  color: var(--bg);
  font-size: 0.75rem;
  padding: 5px 9px;
  border-radius: var(--radius-sm);
  white-space: nowrap;
  pointer-events: none;
  opacity: 0;
  transition: opacity var(--duration-shortest) var(--ease-standard);
  z-index: 200;
}
.cs-tooltip:hover .cs-tooltip__text { opacity: 1; }
.cs-tooltip--top .cs-tooltip__text    { bottom: calc(100% + 6px); left: 50%; transform: translateX(-50%); }
.cs-tooltip--bottom .cs-tooltip__text { top: calc(100% + 6px); left: 50%; transform: translateX(-50%); }
.cs-tooltip--left .cs-tooltip__text   { right: calc(100% + 6px); top: 50%; transform: translateY(-50%); }
.cs-tooltip--right .cs-tooltip__text  { left: calc(100% + 6px); top: 50%; transform: translateY(-50%); }

/* ── Divider ────────────────────────────────────────────────────────────── */
hr.cs-divider {
  border: none;
  border-top: 1px solid var(--border-subtle);
  margin: 16px 0;
}
.cs-divider--labeled {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 16px 0;
}
.cs-divider__line { flex: 1; height: 1px; background: var(--border-subtle); }
.cs-divider__label { font-size: 0.75rem; color: var(--text-tertiary); white-space: nowrap; }

/* ── Tabs ───────────────────────────────────────────────────────────────── */
.cs-tabs {
  display: flex;
  flex-direction: column;
}
.cs-tab {
  display: inline-flex;
  align-items: center;
  padding: 8px 16px;
  font-size: 0.75rem;
  font-weight: var(--font-medium);
  color: var(--text-tertiary);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition: all var(--duration-shortest) var(--ease-standard);
  font-family: inherit;
}
.cs-tab:hover { color: var(--text-primary); }
.cs-tab--active { color: var(--accent); border-bottom-color: var(--accent); }

/* Tab trigger bar — all buttons in a row */
.cs-tabs { overflow: hidden; }
.cs-tab-panel { display: none; padding: 16px 0; }
.cs-tab-panel--active { display: block; }

/* Tabs row — wrap all cs-tab buttons in a flex row */
.cs-tabs > .cs-tab { order: -1; }
.cs-tabs {
  display: flex;
  flex-wrap: wrap;
  border-bottom: 1px solid var(--border-subtle);
}
.cs-tab { flex-shrink: 0; }
.cs-tab-panel { width: 100%; border-top: none; }

/* ── Accordion ──────────────────────────────────────────────────────────── */
.cs-accordion { display: flex; flex-direction: column; }
.cs-accordion-item { border-bottom: 1px solid var(--border-subtle); }
.cs-accordion-trigger {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  background: none;
  border: none;
  color: var(--text-primary);
  font-size: 0.875rem;
  font-weight: var(--font-medium);
  cursor: pointer;
  text-align: left;
  font-family: inherit;
  transition: color var(--duration-shortest) var(--ease-standard);
}
.cs-accordion-trigger:hover { color: var(--accent); }
.cs-accordion-icon {
  font-size: 11px;
  color: var(--text-tertiary);
  transition: transform var(--duration-shorter) var(--ease-standard);
  flex-shrink: 0;
}
.cs-accordion-item--open .cs-accordion-icon { transform: rotate(180deg); }
.cs-accordion-body {
  max-height: 0;
  overflow: hidden;
  transition: max-height var(--duration-shorter) var(--ease-decelerate);
}
.cs-accordion-content {
  padding-bottom: 14px;
  color: var(--text-secondary);
  font-size: 0.875rem;
  line-height: 1.6;
}

/* ── Drawer ─────────────────────────────────────────────────────────────── */
.cs-drawer {
  display: none;
  position: fixed;
  inset: 0;
  z-index: 1000;
}
.cs-drawer--open { display: block; }
.cs-drawer__backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(2px);
}
.cs-drawer__panel {
  position: absolute;
  top: 0; bottom: 0;
  width: 360px;
  max-width: 90vw;
  background: rgba(12, 12, 18, 0.85);
  backdrop-filter: blur(32px);
  -webkit-backdrop-filter: blur(32px);
  border: 1px solid var(--glass-border);
  box-shadow: var(--elevation-6);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}
.cs-drawer--right .cs-drawer__panel { right: 0; border-right: none; border-radius: var(--radius-md) 0 0 var(--radius-md); }
.cs-drawer--left  .cs-drawer__panel { left: 0;  border-left: none;  border-radius: 0 var(--radius-md) var(--radius-md) 0; }
.cs-drawer__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--border-subtle);
  flex-shrink: 0;
}
.cs-drawer__title { font-size: 1rem; font-weight: var(--font-semibold); color: var(--text-primary); }
.cs-drawer__close {
  background: none; border: none; cursor: pointer;
  color: var(--text-tertiary); font-size: 18px; padding: 4px;
}
.cs-drawer__close:hover { color: var(--text-primary); }
.cs-drawer__body { padding: 20px 24px; flex: 1; }

/* ── Snackbar ───────────────────────────────────────────────────────────── */
.cs-snackbar {
  display: inline-flex;
  align-items: center;
  padding: 12px 20px;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: var(--font-medium);
  box-shadow: var(--elevation-4);
  opacity: 0;
  transform: translateY(16px);
  transition:all var(--duration-standard) var(--ease-decelerate);
  max-width: 420px;
}
.cs-snackbar--show { opacity: 1; transform: translateY(0); }
.cs-snackbar--info    { background: var(--text-primary); color: var(--bg); }
.cs-snackbar--success { background: var(--color-success); color: #fff; }
.cs-snackbar--warning { background: var(--color-warning); color: #fff; }
.cs-snackbar--error   { background: var(--color-danger);  color: #fff; }

/* ── Breadcrumb ─────────────────────────────────────────────────────────── */
.cs-breadcrumb__list {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0;
  list-style: none;
  padding: 0;
}
.cs-breadcrumb__item {
  display: flex;
  align-items: center;
  font-size: 0.75rem;
  color: var(--text-tertiary);
}
.cs-breadcrumb__item:not(:last-child)::after {
  content: '/';
  margin: 0 8px;
  color: var(--border);
}
.cs-breadcrumb__item:last-child { color: var(--text-primary); font-weight: var(--font-medium); }
.cs-breadcrumb__link {
  color: var(--text-tertiary);
  text-decoration: none;
  transition: color var(--duration-shortest) var(--ease-standard);
}
.cs-breadcrumb__link:hover { color: var(--accent); }

/* ── Pagination ─────────────────────────────────────────────────────────── */
.cs-pagination {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
}
.cs-pagination__btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  height: 32px;
  padding: 0 6px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--bg-surface);
  color: var(--text-secondary);
  font-size: 0.75rem;
  cursor: pointer;
  transition: all var(--duration-shortest) var(--ease-standard);
  font-family: inherit;
}
.cs-pagination__btn:hover { border-color: var(--accent); color: var(--accent); }
.cs-pagination__btn--active { background: var(--accent); border-color: var(--accent); color: var(--accent-text); font-weight: var(--font-semibold); }
.cs-pagination__btn--disabled { opacity: 0.4; pointer-events: none; }
.cs-pagination__ellipsis { padding: 0 4px; color: var(--text-tertiary); font-size: 0.75rem; }

/* ── Table ──────────────────────────────────────────────────────────────── */
.cs-table-wrap {
  width: 100%;
  overflow-x: auto;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}
.cs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.75rem;
}
.cs-table th {
  padding: 10px 16px;
  text-align: left;
  font-weight: var(--font-semibold);
  color: var(--text-tertiary);
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border);
  white-space: nowrap;
}
.cs-table td {
  padding: 10px 16px;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-subtle);
  vertical-align: middle;
}
.cs-table tr:last-child td { border-bottom: none; }
.cs-table--striped tbody tr:nth-child(even) td { background: var(--bg-surface); }
.cs-table--hoverable tbody tr:hover td { background: var(--bg-surface); }

/* ── List ───────────────────────────────────────────────────────────────── */
.cs-list { list-style: none; padding: 0; margin: 0; }
.cs-list--divided .cs-list-item { border-bottom: 1px solid var(--border-subtle); }
.cs-list--divided .cs-list-item:last-child { border-bottom: none; }
.cs-list-item {
  padding: 10px 0;
  font-size: 0.875rem;
  color: var(--text-primary);
  cursor: default;
}
.cs-list-item:has([onclick]) { cursor: pointer; }
.cs-list__link {
  color: var(--text-primary);
  text-decoration: none;
  display: block;
}
.cs-list__link:hover { color: var(--accent); }

/* ── Avatar ──────────────────────────────────────────────────────────────── */
.cs-avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  font-weight: var(--font-semibold);
  background: var(--glass);
  border: 1px solid var(--glass-border);
  color: var(--text-primary);
}
.cs-avatar--xs { width: 24px; height: 24px; font-size: 10px; }
.cs-avatar--sm { width: 32px; height: 32px; font-size: 12px; }
.cs-avatar--md { width: 40px; height: 40px; font-size: 14px; }
.cs-avatar--lg { width: 56px; height: 56px; font-size: 18px; }
.cs-avatar--color-success { background: var(--color-success-glass); border-color: rgba(52,211,153,0.3); color: var(--color-success); }
.cs-avatar--color-danger  { background: var(--color-danger-glass);  border-color: rgba(220,38,38,0.3);  color: var(--color-danger);  }
.cs-avatar--color-info    { background: var(--color-info-glass);    border-color: rgba(65,105,225,0.3);  color: var(--color-info);    }
.cs-avatar--color-warning { background: var(--color-warning-glass); border-color: rgba(234,179,8,0.3);  color: var(--color-warning); }
.cs-avatar__img { width: 100%; height: 100%; object-fit: cover; }
.cs-avatar__icon { width: 55%; height: 55%; opacity: 0.6; }
.cs-avatar-group { display: flex; }
.cs-avatar-group .cs-avatar { margin-left: -8px; border: 2px solid var(--bg); }
.cs-avatar-group .cs-avatar:first-child { margin-left: 0; }

/* ── EmptyState ──────────────────────────────────────────────────────────── */
.cs-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 48px 24px;
  text-align: center;
}
.cs-empty-state__icon { color: var(--text-tertiary); opacity: 0.6; }
.cs-empty-state__title { font-size: 16px; font-weight: var(--font-semibold); color: var(--text-primary); }
.cs-empty-state__desc  { font-size: 13px; color: var(--text-tertiary); max-width: 320px; line-height: 1.5; }

/* ── Kbd ─────────────────────────────────────────────────────────────────── */
.cs-kbd-group { display: inline-flex; align-items: center; gap: 2px; }
.cs-kbd {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 1px 6px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-bottom-width: 2px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-family: 'Segoe UI', monospace;
  color: var(--text-secondary);
  line-height: 1.6;
}
.cs-kbd__sep { font-size: 10px; color: var(--text-tertiary); padding: 0 1px; }

/* ── Code ────────────────────────────────────────────────────────────────── */
.cs-code {
  display: inline;
  padding: 1px 5px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
  font-size: 0.875em;
  color: var(--copper);
}

/* ── CodeBlock ───────────────────────────────────────────────────────────── */
.cs-code-block {
  background: rgba(0,0,0,0.4);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}
.cs-code-block__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-bottom: 1px solid var(--glass-border);
  background: var(--glass);
}
.cs-code-block__lang {
  font-size: 11px;
  font-weight: var(--font-semibold);
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}
.cs-code-block__copy {
  font-size: 11px;
  color: var(--text-tertiary);
  background: none;
  border: none;
  cursor: pointer;
  padding: 2px 8px;
  border-radius: var(--radius-md);
  transition: color var(--duration-shortest);
}
.cs-code-block__copy:hover { color: var(--text-primary); background: var(--glass); }
.cs-code-block__pre {
  margin: 0;
  padding: 16px;
  overflow-x: auto;
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-primary);
}

/* ── Timeline ────────────────────────────────────────────────────────────── */
.cs-timeline { display: flex; flex-direction: column; }
.cs-timeline-item { display: flex; gap: 16px; }
.cs-timeline-item__track {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
  width: 16px;
}
.cs-timeline-item__dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--glass-border);
  border: 2px solid var(--glass-border-hover);
  flex-shrink: 0;
  margin-top: 4px;
}
.cs-timeline-item__dot--success { background: var(--color-success); border-color: rgba(52,211,153,0.4); }
.cs-timeline-item__dot--danger  { background: var(--color-danger);  border-color: rgba(220,38,38,0.4);  }
.cs-timeline-item__dot--info    { background: var(--color-info);    border-color: rgba(65,105,225,0.4);  }
.cs-timeline-item__dot--warning { background: var(--color-warning); border-color: rgba(234,179,8,0.4);  }
.cs-timeline-item__line {
  flex: 1;
  width: 1px;
  background: var(--border-subtle);
  margin: 4px 0;
  min-height: 16px;
}
.cs-timeline-item:last-child .cs-timeline-item__line { display: none; }
.cs-timeline-item__content { padding-bottom: 20px; flex: 1; }
.cs-timeline-item__time  { font-size: 11px; color: var(--text-tertiary); margin-bottom: 2px; }
.cs-timeline-item__title { font-size: 13px; font-weight: var(--font-medium); color: var(--text-primary); }
.cs-timeline-item__desc  { font-size: 12px; color: var(--text-secondary); margin-top: 2px; line-height: 1.5; }

/* ── Rating ──────────────────────────────────────────────────────────────── */
.cs-rating { display: inline-flex; gap: 2px; align-items: center; }
.cs-rating__star {
  font-size: 20px;
  color: var(--text-tertiary);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: color var(--duration-shortest), transform var(--duration-shortest);
}
.cs-rating__star--filled { color: var(--color-warning); }
.cs-rating[data-rating-readonly] .cs-rating__star { cursor: default; }
.cs-rating__star:not([disabled]):hover { color: var(--color-warning); transform: scale(1.15); }

/* ── Slider ──────────────────────────────────────────────────────────────── */
.cs-slider-wrap { display: flex; flex-direction: column; gap: 8px; }
.cs-slider__header { display: flex; justify-content: space-between; align-items: center; }
.cs-slider__label { font-size: 13px; color: var(--text-secondary); }
.cs-slider__value { font-size: 12px; font-weight: var(--font-semibold); color: var(--text-primary); }
.cs-slider {
  -webkit-appearance: none;
  appearance: none;
  width: 100%;
  height: 4px;
  border-radius: 2px;
  background: var(--glass-border);
  outline: none;
  cursor: pointer;
}
.cs-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--text-primary);
  border: 2px solid var(--bg);
  cursor: pointer;
  transition: transform var(--duration-shortest);
}
.cs-slider::-webkit-slider-thumb:hover { transform: scale(1.2); }
.cs-slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--text-primary);
  border: 2px solid var(--bg);
  cursor: pointer;
}

/* ── NumberInput ─────────────────────────────────────────────────────────── */
.cs-number-input-wrap { display: flex; flex-direction: column; gap: 6px; }
.cs-number-input__label { font-size: 13px; color: var(--text-secondary); }
.cs-number-input {
  display: inline-flex;
  align-items: center;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  overflow: hidden;
  transition: border-color var(--duration-shortest);
}
.cs-number-input:focus-within { border-color: var(--glass-border-hover); }
.cs-number-input__btn {
  padding: 8px 14px;
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  transition: background var(--duration-shortest), color var(--duration-shortest);
  flex-shrink: 0;
}
.cs-number-input__btn:hover { background: var(--glass-hover); color: var(--text-primary); }
.cs-number-input__field {
  width: 60px;
  text-align: center;
  background: none;
  border: none;
  border-left: 1px solid var(--glass-border);
  border-right: 1px solid var(--glass-border);
  color: var(--text-primary);
  font-size: 14px;
  padding: 8px 4px;
  outline: none;
  -moz-appearance: textfield;
}
.cs-number-input__field::-webkit-outer-spin-button,
.cs-number-input__field::-webkit-inner-spin-button { -webkit-appearance: none; }

/* ── FileUpload ──────────────────────────────────────────────────────────── */
.cs-file-upload { display: flex; flex-direction: column; gap: 8px; }
.cs-file-upload__input { display: none; }
.cs-file-upload__zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 32px 24px;
  background: var(--glass);
  border: 2px dashed var(--glass-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--duration-shorter), border-color var(--duration-shorter);
  text-align: center;
}
.cs-file-upload__zone:hover,
.cs-file-upload__zone--drag { background: var(--glass-hover); border-color: var(--glass-border-hover); }
.cs-file-upload__icon { color: var(--text-tertiary); }
.cs-file-upload__text { font-size: 13px; color: var(--text-secondary); font-weight: var(--font-medium); }
.cs-file-upload__hint { font-size: 11px; color: var(--text-tertiary); }
.cs-file-upload__list { display: flex; flex-direction: column; gap: 6px; }
.cs-file-upload__file {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  font-size: 12px;
  color: var(--text-secondary);
}

/* ── TagInput ────────────────────────────────────────────────────────────── */
.cs-tag-input-wrap { display: flex; flex-direction: column; gap: 6px; }
.cs-tag-input__label { font-size: 13px; color: var(--text-secondary); }
.cs-tag-input {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px 12px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: border-color var(--duration-shortest);
  min-height: 42px;
  align-items: center;
}
.cs-tag-input:focus-within { border-color: var(--glass-border-hover); }
.cs-tag-input__tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background: var(--color-info-glass);
  border: 1px solid rgba(65,105,225,0.3);
  border-radius: var(--radius-sm);
  font-size: 12px;
  color: var(--color-info);
}
.cs-tag-input__remove {
  background: none;
  border: none;
  cursor: pointer;
  color: inherit;
  opacity: 0.7;
  font-size: 14px;
  line-height: 1;
  padding: 0;
}
.cs-tag-input__remove:hover { opacity: 1; }
.cs-tag-input__field {
  flex: 1;
  min-width: 80px;
  background: none;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-size: 13px;
  font-family: inherit;
}
.cs-tag-input__field::placeholder { color: var(--text-tertiary); }

/* ── DateInput ───────────────────────────────────────────────────────────── */
.cs-date-input-wrap { display: flex; flex-direction: column; gap: 6px; }
.cs-date-input__label { font-size: 13px; color: var(--text-secondary); }
.cs-date-input__field {
  padding: 8px 12px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  font-family: inherit;
  outline: none;
  transition: border-color var(--duration-shortest);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}
.cs-date-input__field:focus { border-color: var(--glass-border-hover); }
.cs-date-input__field::-webkit-calendar-picker-indicator { filter: invert(0.7); cursor: pointer; }

/* ── Menu ────────────────────────────────────────────────────────────────── */
.cs-menu { position: relative; display: inline-block; }
.cs-menu__arrow { font-size: 10px; opacity: 0.6; }
.cs-menu__dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  min-width: 160px;
  background: rgba(12, 12, 18, 0.92);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  padding: 4px;
  z-index: 200;
  box-shadow: var(--elevation-4);
  display: none;
}
.cs-menu__dropdown.cs-menu--open { display: block; }
.cs-menu__item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 12px;
  background: none;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  font-size: 13px;
  font-family: inherit;
  cursor: pointer;
  text-align: left;
  transition: background var(--duration-shortest), color var(--duration-shortest);
}
.cs-menu__item:hover { background: var(--glass-hover); color: var(--text-primary); }
.cs-menu__item--disabled { opacity: 0.4; cursor: not-allowed; }
.cs-menu__item--disabled:hover { background: none; color: var(--text-secondary); }
.cs-menu__item--color-danger { color: var(--color-danger); }
.cs-menu__item--color-danger:hover { background: var(--color-danger-glass); color: var(--color-danger); }
.cs-menu__item--color-success { color: var(--color-success); }
.cs-menu__item--color-success:hover { background: var(--color-success-glass); color: var(--color-success); }

/* ── Popover ─────────────────────────────────────────────────────────────── */
.cs-popover { position: relative; display: inline-block; }
.cs-popover__panel {
  position: absolute;
  z-index: 200;
  min-width: 200px;
  background: rgba(12, 12, 18, 0.92);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  box-shadow: var(--elevation-4);
  display: none;
}
.cs-popover__panel.cs-popover--open { display: block; }
.cs-popover__panel--bottom { top: calc(100% + 8px); left: 50%; transform: translateX(-50%); }
.cs-popover__panel--top    { bottom: calc(100% + 8px); left: 50%; transform: translateX(-50%); }
.cs-popover__panel--right  { left: calc(100% + 8px); top: 50%; transform: translateY(-50%); }
.cs-popover__panel--left   { right: calc(100% + 8px); top: 50%; transform: translateY(-50%); }
.cs-popover__inner { padding: 16px; }

/* ── Stepper ─────────────────────────────────────────────────────────────── */
.cs-stepper { display: flex; }
.cs-stepper--horizontal { flex-direction: row; align-items: flex-start; }
.cs-stepper--vertical   { flex-direction: column; }
.cs-stepper-step { display: flex; flex: 1; }
.cs-stepper--horizontal .cs-stepper-step { flex-direction: column; align-items: center; text-align: center; }
.cs-stepper--vertical   .cs-stepper-step { flex-direction: row; gap: 16px; align-items: flex-start; }
.cs-stepper-step__track { display: flex; align-items: center; }
.cs-stepper--horizontal .cs-stepper-step__track { flex-direction: row; width: 100%; justify-content: center; margin-bottom: 8px; }
.cs-stepper--vertical   .cs-stepper-step__track { flex-direction: column; align-items: center; }
.cs-stepper-step__circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 2px solid var(--glass-border);
  background: var(--glass);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: var(--font-semibold);
  color: var(--text-tertiary);
  flex-shrink: 0;
  transition: all var(--duration-shorter);
}
.cs-stepper-step--active .cs-stepper-step__circle { border-color: var(--accent); background: var(--accent); color: var(--accent-text); }
.cs-stepper-step--done   .cs-stepper-step__circle { border-color: var(--color-success); background: var(--color-success-glass); color: var(--color-success); }
.cs-stepper-step__connector { flex: 1; height: 1px; background: var(--border-subtle); }
.cs-stepper--vertical .cs-stepper-step__connector { width: 1px; height: auto; min-height: 24px; margin: 4px 0; flex: 1; }
.cs-stepper-step:last-child .cs-stepper-step__connector { display: none; }
.cs-stepper-step__content { padding-bottom: 8px; }
.cs-stepper--vertical .cs-stepper-step__content { padding-bottom: 24px; }
.cs-stepper-step__title { font-size: 13px; font-weight: var(--font-semibold); color: var(--text-secondary); }
.cs-stepper-step--active .cs-stepper-step__title { color: var(--text-primary); }
.cs-stepper-step__desc   { font-size: 12px; color: var(--text-tertiary); margin-top: 2px; }

/* ── Toolbar ─────────────────────────────────────────────────────────────── */
.cs-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 20px;
  background: var(--glass);
  border-radius: var(--radius-md);
}
.cs-toolbar--bordered { border-top: 1px solid var(--glass-border); border-bottom: 1px solid var(--glass-border); border-radius: 0; padding: 8px 24px; }
.cs-toolbar__start { display: flex; align-items: center; gap: 8px; }
.cs-toolbar__end   { display: flex; align-items: center; gap: 8px; }
.cs-toolbar__title { font-size: 18px; font-weight: var(--font-bold); color: var(--text-primary); }

/* ── Chat Widget ─────────────────────────────────────────────────────────── */
.cs-chat-bubble {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: var(--accent);
  color: var(--bg-base);
  border: none;
  cursor: pointer;
  font-size: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1001;
  box-shadow: 0 4px 16px rgba(0,0,0,0.4);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.cs-chat-bubble:hover {
  transform: scale(1.08);
  box-shadow: 0 6px 20px rgba(0,0,0,0.5);
}
.cs-chat-container {
  position: fixed;
  bottom: 88px;
  right: 24px;
  width: 360px;
  height: 520px;
  display: none;
  flex-direction: column;
  background: var(--glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  box-shadow: 0 8px 32px rgba(0,0,0,0.5);
  z-index: 1000;
  overflow: hidden;
  font-family: var(--font-sans);
}
.cs-chat-container.cs-chat--open { display: flex; }
.cs-chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  background: rgba(255,255,255,0.06);
  border-bottom: 1px solid var(--glass-border);
  font-size: 15px;
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}
.cs-chat-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 14px;
  padding: 2px 4px;
  line-height: 1;
  border-radius: var(--radius-sm);
  transition: color 0.15s;
}
.cs-chat-close:hover { color: var(--text-primary); }
.cs-chat-body {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.cs-chat-body::-webkit-scrollbar { width: 4px; }
.cs-chat-body::-webkit-scrollbar-track { background: transparent; }
.cs-chat-body::-webkit-scrollbar-thumb { background: var(--glass-border); border-radius: 2px; }
.cs-chat-msg {
  max-width: 82%;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
}
.cs-chat-msg--user {
  align-self: flex-end;
  background: var(--accent);
  color: var(--bg-base);
  border-bottom-right-radius: var(--radius-sm);
}
.cs-chat-msg--bot {
  align-self: flex-start;
  background: rgba(255,255,255,0.08);
  color: var(--text-primary);
  border: 1px solid var(--glass-border);
  border-bottom-left-radius: var(--radius-sm);
}
.cs-chat-typing {
  align-self: flex-start;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 10px 14px;
  background: rgba(255,255,255,0.08);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  border-bottom-left-radius: var(--radius-sm);
}
.cs-chat-typing span {
  width: 7px;
  height: 7px;
  background: var(--text-muted);
  border-radius: 50%;
  animation: cs-chat-bounce 1.4s infinite;
}
.cs-chat-typing span:nth-child(2) { animation-delay: 0.2s; }
.cs-chat-typing span:nth-child(3) { animation-delay: 0.4s; }
@keyframes cs-chat-bounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.5; }
  30%           { transform: translateY(-6px); opacity: 1; }
}
.cs-chat-footer {
  display: flex;
  gap: 8px;
  padding: 12px;
  border-top: 1px solid var(--glass-border);
  background: rgba(255,255,255,0.03);
}
.cs-chat-input {
  flex: 1;
  background: rgba(255,255,255,0.06);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  padding: 8px 12px;
  outline: none;
  font-family: var(--font-sans);
  transition: border-color 0.15s;
}
.cs-chat-input::placeholder { color: var(--text-muted); }
.cs-chat-input:focus { border-color: var(--accent); }
.cs-chat-send {
  background: var(--accent);
  color: var(--bg-base);
  border: none;
  border-radius: var(--radius-sm);
  padding: 8px 14px;
  font-size: 13px;
  font-weight: var(--font-semibold);
  cursor: pointer;
  font-family: var(--font-sans);
  transition: opacity 0.15s;
}
.cs-chat-send:hover { opacity: 0.85; }

/* ── Form ────────────────────────────────────────────────────────────────── */
.cs-form { display: flex; flex-direction: column; gap: 16px; }

/* ── Sidebar ─────────────────────────────────────────────────────────────── */
.cs-sidebar {
  display: flex;
  flex-direction: column;
  width: 280px;
  height: 100vh;
  background: var(--glass);
  border-right: 1px solid var(--glass-border);
  padding: 16px 12px;
  flex-shrink: 0;
  overflow: hidden;
}
.cs-sidebar__brand {
  font-size: 15px;
  font-weight: var(--font-bold);
  color: var(--accent);
  padding: 8px 12px 16px;
  letter-spacing: 0.02em;
}
.cs-sidebar__nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
.cs-sidebar__nav .cs-nav-link {
  display: block;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: background var(--duration-shortest), color var(--duration-shortest);
}
.cs-sidebar__nav .cs-nav-link:hover { background: var(--glass-hover); color: var(--text-primary); }

/* ── Section ─────────────────────────────────────────────────────────────── */
.cs-section { display: flex; flex-direction: column; gap: 16px; }
.cs-section__header { display: flex; flex-direction: column; gap: 4px; }
.cs-section__title { font-size: 16px; font-weight: var(--font-semibold); color: var(--text-primary); margin: 0; }
.cs-section__desc  { font-size: 13px; color: var(--text-secondary); margin: 0; }
.cs-section__body  { display: flex; flex-direction: column; gap: 12px; }

/* ── Callout ─────────────────────────────────────────────────────────────── */
.cs-callout {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 14px 16px;
  border-radius: var(--radius-md);
  border-left: 3px solid;
  font-size: 13px;
  line-height: 1.5;
}
.cs-callout--info    { background: var(--color-info-glass);    border-color: var(--color-info);    color: var(--color-info); }
.cs-callout--warning { background: var(--color-warning-glass); border-color: var(--color-warning); color: var(--color-warning); }
.cs-callout--tip     { background: var(--color-success-glass); border-color: var(--color-success); color: var(--color-success); }
.cs-callout--danger  { background: var(--color-danger-glass);  border-color: var(--color-danger);  color: var(--color-danger); }
.cs-callout__icon { font-size: 16px; flex-shrink: 0; line-height: 1.4; }
.cs-callout__body { display: flex; flex-direction: column; gap: 4px; color: var(--text-primary); }
.cs-callout__title { font-weight: var(--font-semibold); font-size: 13px; }

/* ── Image ───────────────────────────────────────────────────────────────── */
.cs-image { display: block; max-width: 100%; height: auto; }
.cs-image--rounded { border-radius: var(--radius-md); }

/* ── Link ────────────────────────────────────────────────────────────────── */
.cs-link { color: var(--accent); text-decoration: none; transition: opacity var(--duration-shortest); font-size: inherit; }
.cs-link:hover { opacity: 0.8; text-decoration: underline; }
.cs-link--muted  { color: var(--text-secondary); }
.cs-link--danger { color: var(--color-danger); }

/* ── Search ──────────────────────────────────────────────────────────────── */
.cs-search {
  position: relative;
  display: flex;
  align-items: center;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: border-color var(--duration-shortest);
}
.cs-search:focus-within { border-color: var(--glass-border-hover); }
.cs-search__icon {
  padding: 0 10px 0 12px;
  color: var(--text-tertiary);
  font-size: 14px;
  pointer-events: none;
  flex-shrink: 0;
}
.cs-search__input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-size: 13px;
  font-family: inherit;
  padding: 9px 0;
  -webkit-appearance: none;
}
.cs-search__input::placeholder { color: var(--text-tertiary); }
.cs-search__input::-webkit-search-cancel-button { display: none; }
.cs-search__clear {
  padding: 0 12px 0 8px;
  background: none;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  font-size: 12px;
  flex-shrink: 0;
  transition: color var(--duration-shortest);
  line-height: 1;
}
.cs-search__clear:hover { color: var(--text-primary); }

/* ── ColorInput ──────────────────────────────────────────────────────────── */
.cs-color-input-wrap { display: flex; flex-direction: column; gap: 6px; }
.cs-color-input__label { font-size: 13px; color: var(--text-secondary); }
.cs-color-input {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  transition: border-color var(--duration-shortest);
}
.cs-color-input:focus-within { border-color: var(--glass-border-hover); }
.cs-color-input__field {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  padding: 0;
  background: none;
  flex-shrink: 0;
}
.cs-color-input__field::-webkit-color-swatch-wrapper { padding: 0; }
.cs-color-input__field::-webkit-color-swatch { border: none; border-radius: var(--radius-sm); }
.cs-color-input__hex { font-size: 12px; color: var(--text-secondary); font-family: var(--font-mono); letter-spacing: 0.05em; }

/* ── Snackbar ────────────────────────────────────────────────────────────── */
.cs-snackbar {
  position: fixed;
  z-index: 9000;
  display: flex;
  flex-direction: column;
  gap: 8px;
  pointer-events: none;
}
.cs-snackbar--bottom-right  { bottom: 24px; right: 24px; align-items: flex-end; }
.cs-snackbar--bottom-left   { bottom: 24px; left: 24px; align-items: flex-start; }
.cs-snackbar--top-right     { top: 24px; right: 24px; align-items: flex-end; }
.cs-snackbar--top-left      { top: 24px; left: 24px; align-items: flex-start; }
.cs-snackbar--bottom-center { bottom: 24px; left: 50%; transform: translateX(-50%); align-items: center; }
.cs-snackbar--center { top: 50%; left: 50%; transform: translate(-50%, -50%); align-items: center; }

/* ── Confirm ─────────────────────────────────────────────────────────────── */
.cs-confirm__header { padding: 20px 24px 0; }
.cs-confirm__message {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0 0 8px;
}
.cs-confirm__actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  padding-top: 8px;
}
.cs-confirm__btn--danger  { background: var(--color-danger); color: #fff; border-color: var(--color-danger); }
.cs-confirm__btn--warning { background: var(--color-warning); color: #000; border-color: var(--color-warning); }

/* ── Carousel ────────────────────────────────────────────────────────────── */
.cs-carousel {
  position: relative;
  overflow: hidden;
  border-radius: var(--radius-md);
}
.cs-carousel__track {
  display: flex;
  overflow-x: scroll;
  scroll-snap-type: x mandatory;
  scrollbar-width: none;
  gap: 0;
}
.cs-carousel__track::-webkit-scrollbar { display: none; }
.cs-carousel__slide {
  flex: 0 0 100%;
  scroll-snap-align: start;
}
.cs-carousel__btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid var(--glass-border);
  color: var(--text-primary);
  border-radius: 50%;
  width: 36px;
  height: 36px;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background var(--duration-shortest);
  z-index: 10;
}
.cs-carousel__btn:hover { background: rgba(0,0,0,0.75); }
.cs-carousel__btn--prev { left: 12px; }
.cs-carousel__btn--next { right: 12px; }

/* ── Banner ──────────────────────────────────────────────────────────────── */
.cs-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 20px;
  font-size: 13px;
  line-height: 1.5;
  position: relative;
}
.cs-banner--info    { background: var(--color-info-glass);    border-bottom: 1px solid rgba(65,105,225,0.25);    color: var(--color-info); }
.cs-banner--success { background: var(--color-success-glass); border-bottom: 1px solid rgba(34,197,94,0.25);   color: var(--color-success); }
.cs-banner--warning { background: var(--color-warning-glass); border-bottom: 1px solid rgba(234,179,8,0.25);   color: var(--color-warning); }
.cs-banner--danger  { background: var(--color-danger-glass);  border-bottom: 1px solid rgba(239,68,68,0.25);   color: var(--color-danger); }
.cs-banner__icon { font-size: 16px; flex-shrink: 0; margin-top: 1px; }
.cs-banner__body { flex: 1; color: var(--text-primary); }
.cs-banner__title { color: inherit; }
.cs-banner__dismiss {
  background: none;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  font-size: 12px;
  padding: 0;
  line-height: 1;
  flex-shrink: 0;
  transition: color var(--duration-shortest);
  margin-top: 2px;
}
.cs-banner__dismiss:hover { color: var(--text-primary); }

/* ── FormField ───────────────────────────────────────────────────────────── */
.cs-form-field { display: flex; flex-direction: column; gap: 6px; }
.cs-form-field__label { font-size: 13px; font-weight: var(--font-medium); color: var(--text-secondary); }
.cs-form-field__required { color: var(--color-danger); margin-left: 2px; }
.cs-form-field__hint { font-size: 12px; color: var(--text-tertiary); }
.cs-form-field__hint--error { color: var(--color-danger); }
.cs-form-field--error .cs-form-field__label { color: var(--color-danger); }

/* ── KvList ──────────────────────────────────────────────────────────────── */
.cs-kv-list { display: flex; flex-direction: column; margin: 0; padding: 0; }
.cs-kv-list--divided .cs-kv-item { border-bottom: 1px solid var(--glass-border); }
.cs-kv-list--divided .cs-kv-item:last-child { border-bottom: none; }
.cs-kv-item {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 16px;
  padding: 10px 0;
}
.cs-kv-item__key {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: var(--font-normal);
  flex-shrink: 0;
}
.cs-kv-item__value {
  font-size: 13px;
  color: var(--text-primary);
  font-weight: var(--font-medium);
  text-align: right;
  margin: 0;
}
.cs-kv-item__value--success { color: var(--color-success); }
.cs-kv-item__value--warning { color: var(--color-warning); }
.cs-kv-item__value--danger  { color: var(--color-danger); }
.cs-kv-item__value--muted   { color: var(--text-tertiary); }
.cs-kv-item__link { color: var(--accent); text-decoration: none; }
.cs-kv-item__link:hover { text-decoration: underline; }

/* ── ButtonGroup ─────────────────────────────────────────────────────────── */
.cs-button-group {
  display: inline-flex;
  align-items: center;
}
.cs-button-group .cs-button {
  border-radius: 0;
  margin-left: -1px;
}
.cs-button-group .cs-button:first-child { border-radius: var(--radius-sm) 0 0 var(--radius-sm); margin-left: 0; }
.cs-button-group .cs-button:last-child  { border-radius: 0 var(--radius-sm) var(--radius-sm) 0; }
.cs-button-group .cs-button:only-child  { border-radius: var(--radius-sm); margin-left: 0; }
.cs-button-group .cs-button:hover,
.cs-button-group .cs-button--active,
.cs-button-group .cs-button[data-active] { position: relative; z-index: 1; }

/* ── CopyButton ──────────────────────────────────────────────────────────── */
.cs-copy-button { transition: all var(--duration-shortest); }

/* ── MultiSelect ─────────────────────────────────────────────────────────── */
.cs-multi-select { display: flex; flex-direction: column; gap: 6px; position: relative; }
.cs-multi-select__label { font-size: 13px; color: var(--text-secondary); }
.cs-multi-select__control {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-height: 42px;
  padding: 6px 12px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: border-color var(--duration-shortest);
}
.cs-multi-select--open .cs-multi-select__control { border-color: var(--glass-border-hover); }
.cs-multi-select__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  flex: 1;
}
.cs-multi-select__placeholder { font-size: 13px; color: var(--text-tertiary); }
.cs-multi-select__tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background: var(--color-info-glass);
  border: 1px solid rgba(65,105,225,0.3);
  border-radius: var(--radius-sm);
  font-size: 12px;
  color: var(--color-info);
}
.cs-multi-select__tag button {
  background: none;
  border: none;
  cursor: pointer;
  color: inherit;
  opacity: 0.7;
  font-size: 14px;
  line-height: 1;
  padding: 0;
}
.cs-multi-select__tag button:hover { opacity: 1; }
.cs-multi-select__arrow { font-size: 10px; opacity: 0.5; flex-shrink: 0; }
.cs-multi-select__dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: rgba(12,12,18,0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  padding: 4px;
  z-index: 300;
  box-shadow: var(--elevation-4);
  max-height: 240px;
  overflow-y: auto;
}
.cs-multi-select__option {
  padding: 8px 12px;
  font-size: 13px;
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--duration-shortest), color var(--duration-shortest);
}
.cs-multi-select__option:hover { background: var(--glass-hover); color: var(--text-primary); }
.cs-multi-select__option--active {
  background: var(--color-info-glass);
  color: var(--color-info);
}
.cs-multi-select__option--active::after { content: " ✓"; float: right; }

/* ── Chart ───────────────────────────────────────────────────────────────── */
.cs-chart { display: flex; flex-direction: column; gap: 8px; }
.cs-chart__title { font-size: 13px; font-weight: var(--font-semibold); color: var(--text-secondary); }
.cs-chart__svg { width: 100%; height: auto; display: block; }
.cs-chart--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 120px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  color: var(--text-tertiary);
  font-size: 13px;
}
.cs-chart__grid { stroke: var(--glass-border); stroke-width: 1; }
.cs-chart__axis-label { fill: var(--text-tertiary); font-size: 11px; font-family: var(--font-mono); }
.cs-chart__label { fill: var(--text-tertiary); font-size: 11px; }
.cs-chart__value { fill: var(--text-secondary); font-size: 11px; font-weight: 600; }
.cs-chart__bar { opacity: 0.9; transition: opacity 0.15s; }
.cs-chart__bar:hover { opacity: 1; }
.cs-chart__line { stroke: var(--accent); stroke-width: 2; }
.cs-chart__area { fill: var(--accent); opacity: 0.12; }
.cs-chart__dot { fill: var(--accent); stroke: var(--bg); stroke-width: 2; }
.cs-chart__slice { opacity: 0.9; stroke: var(--bg); stroke-width: 1; transition: opacity 0.15s; }
.cs-chart__slice:hover { opacity: 1; }
.cs-chart__legend-label { fill: var(--text-secondary); font-size: 11px; dominant-baseline: middle; }

/* ── IconButton ──────────────────────────────────────────────────────────── */
.cs-icon-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  border: 1px solid transparent;
  cursor: pointer;
  transition: background var(--duration-shortest), color var(--duration-shortest), border-color var(--duration-shortest);
  flex-shrink: 0;
}
.cs-icon-button--sm { width: 28px; height: 28px; }
.cs-icon-button--md { width: 34px; height: 34px; }
.cs-icon-button--lg { width: 42px; height: 42px; }
.cs-icon-button--ghost { background: none; color: var(--text-tertiary); }
.cs-icon-button--ghost:hover { background: var(--glass-hover); color: var(--text-primary); }
.cs-icon-button--outline { background: none; border-color: var(--glass-border); color: var(--text-secondary); }
.cs-icon-button--outline:hover { border-color: var(--glass-border-hover); color: var(--text-primary); }
.cs-icon-button--solid { background: var(--glass); color: var(--text-primary); border-color: var(--glass-border); }
.cs-icon-button--solid:hover { background: var(--glass-hover); }
.cs-icon-button--color-danger { color: var(--color-danger); }
.cs-icon-button--color-danger:hover { background: var(--color-danger-glass); }

/* ── Tag ─────────────────────────────────────────────────────────────────── */
.cs-tag {
  display: inline-flex;
  align-items: center;
  font-weight: var(--font-medium);
  border-radius: var(--radius-sm);
  letter-spacing: 0.02em;
  white-space: nowrap;
}
.cs-tag--size-xs { font-size: 10px; padding: 1px 6px; }
.cs-tag--size-sm { font-size: 11px; padding: 2px 8px; }
.cs-tag--size-md { font-size: 12px; padding: 3px 10px; }
.cs-tag--default { background: var(--glass); color: var(--text-secondary); border: 1px solid var(--glass-border); }
.cs-tag--cyan    { background: rgba(0,180,216,0.12); color: #00b4d8; border: 1px solid rgba(0,180,216,0.25); }
.cs-tag--purple  { background: rgba(139,92,246,0.12); color: #8b5cf6; border: 1px solid rgba(139,92,246,0.25); }
.cs-tag--orange  { background: rgba(249,115,22,0.12); color: #f97316; border: 1px solid rgba(249,115,22,0.25); }
.cs-tag--pink    { background: rgba(236,72,153,0.12); color: #ec4899; border: 1px solid rgba(236,72,153,0.25); }
.cs-tag--green   { background: var(--color-success-glass); color: var(--color-success); border: 1px solid rgba(34,197,94,0.25); }
.cs-tag--red     { background: var(--color-danger-glass);  color: var(--color-danger);  border: 1px solid rgba(239,68,68,0.25); }
.cs-tag--yellow  { background: var(--color-warning-glass); color: var(--color-warning); border: 1px solid rgba(234,179,8,0.25); }

/* ── DataGrid ────────────────────────────────────────────────────────────── */
.cs-data-grid { display: flex; flex-direction: column; gap: 8px; }
.cs-data-grid__filter { display: flex; }
.cs-data-grid__filter-input {
  flex: 1;
  padding: 8px 12px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  font-family: inherit;
  outline: none;
  transition: border-color var(--duration-shortest);
}
.cs-data-grid__filter-input:focus { border-color: var(--glass-border-hover); }
.cs-data-grid__wrap { overflow-x: auto; }
.cs-data-grid__table { width: 100%; }
.cs-data-grid__th { cursor: pointer; user-select: none; white-space: nowrap; }
.cs-data-grid__th:hover { color: var(--text-primary); }
.cs-data-grid__sort-icon { opacity: 0.3; font-size: 10px; }
.cs-data-grid__th--asc  .cs-data-grid__sort-icon,
.cs-data-grid__th--desc .cs-data-grid__sort-icon { opacity: 1; color: var(--accent); }
.cs-data-grid__empty { padding: 24px; text-align: center; color: var(--text-tertiary); font-size: 13px; display: none; }

/* ── Tree ────────────────────────────────────────────────────────────────── */
.cs-tree { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 1px; }
.cs-tree-item { list-style: none; }
.cs-tree-item__row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: background var(--duration-shortest), color var(--duration-shortest);
  user-select: none;
}
.cs-tree-item__row:hover { background: var(--glass-hover); color: var(--text-primary); }
.cs-tree-item--active > .cs-tree-item__row { background: var(--glass); color: var(--text-primary); }
.cs-tree-item__chevron { font-size: 9px; opacity: 0.5; transition: transform var(--duration-shortest); flex-shrink: 0; }
.cs-tree-item--open > .cs-tree-item__row .cs-tree-item__chevron { transform: rotate(0deg); }
.cs-tree-item:not(.cs-tree-item--open) > .cs-tree-item__row .cs-tree-item__chevron { transform: rotate(-90deg); }
.cs-tree-item__label { flex: 1; }
.cs-tree-item__children { list-style: none; margin: 0; padding-left: 18px; padding-top: 1px; }

/* ── VirtualList ─────────────────────────────────────────────────────────── */
.cs-virtual-list { border: 1px solid var(--glass-border); border-radius: var(--radius-md); background: var(--glass); }
.cs-virtual-list::-webkit-scrollbar { width: 6px; }
.cs-virtual-list::-webkit-scrollbar-track { background: transparent; }
.cs-virtual-list::-webkit-scrollbar-thumb { background: var(--glass-border); border-radius: 3px; }
.cs-virtual-list__row { display: flex; align-items: center; border-bottom: 1px solid var(--glass-border); padding: 0 16px; box-sizing: border-box; }
.cs-virtual-list__row:last-child { border-bottom: none; }
.cs-virtual-list__cell { flex: 1; font-size: 13px; color: var(--text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; padding-right: 16px; }

/* ── Notification ────────────────────────────────────────────────────────── */
.cs-notification { position: relative; display: inline-block; }
.cs-notification__bell {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 6px;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  position: relative;
  transition: color var(--duration-shortest);
  line-height: 1;
}
.cs-notification__bell:hover { color: var(--text-primary); }
.cs-notification__badge {
  position: absolute;
  top: 2px;
  right: 2px;
  background: var(--color-danger);
  color: #fff;
  font-size: 10px;
  font-weight: var(--font-bold);
  line-height: 1;
  padding: 2px 4px;
  border-radius: 999px;
  min-width: 16px;
  text-align: center;
}
.cs-notification__panel {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  width: 320px;
  background: rgba(12,12,18,0.96);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  box-shadow: var(--elevation-4);
  z-index: 500;
  overflow: hidden;
}
.cs-notification__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--glass-border);
}
.cs-notification__title { font-size: 13px; font-weight: var(--font-semibold); color: var(--text-primary); }
.cs-notification__mark-all { background: none; border: none; font-size: 12px; color: var(--accent); cursor: pointer; padding: 0; }
.cs-notification__list { max-height: 360px; overflow-y: auto; }
.cs-notification-item {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  padding: 12px 16px;
  cursor: pointer;
  transition: background var(--duration-shortest);
  border-bottom: 1px solid var(--glass-border);
}
.cs-notification-item:last-child { border-bottom: none; }
.cs-notification-item:hover { background: var(--glass-hover); }
.cs-notification-item__dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: transparent;
  flex-shrink: 0;
  margin-top: 5px;
}
.cs-notification-item--unread .cs-notification-item__dot { background: var(--accent); }
.cs-notification-item__content { flex: 1; }
.cs-notification-item__title { font-size: 13px; font-weight: var(--font-medium); color: var(--text-primary); }
.cs-notification-item--unread .cs-notification-item__title { font-weight: var(--font-semibold); }
.cs-notification-item__body { font-size: 12px; color: var(--text-secondary); margin-top: 2px; }
.cs-notification-item__time { font-size: 11px; color: var(--text-tertiary); margin-top: 4px; }

/* ── Command ─────────────────────────────────────────────────────────────── */
.cs-command {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 15vh;
  z-index: 9999;
}
.cs-command__backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}
.cs-command__dialog {
  position: relative;
  width: 100%;
  max-width: 560px;
  background: rgba(12,12,18,0.98);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-sm);
  box-shadow: 0 20px 60px rgba(0,0,0,0.7);
  overflow: hidden;
  z-index: 1;
}
.cs-command__search-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--glass-border);
}
.cs-command__search-icon { font-size: 16px; color: var(--text-tertiary); flex-shrink: 0; }
.cs-command__input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  font-size: 15px;
  color: var(--text-primary);
  font-family: inherit;
}
.cs-command__input::placeholder { color: var(--text-tertiary); }
.cs-command__results { max-height: 320px; overflow-y: auto; padding: 4px; }
.cs-command__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--duration-shortest);
}
.cs-command__item:hover,
.cs-command__item--active { background: var(--glass-hover); }
.cs-command__item-icon { width: 20px; display: flex; align-items: center; justify-content: center; color: var(--text-tertiary); flex-shrink: 0; }
.cs-command__item-body { display: flex; flex-direction: column; flex: 1; }
.cs-command__item-label { font-size: 13px; color: var(--text-primary); font-weight: var(--font-medium); }
.cs-command__item-desc { font-size: 11px; color: var(--text-tertiary); margin-top: 1px; }

/* ── ContextMenu ─────────────────────────────────────────────────────────── */
.cs-context-menu__inner {
  min-width: 160px;
  background: rgba(12,12,18,0.96);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  padding: 4px;
  box-shadow: var(--elevation-4);
}

/* ── HoverCard ───────────────────────────────────────────────────────────── */
.cs-hover-card { position: relative; display: inline-block; }
.cs-hover-card > *:not(:first-child) {
  position: absolute;
  min-width: 220px;
  background: rgba(12,12,18,0.96);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  box-shadow: var(--elevation-4);
  z-index: 400;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  transition: opacity var(--duration-shortest), visibility var(--duration-shortest);
}
.cs-hover-card:hover > *:not(:first-child) { opacity: 1; visibility: visible; pointer-events: auto; }
.cs-hover-card--top    > *:not(:first-child) { bottom: calc(100% + 8px); left: 50%; transform: translateX(-50%); }
.cs-hover-card--bottom > *:not(:first-child) { top: calc(100% + 8px); left: 50%; transform: translateX(-50%); }
.cs-hover-card--right  > *:not(:first-child) { left: calc(100% + 8px); top: 0; }
.cs-hover-card--left   > *:not(:first-child) { right: calc(100% + 8px); top: 0; }

/* ── SplitView ───────────────────────────────────────────────────────────── */
.cs-split-view { display: flex; width: 100%; height: 100%; overflow: hidden; }
.cs-split-view--vertical { flex-direction: column; }
.cs-split-view__pane { overflow: auto; min-width: 0; min-height: 0; flex-shrink: 0; }
.cs-split-view__pane--first { flex-shrink: 0; }
.cs-split-view__pane:last-child { flex: 1; }
.cs-split-view__divider {
  flex-shrink: 0;
  background: var(--glass-border);
  transition: background var(--duration-shortest);
  cursor: col-resize;
}
.cs-split-view--horizontal .cs-split-view__divider { width: 4px; cursor: col-resize; }
.cs-split-view--vertical   .cs-split-view__divider { height: 4px; cursor: row-resize; }
.cs-split-view__divider:hover { background: var(--accent); }

/* ── Calendar ────────────────────────────────────────────────────────────── */
.cs-calendar {
  display: inline-flex;
  flex-direction: column;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  padding: 16px;
  min-width: 280px;
  user-select: none;
}
.cs-calendar__header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.cs-calendar__nav {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 20px;
  padding: 2px 8px;
  border-radius: var(--radius-md);
  line-height: 1;
  transition: background var(--duration-shortest), color var(--duration-shortest);
}
.cs-calendar__nav:hover { background: var(--glass-hover); color: var(--text-primary); }
.cs-calendar__label { font-size: 14px; font-weight: var(--font-semibold); color: var(--text-primary); }
.cs-calendar__days-header { display: grid; grid-template-columns: repeat(7,1fr); gap: 2px; margin-bottom: 4px; }
.cs-calendar__days-header span { text-align: center; font-size: 11px; color: var(--text-tertiary); padding: 4px 0; }
.cs-calendar__days { display: grid; grid-template-columns: repeat(7,1fr); gap: 2px; }
.cs-calendar__day {
  text-align: center;
  padding: 6px 4px;
  font-size: 13px;
  color: var(--text-secondary);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--duration-shortest), color var(--duration-shortest);
}
.cs-calendar__day:hover { background: var(--glass-hover); color: var(--text-primary); }
.cs-calendar__day--today { color: var(--accent); font-weight: var(--font-semibold); }
.cs-calendar__day--selected { background: var(--accent); color: var(--accent-text) !important; font-weight: var(--font-semibold); }
.cs-calendar__day--empty { cursor: default; }

/* ── Video / Audio / Iframe ──────────────────────────────────────────────── */
.cs-video { width: 100%; height: auto; display: block; border-radius: var(--radius-md); background: #000; }
.cs-audio { width: 100%; display: block; }
.cs-iframe { width: 100%; display: block; border: none; border-radius: var(--radius-md); }

/* ── AspectRatio ─────────────────────────────────────────────────────────── */
.cs-aspect-ratio { width: 100%; overflow: hidden; }
.cs-aspect-ratio > * { width: 100%; height: 100%; object-fit: cover; }

/* ── RichText ────────────────────────────────────────────────────────────── */
.cs-rich-text { color: var(--text-primary); line-height: 1.7; font-size: 14px; }
.cs-rich-text h1,.cs-rich-text h2,.cs-rich-text h3,.cs-rich-text h4 { color: var(--text-primary); margin: 1.2em 0 0.5em; font-weight: var(--font-semibold); }
.cs-rich-text p  { margin: 0 0 0.8em; color: var(--text-secondary); }
.cs-rich-text a  { color: var(--accent); text-decoration: none; }
.cs-rich-text a:hover { text-decoration: underline; }
.cs-rich-text ul,.cs-rich-text ol { padding-left: 1.4em; margin: 0 0 0.8em; color: var(--text-secondary); }
.cs-rich-text code { background: var(--glass); padding: 2px 6px; border-radius: var(--radius-sm); font-family: var(--font-mono); font-size: 0.9em; }
.cs-rich-text blockquote { border-left: 3px solid var(--accent); margin: 0 0 0.8em; padding: 4px 16px; color: var(--text-secondary); }
`
