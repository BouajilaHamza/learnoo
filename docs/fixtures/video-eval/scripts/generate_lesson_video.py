"""Generate unreviewed slides and optionally Derja narration/video."""
from pipeline import narration, options, render, require_tool, screenshots, write_slides


def main():
    parser = options(__doc__)
    parser.add_argument('--html-only', action='store_true', help='Write offline HTML without browser, TTS or video')
    parser.add_argument('--chrome', default='google-chrome', help='Chrome/Chromium executable or absolute path')
    args = parser.parse_args()
    root = args.output_dir.resolve()
    write_slides(root)
    if args.html_only:
        return
    screenshots(root, require_tool(args.chrome))
    narration(root, 'derja')
    render(root, 'derja')


if __name__ == '__main__':
    main()
