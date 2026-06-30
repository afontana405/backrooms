package engine

// atomCSS returns CSS for ChefScript Atoms.
// Appended to the component CSS pipeline.
const atomCSS = `
/* ── Floating Action Button ──────────────────────────────────────────── */
.cs-fab {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
  text-decoration: none;
  font-family: var(--font-family);
  font-weight: 500;
  transition: background var(--duration-shorter) var(--ease-standard),
              box-shadow var(--duration-shorter) var(--ease-standard),
              transform var(--duration-shortest) var(--ease-standard);
}
.cs-fab:active { transform: scale(0.95); }
.cs-fab--circular { border-radius: var(--radius-full); }
.cs-fab--extended { border-radius: var(--radius-chip); gap: var(--spacing-2); padding: 0 var(--spacing-4); }
.cs-fab--small { width: 40px; height: 40px; }
.cs-fab--small.cs-fab--extended { width: auto; height: 34px; font-size: 0.8125rem; }
.cs-fab--medium { width: 56px; height: 56px; }
.cs-fab--medium.cs-fab--extended { width: auto; height: 48px; font-size: 0.875rem; }
.cs-fab--large { width: 72px; height: 72px; }
.cs-fab--large.cs-fab--extended { width: auto; height: 56px; font-size: 0.9375rem; }
.cs-fab--primary { background: var(--color-primary); color: var(--color-primary-contrast); box-shadow: var(--elevation-4); }
.cs-fab--primary:hover { box-shadow: var(--elevation-6); }
.cs-fab--secondary { background: var(--color-secondary); color: #fff; box-shadow: var(--elevation-4); }
.cs-fab--secondary:hover { box-shadow: var(--elevation-6); }
.cs-fab--success { background: var(--color-success); color: #fff; box-shadow: var(--elevation-4); }
.cs-fab--error { background: var(--color-error); color: #fff; box-shadow: var(--elevation-4); }
.cs-fab--warning { background: var(--color-warning); color: #fff; box-shadow: var(--elevation-4); }
.cs-fab--info { background: var(--color-info); color: #fff; box-shadow: var(--elevation-4); }
.cs-fab__label { white-space: nowrap; }
.cs-fab:disabled { opacity: 0.38; cursor: not-allowed; box-shadow: var(--elevation-0); }

/* ── Toggle Button Group ─────────────────────────────────────────────── */
.cs-toggle-group {
  display: inline-flex;
  border: 1px solid var(--color-divider);
  border-radius: var(--radius-sm);
  overflow: hidden;
}
.cs-toggle-group--vertical { flex-direction: column; }
.cs-toggle-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-1);
  padding: var(--spacing-2) var(--spacing-3);
  border: none;
  border-right: 1px solid var(--color-divider);
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  font-family: var(--font-family);
  font-size: 0.875rem;
  transition: background var(--duration-shortest) var(--ease-standard),
              color var(--duration-shortest) var(--ease-standard);
}
.cs-toggle-group--vertical .cs-toggle-btn { border-right: none; border-bottom: 1px solid var(--color-divider); }
.cs-toggle-btn:last-child { border-right: none; border-bottom: none; }
.cs-toggle-btn:hover { background: var(--color-hover); }
.cs-toggle-btn--selected { background: var(--color-selected); color: var(--color-primary); }
.cs-toggle-btn--selected:hover { background: var(--color-selected); }
.cs-toggle-btn:disabled { opacity: 0.38; cursor: not-allowed; }
.cs-toggle-group--size-small .cs-toggle-btn { padding: var(--spacing-1) var(--spacing-2); font-size: 0.75rem; }
.cs-toggle-group--size-large .cs-toggle-btn { padding: var(--spacing-3) var(--spacing-4); font-size: 1rem; }

/* ── Transfer List ───────────────────────────────────────────────────── */
.cs-transfer {
  display: flex;
  align-items: stretch;
  gap: var(--spacing-3);
}
.cs-transfer__panel {
  flex: 1;
  border: 1px solid var(--color-divider);
  border-radius: var(--radius-md);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.cs-transfer__header {
  padding: var(--spacing-2) var(--spacing-3);
  font-weight: 600;
  font-size: 0.875rem;
  border-bottom: 1px solid var(--color-divider);
  background: var(--color-hover);
}
.cs-transfer__list {
  flex: 1;
  overflow-y: auto;
  max-height: 300px;
  min-height: 200px;
}
.cs-transfer__item {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-3);
  cursor: pointer;
  font-size: 0.875rem;
  transition: background var(--duration-shortest) var(--ease-standard);
}
.cs-transfer__item:hover { background: var(--color-hover); }
.cs-transfer__check { flex-shrink: 0; }
.cs-transfer__controls {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: var(--spacing-2);
}
.cs-transfer__btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-divider);
  border-radius: var(--radius-full);
  background: transparent;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background var(--duration-shortest) var(--ease-standard);
}
.cs-transfer__btn:hover { background: var(--color-hover); }

/* ── Paper ───────────────────────────────────────────────────────────── */
.cs-paper {
  background: var(--color-surface);
  color: var(--text-primary);
  border-radius: var(--radius-sm);
  transition: box-shadow var(--duration-standard) var(--ease-standard);
}
.cs-paper--square { border-radius: 0; }
.cs-paper--outlined { border: 1px solid var(--color-divider); box-shadow: none; }
.cs-paper--elevation-0 { box-shadow: var(--elevation-0); }
.cs-paper--elevation-1 { box-shadow: var(--elevation-1); }
.cs-paper--elevation-2 { box-shadow: var(--elevation-2); }
.cs-paper--elevation-3 { box-shadow: var(--elevation-3); }
.cs-paper--elevation-4 { box-shadow: var(--elevation-4); }
.cs-paper--elevation-6 { box-shadow: var(--elevation-6); }
.cs-paper--elevation-8 { box-shadow: var(--elevation-8); }
.cs-paper--elevation-12 { box-shadow: var(--elevation-12); }
.cs-paper--elevation-24 { box-shadow: var(--elevation-24); }

/* ── App Bar ─────────────────────────────────────────────────────────── */
.cs-app-bar {
  display: flex;
  flex-direction: column;
  width: 100%;
  box-sizing: border-box;
  flex-shrink: 0;
  z-index: 50;
}
.cs-app-bar--position-fixed { position: fixed; top: 0; left: 0; right: 0; }
.cs-app-bar--position-absolute { position: absolute; top: 0; left: 0; right: 0; }
.cs-app-bar--position-sticky { position: sticky; top: 0; }
.cs-app-bar--position-static { position: static; }
.cs-app-bar--position-relative { position: relative; }
.cs-app-bar--color-primary { background: var(--color-primary); color: var(--color-primary-contrast); }
.cs-app-bar--color-secondary { background: var(--color-secondary); color: #fff; }
.cs-app-bar--color-transparent { background: transparent; color: inherit; }
.cs-app-bar--color-inherit { background: inherit; color: inherit; }

/* ── Backdrop ────────────────────────────────────────────────────────── */
.cs-backdrop {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0,0,0,0.5);
  z-index: var(--z-modal);
  opacity: 0;
  pointer-events: none;
  transition: opacity var(--duration-shorter) var(--ease-standard);
}
.cs-backdrop--open { opacity: 1; pointer-events: auto; }
.cs-backdrop--invisible { background: transparent; }

/* ── Dialog ──────────────────────────────────────────────────────────── */
.cs-dialog {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
}
.cs-dialog__backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5);
}
.cs-dialog__paper {
  position: relative;
  z-index: 1;
  background: var(--color-surface);
  color: var(--text-primary);
  border-radius: var(--radius-sm);
  box-shadow: var(--elevation-24);
  display: flex;
  flex-direction: column;
  max-height: calc(100% - 64px);
  margin: var(--spacing-8);
  overflow-y: auto;
}
.cs-dialog__paper--xs { max-width: 444px; width: 100%; }
.cs-dialog__paper--sm { max-width: 600px; width: 100%; }
.cs-dialog__paper--md { max-width: 900px; width: 100%; }
.cs-dialog__paper--lg { max-width: 1200px; width: 100%; }
.cs-dialog__paper--xl { max-width: 1536px; width: 100%; }
.cs-dialog--fullscreen .cs-dialog__paper { max-width: none; max-height: none; margin: 0; width: 100%; height: 100%; border-radius: 0; }
.cs-dialog--fullwidth .cs-dialog__paper { width: calc(100% - 64px); }
.cs-dialog__title {
  padding: var(--spacing-4) var(--spacing-6);
  font-size: 1.25rem;
  font-weight: 500;
  line-height: 1.6;
}
.cs-dialog__content {
  padding: var(--spacing-1) var(--spacing-6) var(--spacing-5);
  font-size: 1rem;
  line-height: 1.5;
  overflow-y: auto;
  flex: 1 1 auto;
}
.cs-dialog__content--dividers {
  padding-top: var(--spacing-4);
  border-top: 1px solid var(--color-divider);
  border-bottom: 1px solid var(--color-divider);
}
.cs-dialog__actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-6) var(--spacing-4);
}

/* ── Bottom Navigation ───────────────────────────────────────────────── */
.cs-bottom-nav {
  display: flex;
  justify-content: center;
  width: 100%;
  height: 56px;
  background: var(--color-surface);
  box-shadow: var(--elevation-8);
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 50;
}
.cs-bottom-nav__action {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  max-width: 168px;
  min-width: 80px;
  padding: var(--spacing-1) var(--spacing-3);
  color: var(--text-secondary);
  background: transparent;
  border: none;
  cursor: pointer;
  text-decoration: none;
  transition: color var(--duration-shortest) var(--ease-standard);
}
.cs-bottom-nav__action--selected { color: var(--color-primary); }
.cs-bottom-nav__label {
  font-size: 0.75rem;
  margin-top: var(--spacing-1);
  transition: font-size var(--duration-shortest) var(--ease-standard);
}
.cs-bottom-nav--show-labels .cs-bottom-nav__label { display: block; }

/* ── Menubar ─────────────────────────────────────────────────────────── */
.cs-menubar {
  display: flex;
  align-items: center;
  gap: 0;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-divider);
}
.cs-menubar__item { position: relative; }
.cs-menubar__trigger {
  padding: var(--spacing-2) var(--spacing-3);
  border: none;
  background: transparent;
  font-family: var(--font-family);
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  cursor: pointer;
  transition: background var(--duration-shortest) var(--ease-standard);
}
.cs-menubar__trigger:hover { background: var(--color-hover); }
.cs-menubar__dropdown {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 200px;
  background: var(--color-surface);
  border-radius: var(--radius-sm);
  box-shadow: var(--elevation-8);
  z-index: var(--z-modal);
  padding: var(--spacing-1) 0;
}
.cs-menubar__item--open .cs-menubar__dropdown { display: block; }

/* ── Speed Dial ──────────────────────────────────────────────────────── */
.cs-speed-dial {
  position: fixed;
  z-index: var(--z-modal);
  display: flex;
  flex-direction: column-reverse;
  align-items: center;
}
.cs-speed-dial--up { bottom: var(--spacing-4); right: var(--spacing-4); flex-direction: column-reverse; }
.cs-speed-dial--down { top: var(--spacing-4); right: var(--spacing-4); flex-direction: column; }
.cs-speed-dial--left { bottom: var(--spacing-4); right: var(--spacing-4); flex-direction: row-reverse; }
.cs-speed-dial--right { bottom: var(--spacing-4); left: var(--spacing-4); flex-direction: row; }
.cs-speed-dial__actions {
  display: flex;
  flex-direction: column-reverse;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-3) 0;
  opacity: 0;
  transform: scale(0.5);
  transition: opacity var(--duration-shorter) var(--ease-standard),
              transform var(--duration-shorter) var(--ease-standard);
  pointer-events: none;
}
.cs-speed-dial--open .cs-speed-dial__actions { opacity: 1; transform: scale(1); pointer-events: auto; }
.cs-speed-dial__action {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}
.cs-speed-dial__tooltip {
  background: var(--text-primary);
  color: var(--color-surface);
  padding: var(--spacing-1) var(--spacing-2);
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  white-space: nowrap;
  box-shadow: var(--elevation-4);
}

/* ── Grid ────────────────────────────────────────────────────────────── */
.cs-grid--container {
  display: flex;
  flex-wrap: wrap;
  width: 100%;
}
.cs-grid--item { box-sizing: border-box; }
.cs-grid--xs-1  { flex: 0 0 8.333%; max-width: 8.333%; }
.cs-grid--xs-2  { flex: 0 0 16.667%; max-width: 16.667%; }
.cs-grid--xs-3  { flex: 0 0 25%; max-width: 25%; }
.cs-grid--xs-4  { flex: 0 0 33.333%; max-width: 33.333%; }
.cs-grid--xs-5  { flex: 0 0 41.667%; max-width: 41.667%; }
.cs-grid--xs-6  { flex: 0 0 50%; max-width: 50%; }
.cs-grid--xs-7  { flex: 0 0 58.333%; max-width: 58.333%; }
.cs-grid--xs-8  { flex: 0 0 66.667%; max-width: 66.667%; }
.cs-grid--xs-9  { flex: 0 0 75%; max-width: 75%; }
.cs-grid--xs-10 { flex: 0 0 83.333%; max-width: 83.333%; }
.cs-grid--xs-11 { flex: 0 0 91.667%; max-width: 91.667%; }
.cs-grid--xs-12 { flex: 0 0 100%; max-width: 100%; }
@media (min-width: 600px) {
  .cs-grid--sm-1  { flex: 0 0 8.333%; max-width: 8.333%; }
  .cs-grid--sm-2  { flex: 0 0 16.667%; max-width: 16.667%; }
  .cs-grid--sm-3  { flex: 0 0 25%; max-width: 25%; }
  .cs-grid--sm-4  { flex: 0 0 33.333%; max-width: 33.333%; }
  .cs-grid--sm-5  { flex: 0 0 41.667%; max-width: 41.667%; }
  .cs-grid--sm-6  { flex: 0 0 50%; max-width: 50%; }
  .cs-grid--sm-7  { flex: 0 0 58.333%; max-width: 58.333%; }
  .cs-grid--sm-8  { flex: 0 0 66.667%; max-width: 66.667%; }
  .cs-grid--sm-9  { flex: 0 0 75%; max-width: 75%; }
  .cs-grid--sm-10 { flex: 0 0 83.333%; max-width: 83.333%; }
  .cs-grid--sm-11 { flex: 0 0 91.667%; max-width: 91.667%; }
  .cs-grid--sm-12 { flex: 0 0 100%; max-width: 100%; }
}
@media (min-width: 900px) {
  .cs-grid--md-1  { flex: 0 0 8.333%; max-width: 8.333%; }
  .cs-grid--md-2  { flex: 0 0 16.667%; max-width: 16.667%; }
  .cs-grid--md-3  { flex: 0 0 25%; max-width: 25%; }
  .cs-grid--md-4  { flex: 0 0 33.333%; max-width: 33.333%; }
  .cs-grid--md-5  { flex: 0 0 41.667%; max-width: 41.667%; }
  .cs-grid--md-6  { flex: 0 0 50%; max-width: 50%; }
  .cs-grid--md-7  { flex: 0 0 58.333%; max-width: 58.333%; }
  .cs-grid--md-8  { flex: 0 0 66.667%; max-width: 66.667%; }
  .cs-grid--md-9  { flex: 0 0 75%; max-width: 75%; }
  .cs-grid--md-10 { flex: 0 0 83.333%; max-width: 83.333%; }
  .cs-grid--md-11 { flex: 0 0 91.667%; max-width: 91.667%; }
  .cs-grid--md-12 { flex: 0 0 100%; max-width: 100%; }
}
@media (min-width: 1200px) {
  .cs-grid--lg-1  { flex: 0 0 8.333%; max-width: 8.333%; }
  .cs-grid--lg-2  { flex: 0 0 16.667%; max-width: 16.667%; }
  .cs-grid--lg-3  { flex: 0 0 25%; max-width: 25%; }
  .cs-grid--lg-4  { flex: 0 0 33.333%; max-width: 33.333%; }
  .cs-grid--lg-5  { flex: 0 0 41.667%; max-width: 41.667%; }
  .cs-grid--lg-6  { flex: 0 0 50%; max-width: 50%; }
  .cs-grid--lg-7  { flex: 0 0 58.333%; max-width: 58.333%; }
  .cs-grid--lg-8  { flex: 0 0 66.667%; max-width: 66.667%; }
  .cs-grid--lg-9  { flex: 0 0 75%; max-width: 75%; }
  .cs-grid--lg-10 { flex: 0 0 83.333%; max-width: 83.333%; }
  .cs-grid--lg-11 { flex: 0 0 91.667%; max-width: 91.667%; }
  .cs-grid--lg-12 { flex: 0 0 100%; max-width: 100%; }
}
@media (min-width: 1536px) {
  .cs-grid--xl-1  { flex: 0 0 8.333%; max-width: 8.333%; }
  .cs-grid--xl-2  { flex: 0 0 16.667%; max-width: 16.667%; }
  .cs-grid--xl-3  { flex: 0 0 25%; max-width: 25%; }
  .cs-grid--xl-4  { flex: 0 0 33.333%; max-width: 33.333%; }
  .cs-grid--xl-5  { flex: 0 0 41.667%; max-width: 41.667%; }
  .cs-grid--xl-6  { flex: 0 0 50%; max-width: 50%; }
  .cs-grid--xl-7  { flex: 0 0 58.333%; max-width: 58.333%; }
  .cs-grid--xl-8  { flex: 0 0 66.667%; max-width: 66.667%; }
  .cs-grid--xl-9  { flex: 0 0 75%; max-width: 75%; }
  .cs-grid--xl-10 { flex: 0 0 83.333%; max-width: 83.333%; }
  .cs-grid--xl-11 { flex: 0 0 91.667%; max-width: 91.667%; }
  .cs-grid--xl-12 { flex: 0 0 100%; max-width: 100%; }
}

/* ── Stack ───────────────────────────────────────────────────────────── */
.cs-stack { display: flex; }

/* ── Image List ──────────────────────────────────────────────────────── */
.cs-image-list { display: grid; overflow: hidden; }
.cs-image-list--standard .cs-image-list__item { overflow: hidden; }
.cs-image-list__img { width: 100%; height: 100%; object-fit: cover; display: block; }
.cs-image-list__item { position: relative; overflow: hidden; }
.cs-image-list__title-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0,0,0,0.6));
  color: #fff;
  padding: var(--spacing-2) var(--spacing-3);
  font-size: 0.875rem;
}

/* ── Masonry ─────────────────────────────────────────────────────────── */
.cs-masonry > * { break-inside: avoid; margin-bottom: var(--spacing-2); }
`
