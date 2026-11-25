import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from pipeline import SLIDES, fingerprint, narration, render, write_slides


class PipelineTests(unittest.TestCase):
    def test_offline_slides_carry_unreviewed_notice(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            write_slides(root)
            for slide in SLIDES:
                html = (root / f"{slide['id']}.html").read_text()
                self.assertIn('Prototype non validé', html)
                self.assertNotIn('Programme Officiel', html)
                self.assertNotIn('https://', html)
                self.assertNotIn('__TITLE__', html)

    def test_stale_inputs_fail_before_external_commands(self):
        with tempfile.TemporaryDirectory() as directory, patch('pipeline.subprocess.run') as run:
            with self.assertRaises(SystemExit):
                render(Path(directory), 'msa')
            run.assert_not_called()

    def test_msa_uses_shared_slide_order_and_selected_directory(self):
        with tempfile.TemporaryDirectory(prefix='video test ') as directory:
            root = Path(directory)
            with patch('pipeline.require_tool', return_value='edge-tts'), patch('pipeline.subprocess.run') as run:
                narration(root, 'msa')
                self.assertEqual(run.call_count, len(SLIDES))
                for slide, call in zip(SLIDES, run.call_args_list):
                    command = call.args[0]
                    self.assertEqual(command[command.index('--text')+1], slide['msa_text'])
                    self.assertEqual(command[-1], str(root / 'msa' / f"{slide['id']}.mp3"))
            self.assertEqual((root / 'msa' / 'audio.sha256').read_text(), fingerprint())

    def test_render_pairs_matching_image_and_audio(self):
        with tempfile.TemporaryDirectory(prefix='video test ') as directory:
            root = Path(directory)
            target = root / 'msa'
            target.mkdir()
            (root / 'images.sha256').write_text(fingerprint())
            (target / 'audio.sha256').write_text(fingerprint())
            with patch('pipeline.require_tool', return_value='ffmpeg'), patch('pipeline.subprocess.run') as run:
                render(root, 'msa')
                for slide, call in zip(SLIDES, run.call_args_list):
                    command = call.args[0]
                    self.assertIn(str(root / f"{slide['id']}.png"), command)
                    self.assertIn(str(target / f"{slide['id']}.mp3"), command)
                self.assertEqual(run.call_args_list[-1].args[0][-1], str(target / 'unreviewed-lesson.mp4'))


if __name__ == '__main__':
    unittest.main()
