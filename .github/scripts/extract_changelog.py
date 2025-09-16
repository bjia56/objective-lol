#!/usr/bin/env python3
"""
Extract changelog entries for GitHub releases
"""
import argparse
import re
import sys
from pathlib import Path


def extract_changelog_section(changelog_path: Path, version: str) -> str:
    """Extract the changelog section for a specific version."""
    if not changelog_path.exists():
        return ""

    content = changelog_path.read_text()
    lines = content.split('\n')

    # Pattern to match version headers like "## [1.0.0]" or "### [1.0.0]"
    version_pattern = rf'^(##+)\s*\[{re.escape(version)}\]'

    result_lines = []
    in_section = False
    section_level = None

    for line in lines:
        # Check if this is a version header
        match = re.match(version_pattern, line)
        if match:
            in_section = True
            section_level = len(match.group(1))  # Number of # characters
            continue

        # Slight hack - the Python readme uses 4 # for section headers
        if line.startswith('#### '):
            line = line.replace('#### ', '### ', 1)

        # If we're in the section, check for next section header
        if in_section:
            # Stop if we hit another version section (semantic versioning pattern)
            semver_pattern = rf'^#{{{section_level}}}\s*\[\d+\.\d+\.\d+.*?\]'
            if re.match(semver_pattern, line):
                break
            result_lines.append(line)

    # Clean up the result
    changelog = '\n'.join(result_lines).strip()

    # Remove any trailing empty lines
    while changelog.endswith('\n\n'):
        changelog = changelog.rstrip('\n')

    return changelog


def main():
    parser = argparse.ArgumentParser(description='Extract changelog for version')
    parser.add_argument('--version', required=True, help='Version to extract')
    parser.add_argument('--component', required=True,
                       choices=['main', 'python', 'vscode'],
                       help='Component to extract changelog for')
    parser.add_argument('--output', help='Output file (default: stdout)')

    args = parser.parse_args()

    # Determine changelog file based on component
    if args.component == 'main':
        changelog_path = Path('CHANGELOG.md')
    elif args.component == 'python':
        changelog_path = Path('python/README.md')
    elif args.component == 'vscode':
        changelog_path = Path('vscode/CHANGELOG.md')

    changelog = extract_changelog_section(changelog_path, args.version)

    if not changelog:
        print(f"No changelog entry found for version {args.version} in {changelog_path}",
              file=sys.stderr)
        sys.exit(1)

    if args.output:
        Path(args.output).write_text(changelog)
    else:
        print(changelog)


if __name__ == '__main__':
    main()