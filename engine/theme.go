package engine

// GetThemeCSS returns the full CSS for a given theme name
func GetThemeCSS(theme string) string {
	return tokenCSS() + resetCSS + componentCSS + variantCSS() + componentLibCSS + atomCSS + themeVars(theme) + desktopCSS
}

const desktopCSS = `
/* ── Desktop shell ─────────────────────────────────────── */
.desktop {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #0a0a0f;
  user-select: none;
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
  position: relative;
}
.desktop__rain {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  opacity: 0.06;
  z-index: 0;
}
.desktop__body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 1;
  padding: var(--spacing-6);
}
.desktop__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, 70px);
  gap: var(--spacing-6);
  max-width: 700px;
  width: 100%;
  justify-content: center;
}

/* ── App icons ─────────────────────────────────────────── */
.desk-icon {
  width: 70px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--duration-shorter) var(--ease-standard);
}
.desk-icon:hover .desk-icon__img {
  color: #4f4;
  filter: drop-shadow(0 0 8px rgba(0,255,0,0.6));
}
.desk-icon:active {
  transform: scale(0.92);
}
.desk-icon__img {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0f0;
  transition: all var(--duration-shorter) var(--ease-standard);
}
.desk-icon--locked {
  cursor: default;
}
.desk-icon--locked .desk-icon__img {
  color: #0a0;
  opacity: 0.35;
}
.desk-icon--locked:hover .desk-icon__img {
  color: #0a0;
  filter: none;
  opacity: 0.35;
}

/* ── Taskbar ───────────────────────────────────────────── */
.desktop__taskbar {
  height: 48px;
  background: rgba(0,0,0,0.85);
  border-top: 1px solid rgba(0,255,0,0.15);
  display: flex;
  align-items: center;
  padding: 0 var(--spacing-5);
  gap: var(--spacing-3);
  position: relative;
  z-index: 1;
  backdrop-filter: blur(10px);
}
.taskbar__orb {
  width: 24px;
  height: 24px;
  background: #0a0;
  border-radius: var(--radius-full);
  opacity: 0.6;
  cursor: pointer;
  box-shadow: 0 0 10px rgba(0,255,0,0.3);
}
.taskbar__brand {
  color: #0f0;
  font-size: 0.85rem;
  font-weight: bold;
  letter-spacing: 3px;
  opacity: 0.5;
  margin-left: var(--spacing-1);
}
.taskbar__divider {
  width: 1px;
  height: var(--spacing-5);
  background: rgba(0,255,0,0.15);
  margin: 0 var(--spacing-2);
}
.taskbar__apps {
  display: flex;
  gap: var(--spacing-2);
  flex: 1;
}
.taskbar__clock {
  color: #070;
  font-size: 0.85rem;
  letter-spacing: 1px;
}
.taskbar-btn {
  padding: var(--spacing-2) var(--spacing-4);
  background: rgba(0,255,0,0.06);
  border: 1px solid rgba(0,255,0,0.15);
  border-radius: var(--radius-sm);
  color: #0f0;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background var(--duration-shortest) var(--ease-standard);
}
.taskbar-btn:hover {
  background: rgba(0,255,0,0.12);
}
.taskbar__logout {
  padding: var(--spacing-2) var(--spacing-3);
  color: #070;
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: color var(--duration-shortest) var(--ease-standard), background var(--duration-shortest) var(--ease-standard);
  display: flex;
  align-items: center;
}
.taskbar__logout:hover {
  color: #0f0;
  background: rgba(0,255,0,0.1);
}
.taskbar__role {
  font-size: 0.75rem;
  letter-spacing: 2px;
  color: #0a0;
  background: rgba(0,255,0,0.08);
  padding: var(--spacing-1) var(--spacing-2);
  border-radius: var(--radius-sm);
  border: 1px solid rgba(0,255,0,0.15);
}

/* ── Per-app headers ─────────────────────────────────────── */
.app-header {
  padding: var(--spacing-3) var(--spacing-5);
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  font-weight: 600;
  font-size: 1rem;
  letter-spacing: 0.02em;
  border-bottom: 1px solid rgba(255,255,255,0.08);
}
.app-header .app-icon { font-size: 1.25rem; }
.app-header-email     { background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%); color: #60a5fa; }
.app-header-group     { background: linear-gradient(135deg, #1a1a2e 0%, #1e2a3a 100%); color: #a78bfa; }
.app-header-chat      { background: linear-gradient(135deg, #0f1a0f 0%, #1a2e1a 100%); color: #4ade80; }
.app-header-tickets   { background: linear-gradient(135deg, #2a1a1a 0%, #2e1616 100%); color: #f87171; }
.app-header-siem      { background: linear-gradient(135deg, #1a1a0f 0%, #2e2a16 100%); color: #facc15; }
.app-header-playbook  { background: linear-gradient(135deg, #1a1a2e 0%, #1e1a2e 100%); color: #c084fc; }
.app-header-directory { background: linear-gradient(135deg, #0f1a1a 0%, #162e2e 100%); color: #2dd4bf; }
`

func themeVars(theme string) string {
	switch theme {
	case "dark":
		return themeDark()
	case "stretch":
		return themeStretch()
	default:
		return themeLight()
	}
}

func themeLight() string {
	return `
:root {
  --bg: #fafbfc;
  --bg-elevated: #ffffff;
  --bg-surface: #f6f8fa;
  --glass: rgba(255, 255, 255, 0.95);
  --glass-hover: rgba(0, 0, 0, 0.04);
  --glass-border: #e1e4e8;
  --glass-border-hover: #d1d5da;
  --text-primary: #1a1a1a;
  --text-secondary: #586069;
  --text-tertiary: #6a737d;
  --text-muted: #959da5;
  --border: #e1e4e8;
  --border-subtle: #eaecef;
  --accent: #667eea;
  --accent-hover: #5a6fd6;
  --accent-text: #ffffff;
  --copper: #667eea;
  --copper-hover: #5a6fd6;
  --success: #4caf50;
  --danger: #d73a49;

  /* Semantic Colors */
  --color-primary: #667eea;
  --color-primary-hover: #5a6fd6;
  --color-primary-glass: rgba(102, 126, 234, 0.1);
  --color-secondary: #f6f8fa;
  --color-secondary-hover: #e1e4e8;
  --color-success: #4caf50;
  --color-success-hover: #43a047;
  --color-success-glass: rgba(76, 175, 80, 0.1);
  --color-warning: #ff9800;
  --color-warning-hover: #fb8c00;
  --color-warning-glass: rgba(255, 152, 0, 0.1);
  --color-danger: #d73a49;
  --color-danger-hover: #cb2431;
  --color-danger-glass: rgba(215, 58, 73, 0.1);
  --color-info: #0366d6;
  --color-info-hover: #0256b9;
  --color-info-glass: rgba(3, 102, 214, 0.1);
  --color-muted: #f6f8fa;
  --color-muted-hover: #e1e4e8;
}`
}

func themeDark() string {
	return `
:root {
  --bg: #0a0a0f;
  --bg-elevated: rgba(255, 255, 255, 0.04);
  --bg-surface: rgba(255, 255, 255, 0.07);
  --glass: rgba(255, 255, 255, 0.06);
  --glass-hover: rgba(255, 255, 255, 0.1);
  --glass-border: rgba(255, 255, 255, 0.12);
  --glass-border-hover: rgba(255, 255, 255, 0.22);
  --text-primary: rgba(255, 255, 255, 0.95);
  --text-secondary: rgba(255, 255, 255, 0.65);
  --text-tertiary: rgba(255, 255, 255, 0.35);
  --border: rgba(255, 255, 255, 0.1);
  --border-subtle: rgba(255, 255, 255, 0.06);
  --accent: rgba(255, 255, 255, 0.9);
  --accent-hover: rgba(255, 255, 255, 1);
  --accent-text: #0a0a0f;
  --copper: #e8a87c;
  --copper-hover: #f0bc96;
  --success: #34d399;
  --danger: #f87171;

  /* Semantic Colors */
  --color-primary: rgba(255, 255, 255, 0.9);
  --color-primary-hover: rgba(255, 255, 255, 1);
  --color-primary-glass: rgba(255, 255, 255, 0.1);
  --color-secondary: rgba(255, 255, 255, 0.08);
  --color-secondary-hover: rgba(255, 255, 255, 0.14);
  --color-success: #34d399;
  --color-success-hover: #10b981;
  --color-success-glass: rgba(52, 211, 153, 0.15);
  --color-warning: #eab308;
  --color-warning-hover: #ca8a04;
  --color-warning-glass: rgba(234, 179, 8, 0.15);
  --color-danger: #dc2626;
  --color-danger-hover: #b91c1c;
  --color-danger-glass: rgba(220, 38, 38, 0.15);
  --color-info: #4169e1;
  --color-info-hover: #2f55c8;
  --color-info-glass: rgba(65, 105, 225, 0.15);
  --color-muted: rgba(255, 255, 255, 0.06);
  --color-muted-hover: rgba(255, 255, 255, 0.1);
}`
}

func themeStretch() string {
	return `
:root {
  --bg: #FAF8F5;
  --bg-elevated: #FFFFFF;
  --bg-surface: #F3F0EC;
  --text-primary: #1B2A4A;
  --text-secondary: #4A5568;
  --text-tertiary: #8C9AAF;
  --border: #E8E4DF;
  --border-subtle: #F0ECE7;
  --accent: #1B2A4A;
  --accent-hover: #253759;
  --accent-text: #FAF8F5;
  --copper: #C4724E;
  --copper-hover: #B56542;
  --success: #2F7D5B;
  --danger: #C0392B;

  /* Semantic Colors */
  --color-primary: #1B2A4A;
  --color-primary-hover: #253759;
  --color-secondary: #F3F0EC;
  --color-secondary-hover: #E8E4DF;
  --color-success: #2F7D5B;
  --color-success-hover: #256A4C;
  --color-warning: #D4850A;
  --color-warning-hover: #B87208;
  --color-danger: #C0392B;
  --color-danger-hover: #A93226;
  --color-info: #2B6CB0;
  --color-info-hover: #245BA0;
  --color-muted: #F3F0EC;
  --color-muted-hover: #E8E4DF;
}`
}

const resetCSS = `
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: var(--font-family);
  background: var(--bg);
  color: var(--text-primary);
  line-height: 1.5;
  min-height: 100vh;
  font-size: 1rem;
  letter-spacing: -0.006em;
}
`

const componentCSS = `
/* --- Components on MUI token system. 8px spacing base. --- */

/* Header */
.cs-header {
  padding: var(--spacing-5) var(--spacing-6);
  border-bottom: 1px solid var(--border-subtle);
}
.header-title {
  font-family: var(--font-family);
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.334;
  letter-spacing: -0.02em;
}
.header-subtitle {
  font-size: 0.875rem;
  color: var(--text-tertiary);
  margin-left: var(--spacing-3);
  font-weight: 400;
}

/* Text */
.cs-text {
  font-size: 0.875rem;
  color: var(--text-secondary);
  line-height: 1.43;
}

/* Button — base styles in variantCSS, legacy compat below */

/* Row */
/* Card */
.cs-card {
  background: var(--glass);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  padding: var(--spacing-6) var(--spacing-8);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
  align-items: flex-start;
  height: 100%;
  transition: border-color var(--duration-shorter) var(--ease-standard);
}
.cs-card:hover {
  border-color: var(--glass-border-hover);
}

/* Heading */
.cs-heading {
  font-family: var(--font-family);
  color: var(--text-primary);
  font-weight: 600;
  line-height: 1.235;
  letter-spacing: -0.02em;
}
h2.cs-heading { font-size: 1.25rem; }
h3.cs-heading { font-size: 1rem; }

/* Footer */
.cs-footer {
  padding: var(--spacing-6) var(--spacing-6) var(--spacing-4);
  text-align: center;
  color: var(--text-tertiary);
  font-size: 0.75rem;
}

/* Nav */
.cs-nav {
  display: flex;
  align-items: center;
  padding: 0 var(--spacing-6);
  height: var(--spacing-12);
  border-bottom: 1px solid var(--glass-border);
  gap: var(--spacing-6);
  background: var(--bg-elevated);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  position: sticky;
  top: 0;
  z-index: 50;
}
.nav-brand {
  font-family: var(--font-family);
  font-size: 1rem;
  font-weight: 600;
  color: var(--copper);
  letter-spacing: -0.02em;
}
.nav-links {
  display: flex;
  gap: var(--spacing-5);
}
.cs-nav-link {
  color: var(--text-tertiary);
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: color var(--duration-shortest) var(--ease-standard);
}
.cs-nav-link:hover {
  color: var(--text-primary);
}

/* Stat Card */
.cs-stat-card {
  background: var(--glass);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  padding: var(--spacing-5) var(--spacing-6);
  flex: 1;
  transition: border-color var(--duration-shorter) var(--ease-standard);
}
.cs-stat-card:hover {
  border-color: var(--glass-border-hover);
}
.stat-label {
  font-size: 0.75rem;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-weight: 500;
  margin-bottom: var(--spacing-1);
}
.stat-value {
  font-family: var(--font-family);
  font-size: 2.125rem;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: -0.02em;
  line-height: 1.235;
}
.trend-up { color: var(--success); font-size: 0.875rem; font-weight: 500; }
.trend-down { color: var(--copper); font-size: 0.875rem; font-weight: 500; }
`
