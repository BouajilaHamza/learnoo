"""Combine current shared slides with matching MSA narration."""
from pipeline import options, render


def main():
    args = options(__doc__).parse_args()
    render(args.output_dir.resolve(), 'msa')


if __name__ == '__main__':
    main()
