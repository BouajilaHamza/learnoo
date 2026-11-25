"""Local, unreviewed slide experiment. Importing this module has no side effects."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess

from slides import HTML_TEMPLATE, SLIDES


def fingerprint():
    return hashlib.sha256(json.dumps([HTML_TEMPLATE, SLIDES], ensure_ascii=False).encode()).hexdigest()


def options(description):
    parser = argparse.ArgumentParser(description=description)
    parser.add_argument('--output-dir', type=Path, required=True, help='Local scratch directory; never published automatically')
    return parser


def require_tool(name):
    result = shutil.which(name)
    if result is None:
        raise SystemExit(f'Missing executable: {name}. Install it before running this optional experiment.')
    return result


def write_slides(root):
    root.mkdir(parents=True, exist_ok=True)
    for slide in SLIDES:
        html = HTML_TEMPLATE
        for token, key in [('BADGE', 'badge'), ('TITLE', 'title'), ('SUBTITLE', 'subtitle'), ('CONTENT', 'content_html')]:
            html = html.replace(f'__{token}__', slide[key])
        (root / f"{slide['id']}.html").write_text(html, encoding='utf-8')


def require_current(root, marker):
    path = root / marker
    if not path.is_file() or path.read_text() != fingerprint():
        raise SystemExit(f'Missing or stale {marker}; regenerate these inputs from the current slides.')


def screenshots(root, chrome):
    marker = root / 'images.sha256'
    marker.unlink(missing_ok=True)
    for slide in SLIDES:
        output = root / f"{slide['id']}.png"
        output.unlink(missing_ok=True)
        subprocess.run([chrome, '--headless=new', '--disable-gpu', '--hide-scrollbars',
                        '--window-size=1920,1080', f'--screenshot={output}',
                        (root / f"{slide['id']}.html").as_uri()], check=True)
        if not output.is_file():
            raise SystemExit(f'Browser did not produce {output}')
    marker.write_text(fingerprint())


def narration(root, variant):
    tts = require_tool('edge-tts')
    target = root / variant
    target.mkdir(parents=True, exist_ok=True)
    marker = target / 'audio.sha256'
    marker.unlink(missing_ok=True)
    for slide in SLIDES:
        text = slide['voice_text' if variant == 'derja' else 'msa_text']
        subprocess.run([tts, '--voice', 'ar-TN-HediNeural', '--text', text,
                        '--write-media', str(target / f"{slide['id']}.mp3")], check=True)
    marker.write_text(fingerprint())


def render(root, variant):
    require_current(root, 'images.sha256')
    target = root / variant
    require_current(target, 'audio.sha256')
    ffmpeg = require_tool('ffmpeg')
    for slide in SLIDES:
        stem = slide['id']
        subprocess.run([ffmpeg, '-y', '-loop', '1', '-framerate', '25', '-i', str(root / f'{stem}.png'),
                        '-i', str(target / f'{stem}.mp3'), '-c:v', 'libx264', '-tune', 'stillimage',
                        '-c:a', 'aac', '-b:a', '192k', '-pix_fmt', 'yuv420p', '-shortest',
                        str(target / f'{stem}.mp4')], check=True)
    # Fixed relative filenames avoid concat escaping problems in user-selected paths.
    listing = target / 'segments.txt'
    listing.write_text(''.join(f"file '{s['id']}.mp4'\n" for s in SLIDES))
    subprocess.run([ffmpeg, '-y', '-f', 'concat', '-safe', '1', '-i', str(listing),
                    '-c', 'copy', str(target / 'unreviewed-lesson.mp4')], check=True)
