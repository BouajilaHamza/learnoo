"""Generate MSA narration matching the shared, unreviewed slides."""
from pipeline import narration, options


def main():
    args = options(__doc__).parse_args()
    narration(args.output_dir.resolve(), 'msa')


if __name__ == '__main__':
    main()
