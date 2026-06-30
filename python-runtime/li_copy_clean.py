#!/usr/bin/env python3
"""
LinkedIn scraper - Core functions for Electron frontend
All tkinter GUI code removed - uses Electron instead
"""

from __future__ import annotations

import asyncio
import json
import random
import sys
from pathlib import Path

import pandas as pd
import requests
from playwright.async_api import async_playwright


# ----------  Core scraping function  ----------
async def copy_one_profile(
    url: str, page, email: str = "", category: str = "", log_callback=None
) -> dict[str, str]:
    """Copy visible fields from a single profile."""
    await page.goto(url, wait_until="domcontentloaded")
    await asyncio.sleep(3)  # let JS settle

    def log(msg: str):
        """Log message via callback"""
        if log_callback:
            log_callback(msg)
        else:
            print(msg, flush=True)

    async def _text(sel: str, default: str = "", label: str = "") -> str:
        try:
            element = page.locator(sel).first
            text = (await element.inner_text(timeout=5_000)).strip()
            if text:
                log(f"  [+] {label}: {text[:100]}{'...' if len(text) > 100 else ''}")
            else:
                log(f"  [!] {label}: (empty)")
            return text
        except Exception as e:
            log(f"  [-] {label}: Not found ({str(e)[:50]})")
            return default

    log(f"\n[>] Copying: {url}")

    # Click all "see more" buttons to expand content
    try:
        see_more_buttons = page.locator(
            'button:has-text("see more"), button:has-text("…see more")'
        )
        count = await see_more_buttons.count()
        if count > 0:
            log(f"  [+] Expanding {count} collapsed sections...")
            for i in range(count):
                try:
                    await see_more_buttons.nth(i).click(timeout=1_000)
                    await asyncio.sleep(0.3)  # Small delay after each click
                except Exception:
                    pass  # Some buttons might not be clickable, that's OK
    except Exception as e:
        log(f"  [!] Could not expand sections: {str(e)[:40]}")

    # Basic Info
    name = await _text("h1", default="", label="Name")

    # If name not found, page doesn't exist or didn't load - skip this profile
    if not name:
        log(f"  [X] SKIPPING: No name found - page may not exist or failed to load")
        return {
            "url": url,
            "email": email,
            "category": category,
            "firstName": "",
            "lastName": "",
            "fullName": "",
            "headline": "",
            "location": "",
            "currentRole": "",
            "currentCompany": "",
            "about": "",
            "experience": "",
            "education": "",
            "skills": "",
            "activity": "",
            "publications": "",
            "interests": "",
            "error": "Profile not found or failed to load",
        }

    # Split name into first and last
    name_parts = name.strip().split(None, 1)  # Split on first whitespace
    first_name = name_parts[0] if len(name_parts) > 0 else ""
    last_name = name_parts[1] if len(name_parts) > 1 else ""

    headline = await _text("div.text-body-medium", default="", label="Headline")
    location = await _text(
        "span.text-body-small.inline.t-black--light", default="", label="Location"
    )

    # About section - get all text from the section
    about = ""
    try:
        # Try multiple strategies to find About text
        # Strategy 1: Get the parent section element
        about_section = (
            page.locator("#about").locator("xpath=ancestor::section[1]").first
        )
        if await about_section.count() > 0:
            about = (await about_section.inner_text(timeout=5_000)).strip()

        # Strategy 2: Navigate up from #about anchor
        if not about:
            about_container = (
                page.locator("#about").locator("xpath=../../div[last()]").first
            )
            if await about_container.count() > 0:
                about = (await about_container.inner_text(timeout=5_000)).strip()

        # Clean up - remove heading
        if about:
            about = about.replace("About\n", "").replace("About", "").strip()
            log(f"  [+] About: {about[:150]}...")
        else:
            log(f"  [!] About: (empty)")
    except Exception as e:
        log(f"  [-] About: Not found ({str(e)[:50]})")

    # Experience section - get all text from the experience list only
    experience_text = ""
    try:
        # Get the experience section
        exp_section = (
            page.locator("#experience").locator("xpath=ancestor::section[1]").first
        )
        if await exp_section.count() > 0:
            # Get just the UL containing experience items
            exp_ul = exp_section.locator("ul").first
            if await exp_ul.count() > 0:
                experience_text = (await exp_ul.inner_text(timeout=5_000)).strip()
            else:
                # Fallback: get text from section but will include headers
                experience_text = (await exp_section.inner_text(timeout=5_000)).strip()
                experience_text = (
                    experience_text.replace("Experience\n", "")
                    .replace("Experience", "")
                    .strip()
                )

        if experience_text:
            log(f"  [+] Experience: {len(experience_text)} characters captured")
        else:
            log(f"  [!] Experience: (empty)")
    except Exception as e:
        log(f"  [-] Experience: {str(e)[:50]}")

    # Current Role and Company - parse from experience text
    current_role = ""
    current_company = ""
    try:
        if experience_text:
            # Split by 'Present' to isolate the current position
            if "Present" in experience_text:
                before_present = experience_text.split("Present")[0]
                lines = before_present.split("\n")

                # Role is the first line
                if len(lines) > 0:
                    current_role = lines[0].strip()

                # Company is typically the 3rd line (index 2)
                if len(lines) > 2:
                    company_line = lines[2].strip()
                    # Remove "Full-time", "Part-time", etc after ·
                    if "·" in company_line:
                        current_company = company_line.split("·")[0].strip()
                    else:
                        current_company = company_line

            if current_role:
                log(f"  [+] Current Role: {current_role[:80]}")
            if current_company:
                log(f"  [+] Current Company: {current_company[:80]}")

            if not current_role and not current_company:
                log(f"  [!] Current role/company: No 'Present' position found")
    except Exception as e:
        log(f"  [-] Current role/company parsing: {str(e)[:50]}")

    # Education section - get all text
    education_text = ""
    try:
        edu_section = (
            page.locator("#education").locator("xpath=ancestor::section[1]").first
        )
        if await edu_section.count() > 0:
            education_text = (await edu_section.inner_text(timeout=5_000)).strip()
            education_text = (
                education_text.replace("Education\n", "")
                .replace("Education", "")
                .strip()
            )

        if education_text:
            log(f"  [+] Education: {len(education_text)} characters captured")
        else:
            log(f"  [!] Education: (empty)")
    except Exception as e:
        log(f"  [-] Education: {str(e)[:50]}")

    # Skills section - get all text
    skills_text = ""
    try:
        skills_section = (
            page.locator("#skills").locator("xpath=ancestor::section[1]").first
        )
        if await skills_section.count() > 0:
            skills_text = (await skills_section.inner_text(timeout=5_000)).strip()
            skills_text = (
                skills_text.replace("Skills\n", "").replace("Skills", "").strip()
            )

        if skills_text:
            log(f"  [+] Skills: {len(skills_text)} characters captured")
        else:
            log(f"  [!] Skills: (empty)")
    except Exception as e:
        log(f"  [-] Skills: {str(e)[:50]}")

    # Activity section - get all text
    activity_text = ""
    try:
        activity_section = page.locator(
            'xpath=//*[@id="profile-content"]/div/div[2]/div/div/main/section[3]/div[4]'
        ).first
        if await activity_section.count() > 0:
            activity_text = (await activity_section.inner_text(timeout=5_000)).strip()

        if activity_text:
            log(f"  [+] Activity: {len(activity_text)} characters captured")
        else:
            log(f"  [!] Activity: (empty)")
    except Exception as e:
        log(f"  [-] Activity: {str(e)[:50]}")

    # Publications section - get all text
    publications_text = ""
    try:
        pub_section = (
            page.locator("#publications").locator("xpath=ancestor::section[1]").first
        )
        if await pub_section.count() > 0:
            publications_text = (await pub_section.inner_text(timeout=5_000)).strip()
            publications_text = (
                publications_text.replace("Publications\n", "")
                .replace("Publications", "")
                .strip()
            )

        if publications_text:
            log(f"  [+] Publications: {len(publications_text)} characters captured")
        else:
            log(f"  [!] Publications: (empty)")
    except Exception as e:
        log(f"  [-] Publications: {str(e)[:50]}")

    # Interests section - get all text
    interests_text = ""
    try:
        int_section = (
            page.locator("#interests").locator("xpath=ancestor::section[1]").first
        )
        if await int_section.count() > 0:
            interests_text = (await int_section.inner_text(timeout=5_000)).strip()
            interests_text = (
                interests_text.replace("Interests\n", "")
                .replace("Interests", "")
                .strip()
            )

        if interests_text:
            log(f"  [+] Interests: {len(interests_text)} characters captured")
        else:
            log(f"  [!] Interests: (empty)")
    except Exception as e:
        log(f"  [-] Interests: {str(e)[:50]}")

    return {
        "url": url,
        "email": email,
        "category": category,
        "firstName": first_name,
        "lastName": last_name,
        "fullName": name,
        "headline": headline,
        "location": location,
        "currentRole": current_role,
        "currentCompany": current_company,
        "about": about,
        "experience": experience_text,
        "education": education_text,
        "skills": skills_text,
        "activity": activity_text,
        "publications": publications_text,
        "interests": interests_text,
    }


# ----------  Main scraper function  ----------
async def run_scraper(
    csv_path, port, min_delay, max_delay, test_limit, output_format, webhook_url
):
    """Main scraper function for command-line/Electron use"""
    csv_path = Path(csv_path)

    # Read CSV
    try:
        df = pd.read_csv(csv_path)
        if "linkedin_url" not in df.columns:
            print("ERROR: CSV must have 'linkedin_url' column", flush=True)
            return False

        urls = df["linkedin_url"].dropna().astype(str).tolist()
        if not urls:
            print("ERROR: No URLs found in linkedin_url column", flush=True)
            return False

        # Get emails if available
        if "Email Address" in df.columns:
            emails = df["Email Address"].fillna("").astype(str).tolist()
        else:
            emails = [""] * len(urls)
            print("[!] No 'Email Address' column found in CSV", flush=True)

        # Get categories if available
        if "category" in df.columns:
            categories = df["category"].fillna("").astype(str).tolist()
        else:
            categories = [""] * len(urls)
            print("[!] No 'category' column found in CSV", flush=True)

        profile_data = list(zip(urls, emails, categories))

        # Apply test limit
        total_urls = len(profile_data)
        if test_limit > 0 and test_limit < total_urls:
            profile_data = profile_data[:test_limit]
            print(
                f"\n[!] TEST MODE: Processing only first {test_limit} of {total_urls} profiles\n",
                flush=True,
            )

    except Exception as e:
        print(f"ERROR: Failed to read CSV: {e}", flush=True)
        return False

    # Connect to Chrome and scrape
    rows = []
    cdp_url = f"http://localhost:{port}"

    print(f"[*] Connecting to Chrome on port {port}...", flush=True)

    async with async_playwright() as p:
        try:
            browser = await p.chromium.connect_over_cdp(cdp_url)
            print("[+] Connected to Chrome!", flush=True)

            contexts = browser.contexts
            if not contexts:
                print("ERROR: No browser contexts found", flush=True)
                return False

            pages = contexts[0].pages
            if not pages:
                page = await contexts[0].new_page()
            else:
                page = pages[0]

            print(
                f"\n[*] Starting to copy {len(profile_data)} profiles...\n", flush=True
            )

            for idx, (url, email, category) in enumerate(profile_data, 1):
                print(f"PROGRESS:{idx}/{len(profile_data)}", flush=True)

                try:
                    row = await copy_one_profile(url, page, email, category)
                except Exception as exc:
                    row = {
                        "url": url,
                        "email": email,
                        "category": category,
                        "error": str(exc),
                    }
                    print(f"  [X] Error: {exc}", flush=True)

                rows.append(row)

                # Delay between profiles
                if idx < len(profile_data):
                    delay = random.randint(min_delay, max_delay)
                    print(f"  [~] Waiting {delay}s...\n", flush=True)
                    await asyncio.sleep(delay)

            print(f"\n[+] Completed all {len(rows)} profiles!", flush=True)

        except Exception as exc:
            print(f"ERROR: Failed to connect to Chrome: {exc}", flush=True)
            return False

    # Save results
    try:
        if output_format == "json":
            output_file = csv_path.with_name(csv_path.stem + "_copied.json")
            with open(output_file, "w", encoding="utf-8") as f:
                json.dump(rows, f, indent=2, ensure_ascii=False)
        else:
            output_file = csv_path.with_name(csv_path.stem + "_copied.csv")
            new_df = pd.DataFrame(rows)
            final_df = df.merge(
                new_df, left_on="linkedin_url", right_on="url", how="left"
            )
            final_df.to_csv(output_file, index=False)

        print(f"\n[+] Saved {len(rows)} profiles to {output_file}", flush=True)

        # Send to webhook if provided
        if webhook_url and output_format == "json":
            try:
                print(f"\n[*] Sending data to webhook...", flush=True)
                response = requests.post(
                    webhook_url,
                    json=rows,
                    headers={"Content-Type": "application/json"},
                    timeout=30,
                )

                if response.status_code in [200, 201, 202, 204]:
                    print(
                        f"[+] Webhook success! Status: {response.status_code}",
                        flush=True,
                    )
                else:
                    print(
                        f"[!] Webhook returned status {response.status_code}",
                        flush=True,
                    )

            except Exception as e:
                print(f"[X] Webhook error: {str(e)[:100]}", flush=True)

        return True

    except Exception as e:
        print(f"ERROR: Failed to save results: {e}", flush=True)
        return False


# ----------  Entry point for command-line use  ----------
if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description="LinkedIn profile scraper")
    parser.add_argument("--csv", required=True, help="Path to CSV file")
    parser.add_argument("--port", type=int, required=True, help="Chrome debug port")
    parser.add_argument(
        "--min-delay", type=int, default=3, help="Min delay between profiles (seconds)"
    )
    parser.add_argument(
        "--max-delay", type=int, default=8, help="Max delay between profiles (seconds)"
    )
    parser.add_argument(
        "--test-limit", type=int, default=0, help="Test limit (0 = all)"
    )
    parser.add_argument(
        "--format", choices=["json", "csv"], default="json", help="Output format"
    )
    parser.add_argument("--webhook", default=None, help="Webhook URL")

    args = parser.parse_args()

    success = asyncio.run(
        run_scraper(
            args.csv,
            args.port,
            args.min_delay,
            args.max_delay,
            args.test_limit,
            args.format,
            args.webhook,
        )
    )

    sys.exit(0 if success else 1)
